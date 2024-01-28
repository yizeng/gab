package api

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/yizeng/gab/gin/crud-gorm/docs"
	v1 "github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/middleware"
	"github.com/yizeng/gab/gin/crud-gorm/internal/config"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository/dao"
	"github.com/yizeng/gab/gin/crud-gorm/internal/service"
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

	articleHandler := initArticleHandler(db)
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
	// Logger and Recovery are needed unless we use gin.Default().
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(requestid.New())
	s.Router.Use(middleware.ConfigCORS(s.Config.API.AllowedCORSDomains))
}

func (s *Server) MountHandlers(articleHandler *v1.ArticleHandler) {
	const basePath = "/api/v1"

	apiV1 := s.Router.Group(basePath)
	{
		apiV1.GET("/articles", middleware.Paginate(), articleHandler.HandleListArticles)
		apiV1.POST("/articles", articleHandler.HandleCreateArticle)
		apiV1.GET("/articles/:articleID", articleHandler.HandleGetArticle)
		apiV1.GET("/articles/search", articleHandler.HandleSearchArticles)
	}

	s.Router.GET("/", v1.HandleHealthcheck)

	// Setup Swagger UI.
	docs.SwaggerInfo.Host = s.Config.API.BaseURL
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Title = "API for gin/crud-gorm"
	docs.SwaggerInfo.Description = "This is an example of Go API with Gin."
	docs.SwaggerInfo.Version = "1.0"
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
