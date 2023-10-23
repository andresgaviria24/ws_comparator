package dto

import "ws_comparator/infrastructure/diff"

type Response struct {
	Status      int          `json:"status"`
	Description []diff.Delta `json:"description"`
	Message     string       `json:"message"`
}
