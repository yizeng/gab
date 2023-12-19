package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/yizeng/gab/chi/crud-gorm/internal/service"
	v1 "github.com/yizeng/gab/chi/crud-gorm/internal/web/handler/v1"
)

type Server struct {
	Address string
	Router  *chi.Mux
}

func NewServer(db *gorm.DB) *Server {
	articleHandler := initArticleHandler(db)

	s := &Server{
		Address: getServerAddress(),
		Router:  chi.NewRouter(),
	}

	s.MountMiddlewares()
	s.MountHandlers(articleHandler)

	return s
}

func initArticleHandler(db *gorm.DB) *v1.ArticleHandler {
	articleSvc := service.NewArticleService()
	articleHandler := v1.NewArticleHandler(articleSvc)

	return articleHandler
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

func (s *Server) MountHandlers(articleHandler *v1.ArticleHandler) {
	apiV1Router := chi.NewRouter()
	apiV1Router.Route("/", func(r chi.Router) {
		r.Get("/articles", articleHandler.HandleListArticles)
		r.Post("/articles", articleHandler.HandleCreateArticle)
		r.Get("/articles/{articleID}", articleHandler.HandleGetArticle)
	})

	s.Router.Mount("/api/v1", apiV1Router)

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
