package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
	"github.com/yizeng/gab/chi/crud-gorm/internal/web/handler/v1/request"
	"github.com/yizeng/gab/chi/crud-gorm/internal/web/handler/v1/response"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error)
	GetArticle(ctx context.Context, ID uint) (domain.Article, error)
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

func (h *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	req := request.CreateArticleRequest{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, response.NewBadRequest(err.Error()))

		return
	}

	article, err := h.svc.CreateArticle(r.Context(), req.Article)
	if err != nil {
		err = fmt.Errorf("v1.HandleCreateArticle -> h.svc.CreateArticle -> %w", err)
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, response.NewArticle(&article))
	if err != nil {
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}
}

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
		err = fmt.Errorf("v1.HandleGetArticle -> h.svc.GetArticle -> %w", err)
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}

	err = render.Render(w, r, response.NewArticle(&article))
	if err != nil {
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}
}

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
