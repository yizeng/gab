package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/service"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error)
	GetArticle(ctx context.Context, id uint) (*domain.Article, error)
	ListArticles(ctx context.Context) ([]domain.Article, error)
}

type ArticleHandler struct {
	svc ArticleService
}

func NewArticleHandler(svc ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

// HandleCreateArticle godoc
// @Summary      Create an article
// @Tags         articles
// @Produce      json
// @Param        request   body      request.CreateArticleRequest true "request body"
// @Success      200      {object}   domain.Article
// @Failure      400      {object}   response.ErrResponse
// @Failure      500      {object}   response.ErrResponse
// @Router       /articles [post]
func (h *ArticleHandler) HandleCreateArticle(ctx *gin.Context) {
	req := request.CreateArticleRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err.Error()))

		return
	}

	if err := req.Validate(); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err.Error()))

		return
	}

	article, err := h.svc.CreateArticle(ctx.Request.Context(), &domain.Article{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		if errors.Is(err, service.ErrArticleDuplicated) {
			response.RenderError(ctx, response.NewBadRequest(service.ErrArticleDuplicated.Error()))

			return
		}

		err = fmt.Errorf("v1.HandleCreateArticle -> h.svc.CreateArticle -> %w", err)
		response.RenderError(ctx, response.NewInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusCreated, article)
}

// HandleGetArticle godoc
// @Summary      Get an article
// @Tags         articles
// @Produce      json
// @Param        articleID   path    int  true "article ID"
// @Success      200      {object}   domain.Article
// @Failure      400      {object}   response.ErrResponse
// @Failure      404      {object}   response.ErrResponse
// @Failure      500      {object}   response.ErrResponse
// @Router       /articles/{articleID} [get]
func (h *ArticleHandler) HandleGetArticle(ctx *gin.Context) {
	rawArticleID := ctx.Param("articleID")
	articleID, err := strconv.Atoi(rawArticleID)
	if err != nil {
		response.RenderError(ctx, response.NewInvalidInput("articleID", rawArticleID))

		return
	}

	if articleID <= 0 {
		response.RenderError(ctx, response.NewNotFound("article", "ID", articleID))

		return
	}

	article, err := h.svc.GetArticle(ctx.Request.Context(), uint(articleID))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.RenderError(ctx, response.NewNotFound("article", "ID", articleID))

			return
		}

		err = fmt.Errorf("v1.HandleGetArticle -> h.svc.GetArticle -> %w", err)
		response.RenderError(ctx, response.NewInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusOK, article)
}

// HandleListArticles godoc
// @Summary      List all articles
// @Tags         articles
// @Produce      json
// @Success      200      {object}   []domain.Article
// @Failure      500      {object}   response.ErrResponse
// @Router       /articles [get]
func (h *ArticleHandler) HandleListArticles(ctx *gin.Context) {
	articles, err := h.svc.ListArticles(ctx.Request.Context())
	if err != nil {
		err = fmt.Errorf("v1.HandleListArticles -> h.svc.ListArticles -> %w", err)
		response.RenderError(ctx, response.NewInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusOK, articles)
}
