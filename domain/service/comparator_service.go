package service

import "ws_comparator/domain/dto"

type RestaurantService interface {
	GetFoods() ([]dto.FoodDto, dto.Response)
}
