package diff

import (
	"container/list"
	"encoding/json"
	"reflect"
	"sort"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
	lcs "github.com/yudai/golcs"
)

type Diff interface {
	Deltas() []Delta

	Modified() bool
}

type diff struct {
	deltas []Delta
}

func (diff *diff) Deltas() []Delta {
	return diff.deltas
}

func (diff *diff) Modified() bool {
	return len(diff.deltas) > 0
}

type Differ struct {
	textDiffMinimumLength int
}

func New() *Differ {
	return &Differ{
		textDiffMinimumLength: 30,
	}
}

func (differ *Differ) Compare(
	left []byte,
	right []byte,
) (Diff, error) {
	var leftMap, rightMap map[string]interface{}
	err := json.Unmarshal(left, &leftMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(right, &rightMap)
	if err != nil {
		return nil, err
	}
	return differ.CompareObjects(leftMap, rightMap), nil
}

func (differ *Differ) CompareObjects(
	left map[string]interface{},
	right map[string]interface{},
) Diff {
	deltas := differ.compareMaps(left, right)
	return &diff{deltas: deltas}
}

func (differ *Differ) CompareArrays(
	left []interface{},
	right []interface{},
) Diff {
	deltas := differ.compareArrays(left, right)
	return &diff{deltas: deltas}
}

func (differ *Differ) compareMaps(
	left map[string]interface{},
	right map[string]interface{},
) (deltas []Delta) {
	deltas = make([]Delta, 0)

	names := sortedKeys(left)
	for _, name := range names {
		if rightValue, ok := right[name]; ok {
			same, delta := differ.compareValues(Name(name), left[name], rightValue)
			if !same {
				deltas = append(deltas, delta)
			}
		} else {
			deltas = append(deltas, NewDeleted(Name(name), left[name]))
		}
	}

	names = sortedKeys(right)
	for _, name := range names {
		if _, ok := left[name]; !ok {
			deltas = append(deltas, NewAdded(Name(name), right[name]))
		}
	}

	return deltas
}

func (differ *Differ) ApplyPatch(json map[string]interface{}, patch Diff) {
	applyDeltas(patch.Deltas(), json)
}

type maybe struct {
	index    int
	lcsIndex int
	item     interface{}
}

func (differ *Differ) compareArrays(
	left []interface{},
	right []interface{},
) (deltas []Delta) {
	deltas = make([]Delta, 0)
	lcsPairs := lcs.New(left, right).IndexPairs()

	maybeDeleted := list.New()
	lcsI := 0
	for i, leftValue := range left {
		if lcsI < len(lcsPairs) && lcsPairs[lcsI].Left == i {
			lcsI++
		} else {
			maybeDeleted.PushBack(maybe{index: i, lcsIndex: lcsI, item: leftValue})
		}
	}

	maybeAdded := list.New()
	lcsI = 0
	for i, rightValue := range right {
		if lcsI < len(lcsPairs) && lcsPairs[lcsI].Right == i {
			lcsI++
		} else {
			maybeAdded.PushBack(maybe{index: i, lcsIndex: lcsI, item: rightValue})
		}
	}

	var delNext *list.Element
	for delCandidate := maybeDeleted.Front(); delCandidate != nil; delCandidate = delNext {
		delCan := delCandidate.Value.(maybe)
		delNext = delCandidate.Next()

		for addCandidate := maybeAdded.Front(); addCandidate != nil; addCandidate = addCandidate.Next() {
			addCan := addCandidate.Value.(maybe)
			if reflect.DeepEqual(delCan.item, addCan.item) {
				deltas = append(deltas, NewMoved(Index(delCan.index), Index(addCan.index), delCan.item, nil))
				maybeAdded.Remove(addCandidate)
				maybeDeleted.Remove(delCandidate)
				break
			}
		}
	}

	prevIndexDel := 0
	prevIndexAdd := 0
	delElement := maybeDeleted.Front()
	addElement := maybeAdded.Front()
	for i := 0; i <= len(lcsPairs); i++ { // not "< len(lcsPairs)"
		var lcsPair lcs.IndexPair
		var delSize, addSize int
		if i < len(lcsPairs) {
			lcsPair = lcsPairs[i]
			delSize = lcsPair.Left - prevIndexDel - 1
			addSize = lcsPair.Right - prevIndexAdd - 1
			prevIndexDel = lcsPair.Left
			prevIndexAdd = lcsPair.Right
		}

		var delSlice []maybe
		if delSize > 0 {
			delSlice = make([]maybe, 0, delSize)
		} else {
			delSlice = make([]maybe, 0, maybeDeleted.Len())
		}
		for ; delElement != nil; delElement = delElement.Next() {
			d := delElement.Value.(maybe)
			if d.lcsIndex != i {
				break
			}
			delSlice = append(delSlice, d)
		}

		var addSlice []maybe
		if addSize > 0 {
			addSlice = make([]maybe, 0, addSize)
		} else {
			addSlice = make([]maybe, 0, maybeAdded.Len())
		}
		for ; addElement != nil; addElement = addElement.Next() {
			a := addElement.Value.(maybe)
			if a.lcsIndex != i {
				break
			}
			addSlice = append(addSlice, a)
		}

		if len(delSlice) > 0 && len(addSlice) > 0 {
			var bestDeltas []Delta
			bestDeltas, delSlice, addSlice = differ.maximizeSimilarities(delSlice, addSlice)
			for _, delta := range bestDeltas {
				deltas = append(deltas, delta)
			}
		}

		for _, del := range delSlice {
			deltas = append(deltas, NewDeleted(Index(del.index), del.item))
		}
		for _, add := range addSlice {
			deltas = append(deltas, NewAdded(Index(add.index), add.item))
		}
	}

	return deltas
}

func (differ *Differ) compareValues(
	position Position,
	left interface{},
	right interface{},
) (same bool, delta Delta) {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false, NewModified(position, left, right)
	}

	switch left.(type) {

	case map[string]interface{}:
		l := left.(map[string]interface{})
		childDeltas := differ.compareMaps(l, right.(map[string]interface{}))
		if len(childDeltas) > 0 {
			return false, NewObject(position, childDeltas)
		}

	case []interface{}:
		l := left.([]interface{})
		childDeltas := differ.compareArrays(l, right.([]interface{}))

		if len(childDeltas) > 0 {
			return false, NewArray(position, childDeltas)
		}

	default:
		if !reflect.DeepEqual(left, right) {

			if reflect.ValueOf(left).Kind() == reflect.String &&
				reflect.ValueOf(right).Kind() == reflect.String &&
				differ.textDiffMinimumLength <= len(left.(string)) {

				textDiff := dmp.New()
				patchs := textDiff.PatchMake(left.(string), right.(string))
				return false, NewTextDiff(position, patchs, left, right)

			} else {
				return false, NewModified(position, left, right)
			}
		}
	}

	return true, nil
}

