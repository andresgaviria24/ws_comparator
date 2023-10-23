package diffjson

import (
	"ws_comparator/domain/dto"
	"ws_comparator/infrastructure/diff"
	"ws_comparator/infrastructure/format"
)

func DiffJson(result map[string]interface{}, comparator dto.ComparatorIn) (diff.Diff, string) {
	differ := diff.New()
	diff := differ.CompareObjects(comparator.ResponseC, result)

	f := format.NewAsciiFormatter(result, format.AsciiFormatterDefaultConfig)
	deltaJSON, _ := f.Format(diff)
	return diff, deltaJSON
}
