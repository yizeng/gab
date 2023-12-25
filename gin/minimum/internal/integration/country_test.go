package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/yizeng/gab/gin/minimum/internal/api"
	"github.com/yizeng/gab/gin/minimum/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/minimum/internal/config"
	"github.com/yizeng/gab/gin/minimum/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCalculateTotalPopulation(t *testing.T) {
	const expectedTotalPopulation = 68085736

	s := api.NewServer(&config.AppConfig{
		API: &config.APIConfig{},
		Gin: &config.GinConfig{
			Mode: gin.TestMode,
		},
	})

	tests := []struct {
		name         string
		buildReqBody func() string
		wantCode     int
		wantBody     string
		wantErr      bool
	}{
		{
			name: "Happy Path",
			buildReqBody: func() string {
				states := request.SumPopulationByState{
					States: []domain.State{
						{
							Name:       "California",
							Population: 38940231,
						},
						{
							Name:       "Texas",
							Population: 29145505,
						},
					},
				}

				body, err := json.Marshal(states)
				require.NoError(t, err)

				return string(body)
			},
			wantCode: http.StatusOK,
			wantBody: fmt.Sprintf("{\"total_population\":%v}", expectedTotalPopulation),
		},
		{
			name: "400 Bad Request - Invalid JSON",
			buildReqBody: func() string {
				return "["
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"status":400,"error":"unexpected EOF"}`,
		},
		{
			name: "400 Bad Request - Missing required values",
			buildReqBody: func() string {
				return `{"states": [{"population": 123}]}`
			},
			wantCode: http.StatusBadRequest,
			wantBody: `{"status":400,"error":"Key: 'SumPopulationByState.States[0].Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a New Request.
			body := tt.buildReqBody()
			req, err := http.NewRequest("POST", "/api/v1/countries/sum-population-by-state", strings.NewReader(body))
			require.NoError(t, err)

			// Execute Request.
			response := executeRequest(req, s)

			// Check the response code and body.
			assert.Equal(t, tt.wantCode, response.Code)
			assert.Equal(t, tt.wantBody, response.Body.String())
		})
	}
}
