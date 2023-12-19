package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yizeng/gab/chi/minimum/docs"
	"github.com/yizeng/gab/chi/minimum/internal/service"
	v1 "github.com/yizeng/gab/chi/minimum/internal/web/handler/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Server struct {
	Address string
	Router  *chi.Mux
}

func NewServer() *Server {
	s := &Server{
		Address: getServerAddress(),
		Router:  chi.NewRouter(),
	}

	s.MountMiddlewares()
	s.MountHandlers()

	return s
}

func getServerAddress() string {
	host := viper.Get("API_HOST")
	port := viper.Get("API_PORT")
	addr := fmt.Sprintf("%v:%v", host, port)

	return addr
}

func (s *Server) MountMiddlewares() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(middleware.Heartbeat("/"))
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
}

func (s *Server) MountHandlers() {
	const basePath = "/api/v1"

	apiV1Router := chi.NewRouter()
	apiV1Router.Route("/", func(r chi.Router) {
		countrySvc := service.NewCountryService()
		countryHandler := v1.NewCountryHandler(countrySvc)
		r.Post("/countries/sum-population-by-state", countryHandler.HandleSumPopulationByState)
	})

	s.Router.Mount(basePath, apiV1Router)

	// Setup Swagger UI.
	docs.SwaggerInfo.Host = s.Address
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "API for chi/minimum"
	docs.SwaggerInfo.Description = "This is an example of Go API with Chi router."
	docs.SwaggerInfo.Version = "1.0"
	s.Router.Get("/swagger/*", httpSwagger.WrapHandler)

	s.printAllRoutes()
}

func (s *Server) printAllRoutes() {
	zap.L().Info("printing all routes...")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)

		zap.L().Info(fmt.Sprintf("%v\t%v", method, route))

		return nil
	}

	if err := chi.Walk(s.Router, walkFunc); err != nil {
		zap.L().Error("printing all routes failed", zap.Error(err))
	}
}
