package gontracts

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

// AuthHandler authentication handler
type AuthHandler struct {
	jwtmiddleware.JWTMiddleware
	secretKey []byte
}

// NewAuthHandler creates new authentication handler
func NewAuthHandler(key []byte) *AuthHandler {
	return &AuthHandler{
		JWTMiddleware: *jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return key, nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		}),
		secretKey: key,
	}
}

// GenerateToken returns new authentication token
func (a *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	// create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// sign token with a key
	tokenString, _ := token.SignedString(a.secretKey)

	w.Write([]byte(tokenString))
}

// HandlerFunc is a function handler adepter that wraps handler function into auth handler object
func (a *AuthHandler) HandlerFunc(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return a.Handler(http.HandlerFunc(f))
}
