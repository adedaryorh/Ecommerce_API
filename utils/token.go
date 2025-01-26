package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTToken struct {
	config *Config
}
type jwtClaim struct {
	jwt.StandardClaims
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func NewJWTToken(config *Config) *JWTToken {
	return &JWTToken{config: config}
}

func (j *JWTToken) CreateToken(user_id int64) (string, error) {
	claims := jwtClaim{
		UserId: user_id,
		Exp:    time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.config.Signing_key))
	if err != nil {
		return "", err
	}
	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Authentication token")
		}
		return []byte(j.config.Signing_key), nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("Invalid Authentication token")
	}
	claims, ok := token.Claims.(*jwtClaim)
	if !ok {
		return 0, "", fmt.Errorf("Invalid Authentication token")
	}
	fmt.Println("Token Claims - UserId:", claims.UserId, "Role:", claims.Role) // Debug log
	if claims.Exp < time.Now().Unix() {
		return 0, "", fmt.Errorf("token expired")
	}
	return claims.UserId, claims.Role, nil // Return userId and role
}
