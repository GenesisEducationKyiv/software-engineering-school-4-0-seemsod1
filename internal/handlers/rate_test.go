package handlers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Repo.Rate))
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	type ApiResponse struct {
		Price float64 `json:"price"`
	}
	var apiResponse ApiResponse
	if err = json.Unmarshal(b, &apiResponse); err != nil {
		t.Error(err)
	}

	if apiResponse.Price == -1 {
		t.Error("expected amount; got empty")
	}

	if _, err = json.Marshal(apiResponse); err != nil {
		t.Error(err)
	}
}
