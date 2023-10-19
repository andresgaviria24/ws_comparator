package application

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"ws_comparator/domain/dto"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
)

type MockComparatorService struct{}

func (s *MockComparatorService) Comparator(comparator dto.ComparatorIn, app *newrelic.Application) dto.Response {
	return dto.Response{
		Status: http.StatusOK,
	}
}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(m.Run())
}

func TestComparatorHandlerSuccess(t *testing.T) {
	w := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, w.Code)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	c.Request, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(comparatorJSON))

	newrelicApp := &newrelic.Application{}

	controller := &ComparatorController{
		comparatorService: &MockComparatorService{},
		newrelicClient:    newrelicApp,
	}

	controller.ComparatorHandler(c)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestComparatorHandlerBody(t *testing.T) {
	w := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, w.Code)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	c.Request, _ = http.NewRequest(http.MethodPost, "", nil)

	newrelicApp := &newrelic.Application{}

	controller := &ComparatorController{
		comparatorService: &MockComparatorService{},
		newrelicClient:    newrelicApp,
	}

	controller.ComparatorHandler(c)
	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
}

/*
func TestComparatorHandler_InvalidJSON(t *testing.T) {
	// Preparar un router de prueba de Gin
	//router := gin.Default()

	// Configurar un contexto de prueba
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Crear un JSON de solicitud no válido
	invalidJSON := `{"invalid_json": "}`

	// Configurar la solicitud con el JSON no válido
	req, _ := http.NewRequest("POST", "/v1/request/comparation", strings.NewReader(invalidJSON))
	c.Request = req

	// Inicializar un mock de New Relic Application
	newrelicApp := &newrelic.Application{}

	// Configurar el controlador con el servicio de mock y la aplicación de New Relic de mock
	controller := &ComparatorController{
		comparatorService: &MockComparatorService{},
		newrelicClient:    newrelicApp,
	}

	// Ejecutar el controlador
	controller.ComparatorHandler(c)

	// Verificar que la respuesta sea un código 422 (Unprocessable Entity) debido al JSON no válido
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
*/

var comparatorJSON = `{
	"method": "POST",
	"body": {
		"country": "CHL",
		"deliveryType": "SMD",
		"originCoordinates": {
			"longitude": -70.6363879,
			"latitude": -33.439852
		},
		"destinationCoordinates": {
			"longitude": -70.6363879,
			"latitude": -33.439852
		},
		"useCoordinates": true,
		"client_id": "5470b131-c395-4d3b-8027-584efee7e8ab"
	},
	"query_params": {
		"param1": "valor1",
		"param2": "valor2"
	},
	"response": {
		"route": {
			"id": "",
			"originZipcode": "",
			"originCoordinates": {
				"latitude": -33.439852,
				"longitude": -70.6363879
			},
			"destinationZipcode": "",
			"destinationCoordinates": {
				"latitude": -33.439852,
				"longitude": -70.6363879
			},
			"stations": [
				{
					"station": "RMRU",
					"coverage": "RM.107",
					"code": "",
					"substationPrefix": "",
					"country": "CHL"
				},
				{
					"station": "RMRU",
					"coverage": "RM.107",
					"code": "",
					"substationPrefix": "",
					"country": "CHL"
				}
			],
			"hasCoverage": false,
			"hasP99Coverage": false,
			"cityToCity": false,
			"p99CoveragePoint": null
		},
		"verbose": {
			"strategy": "StandardCoverage"
		}
	},
	"url": "https://coverage-microservice-qndxoltwga-uc.a.run.app/coverage.CoverageService/GetRoute"
}`
