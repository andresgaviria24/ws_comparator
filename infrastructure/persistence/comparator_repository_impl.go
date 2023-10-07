package persistence

import (
	"ws_comparator/domain/entity"

	"github.com/jinzhu/gorm"
)

type ComparatorRepositoryImpl struct {
	db *gorm.DB
}

func (r *ComparatorRepositoryImpl) GetConfiguration() ([]entity.Food, error) {

	var listUserDetail []entity.Food
	err := r.db.Find(&listUserDetail).Error

	return listUserDetail, err
}
