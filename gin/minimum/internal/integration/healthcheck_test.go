package integration

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yizeng/gab/gin/minimum/internal/api"
	"github.com/yizeng/gab/gin/minimum/internal/config"
)

func TestHandleHealthcheck(t *testing.T) {
	s := api.NewServer(&config.AppConfig{
		API: &config.APIConfig{},
		Gin: &config.GinConfig{
			Mode: gin.TestMode,
		},
	})

	// Create a New Request.
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	// Execute Request.
	response := executeRequest(req, s)

	// Check the response code.
	assert.Equal(t, http.StatusOK, response.Code)

	assert.Equal(t, `{"message":"pong"}`, response.Body.String())
}
