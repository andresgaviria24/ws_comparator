package dto

import (
	"errors"
	"net/http"
)

type ComparatorIn struct {
	Method      string                 `json:"method"`
	Body        map[string]interface{} `json:"body"`
	QueryParams map[string]string      `json:"query_params"`
	ResponseC   map[string]interface{} `json:"response"`
	Url         string                 `json:"url"`
}

func (c *ComparatorIn) Validator() error {
	if c.Method != http.MethodPost && c.Method != http.MethodGet {
		return errors.New("método HTTP no permitido")
	}

	/*if len(c.Body) == 0 {
		return errors.New("la solicitud debe contener un cuerpo")
	}*/
	return nil
}
