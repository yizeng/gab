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
	"gorm.io/gorm"

	"github.com/yizeng/gab/chi/gorm/crud/docs"
	v1 "github.com/yizeng/gab/chi/gorm/crud/internal/api/handler/v1"
	"github.com/yizeng/gab/chi/gorm/crud/internal/api/middleware"
	"github.com/yizeng/gab/chi/gorm/crud/internal/config"
	"github.com/yizeng/gab/chi/gorm/crud/internal/repository"
	"github.com/yizeng/gab/chi/gorm/crud/internal/repository/dao"
	"github.com/yizeng/gab/chi/gorm/crud/internal/service"
)

type Server struct {
	Config *config.AppConfig
	Router *chi.Mux
}

func NewServer(conf *config.AppConfig, db *gorm.DB) *Server {
	s := &Server{
		Config: conf,
		Router: chi.NewRouter(),
	}

	s.MountMiddlewares()

	articleHandler := s.initArticleHandler(db)
	s.MountHandlers(articleHandler)

	return s
}

func (s *Server) MountMiddlewares() {
	s.Router.Use(chimiddleware.RequestID)
	s.Router.Use(chimiddleware.Logger)
	s.Router.Use(chimiddleware.Recoverer)
	s.Router.Use(chimiddleware.CleanPath)
	s.Router.Use(chimiddleware.Heartbeat("/"))

	s.Router.Use(middleware.ConfigCORS(s.Config.API.Environment, s.Config.API.AllowedCORSDomains))
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
}

func (s *Server) MountHandlers(articleHandler *v1.ArticleHandler) {
	const basePath = "/api/v1"

	apiV1Router := chi.NewRouter()
	apiV1Router.Route("/", func(r chi.Router) {
		r.With(middleware.Pagination).Get("/articles", articleHandler.HandleListArticles)
		r.Post("/articles", articleHandler.HandleCreateArticle)
		r.Get("/articles/{articleID}", articleHandler.HandleGetArticle)
		r.Get("/articles/search", articleHandler.HandleSearchArticles)
	})

	s.Router.Mount(basePath, apiV1Router)

	// Setup Swagger UI.
	docs.SwaggerInfo.Host = s.Config.API.BaseURL
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "API for chi/crud-gorm"
	docs.SwaggerInfo.Description = "This is an example of Go API with Chi router."
	docs.SwaggerInfo.Version = "1.0"
	s.Router.Get("/swagger/*", httpSwagger.WrapHandler)

	s.printAllRoutes()
}

func (s *Server) initArticleHandler(db *gorm.DB) *v1.ArticleHandler {
	articleDAO := dao.NewArticleDAO(db)
	articleRepo := repository.NewArticleRepository(articleDAO)
	articleSvc := service.NewArticleService(articleRepo)
	articleHandler := v1.NewArticleHandler(articleSvc)

	return articleHandler
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
