package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSubscribe(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fieldWriter, err := writer.CreateFormField("email")
	assert.NoError(t, err)
	_, err = io.Copy(fieldWriter, strings.NewReader("test@mail.com"))
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/subscribe", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.Subscribe)
	handler.ServeHTTP(rr, req)
}
