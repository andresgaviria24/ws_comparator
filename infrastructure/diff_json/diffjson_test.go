package diffjson

import (
	"encoding/json"
	"fmt"
	"testing"
	"ws_comparator/domain/dto"

	"github.com/stretchr/testify/assert"
)

func TestObjects(t *testing.T) {

	tests := []struct {
		name       string
		argsResult map[string]interface{}
		argsIn     dto.ComparatorIn
		want       string
	}{
		{
			name:       "Identical",
			argsResult: DataInJson(),
			argsIn:     dto.ComparatorIn{ResponseC: DataInJson()},
			want:       diff_string_answer,
		},
		{
			name:       "Diff string",
			argsResult: DataInJson(),
			argsIn:     dto.ComparatorIn{ResponseC: DataStringJson()},
			want:       identical_answer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, deltaJson := DiffJson(tt.argsResult, tt.argsIn)
			assert.Equal(t, deltaJson, tt.want)
		})
	}
}

func DataInJson() map[string]interface{} {
	var data map[string]interface{}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error al convertir JSON:", err)
		return nil
	}

	return data
}

var jsonData = []byte(`
{
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
	}
}
`)

var identical_answer = " {\n   \"response\": {\n     \"route\": {\n       \"cityToCity\": false,\n       \"destinationCoordinates\": {\n         \"latitude\": -33.439852,\n         \"longitude\": -70.6363879\n       },\n       \"destinationZipcode\": \"\",\n       \"hasCoverage\": false,\n       \"hasP99Coverage\": false,\n       \"id\": \"\",\n       \"originCoordinates\": {\n         \"latitude\": -33.439852,\n         \"longitude\": -70.6363879\n       },\n       \"originZipcode\": \"\",\n       \"p99CoveragePoint\": null,\n       \"stations\": [\n         {\n           \"code\": \"\",\n           \"country\": \"CHL\",\n           \"coverage\": \"RM.107\",\n           \"station\": \"RMRU\",\n           \"substationPrefix\": \"\"\n         },\n         {\n           \"code\": \"\",\n           \"country\": \"CHL\",\n           \"coverage\": \"RM.107\",\n           \"station\": \"RMRU\",\n           \"substationPrefix\": \"\"\n         }\n       ]\n     },\n     \"verbose\": {\n       \"strategy\": \"StandardCoverage\"\n     }\n   }\n }\n"

func DataStringJson() map[string]interface{} {
	var data map[string]interface{}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error al convertir JSON:", err)
		return nil
	}

	return data
}

var jsonStringData = []byte(`
{
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
					"station": "BBBB",
					"coverage": "RM.107",
					"code": "",
					"substationPrefix": "",
					"country": "CHL"
				},
				{
					"station": "AAAA",
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
	}
}
`)

var diff_string_answer = " {\n   \"response\": {\n     \"route\": {\n       \"cityToCity\": false,\n       \"destinationCoordinates\": {\n         \"latitude\": -33.439852,\n         \"longitude\": -70.6363879\n       },\n       \"destinationZipcode\": \"\",\n       \"hasCoverage\": false,\n       \"hasP99Coverage\": false,\n       \"id\": \"\",\n       \"originCoordinates\": {\n         \"latitude\": -33.439852,\n         \"longitude\": -70.6363879\n       },\n       \"originZipcode\": \"\",\n       \"p99CoveragePoint\": null,\n       \"stations\": [\n         {\n           \"code\": \"\",\n           \"country\": \"CHL\",\n           \"coverage\": \"RM.107\",\n           \"station\": \"RMRU\",\n           \"substationPrefix\": \"\"\n         },\n         {\n           \"code\": \"\",\n           \"country\": \"CHL\",\n           \"coverage\": \"RM.107\",\n           \"station\": \"RMRU\",\n           \"substationPrefix\": \"\"\n         }\n       ]\n     },\n     \"verbose\": {\n       \"strategy\": \"StandardCoverage\"\n     }\n   }\n }\n"
