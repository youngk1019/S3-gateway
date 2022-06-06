package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"s3-gateway/command/vars"
	"strings"
)

func ParseJWT(token string) (jwt.MapClaims, error) {
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err := parser.ParseUnverified(token, claims)
	return claims, err
}

func FetchJWTToken(r *http.Request) string {
	defer func() {
		r.Header.Del(vars.JWTHeader)
		r.URL.Query().Del(vars.JWTQuery)
	}()
	token := r.Header.Get(vars.JWTHeader)
	if strings.HasPrefix(token, "Bearer ") || strings.HasPrefix(token, "bearer ") {
		return token[7:]
	}

	token = r.URL.Query().Get(vars.JWTQuery)
	if token != "" {
		return token
	}

	cookie, err := r.Cookie(vars.JWTCookie)
	if err != nil {
		return ""
	}
	return cookie.Value
}
