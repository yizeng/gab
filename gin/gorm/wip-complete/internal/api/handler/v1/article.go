package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yizeng/gab/gin/wip-complete/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/wip-complete/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/wip-complete/internal/api/middleware"
	"github.com/yizeng/gab/gin/wip-complete/internal/domain"
	"github.com/yizeng/gab/gin/wip-complete/internal/service"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error)
	GetArticle(ctx context.Context, id uint) (domain.Article, error)
	ListArticles(ctx context.Context, page, perPage uint) ([]domain.Article, error)
	SearchArticles(ctx context.Context, title, content string) ([]domain.Article, error)
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
// @Failure      400      {object}   response.Err
// @Success      401      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /articles [post]
func (h *ArticleHandler) HandleCreateArticle(ctx *gin.Context) {
	req := request.CreateArticleRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RenderErr(ctx, response.ErrBadRequest(err))

		return
	}

	if err := req.Validate(); err != nil {
		response.RenderErr(ctx, response.ErrBadRequest(err))

		return
	}

	article, err := h.svc.CreateArticle(ctx.Request.Context(), domain.Article{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		if errors.Is(err, service.ErrArticleDuplicated) {
			response.RenderErr(ctx, response.ErrBadRequest(service.ErrArticleDuplicated))

			return
		}

		err = fmt.Errorf("v1.HandleCreateArticle -> h.svc.CreateArticle -> %w", err)
		response.RenderErr(ctx, response.ErrInternalServerError(err))

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
// @Failure      400      {object}   response.Err
// @Success      401      {object}   response.Err
// @Failure      404      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /articles/{articleID} [get]
func (h *ArticleHandler) HandleGetArticle(ctx *gin.Context) {
	rawArticleID := ctx.Param("articleID")
	articleID, err := strconv.Atoi(rawArticleID)
	if err != nil {
		response.RenderErr(ctx, response.ErrInvalidInput("articleID", rawArticleID))

		return
	}

	if articleID <= 0 {
		response.RenderErr(ctx, response.ErrNotFound("article", "ID", articleID))

		return
	}

	article, err := h.svc.GetArticle(ctx.Request.Context(), uint(articleID))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.RenderErr(ctx, response.ErrNotFound("article", "ID", articleID))

			return
		}

		err = fmt.Errorf("v1.HandleGetArticle -> h.svc.GetArticle -> %w", err)
		response.RenderErr(ctx, response.ErrInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusOK, article)
}

// HandleListArticles godoc
// @Summary      List all articles
// @Tags         articles
// @Produce      json
// @Param        page     query      int  false  "which page to load. Default to 1 if empty."
// @Param        per_page query      int  false  "how many items per page. Default to 10 if empty."
// @Success      200      {object}   []domain.Article
// @Success      401      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /articles [get]
func (h *ArticleHandler) HandleListArticles(ctx *gin.Context) {
	page, err := parsePaginationQuery(ctx, middleware.PageQueryKey)
	if err != nil {
		response.RenderErr(ctx, response.ErrInternalServerError(err))

		return
	}
	perPage, err := parsePaginationQuery(ctx, middleware.PerPageQueryKey)
	if err != nil {
		response.RenderErr(ctx, response.ErrInternalServerError(err))

		return
	}

	articles, err := h.svc.ListArticles(ctx.Request.Context(), page, perPage)
	if err != nil {
		err = fmt.Errorf("v1.HandleListArticles -> h.svc.ListArticles -> %w", err)
		response.RenderErr(ctx, response.ErrInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusOK, articles)
}

// HandleSearchArticles godoc
// @Summary      Search articles
// @Tags         articles
// @Produce      json
// @Param        title    query     string  false  "search by title"
// @Param        content  query     string  false  "search by content"
// @Success      200      {object}   []domain.Article
// @Success      401      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /articles/search [get]
func (h *ArticleHandler) HandleSearchArticles(ctx *gin.Context) {
	titleParam := ctx.Query("title")
	contentParam := ctx.Query("content")

	articles, err := h.svc.SearchArticles(ctx.Request.Context(), titleParam, contentParam)
	if err != nil {
		err = fmt.Errorf("v1.HandleSearchArticles -> h.svc.SearchArticles -> %w", err)
		response.RenderErr(ctx, response.ErrInternalServerError(err))

		return
	}

	ctx.JSON(http.StatusOK, articles)
}

func parsePaginationQuery(ctx *gin.Context, key string) (uint, error) {
	val, exists := ctx.Get(key)
	if !exists {
		return 0, fmt.Errorf("key %q is not found in context", key)
	}

	result, ok := val.(uint)
	if !ok {
		return 0, fmt.Errorf("key %q's value %v cannot be casted into uint", key, val)
	}

	return result, nil
}