func applyDeltas(deltas []Delta, object interface{}) interface{} {
	preDeltas := make(preDeltas, 0)
	for _, delta := range deltas {
		switch delta.(type) {
		case PreDelta:
			preDeltas = append(preDeltas, delta.(PreDelta))
		}
	}
	sort.Sort(preDeltas)
	for _, delta := range preDeltas {
		object = delta.PreApply(object)
	}

	postDeltas := make(postDeltas, 0, len(deltas)-len(preDeltas))
	for _, delta := range deltas {
		switch delta.(type) {
		case PostDelta:
			postDeltas = append(postDeltas, delta.(PostDelta))
		}
	}
	sort.Sort(postDeltas)

	for _, delta := range postDeltas {
		object = delta.PostApply(object)
	}

	return object
}

func (differ *Differ) maximizeSimilarities(left []maybe, right []maybe) (resultDeltas []Delta, freeLeft, freeRight []maybe) {
	deltaTable := make([][]Delta, len(left))
	for i := 0; i < len(left); i++ {
		deltaTable[i] = make([]Delta, len(right))
	}
	for i, leftValue := range left {
		for j, rightValue := range right {
			_, delta := differ.compareValues(Index(rightValue.index), leftValue.item, rightValue.item)
			deltaTable[i][j] = delta
		}
	}

	sizeX := len(left) + 1 // margins for both sides
	sizeY := len(right) + 1

	dpTable := make([][]float64, sizeX)
	for i := 0; i < sizeX; i++ {
		dpTable[i] = make([]float64, sizeY)
	}
	for x := sizeX - 2; x >= 0; x-- {
		for y := sizeY - 2; y >= 0; y-- {
			prevX := dpTable[x+1][y]
			prevY := dpTable[x][y+1]
			score := deltaTable[x][y].Similarity() + dpTable[x+1][y+1]

			dpTable[x][y] = max(prevX, prevY, score)
		}
	}

	minLength := len(left)
	if minLength > len(right) {
		minLength = len(right)
	}
	maxInvalidLength := minLength - 1

	freeLeft = make([]maybe, 0, len(left)-minLength)
	freeRight = make([]maybe, 0, len(right)-minLength)

	resultDeltas = make([]Delta, 0, minLength)
	var x, y int
	for x, y = 0, 0; x <= sizeX-2 && y <= sizeY-2; {
		current := dpTable[x][y]
		nextX := dpTable[x+1][y]
		nextY := dpTable[x][y+1]

		xValidLength := len(left) - maxInvalidLength + y
		yValidLength := len(right) - maxInvalidLength + x

		if x+1 < xValidLength && current == nextX {
			freeLeft = append(freeLeft, left[x])
			x++
		} else if y+1 < yValidLength && current == nextY {
			freeRight = append(freeRight, right[y])
			y++
		} else {
			resultDeltas = append(resultDeltas, deltaTable[x][y])
			x++
			y++
		}
	}
	for ; x < sizeX-1; x++ {
		freeLeft = append(freeLeft, left[x-1])
	}
	for ; y < sizeY-1; y++ {
		freeRight = append(freeRight, right[y-1])
	}

	return resultDeltas, freeLeft, freeRight
}

func deltasSimilarity(deltas []Delta) (similarity float64) {
	for _, delta := range deltas {
		similarity += delta.Similarity()
	}
	similarity = similarity / float64(len(deltas))
	return
}

func stringSimilarity(left, right string) (similarity float64) {
	matchingLength := float64(
		lcs.New(
			stringToInterfaceSlice(left),
			stringToInterfaceSlice(right),
		).Length(),
	)
	similarity =
		(matchingLength / float64(len(left))) * (matchingLength / float64(len(right)))
	return
}

func stringToInterfaceSlice(str string) []interface{} {
	s := make([]interface{}, len(str))
	for i, v := range str {
		s[i] = v
	}
	return s
}

func sortedKeys(m map[string]interface{}) (keys []string) {
	keys = make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func max(first float64, rest ...float64) (max float64) {
	max = first
	for _, value := range rest {
		if max < value {
			max = value
		}
	}
	return max
}
