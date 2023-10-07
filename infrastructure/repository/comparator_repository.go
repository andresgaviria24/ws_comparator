package repository

import (
	"ws_comparator/domain/entity"
)

type ComparatorRepository interface {
	GetConfiguration() ([]entity.Food, error)
}
