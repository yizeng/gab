package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yizeng/gab/chi/crud-gorm/internal/api/handler/v1/request"
	"github.com/yizeng/gab/chi/crud-gorm/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
	"github.com/yizeng/gab/chi/crud-gorm/internal/service"
)

func TestArticleHandler_HandleCreateArticle(t *testing.T) {
	testArticle := &domain.Article{
		ID:      999,
		UserID:  123,
		Title:   "title",
		Content: "content",
	}
	testError := errors.New("test error")

	tests := []struct {
		name         string
		setupService func() ArticleService
		buildReqBody func() string
		respCode     int
		want         *domain.Article
		wantErr      bool
		err          *response.ErrResponse
	}{
		{
			name: "201 Created",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockCreate = func(ctx context.Context, article *domain.Article) (*domain.Article, error) {
					return testArticle, nil
				}
				return mock
			},
			buildReqBody: func() string {
				article := request.CreateArticleRequest{
					UserID:  123,
					Title:   "title",
					Content: "content",
				}

				body, err := json.Marshal(article)
				require.NoError(t, err)

				return string(body)
			},
			respCode: http.StatusCreated,
			wantErr:  false,
			want:     testArticle,
			err:      nil,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockCreate = func(ctx context.Context, article *domain.Article) (*domain.Article, error) {
					return nil, testError
				}
				return mock
			},
			buildReqBody: func() string {
				article := request.CreateArticleRequest{
					UserID:  123,
					Title:   "title",
					Content: "content",
				}

				body, err := json.Marshal(article)
				require.NoError(t, err)

				return string(body)
			},
			respCode: http.StatusInternalServerError,
			wantErr:  true,
			want:     nil,
			err:      response.NewInternalServerError(testError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			r := chi.NewRouter()
			r.Post("/", h.HandleCreateArticle)

			// Prepare request.
			body := tt.buildReqBody()
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.err.ErrorCode, result.ErrorCode)
			} else {
				var result domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.UserID, result.UserID)
				assert.Equal(t, tt.want.Title, result.Title)
				assert.Equal(t, tt.want.Content, result.Content)
			}
		})
	}
}

func TestArticleHandler_HandleGetArticle(t *testing.T) {
	testArticle := &domain.Article{
		ID:      999,
		UserID:  123,
		Title:   "title",
		Content: "content",
	}
	testError := errors.New("test error")

	tests := []struct {
		name         string
		setupService func() ArticleService
		articleID    string
		respCode     int
		want         *domain.Article
		wantErr      bool
		err          *response.ErrResponse
	}{
		{
			name: "200 OK",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockGetArticle = func(ctx context.Context, id uint) (*domain.Article, error) {
					return testArticle, nil
				}
				return mock
			},
			articleID: "999",
			respCode:  http.StatusOK,
			wantErr:   false,
			want:      testArticle,
			err:       nil,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockGetArticle = func(ctx context.Context, id uint) (*domain.Article, error) {
					return nil, testError
				}
				return mock
			},
			articleID: "999",
			respCode:  http.StatusInternalServerError,
			wantErr:   true,
			want:      nil,
			err:       response.NewInternalServerError(testError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			r := chi.NewRouter()
			r.Get("/{articleID}", h.HandleGetArticle)

			// Prepare request.
			url := fmt.Sprintf("/%v", tt.articleID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.err.ErrorCode, result.ErrorCode)
			} else {
				var result domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.UserID, result.UserID)
				assert.Equal(t, tt.want.Title, result.Title)
				assert.Equal(t, tt.want.Content, result.Content)
			}
		})
	}
}

func TestArticleHandler_HandleListArticles(t *testing.T) {
	testArticles := []domain.Article{
		{
			ID:      999,
			UserID:  123,
			Title:   "title 999",
			Content: "content 999",
		}, {
			ID:      888,
			UserID:  123,
			Title:   "title 888",
			Content: "content 888",
		},
	}
	testError := errors.New("test error")

	tests := []struct {
		name         string
		setupService func() ArticleService
		respCode     int
		want         []domain.Article
		wantErr      bool
		err          *response.ErrResponse
	}{
		{
			name: "200 OK",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockListArticles = func(ctx context.Context) ([]domain.Article, error) {
					return testArticles, nil
				}
				return mock
			},
			respCode: http.StatusOK,
			wantErr:  false,
			want:     testArticles,
			err:      nil,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockListArticles = func(ctx context.Context) ([]domain.Article, error) {
					return nil, testError
				}
				return mock
			},
			respCode: http.StatusInternalServerError,
			wantErr:  true,
			want:     nil,
			err:      response.NewInternalServerError(testError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			r := chi.NewRouter()
			r.Get("/", h.HandleListArticles)

			// Prepare request.
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.err.ErrorCode, result.ErrorCode)
			} else {
				var result []domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(result))

				for i, v := range result {
					assert.Equal(t, tt.want[i].UserID, v.UserID)
					assert.Equal(t, tt.want[i].Title, v.Title)
					assert.Equal(t, tt.want[i].Content, v.Content)
				}
			}
		})
	}
}

func TestArticleHandler_HandleSearchArticles(t *testing.T) {
	testArticles := []domain.Article{
		{
			ID:      999,
			UserID:  123,
			Title:   "title 999",
			Content: "content 999",
		},
	}
	testError := errors.New("test error")

	tests := []struct {
		name         string
		setupService func() ArticleService
		respCode     int
		want         []domain.Article
		wantErr      bool
		err          *response.ErrResponse
	}{
		{
			name: "200 OK",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockSearchArticles = func(ctx context.Context, title, content string) ([]domain.Article, error) {
					return testArticles, nil
				}
				return mock
			},
			respCode: http.StatusOK,
			wantErr:  false,
			want:     testArticles,
			err:      nil,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			setupService: func() ArticleService {
				mock := service.NewArticleServiceMock()
				mock.MockSearchArticles = func(ctx context.Context, title, content string) ([]domain.Article, error) {
					return nil, testError
				}
				return mock
			},
			respCode: http.StatusInternalServerError,
			wantErr:  true,
			want:     nil,
			err:      response.NewInternalServerError(testError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			r := chi.NewRouter()
			r.Get("/", h.HandleSearchArticles)

			// Prepare request.
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.err.ErrorCode, result.ErrorCode)
			} else {
				var result []domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(result))

				for i, v := range result {
					assert.Equal(t, tt.want[i].UserID, v.UserID)
					assert.Equal(t, tt.want[i].Title, v.Title)
					assert.Equal(t, tt.want[i].Content, v.Content)
				}
			}
		})
	}
}
