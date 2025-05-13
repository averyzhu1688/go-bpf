package utils

import (
	"errors"
	"time"

	"bpf.com/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

// custom the jwt claim
type JWTClaims struct {
	UserId   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// generator access token
func GenerateAccessToken(userId uint64) (string, error) {
	cfg := config.GetAppConfig().JWT
	claims := JWTClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.AccessTokenExp) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.TokenIssuer,
		},
	}
	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign token
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// generate the refreshtoken
func GenerateRefreshToken(userId uint64) (string, error) {
	cfg := config.GetAppConfig().JWT
	claims := JWTClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.RefreshTokenExp) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.TokenIssuer,
		},
	}
	//create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign the token
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// parse token
func ParseAccessToken(tokenString string) (*JWTClaims, error) {
	return parseToken(tokenString)
}

// parse refreshtoken
func ParseRefreshToken(refreshTokenString string) (*JWTClaims, error) {
	return parseToken(refreshTokenString)
}

// parse token
func parseToken(tokenString string) (*JWTClaims, error) {
	cfg := config.GetAppConfig().JWT

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// check sign
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	// verify token and convert to custom claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// validate token
func ValidateToken(tokenString string) bool {
	_, err := ParseAccessToken(tokenString)
	return err == nil
}
