package newrelic_metrics

import (
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/yudai/gojsondiff"
)

func SendMetric(diff gojsondiff.Diff, app *newrelic.Application, diffJson string) {
	if len(diff.Deltas()) == 0 {
		app.RecordCustomMetric(os.Getenv("SERVICE_NAME")+".identical_answer", 1)
	} else {
		app.RecordLog(newrelic.LogData{
			Severity: "error",
			Message:  diffJson,
		})
		app.RecordCustomMetric(os.Getenv("SERVICE_NAME")+".diff_answer", 1)
	}

	app.RecordCustomMetric(os.Getenv("SERVICE_NAME")+".total_request", 1)
}
