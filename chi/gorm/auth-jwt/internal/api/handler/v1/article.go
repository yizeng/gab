package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/api/handler/v1/request"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/api/middleware"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/domain"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/service"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error)
	GetArticle(ctx context.Context, id uint) (domain.Article, error)
	ListArticles(ctx context.Context, page uint, perPage uint) ([]domain.Article, error)
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
// @Failure      500      {object}   response.Err
// @Router       /articles [post]
func (h *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	req := request.CreateArticleRequest{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, response.ErrBadRequest(err))

		return
	}

	article, err := h.svc.CreateArticle(r.Context(), domain.Article{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		if errors.Is(err, service.ErrArticleDuplicated) {
			_ = render.Render(w, r, response.ErrBadRequest(service.ErrArticleDuplicated))

			return
		}

		err = fmt.Errorf("v1.HandleCreateArticle -> h.svc.CreateArticle -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	render.Status(r, http.StatusCreated)
	if err = render.Render(w, r, response.NewArticle(&article)); err != nil {
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}
}

// HandleGetArticle godoc
// @Summary      Get an article
// @Tags         articles
// @Produce      json
// @Param        articleID   path    int  true "article ID"
// @Success      200      {object}   domain.Article
// @Failure      400      {object}   response.Err
// @Failure      404      {object}   response.Err
// @Failure      500      {object}   response.Err
// @Router       /articles/{articleID} [get]
func (h *ArticleHandler) HandleGetArticle(w http.ResponseWriter, r *http.Request) {
	rawArticleID := chi.URLParam(r, "articleID")
	articleID, err := strconv.Atoi(rawArticleID)
	if err != nil {
		_ = render.Render(w, r, response.ErrInvalidInput("articleID", rawArticleID))

		return
	}

	if articleID <= 0 {
		_ = render.Render(w, r, response.ErrNotFound("article", "ID", articleID))

		return
	}

	article, err := h.svc.GetArticle(r.Context(), uint(articleID))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			_ = render.Render(w, r, response.ErrNotFound("article", "ID", articleID))

			return
		}

		err = fmt.Errorf("v1.HandleGetArticle -> h.svc.GetArticle -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	if err = render.Render(w, r, response.NewArticle(&article)); err != nil {
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}
}

// HandleListArticles godoc
// @Summary      List all articles
// @Tags         articles
// @Produce      json
// @Param        page     query      int  false  "which page to load. Default to 1 if empty."
// @Param        per_page query      int  false  "how many items per page. Default to 10 if empty."
// @Success      200      {object}   []domain.Article
// @Failure      500      {object}   response.Err
// @Router       /articles [get]
func (h *ArticleHandler) HandleListArticles(w http.ResponseWriter, r *http.Request) {
	pageVal := r.Context().Value(middleware.PageQueryKey)
	page, ok := pageVal.(uint)
	if !ok {
		err := fmt.Errorf("key %q's value %v cannot be casted into uint", middleware.PageQueryKey, pageVal)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}
	perPageVal := r.Context().Value(middleware.PerPageQueryKey)
	perPage, ok := perPageVal.(uint)
	if !ok {
		err := fmt.Errorf("key %q's value %v cannot be casted into uint", middleware.PerPageQueryKey, perPage)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	articles, err := h.svc.ListArticles(r.Context(), page, perPage)
	if err != nil {
		err = fmt.Errorf("v1.HandleListArticles -> h.svc.ListArticles -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	if err = render.RenderList(w, r, response.NewArticles(articles)); err != nil {
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}
}

// HandleSearchArticles godoc
// @Summary      Search articles
// @Tags         articles
// @Produce      json
// @Param        title    query     string  false  "search by title"
// @Param        content  query     string  false  "search by content"
// @Success      200      {object}   []domain.Article
// @Failure      500      {object}   response.Err
// @Router       /articles/search [get]
func (h *ArticleHandler) HandleSearchArticles(w http.ResponseWriter, r *http.Request) {
	titleParam := r.URL.Query().Get("title")
	contentParam := r.URL.Query().Get("content")

	articles, err := h.svc.SearchArticles(r.Context(), titleParam, contentParam)
	if err != nil {
		err = fmt.Errorf("v1.HandleSearchArticles -> h.svc.SearchArticles -> %w", err)
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}

	if err = render.RenderList(w, r, response.NewArticles(articles)); err != nil {
		_ = render.Render(w, r, response.ErrInternalServerError(err))

		return
	}
}
