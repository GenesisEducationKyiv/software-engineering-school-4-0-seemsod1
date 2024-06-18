package handlers_test

import (
	"testing"

	"github.com/seemsod1/api-project/internal/rateapi"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	provider := rateapi.NewMockProvider()

	repo := handlers.NewRepo(mockDB, provider)

	assert.NotNil(t, repo)
	assert.Equal(t, mockDB, repo.DB)
}

func TestNewHandlers(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	provider := rateapi.NewMockProvider()

	repo := handlers.NewRepo(mockDB, provider)

	handlers.NewHandlers(repo)

	assert.Equal(t, repo, handlers.Repo)
}
