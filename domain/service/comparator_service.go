package service

import "ws_comparator/domain/dto"

type ComparatorService interface {
	Comparator(dto.ComparatorIn) dto.Response
}
