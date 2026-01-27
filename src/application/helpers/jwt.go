package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("minha_chave_secreta") // troque por algo seguro

func GenerateJWT(username string, tenant_id int64, role string, userId int64) (string, error) {
	claims := jwt.MapClaims{
		"role":      role,
		"user_id":   userId,
		"username":  username,
		"tenant_id": tenant_id,
		"exp":       time.Now().Add(time.Hour * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type BillingClaims struct {
	CanWrite bool
}

func GenerateJWTWithBilling(username string, tenant_id int64, role string, userId int64, billing BillingClaims) (string, error) {
	claims := jwt.MapClaims{
		"role":              role,
		"user_id":           userId,
		"username":          username,
		"tenant_id":         tenant_id,
		"billing_can_write": billing.CanWrite,
		"exp":               time.Now().Add(time.Hour * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (string, float64, string, float64, error) {
	claims, err := parseJWTClaims(tokenString)
	if err != nil {
		return "", 0, "", 0, err
	}
	username := claims["username"].(string)
	tenant_id := claims["tenant_id"].(float64)
	role := claims["role"].(string)
	user_id := claims["user_id"].(float64)
	return username, tenant_id, role, user_id, nil
}

func ParseJWTWithBilling(tokenString string) (string, float64, string, float64, BillingClaims, error) {
	claims, err := parseJWTClaims(tokenString)
	if err != nil {
		return "", 0, "", 0, BillingClaims{}, err
	}

	username := claims["username"].(string)
	tenant_id := claims["tenant_id"].(float64)
	role := claims["role"].(string)
	user_id := claims["user_id"].(float64)

	billing := BillingClaims{}
	if canWrite, ok := claims["billing_can_write"].(bool); ok {
		billing.CanWrite = canWrite
	}

	return username, tenant_id, role, user_id, billing, nil
}

func parseJWTClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || token == nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
