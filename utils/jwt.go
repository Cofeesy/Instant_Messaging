package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"gin_chat/utils/setting"
)

var JwtSecret = []byte(setting.JwtSecret)

// 自定义的
type Claims struct {
	ID uint
	Username string `json:"username"`
	jwt.StandardClaims
}


// 生成token
func GenerateToken(id uint,username string)(string, error){
	Claims:=Claims{
		ID:id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			//30天后过期
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			// ExpiresAt: time.Now().Add(24 *time.Hour ).Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "gin_chat",
		},
	}

	// 创建jwt_token
	tokenClaims:=jwt.NewWithClaims(jwt.SigningMethodHS256,Claims)

	// 使用密钥签名
	token, err := tokenClaims.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 解析token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
}




