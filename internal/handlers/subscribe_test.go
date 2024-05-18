package handlers

import (
	"bytes"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/storage/dbrepo"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubscribe(t *testing.T) {
	var body bytes.Buffer

	mockDB := dbrepo.NewMockDB()
	repo := Repository{DB: mockDB}

	// no form data
	req := httptest.NewRequest(http.MethodPost, "/subscribe", &body)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repo.Subscribe)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// empty form data
	body.Reset()
	writer := multipart.NewWriter(&body)
	err := writer.Close()
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/subscribe", nil)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// invalid form data
	body.Reset()
	writer = multipart.NewWriter(&body)
	_, err = writer.CreateFormFile("abc", "123")
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/subscribe", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// empty email
	body.Reset()
	writer = multipart.NewWriter(&body)
	err = writer.WriteField("email", "")
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/subscribe", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// invalid email
	body.Reset()
	writer = multipart.NewWriter(&body)
	err = writer.WriteField("email", "abc")
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/subscribe", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// correct email
	body.Reset()
	writer = multipart.NewWriter(&body)
	err = writer.WriteField("email", "test@mail.com")
	assert.NoError(t, err)
	err = writer.Close()

	req = httptest.NewRequest(http.MethodPost, "/subscribe", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()

	mockDB.On("AddSubscriber", "test@mail.com").Return(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// duplicated email
	rr = httptest.NewRecorder()

	mockDB.ExpectedCalls = nil
	mockDB.On("AddSubscriber", "test@mail.com").Return(customerrors.DuplicatedKey)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)

	// internal error
	rr = httptest.NewRecorder()

	mockDB.ExpectedCalls = nil
	mockDB.On("AddSubscriber", "test1@mail.com").Return(assert.AnError)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
