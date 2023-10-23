package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"ws_comparator/domain/dto"
	diffjson "ws_comparator/infrastructure/diff_json"
	newrelic_metrics "ws_comparator/infrastructure/newrelic"
	"ws_comparator/infrastructure/persistence"
	"ws_comparator/infrastructure/repository"

	"github.com/go-resty/resty/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ComparatorServiceImpl struct {
	comparatorRepository repository.ComparatorRepository
}

func InitComparatorServiceImpl() *ComparatorServiceImpl {
	dbHelper, err := persistence.InitDbHelper()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &ComparatorServiceImpl{
		comparatorRepository: dbHelper.ComparatorRepository,
	}
}

func (r *ComparatorServiceImpl) Comparator(comparator dto.ComparatorIn, app *newrelic.Application) dto.Response {
	var response dto.Response

	payload := make(map[string]interface{})
	for k, b := range comparator.Body {
		payload[k] = b
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(comparator.Url)

	if err != nil {
		fmt.Println("Error realizando la solicitud:", err)
		return responseWithError(response, app)
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return responseWithError(response, app)
	}

	diff, deltaJSON := diffjson.DiffJson(result, comparator)

	newrelic_metrics.SendMetric(diff, app, deltaJSON)

	response.Status = http.StatusOK
	response.Description = diff.Deltas()
	return response
}

func responseWithError(response dto.Response, app *newrelic.Application) dto.Response {
	response.Status = http.StatusInternalServerError
	app.RecordCustomMetric(os.Getenv("SERVICE_NAME")+".error_request", 1)
	return response
}
