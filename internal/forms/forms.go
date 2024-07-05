package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(f.Get(field), "email"); err != nil {
		f.Errors.Add(field, "Invalid email address")
	}
}

func ParseEmail(r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return "", fmt.Errorf("unable to parse form")
	}

	email := r.Form.Get("email")

	form := New(r.PostForm)
	form.Required("email")
	form.IsEmail("email")

	if !form.Valid() {
		return "", fmt.Errorf("invalid email")
	}

	return email, nil
}
