package dto

import "github.com/yudai/gojsondiff"

type Response struct {
	Status      int                `json:"status"`
	Description []gojsondiff.Delta `json:"description"`
	Message     string             `json:"message"`
}
