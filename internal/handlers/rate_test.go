package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestRate(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	server := httptest.NewServer(http.HandlerFunc(handlers.Repo.Rate))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	type APIResponse struct {
		Price float64 `json:"price"`
	}
	var apiResponse APIResponse
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
