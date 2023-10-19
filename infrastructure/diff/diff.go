package diff

import (
	"ws_comparator/domain/dto"

	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func Diff(result map[string]interface{}, comparator dto.ComparatorIn) (gojsondiff.Diff, string) {
	differ := gojsondiff.New()
	diff := differ.CompareObjects(comparator.ResponseC, result)

	f := formatter.NewAsciiFormatter(result, formatter.AsciiFormatterDefaultConfig)
	deltaJSON, _ := f.Format(diff)
	return diff, deltaJSON
}
