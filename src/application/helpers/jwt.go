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

func ParseJWT(tokenString string) (string, float64, string, float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		tenant_id := claims["tenant_id"].(float64)
		role := claims["role"].(string)
		user_id := claims["user_id"].(float64)
		return username, tenant_id, role, user_id, nil
	}

	return "", 0, "", 0, err
}
