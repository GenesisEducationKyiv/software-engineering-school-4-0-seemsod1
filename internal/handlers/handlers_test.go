package handlers_test

import (
	"context"
	"testing"

	"github.com/seemsod1/api-project/pkg/logger"

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

func (mp *mockProvider) GetRate(ctx context.Context, base, target string) (float64, error) {
	args := mp.Called(ctx, base, target)
	return args.Get(0).(float64), args.Error(1)
}

type mockNotifier struct {
	mock.Mock
}

func newMockNotifier() *mockNotifier {
	return &mockNotifier{}
}

func (mn *mockNotifier) SendRate(subscribers []string) {
	mn.Called(subscribers)
}

func TestNewRepo(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	provider := newMockProvider()
	not := newMockNotifier()
	log, _ := logger.NewLogger("test")

	repo := handlers.NewRepo(mockDB, provider, not, log)

	assert.NotNil(t, repo)
	assert.Equal(t, mockDB, repo.Subscriber)
}

func TestNewHandlers(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	provider := newMockProvider()
	not := newMockNotifier()
	log, _ := logger.NewLogger("test")

	repo := handlers.NewRepo(mockDB, provider, not, log)

	handlers.NewHandlers(repo)

	assert.Equal(t, repo, handlers.Repo)
}
