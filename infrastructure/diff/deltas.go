package diff

import (
	"errors"
	"reflect"
	"strconv"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

type Delta interface {
	Similarity() (similarity float64)
}

type similariter interface {
	similarity() (similarity float64)
}

type similarityCache struct {
	similariter
	value float64
}

func newSimilarityCache(sim similariter) similarityCache {
	cache := similarityCache{similariter: sim, value: -1}
	return cache
}

func (cache similarityCache) Similarity() (similarity float64) {
	if cache.value < 0 {
		cache.value = cache.similariter.similarity()
	}
	return cache.value
}

type Position interface {
	String() (name string)

	CompareTo(another Position) bool
}

type Name string

func (n Name) String() (name string) {
	return string(n)
}

func (n Name) CompareTo(another Position) bool {
	return n < another.(Name)
}

type Index int

func (i Index) String() (name string) {
	return strconv.Itoa(int(i))
}

func (i Index) CompareTo(another Position) bool {
	return i < another.(Index)
}

type PreDelta interface {
	PrePosition() Position

	PreApply(object interface{}) interface{}
}

type preDelta struct{ Position }

func (i preDelta) PrePosition() Position {
	return Position(i.Position)
}

type preDeltas []PreDelta

func (s preDeltas) Len() int {
	return len(s)
}

func (s preDeltas) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s preDeltas) Less(i, j int) bool {
	return !s[i].PrePosition().CompareTo(s[j].PrePosition())
}

type PostDelta interface {
	PostPosition() Position

	PostApply(object interface{}) interface{}
}

type postDelta struct{ Position }

func (i postDelta) PostPosition() Position {
	return Position(i.Position)
}

type postDeltas []PostDelta

func (s postDeltas) Len() int {
	return len(s)
}

func (s postDeltas) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s postDeltas) Less(i, j int) bool {
	return s[i].PostPosition().CompareTo(s[j].PostPosition())
}

type Object struct {
	postDelta
	similarityCache

	Deltas []Delta
}

func NewObject(position Position, deltas []Delta) *Object {
	d := Object{postDelta: postDelta{position}, Deltas: deltas}
	d.similarityCache = newSimilarityCache(&d)
	return &d
}

func (d *Object) PostApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		o := object.(map[string]interface{})
		n := string(d.PostPosition().(Name))
		o[n] = applyDeltas(d.Deltas, o[n])
	case []interface{}:
		o := object.([]interface{})
		n := int(d.PostPosition().(Index))
		o[n] = applyDeltas(d.Deltas, o[n])
	}
	return object
}

func (d *Object) similarity() (similarity float64) {
	similarity = deltasSimilarity(d.Deltas)
	return
}

type Array struct {
	postDelta
	similarityCache

	Deltas []Delta
}

func NewArray(position Position, deltas []Delta) *Array {
	d := Array{postDelta: postDelta{position}, Deltas: deltas}
	d.similarityCache = newSimilarityCache(&d)
	return &d
}

func (d *Array) PostApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		o := object.(map[string]interface{})
		n := string(d.PostPosition().(Name))
		o[n] = applyDeltas(d.Deltas, o[n])
	case []interface{}:
		o := object.([]interface{})
		n := int(d.PostPosition().(Index))
		o[n] = applyDeltas(d.Deltas, o[n])
	}
	return object
}

func (d *Array) similarity() (similarity float64) {
	similarity = deltasSimilarity(d.Deltas)
	return
}

type Added struct {
	postDelta
	similarityCache

	Value interface{}
}

func NewAdded(position Position, value interface{}) *Added {
	d := Added{postDelta: postDelta{position}, Value: value}
	return &d
}

func (d *Added) PostApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		object.(map[string]interface{})[string(d.PostPosition().(Name))] = d.Value
	case []interface{}:
		i := int(d.PostPosition().(Index))
		o := object.([]interface{})
		if i < len(o) {
			o = append(o, 0) //dummy
			copy(o[i+1:], o[i:])
			o[i] = d.Value
			object = o
		} else {
			object = append(o, d.Value)
		}
	}

	return object
}

func (d *Added) similarity() (similarity float64) {
	return 0
}

type Modified struct {
	postDelta
	similarityCache

	OldValue interface{}

	NewValue interface{}
}

