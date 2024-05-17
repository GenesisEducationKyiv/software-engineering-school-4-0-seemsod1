package main

import (
	"testing"
)

func TestSetup(t *testing.T) {
	if err := setup(&app); err != nil {
		t.Error(err)
	}
}
