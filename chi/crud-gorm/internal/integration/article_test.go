package integration

//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
//
// 	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
// 	"github.com/yizeng/gab/chi/crud-gorm/internal/web"
// 	"github.com/yizeng/gab/chi/crud-gorm/internal/web/handler/v1/request"
// )
//
// func TestArticleHandler_HandleCreateArticle(t *testing.T) {
// 	s := web.NewServer()
//
// 	tests := []struct {
// 		name         string
// 		buildReqBody func() string
// 		respCode     int
// 		respBody     string
// 	}{
// 		{
// 			name: "201 Created",
// 			buildReqBody: func() string {
// 				article := request.CreateArticleRequest{
// 					Article: domain.Article{
// 						ID:      1,
// 						Title:   "title 1",
// 						Content: "content 1",
// 					},
// 				}
//
// 				body, err := json.Marshal(article)
// 				require.NoError(t, err)
//
// 				return string(body)
// 			},
// 			respCode: http.StatusCreated,
// 			respBody: `{"id":1,"title":"title 1","content":"content 1"}`,
// 		},
// 		{
// 			name: "400 Bad Request",
// 			buildReqBody: func() string {
// 				return "["
// 			},
// 			respCode: http.StatusBadRequest,
// 			respBody: `{"status":400,"error":"unexpected EOF"}`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Prepare Request.
// 			body := tt.buildReqBody()
// 			req, err := http.NewRequest("POST", "/api/v1/articles", strings.NewReader(body))
// 			require.NoError(t, err)
//
// 			// Execute Request.
// 			response := executeRequest(req, s)
//
// 			// Check the response code.
// 			checkResponseCode(t, tt.respCode, response.Code)
//
// 			assert.Equal(t, fmt.Sprintf("%v\n", tt.respBody), response.Body.String())
// 		})
// 	}
// }
//
// func TestArticleHandler_HandleGetArticle(t *testing.T) {
// 	s := web.NewServer()
//
// 	tests := []struct {
// 		name      string
// 		articleID string
// 		respCode  int
// 		respBody  string
// 	}{
// 		{
// 			name:      "200 OK",
// 			articleID: "123",
// 			respCode:  http.StatusOK,
// 			respBody:  `{"id":123,"title":"title 123","content":"content 123"}`,
// 		},
// 		{
// 			name:      "404 Not Found - when articleID is negative",
// 			articleID: "-123",
// 			respCode:  http.StatusNotFound,
// 			respBody:  `{"status":404,"error":"article not found (ID=-123)"}`,
// 		},
// 		{
// 			name:      "400 Bad Request",
// 			articleID: "abc",
// 			respCode:  http.StatusBadRequest,
// 			respBody:  `{"status":400,"error":"invalid input field articleID=abc"}`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Prepare Request.
// 			url := fmt.Sprintf("/api/v1/articles/%v", tt.articleID)
// 			req, err := http.NewRequest("GET", url, strings.NewReader(""))
// 			require.NoError(t, err)
//
// 			// Execute Request.
// 			response := executeRequest(req, s)
//
// 			// Check the response code.
// 			checkResponseCode(t, tt.respCode, response.Code)
//
// 			assert.Equal(t, fmt.Sprintf("%v\n", tt.respBody), response.Body.String())
// 		})
// 	}
// }
//
// func TestArticleHandler_HandleListArticles(t *testing.T) {
// 	s := web.NewServer()
//
// 	tests := []struct {
// 		name     string
// 		respCode int
// 		respBody string
// 	}{
// 		{
// 			name:     "200 OK",
// 			respCode: http.StatusOK,
// 			respBody: `[{"id":1,"title":"title 1","content":"content 1"},{"id":2,"title":"title 2","content":"content 2"}]`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Prepare Request.
// 			req, err := http.NewRequest("GET", "/api/v1/articles", strings.NewReader(""))
// 			require.NoError(t, err)
//
// 			// Execute Request.
// 			response := executeRequest(req, s)
//
// 			// Check the response code.
// 			checkResponseCode(t, tt.respCode, response.Code)
//
// 			assert.Equal(t, fmt.Sprintf("%v\n", tt.respBody), response.Body.String())
// 		})
// 	}
// }
