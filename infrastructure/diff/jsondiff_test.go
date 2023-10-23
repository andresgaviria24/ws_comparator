package diff

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	differ := New()
	if differ == nil {
		t.Error("Expected a non-nil Differ, but got nil")
	}
}

func TestDiffer_CompareObjects(t *testing.T) {
	differ := New()
	left := map[string]interface{}{"key1": "value1", "key2": "value2"}
	right := map[string]interface{}{"key1": "value1", "key3": "value3"}

	diff := differ.CompareObjects(left, right)
	if !diff.Modified() {
		t.Error("Expected Modified() to be true, but it's false")
	}
}

func TestDiffer_CompareArrays(t *testing.T) {
	differ := New()
	left := []interface{}{1, 2, 3}
	right := []interface{}{1, 2, 4}

	diff := differ.CompareArrays(left, right)
	if !diff.Modified() {
		t.Error("Expected Modified() to be true, but it's false")
	}
}

func TestDiffer_Compare(t *testing.T) {
	differ := New()
	leftJSON := []byte(`{"key1": "value1", "key2": "value2"}`)
	rightJSON := []byte(`{"key1": "value1", "key3": "value3"}`)

	diff, err := differ.Compare(leftJSON, rightJSON)
	if err != nil {
		t.Errorf("Error comparing JSON: %v", err)
	}

	if !diff.Modified() {
		t.Error("Expected Modified() to be true, but it's false")
	}
}

func TestStringToInterfaceSlice(t *testing.T) {
	str := "hello"
	expected := []interface{}{'h', 'e', 'l', 'l', 'o'}
	result := stringToInterfaceSlice(str)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]interface{}{"b": 2, "a": 1, "c": 3}
	expected := []string{"a", "b", "c"}
	keys := sortedKeys(m)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, but got %v", expected, keys)
	}
}
