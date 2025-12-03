package auth

import (
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secretKey string, iat, seconds int64, accountID string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[vars.AccountKey] = accountID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
