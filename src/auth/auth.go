package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
)

type Auth struct {
	*jwtauth.JWTAuth
	devMode bool
}

func NewAuth(secret string, devMode bool) *Auth {
	return &Auth{jwtauth.New("HS256", []byte(secret), nil), devMode}
}

type TokenClaims struct {
	UserID    int64  `json:"userID"`
	UserEmail string `json:"userEmail"`
	ExpireIn  int64  `json:"exp"`
}

func TokenClaimsFromRequest(r *http.Request) *TokenClaims {
	_, claims, _ := jwtauth.FromContext(r.Context())

	tc := TokenClaims{
		UserID:    int64(claims["userID"].(float64)),
		UserEmail: claims["userEmail"].(string),
		ExpireIn:  claims["exp"].(time.Time).Unix(),
	}

	return &tc
}

func NewTokenClaims() *TokenClaims {
	t := &TokenClaims{}

	return t
}

func NewTokenClaimsForUser(userID int64, userEmail string) (*TokenClaims, error) {
	t := NewTokenClaims()

	t.UserID = userID
	t.UserEmail = userEmail

	return t, nil
}

func (tc *TokenClaims) asMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"userID":    tc.UserID,
		"userEmail": tc.UserEmail,
		"exp":       time.Now().Add(time.Minute * 3).Unix(),
	}
}

func (tc *TokenClaims) Encode(t *jwtauth.JWTAuth) string {
	mc := tc.asMapClaims()
	_, tokenString, _ := t.Encode(mc)
	return tokenString
}