func NewModified(position Position, oldValue, newValue interface{}) *Modified {
	d := Modified{
		postDelta: postDelta{position},
		OldValue:  oldValue,
		NewValue:  newValue,
	}
	d.similarityCache = newSimilarityCache(&d)
	return &d

}

func (d *Modified) PostApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		object.(map[string]interface{})[string(d.PostPosition().(Name))] = d.NewValue
	case []interface{}:
		object.([]interface{})[int(d.PostPosition().(Index))] = d.NewValue
	}
	return object
}

func (d *Modified) similarity() (similarity float64) {
	similarity += 0.3
	if reflect.TypeOf(d.OldValue) == reflect.TypeOf(d.NewValue) {
		similarity += 0.3

		switch d.OldValue.(type) {
		case string:
			similarity += 0.4 * stringSimilarity(d.OldValue.(string), d.NewValue.(string))
		case float64:
			ratio := d.OldValue.(float64) / d.NewValue.(float64)
			if ratio > 1 {
				ratio = 1 / ratio
			}
			similarity += 0.4 * ratio
		}
	}
	return
}

type TextDiff struct {
	Modified

	Diff []dmp.Patch
}

func NewTextDiff(position Position, diff []dmp.Patch, oldValue, newValue interface{}) *TextDiff {
	d := TextDiff{
		Modified: *NewModified(position, oldValue, newValue),
		Diff:     diff,
	}
	return &d
}

func (d *TextDiff) PostApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		o := object.(map[string]interface{})
		i := string(d.PostPosition().(Name))
		d.OldValue = o[i]
		d.patch()
		o[i] = d.NewValue
	case []interface{}:
		o := object.([]interface{})
		i := d.PostPosition().(Index)
		d.OldValue = o[i]
		d.patch()
		o[i] = d.NewValue
	}
	return object
}

func (d *TextDiff) patch() error {
	if d.OldValue == nil {
		return errors.New("Old Value is not set")
	}
	patcher := dmp.New()
	patched, successes := patcher.PatchApply(d.Diff, d.OldValue.(string))
	for _, success := range successes {
		if !success {
			return errors.New("Failed to apply a patch")
		}
	}
	d.NewValue = patched
	return nil
}

func (d *TextDiff) DiffString() string {
	dmp := dmp.New()
	return dmp.PatchToText(d.Diff)
}

type Deleted struct {
	preDelta

	Value interface{}
}

func NewDeleted(position Position, value interface{}) *Deleted {
	d := Deleted{
		preDelta: preDelta{position},
		Value:    value,
	}
	return &d

}

func (d *Deleted) PreApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
		delete(object.(map[string]interface{}), string(d.PrePosition().(Name)))
	case []interface{}:
		i := int(d.PrePosition().(Index))
		o := object.([]interface{})
		object = append(o[:i], o[i+1:]...)
	}
	return object
}

func (d Deleted) Similarity() (similarity float64) {
	return 0
}

type Moved struct {
	preDelta
	postDelta
	similarityCache

	Value interface{}
	Delta interface{}
}

func NewMoved(oldPosition Position, newPosition Position, value interface{}, delta Delta) *Moved {
	d := Moved{
		preDelta:  preDelta{oldPosition},
		postDelta: postDelta{newPosition},
		Value:     value,
		Delta:     delta,
	}
	d.similarityCache = newSimilarityCache(&d)
	return &d
}

func (d *Moved) PreApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
	case []interface{}:
		i := int(d.PrePosition().(Index))
		o := object.([]interface{})
		d.Value = o[i]
		object = append(o[:i], o[i+1:]...)
	}
	return object
}

func (d *Moved) PostApply(object interface{}) interface{} {
	switch object.(type) {
	case map[string]interface{}:
	case []interface{}:
		i := int(d.PostPosition().(Index))
		o := object.([]interface{})
		o = append(o, 0) //dummy
		copy(o[i+1:], o[i:])
		o[i] = d.Value
		object = o
	}

	if d.Delta != nil {
		d.Delta.(PostDelta).PostApply(object)
	}

	return object
}

func (d *Moved) similarity() (similarity float64) {
	similarity = 0.6
	ratio := float64(d.PrePosition().(Index)) / float64(d.PostPosition().(Index))
	if ratio > 1 {
		ratio = 1 / ratio
	}
	similarity += 0.4 * ratio
	return
}
