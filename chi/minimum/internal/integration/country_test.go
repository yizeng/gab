package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/yizeng/gab/chi/minimum/internal/api"
	"github.com/yizeng/gab/chi/minimum/internal/api/handler/v1/request"
	"github.com/yizeng/gab/chi/minimum/internal/config"
	"github.com/yizeng/gab/chi/minimum/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCalculateTotalPopulation(t *testing.T) {
	s := api.NewServer(&config.APIConfig{})

	// Create a New Request.
	expectedTotalPopulation := 68085736
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

	req, err := http.NewRequest("POST", "/api/v1/countries/sum-population-by-state", strings.NewReader(string(body)))
	require.NoError(t, err)

	// Execute Request.
	response := executeRequest(req, s)

	// Check the response code.
	assert.Equal(t, http.StatusOK, response.Code)

	assert.Equal(t, fmt.Sprintf("{\"total_population\":%v}\n", expectedTotalPopulation), response.Body.String())
}

func TestHandleCalculateTotalPopulation_BadRequest(t *testing.T) {
	s := api.NewServer(&config.APIConfig{})

	// Create a New Request.
	body := strings.NewReader(`[`)
	req, err := http.NewRequest("POST", "/api/v1/countries/sum-population-by-state", body)
	require.NoError(t, err)

	// Execute Request.
	response := executeRequest(req, s)

	// Check the response code.
	assert.Equal(t, http.StatusBadRequest, response.Code)

	msg := `{"status":400,"error":"unexpected EOF"}`
	assert.Equal(t, fmt.Sprintf("%v\n", msg), response.Body.String())
}
