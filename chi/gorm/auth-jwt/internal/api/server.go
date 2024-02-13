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

	"github.com/yizeng/gab/chi/gorm/auth-jwt/docs"
	v1 "github.com/yizeng/gab/chi/gorm/auth-jwt/internal/api/handler/v1"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/api/middleware"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/config"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/repository"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/repository/dao"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/service"
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

	authHandler := s.initAuthHandler(db)
	userHandler := s.initUserHandler(db)
	s.MountHandlers(authHandler, userHandler)

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

func (s *Server) MountHandlers(authHandler *v1.AuthHandler, userHandler *v1.UserHandler) {
	const basePath = "/api/v1"

	apiV1Router := chi.NewRouter()
	apiV1Router.Route("/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/auth/signup", authHandler.HandleSignup)
			r.Post("/auth/login", authHandler.HandleLogin)
		})

		r.Group(func(r chi.Router) {
			authenticator := middleware.NewAuthenticator(s.Config.API.JWTSigningKey)
			r.Use(authenticator.VerifyJWT)

			r.Get("/users/{userID}", userHandler.HandleGetUser)
		})
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

func (s *Server) initAuthHandler(db *gorm.DB) *v1.AuthHandler {
	userDAO := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(userDAO)
	svc := service.NewAuthService(repo)
	handler := v1.NewAuthHandler(s.Config.API, svc)

	return handler
}

func (s *Server) initUserHandler(db *gorm.DB) *v1.UserHandler {
	userDAO := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(userDAO)
	svc := service.NewUserService(repo)
	handler := v1.NewUserHandler(svc)

	return handler
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
