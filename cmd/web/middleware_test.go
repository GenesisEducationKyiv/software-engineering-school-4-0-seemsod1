package main

import (
	"fmt"
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
		t.Error(fmt.Sprintf("type is not http.Handler, got %v", v))
	}
}
