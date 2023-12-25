package api

import (
	"fmt"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/yizeng/gab/gin/minimum/docs"
	v1 "github.com/yizeng/gab/gin/minimum/internal/api/handler/v1"
	"github.com/yizeng/gab/gin/minimum/internal/api/middleware"
	"github.com/yizeng/gab/gin/minimum/internal/config"
	"github.com/yizeng/gab/gin/minimum/internal/service"
)

type Server struct {
	Address string
	Config  *config.AppConfig
	Router  *gin.Engine
}

func NewServer(conf *config.AppConfig) *Server {
	address := fmt.Sprintf("%v:%v", conf.API.Host, conf.API.Port)

	gin.SetMode(conf.Gin.Mode)
	engine := gin.New()

	s := &Server{
		Address: address,
		Config:  conf,
		Router:  engine,
	}

	s.MountMiddlewares()

	countrySvc := service.NewCountryService()
	countryHandler := v1.NewCountryHandler(countrySvc)
	s.MountHandlers(countryHandler)

	return s
}

func (s *Server) MountMiddlewares() {
	// Logger and Recovery are needed unless we use gin.Default().
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(requestid.New())
	s.Router.Use(middleware.ConfigCORS(s.Config.API.Host))
}

func (s *Server) MountHandlers(countryHandler *v1.CountryHandler) {
	const basePath = "/api/v1"

	apiV1 := s.Router.Group(basePath)
	{
		apiV1.POST("/countries/sum-population-by-state", countryHandler.HandleSumPopulationByState)
	}

	s.Router.GET("/", v1.HandleHealthcheck)

	// Setup Swagger UI.
	docs.SwaggerInfo.Host = s.Address
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "API for gin/minimum"
	docs.SwaggerInfo.Description = "This is an example of Go API with Gin."
	docs.SwaggerInfo.Version = "1.0"
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
