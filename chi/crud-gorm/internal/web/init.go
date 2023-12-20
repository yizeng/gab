package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/yizeng/gab/chi/crud-gorm/docs"
	"github.com/yizeng/gab/chi/crud-gorm/internal/config"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository/dao"
	"github.com/yizeng/gab/chi/crud-gorm/internal/service"
	v1 "github.com/yizeng/gab/chi/crud-gorm/internal/web/handler/v1"
)

type Server struct {
	Address string
	Config  *config.APIConfig
	Router  *chi.Mux
}

func NewServer(conf *config.APIConfig, db *gorm.DB) *Server {
	address := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	articleHandler := initArticleHandler(db)

	s := &Server{
		Address: address,
		Config:  conf,
		Router:  chi.NewRouter(),
	}

	s.MountMiddlewares()
	s.MountHandlers(articleHandler)

	return s
}

func initArticleHandler(db *gorm.DB) *v1.ArticleHandler {
	articleDAO := dao.NewArticleDAO(db)
	articleRepo := repository.NewArticleRepository(articleDAO)
	articleSvc := service.NewArticleService(articleRepo)
	articleHandler := v1.NewArticleHandler(articleSvc)

	return articleHandler
}

func (s *Server) MountMiddlewares() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://0.0.0.0") {
				return true
			}

			return strings.HasPrefix(origin, s.Config.Host)
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            strings.EqualFold(s.Config.Environment, "development"),
	}))
	s.Router.Use(middleware.Heartbeat("/"))
	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
}

func (s *Server) MountHandlers(articleHandler *v1.ArticleHandler) {
	const basePath = "/api/v1"

	apiV1Router := chi.NewRouter()
	apiV1Router.Route("/", func(r chi.Router) {
		r.Get("/articles", articleHandler.HandleListArticles)
		r.Post("/articles", articleHandler.HandleCreateArticle)
		r.Get("/articles/{articleID}", articleHandler.HandleGetArticle)
	})

	s.Router.Mount(basePath, apiV1Router)

	// Setup Swagger UI.
	docs.SwaggerInfo.Host = s.Address
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "API for chi/crud-gorm"
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
