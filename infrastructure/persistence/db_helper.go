package persistence

import (
	"os"
	"ws_comparator/infrastructure/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DbHelper struct {
	ComparatorRepository repository.ComparatorRepository
	db                   *gorm.DB
}

func InitDbHelper() (*DbHelper, error) {

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate()
	return &DbHelper{
		ComparatorRepository: &ComparatorRepositoryImpl{db},
		db:                   db,
	}, nil
}
