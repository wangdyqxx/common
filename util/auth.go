package util

import "github.com/dgrijalva/jwt-go"

var jwtKey = "aaaaa"

// 生成awt的token
func CreateToken() (string, error) {
	info := jwt.MapClaims{
		"test1:": "test1",
		"test2":  "test2",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
	return token.SignedString([]byte(jwtKey))
}

// 解析token
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
