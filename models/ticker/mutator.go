package ticker

import (
	"errors"
	"net/url"
)

// MutatorFn are used to validate and update Ticker properties
type MutatorFn func(*Ticker) error

// ErrBadScheme indicates a URL Scheme is neither HTTP nor HTTPS
var ErrBadScheme = errors.New("URL scheme must be HTTP or HTTPS")

// WithSymbol sets a Ticker's Symbol property
func WithSymbol(symbol string) MutatorFn {
	return func(t *Ticker) error {
		t.Symbol = symbol

		return nil
	}
}

// WithHref validates that the given href is a valid URL. Valid URLs must be an
// absolute paths and the scheme must be HTTP/S. If the URL is valid, the Ticker's
// Href property will be updated
func WithHref(href string) MutatorFn {
	return func(t *Ticker) error {
		parsed, err := url.ParseRequestURI(href)
		if err != nil {
			return err
		}

		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return ErrBadScheme
		}

		t.Href = parsed.String()
		return nil
	}
}
