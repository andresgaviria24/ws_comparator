package service

import (
	"ws_comparator/domain/dto"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type ComparatorService interface {
	Comparator(dto.ComparatorIn, *newrelic.Application) dto.Response
}
