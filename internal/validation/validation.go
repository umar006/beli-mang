package validation

import (
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

func URL(fl validator.FieldLevel) bool {
	urlString := fl.Field().String()
	if !strings.Contains(urlString, ".") {
		return false
	}
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	u, err := url.Parse(urlString)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
