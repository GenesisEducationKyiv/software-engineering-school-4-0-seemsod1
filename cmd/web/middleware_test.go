package main

import (
	"net/http"
	"testing"
)

func TestEnableCORS(t *testing.T) {
	var mH myHandler

	h := EnableCORS(&mH)

	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Errorf("type is not http.Handler, got %w", v)
	}
}
