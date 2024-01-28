package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/yizeng/gab/chi/minimum/docs"
	v1 "github.com/yizeng/gab/chi/minimum/internal/api/handler/v1"
	"github.com/yizeng/gab/chi/minimum/internal/api/middleware"
	"github.com/yizeng/gab/chi/minimum/internal/config"
	"github.com/yizeng/gab/chi/minimum/internal/service"
)

type Server struct {
	Config *config.APIConfig
	Router *chi.Mux
}

func NewServer(conf *config.APIConfig) *Server {
	s := &Server{
		Config: conf,
		Router: chi.NewRouter(),
	}

	s.MountMiddlewares()
	s.MountHandlers()

	return s
}

func (s *Server) MountMiddlewares() {
	s.Router.Use(chimiddleware.RequestID)
	s.Router.Use(chimiddleware.Logger)
	s.Router.Use(chimiddleware.Recoverer)
	s.Router.Use(chimiddleware.CleanPath)
	s.Router.Use(chimiddleware.Heartbeat("/"))

	s.Router.Use(middleware.ConfigCORS(s.Config.Environment, s.Config.AllowedCORSDomains))
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
	docs.SwaggerInfo.Host = s.Config.BaseURL
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
