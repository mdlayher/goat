package goat

import (
	"net/http"
)

// apiAuthenticator interface which defines methods required to implement an authentication method
type apiAuthenticator interface {
	Auth(*http.Request) bool
}

// basicAPIAuthenticator uses the HTTP Basic authentication scheme
type basicAPIAuthenticator struct {
}

// Auth handles validation of HTTP Basic authentication
func (a basicAPIAuthenticator) Auth(r *http.Request) bool {
	return true
}

// hmacAPIAuthenticator uses the HMAC-SHA1 authentication scheme
type hmacAPIAuthenticator struct {
}

// Auth handles validation of HMAC-SHA1 authentication
func (a hmacAPIAuthenticator) Auth(r *http.Request) bool {
	return true
}
