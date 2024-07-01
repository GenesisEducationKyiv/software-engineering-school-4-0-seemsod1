package forms_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/seemsod1/api-project/internal/forms"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/smth", http.NoBody)

	form := forms.New(r.PostForm)

	if !form.Valid() {
		t.Error("expected valid form; got invalid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/smth", http.NoBody)
	form := forms.New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/smth", http.NoBody)
	r.PostForm = postedData

	form = forms.New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when required fields exist")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := forms.New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "test@mail.com")
	form = forms.New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = forms.New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}
