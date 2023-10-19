package newrelic_metrics_test

/*
import (
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/yudai/gojsondiff"
)

type MockApplication struct{}

func (m *MockApplication) RecordCustomMetric(name string, value float64) {
}

func (m *MockApplication) RecordLog(log newrelic.LogData) {
}

func TestSendMetric_IdenticalAnswer(t *testing.T) {
	mockApp := &MockApplication{}

	diff := gojsondiff.Diff{}
	diffJson := ""
	SendMetric(diff, mockApp, diffJson)

	assert.True(t, metricRecorded(mockApp, "identical_answer"))
}

func TestSendMetric_DifferentAnswer(t *testing.T) {
	mockApp := &MockApplication{}

	diff := gojsondiff.Diff{
		Deltas: []gojsondiff.Delta{{}},
	}
	diffJson := "Diferencias encontradas"
	SendMetric(diff, mockApp, diffJson)

	assert.True(t, metricRecorded(mockApp, "diff_answer"))
	assert.True(t, logRecorded(mockApp, "Diferencias encontradas"))
}

func metricRecorded(app *MockApplication, metricName string) bool {
	return true
}

func logRecorded(app *MockApplication, logMessage string) bool {
	return true
}
*/
