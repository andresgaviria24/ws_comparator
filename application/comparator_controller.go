package application

import (
	"net/http"
	"ws_comparator/domain/dto"
	"ws_comparator/domain/service"

	"github.com/gin-gonic/gin"
)

type ComparatorController struct {
	restaurantService service.RestaurantService
}

func InitComparatorController(router *gin.Engine) {
	comparatorController := ComparatorController{
		restaurantService: service.InitComparatorServiceImpl(),
	}
	router.POST("v1/request/configuration", comparatorController.ComparatorHandler)
}

func (r *ComparatorController) ComparatorHandler(c *gin.Context) {

	var response dto.Response

	foods, response := r.restaurantService.GetFoods()

	if response.Status != http.StatusOK {
		c.JSON(response.Status, response)
		return
	}
	c.JSON(response.Status, foods)
}
