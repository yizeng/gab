package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/crud-gorm/internal/api/handler/v1/request"
	"github.com/yizeng/gab/chi/crud-gorm/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
	"github.com/yizeng/gab/chi/crud-gorm/internal/service"
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
func (h *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	req := request.CreateArticleRequest{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, response.NewBadRequest(err.Error()))

		return
	}

	article, err := h.svc.CreateArticle(r.Context(), &domain.Article{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		if errors.Is(err, service.ErrArticleDuplicated) {
			render.Render(w, r, response.NewBadRequest(service.ErrArticleDuplicated.Error()))

			return
		}

		err = fmt.Errorf("v1.HandleCreateArticle -> h.svc.CreateArticle -> %w", err)
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, response.NewArticle(article))
	if err != nil {
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}
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
func (h *ArticleHandler) HandleGetArticle(w http.ResponseWriter, r *http.Request) {
	rawArticleID := chi.URLParam(r, "articleID")
	articleID, err := strconv.Atoi(rawArticleID)
	if err != nil {
		render.Render(w, r, response.NewInvalidInput("articleID", rawArticleID))

		return
	}

	if articleID <= 0 {
		render.Render(w, r, response.NewNotFound("article", "ID", articleID))

		return
	}

	article, err := h.svc.GetArticle(r.Context(), uint(articleID))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			render.Render(w, r, response.NewNotFound("article", "ID", articleID))

			return
		}

		err = fmt.Errorf("v1.HandleGetArticle -> h.svc.GetArticle -> %w", err)
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}

	err = render.Render(w, r, response.NewArticle(article))
	if err != nil {
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}
}

// HandleListArticles godoc
// @Summary      List all articles
// @Tags         articles
// @Produce      json
// @Success      200      {object}   []domain.Article
// @Failure      500      {object}   response.ErrResponse
// @Router       /articles [get]
func (h *ArticleHandler) HandleListArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.svc.ListArticles(r.Context())
	if err != nil {
		err = fmt.Errorf("v1.HandleListArticles -> h.svc.ListArticles -> %w", err)
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}

	err = render.RenderList(w, r, response.NewArticles(articles))
	if err != nil {
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}
}
