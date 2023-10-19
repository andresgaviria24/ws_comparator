package server

import (
	"log"
	"os"
	"ws_comparator/application"
	"ws_comparator/docs"
	"ws_comparator/interfaces/middleware"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerImpl struct {
	router *gin.Engine
}

func InitServer() Server {
	serverImpl := &ServerImpl{}
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("SERVICE_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		//newrelic.ConfigDebugLogger(os.Stdout),
	)

	router.Use(nrgin.Middleware(app))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	application.InitComparatorController(router, app)

	swaggerDocs()
	serverImpl.router = router
	return serverImpl
}

func swaggerDocs() {
	docs.SwaggerInfo.Title = os.Getenv("SWAGGER_TITLE")
	docs.SwaggerInfo.Description = os.Getenv("SWAGGER_DESCRIPTION")
	docs.SwaggerInfo.Version = os.Getenv("SWAGGER_VERSION")
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.BasePath = os.Getenv("SWAGGER_BASEPATH")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func (api *ServerImpl) RunServer() {
	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = os.Getenv("LOCAL_PORT") //localhost
	}
	log.Fatal(api.router.Run(":" + appPort))
}
