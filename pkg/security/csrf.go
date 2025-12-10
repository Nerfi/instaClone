package security

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func NewCSRF(secret []byte, secure bool) func(http.Handler) http.Handler {
	return csrf.Protect(secret, csrf.Secure(secure))
}
