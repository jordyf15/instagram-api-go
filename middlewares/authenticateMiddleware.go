package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type AuthenticateMiddleware struct {
	handler http.Handler
}

func (am *AuthenticateMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) == 2 && (parts[1] == "users" || parts[1] == "authentications") {
		am.handler.ServeHTTP(w, r)
		return
	}
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != t.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if token != nil && err == nil {
		am.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	}
}

func NewAuthenticateMiddleware(handlerToWrap http.Handler) *AuthenticateMiddleware {
	return &AuthenticateMiddleware{handlerToWrap}
}
