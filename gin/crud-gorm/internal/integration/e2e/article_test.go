package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/yizeng/gab/gin/crud-gorm/internal/api"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/crud-gorm/internal/config"
	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository/dao"
	"github.com/yizeng/gab/gin/crud-gorm/pkg/dockertester"
)

type ArticleHandlersTestSuite struct {
	suite.Suite

	db       *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	server   *api.Server
}

func (s *ArticleHandlersTestSuite) SetupSuite() {
	// Initialize container.
	dt := dockertester.InitPostgres()
	s.pool = dt.Pool
	s.resource = dt.Resource

	// Open connection.
	db, err := dockertester.OpenPostgres(dt.Resource, dt.HostPort)
	require.NoError(s.T(), err)

	s.db = db
}

func (s *ArticleHandlersTestSuite) TearDownSuite() {
	err := s.pool.Purge(s.resource) // Destroy the container.
	require.NoError(s.T(), err)
}

func (s *ArticleHandlersTestSuite) SetupTest() {
	// Run migrations.
	err := dao.InitTables(s.db)
	require.NoError(s.T(), err)

	// Seed database.
	script, err := os.ReadFile("../scripts/seed_articles.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)

	// Create API server.
	s.server = api.NewServer(&config.AppConfig{
		API: &config.APIConfig{},
		Gin: &config.GinConfig{
			Mode: gin.TestMode,
		},
		Postgres: &config.PostgresConfig{},
	}, s.db)
}

func (s *ArticleHandlersTestSuite) TearDownTest() {
	s.deleteAllArticles()
}

func (s *ArticleHandlersTestSuite) deleteAllArticles() {
	script, err := os.ReadFile("../scripts/delete_articles.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)
}

func TestArticleHandlers(t *testing.T) {
	suite.Run(t, new(ArticleHandlersTestSuite))
}

func (s *ArticleHandlersTestSuite) TestArticleHandlers_HandleCreateArticle() {
	tests := []struct {
		name         string
		buildReqBody func() string
		respCode     int
		want         *domain.Article
		wantErr      bool
		err          *response.ErrResponse
	}{
		{
			name: "201 Created",
			buildReqBody: func() string {
				article := request.CreateArticleRequest{
					UserID:  123,
					Title:   "title 1",
					Content: "content 1",
				}

				body, err := json.Marshal(article)
				require.NoError(s.T(), err)

				return string(body)
			},
			respCode: http.StatusCreated,
			wantErr:  false,
			err:      nil,
			want: &domain.Article{
				UserID:  123,
				Title:   "title 1",
				Content: "content 1",
			},
		},
		{
			name: "400 Bad Request - Missing user_id, Title too long, Content too long",
			buildReqBody: func() string {
				article := request.CreateArticleRequest{
					Title:   uniuri.NewLen(200),
					Content: uniuri.NewLen(10000),
				}

				body, err := json.Marshal(article)
				require.NoError(s.T(), err)

				return string(body)
			},
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewBadRequest("content: the length must be between 1 and 5000; title: the length must be between 1 and 128; user_id: cannot be blank."),
			want:     nil,
		},
		{
			name: "400 Bad Request - invalid JSON",
			buildReqBody: func() string {
				return "["
			},
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewBadRequest("unexpected EOF"),
			want:     nil,
		},
		{
			name: "400 Bad Request - Already exists",
			buildReqBody: func() string {
				article := request.CreateArticleRequest{
					UserID:  123,
					Title:   "seeded title 999",
					Content: "seeded content 999",
				}

				body, err := json.Marshal(article)
				require.NoError(s.T(), err)

				return string(body)
			},
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewBadRequest("article already exists"),
			want:     nil,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Prepare Request.
			body := tt.buildReqBody()
			req, err := http.NewRequest("POST", "/api/v1/articles", strings.NewReader(body))
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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

func (s *ArticleHandlersTestSuite) TestArticleHandlers_HandleGetArticle() {
	tests := []struct {
		name      string
		articleID string
		respCode  int
		want      *domain.Article
		wantErr   bool
		err       *response.ErrResponse
	}{
		{
			name:      "200 OK",
			articleID: "999",
			respCode:  http.StatusOK,
			wantErr:   false,
			err:       nil,
			want: &domain.Article{
				ID:      999,
				UserID:  123,
				Title:   "seeded title 999",
				Content: "seeded content 999",
			},
		},
		{
			name:      "404 Not Found - articleID is not found",
			articleID: "1",
			respCode:  http.StatusNotFound,
			wantErr:   true,
			err:       response.NewNotFound("article", "ID", "1"),
			want:      nil,
		},
		{
			name:      "404 Not Found - articleID is negative",
			articleID: "-1",
			respCode:  http.StatusNotFound,
			wantErr:   true,
			err:       response.NewNotFound("article", "ID", "-1"),
			want:      nil,
		},
		{
			name:      "400 Bad Request - invalid articleID",
			articleID: "abc",
			respCode:  http.StatusBadRequest,
			wantErr:   true,
			err:       response.NewInvalidInput("articleID", "abc"),
			want:      nil,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Prepare Request.
			url := fmt.Sprintf("/api/v1/articles/%v", tt.articleID)
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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

func (s *ArticleHandlersTestSuite) TestArticleHandlers_HandleListArticles() {
	tests := []struct {
		name     string
		setup    func()
		query    string
		respCode int
		want     []domain.Article
		wantErr  bool
		err      *response.ErrResponse
	}{
		{
			name:     "200 OK",
			setup:    func() {},
			query:    "",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want: []domain.Article{
				{
					ID:      999,
					UserID:  123,
					Title:   "seeded title 999",
					Content: "seeded content 999",
				}, {
					ID:      888,
					UserID:  123,
					Title:   "seeded title 888",
					Content: "seeded content 888",
				},
			},
		},
		{
			name:     "200 OK - With pagination",
			setup:    func() {},
			query:    "?page=2&per_page=1",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want: []domain.Article{
				{
					ID:      888,
					UserID:  123,
					Title:   "seeded title 888",
					Content: "seeded content 888",
				},
			},
		},
		{
			name: "200 OK - When there are no articles",
			setup: func() {
				s.deleteAllArticles()
			},
			query:    "?page=1&per_page=2",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want:     []domain.Article{},
		},
		{
			name:     "400 Bad Request - Invalid page query",
			setup:    func() {},
			query:    "?page=abc&per_page=2",
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewInvalidInput("page", "abc"),
			want:     []domain.Article{},
		},
		{
			name:     "400 Bad Request - Negative page query",
			setup:    func() {},
			query:    "?page=-123&per_page=2",
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewInvalidInput("page", "-123"),
			want:     []domain.Article{},
		},
		{
			name:     "400 Bad Request - Invalid per_page query",
			setup:    func() {},
			query:    "?page=1&per_page=abc",
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewInvalidInput("per_page", "abc"),
			want:     []domain.Article{},
		},
		{
			name:     "400 Bad Request - Negative per_page query",
			setup:    func() {},
			query:    "?page=1&per_page=-123",
			respCode: http.StatusBadRequest,
			wantErr:  true,
			err:      response.NewInvalidInput("per_page", "-123"),
			want:     []domain.Article{},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup the tests.
			tt.setup()

			// Prepare Request.
			req, err := http.NewRequest("GET", "/api/v1/articles"+tt.query, nil)
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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

func (s *ArticleHandlersTestSuite) TestArticleHandlers_HandleSearchArticles() {
	tests := []struct {
		name     string
		query    string
		respCode int
		want     []domain.Article
		wantErr  bool
		err      *response.ErrResponse
	}{
		{
			name:     "200 OK - by title",
			query:    "title=999",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want: []domain.Article{
				{
					ID:      999,
					UserID:  123,
					Title:   "seeded title 999",
					Content: "seeded content 999",
				},
			},
		},
		{
			name:     "200 OK - by content",
			query:    "content=999",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want: []domain.Article{
				{
					ID:      999,
					UserID:  123,
					Title:   "seeded title 999",
					Content: "seeded content 999",
				},
			},
		},
		{
			name:     "200 OK - When there are no results",
			query:    "title=no-title&content=no-content",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want:     []domain.Article{},
		},
		{
			name:     "200 OK - No query parameters",
			query:    "",
			respCode: http.StatusOK,
			wantErr:  false,
			err:      nil,
			want: []domain.Article{
				{
					ID:      999,
					UserID:  123,
					Title:   "seeded title 999",
					Content: "seeded content 999",
				},
				{
					ID:      888,
					UserID:  123,
					Title:   "seeded title 888",
					Content: "seeded content 888",
				},
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Prepare Request.
			req, err := http.NewRequest("GET", "/api/v1/articles/search?"+tt.query, nil)
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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
