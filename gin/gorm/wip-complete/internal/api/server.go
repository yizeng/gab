package api

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/yizeng/gab/gin/wip-complete/docs"
	v1 "github.com/yizeng/gab/gin/wip-complete/internal/api/handler/v1"
	"github.com/yizeng/gab/gin/wip-complete/internal/api/middleware"
	"github.com/yizeng/gab/gin/wip-complete/internal/config"
	"github.com/yizeng/gab/gin/wip-complete/internal/repository"
	"github.com/yizeng/gab/gin/wip-complete/internal/repository/dao"
	"github.com/yizeng/gab/gin/wip-complete/internal/service"
)

type Server struct {
	Config *config.AppConfig
	Router *gin.Engine
}

func NewServer(conf *config.AppConfig, db *gorm.DB) *Server {
	gin.SetMode(conf.Gin.Mode)
	engine := gin.New()

	s := &Server{
		Config: conf,
		Router: engine,
	}

	s.MountMiddlewares()

	authHandler := s.initAuthHandler(db)
	userHandler := s.initUserHandler(db)
	articleHandler := s.initArticleHandler(db)
	s.MountHandlers(authHandler, userHandler, articleHandler)

	return s
}

func (s *Server) initArticleHandler(db *gorm.DB) *v1.ArticleHandler {
	articleDAO := dao.NewArticleDAO(db)
	repo := repository.NewArticleRepository(articleDAO)
	svc := service.NewArticleService(repo)
	handler := v1.NewArticleHandler(svc)

	return handler
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

func (s *Server) MountMiddlewares() {
	// Logger and Recovery are needed unless we use gin.Default().
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(requestid.New())
	s.Router.Use(middleware.ConfigCORS(s.Config.API.AllowedCORSDomains))
}

func (s *Server) MountHandlers(authHandler *v1.AuthHandler, userHandler *v1.UserHandler, articleHandler *v1.ArticleHandler) {
	const basePath = "/api/v1"

	auth := s.Router.Group(basePath)
	{
		auth.POST("/auth/signup", authHandler.HandleSignup)
		auth.POST("/auth/login", authHandler.HandleLogin)
	}

	users := s.Router.Group(basePath, middleware.NewAuthenticator(s.Config.API.JWTSigningKey).VerifyJWT())
	{
		users.GET("/users/:userID", userHandler.HandleGetUser)
	}

	articles := s.Router.Group(basePath)
	{
		articles.GET("/articles", middleware.Paginate(), articleHandler.HandleListArticles)
		articles.POST("/articles", articleHandler.HandleCreateArticle)
		articles.GET("/articles/:articleID", articleHandler.HandleGetArticle)
		articles.GET("/articles/search", articleHandler.HandleSearchArticles)
	}

	s.Router.GET("/", v1.HandleHealthcheck)

	// Setup Swagger UI.
	docs.SwaggerInfo.Host = s.Config.API.BaseURL
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "API for gin/complete"
	docs.SwaggerInfo.Description = "This is an example of Go API with Gin."
	docs.SwaggerInfo.Version = "1.0"
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
