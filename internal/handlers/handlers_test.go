package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	mock.Mock
}

func newMockProvider() *mockProvider {
	return &mockProvider{}
}

func (mp *mockProvider) GetRate(base, target string) (float64, error) {
	args := mp.Called(base, target)
	return args.Get(0).(float64), args.Error(1)
}

func TestNewRepo(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	provider := newMockProvider()

	repo := handlers.NewRepo(mockDB, provider)

	assert.NotNil(t, repo)
	assert.Equal(t, mockDB, repo.DB)
}

func TestNewHandlers(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	provider := newMockProvider()

	repo := handlers.NewRepo(mockDB, provider)

	handlers.NewHandlers(repo)

	assert.Equal(t, repo, handlers.Repo)
}
