package application

import (
	"net/http"
	"ws_comparator/domain/drip"
	"ws_comparator/domain/dto"
	"ws_comparator/domain/service"
	"ws_comparator/utils"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ComparatorController struct {
	comparatorService service.ComparatorService
	newrelicClient    *newrelic.Application
}

func InitComparatorController(router *gin.Engine, newrelicClient *newrelic.Application) {
	comparatorController := ComparatorController{
		comparatorService: service.InitComparatorServiceImpl(),
		newrelicClient:    newrelicClient,
	}
	router.POST("v1/request/comparation", comparatorController.ComparatorHandler)
}

func (r *ComparatorController) ComparatorHandler(c *gin.Context) {
	var (
		comparator dto.ComparatorIn
		response   = dto.Response{
			Status: http.StatusOK,
		}
	)

	if !drip.Drip(c, &utils.RealUtils{}) {
		return
	}

	if err := c.ShouldBindJSON(&comparator); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.Response{})
		return
	}

	if response = r.comparatorService.Comparator(comparator, r.newrelicClient); response.Status != http.StatusOK {
		c.JSON(response.Status, response)
	}

	c.JSON(http.StatusOK, response)

}
