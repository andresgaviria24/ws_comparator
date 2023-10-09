package application

import (
	"net/http"
	"ws_comparator/domain/dto"
	"ws_comparator/domain/service"

	"github.com/gin-gonic/gin"
)

type ComparatorController struct {
	comparatorService service.ComparatorService
}

func InitComparatorController(router *gin.Engine) {
	comparatorController := ComparatorController{
		comparatorService: service.InitComparatorServiceImpl(),
	}
	router.POST("v1/request/comparation", comparatorController.ComparatorHandler)
}

func (r *ComparatorController) ComparatorHandler(c *gin.Context) {

	var (
		response   dto.Response
		comparator dto.ComparatorIn
	)

	if err := c.ShouldBindJSON(&comparator); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.Response{})
		return
	}

	response = r.comparatorService.Comparator(comparator)

	if response.Status != http.StatusOK {
		c.JSON(response.Status, response)
		return
	}
	c.JSON(response.Status, response)
}
