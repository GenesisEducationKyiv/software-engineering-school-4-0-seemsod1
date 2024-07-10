package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seemsod1/api-project/pkg/timezone"

	"github.com/seemsod1/api-project/pkg/logger"

	"github.com/seemsod1/api-project/internal/handlers"
	"github.com/seemsod1/api-project/internal/models"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribeIntegration(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	logg, _ := logger.NewLogger("test")
	repo := handlers.Repository{Subscriber: mockDB, Logger: logg}
	handler := http.HandlerFunc(repo.Subscribe)

	testCases := []struct {
		name       string
		formFields map[string]string
		headers    map[string]string
		mockReturn error
		wantStatus int
	}{
		{
			name:       "no form data",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty form data",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid form data",
			formFields: map[string]string{
				"abc": "123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty email",
			formFields: map[string]string{
				"email": "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			formFields: map[string]string{
				"email": "abc",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "correct email but invalid timezone",
			formFields: map[string]string{
				"email": "test@mail.com",
			},
			headers: map[string]string{
				"Accept-Timezone": "abc",
			},
			mockReturn: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "correct email and timezone",
			formFields: map[string]string{
				"email": "test@mail.com",
			},
			headers: map[string]string{
				"Accept-Timezone": "UTC",
			},
			mockReturn: nil,
			wantStatus: http.StatusOK,
		},
		{
			name: "duplicated email",
			formFields: map[string]string{
				"email": "test@mail.com",
			},
			headers: map[string]string{
				"Accept-Timezone": "UTC",
			},
			mockReturn: dbrepo.ErrorDuplicateSubscription,
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			formFields: map[string]string{
				"email": "test@mail.com",
			},
			mockReturn: assert.AnError,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			for key, value := range tt.formFields {
				err := writer.WriteField(key, value)
				require.NoError(t, err)
			}

			err := writer.Close()
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/subscribe", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			rr := httptest.NewRecorder()

			if tt.wantStatus != http.StatusBadRequest {
				offset := 0
				if timezoneStr, ok := tt.headers["Accept-Timezone"]; ok {
					req.Header.Set("Accept-Timezone", timezoneStr)
					offset, err = timezone.ProcessTimezoneHeader(req)
					require.NoError(t, err)
				}
				mockDB.ExpectedCalls = nil
				mockDB.On("AddSubscriber", models.Subscriber{
					Email:    tt.formFields["email"],
					Timezone: offset,
				}).Return(tt.mockReturn)
			}

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestSubscribe_ErrorParsingForm(t *testing.T) {
	mockDB := dbrepo.NewMockDB()
	logg, _ := logger.NewLogger("test")
	repo := handlers.Repository{Subscriber: mockDB, Logger: logg}
	handler := http.HandlerFunc(repo.Subscribe)

	req := httptest.NewRequest(http.MethodPost, "/subscribe", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
