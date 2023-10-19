package drip

import (
	"math/rand"
	"net/http"
	"time"
	"ws_comparator/domain/dto"
	"ws_comparator/utils"

	"github.com/gin-gonic/gin"
)

func isDropActive(envGetter utils.EnvGetter) bool {
	return envGetter.GetBoolEnv("DRIP_STRATEGY_ACTIVE")
}

func Drip(c *gin.Context, envGetter utils.EnvGetter) bool {

	if !isDropActive(envGetter) {
		return true
	}

	rand.Seed(time.Now().UnixNano())

	threshold := envGetter.GetDoubleEnv("DRIP_STRATEGY_PORCENTAJE") / 100

	random := rand.Float64()

	if random < threshold {
		return true
	}

	c.JSON(http.StatusOK, dto.Response{})
	return false
}
