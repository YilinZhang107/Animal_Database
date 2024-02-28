/*
* @Author: Oatmeal107
* @Date:   2023/6/15 21:15
 */

package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("my-secret-key") // 设置密钥，用于签名验证

type Claims struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	//Grade     uint8   `json:"grade"`
	Authority int `json:"authority"`
	jwt.StandardClaims
}

// GenerateToken 生成 Token
func GenerateToken(id uint, username string, authority int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(7 * 24 * time.Hour) // 过期时间为7* 24 小时
	Claims := Claims{
		ID:       id,
		UserName: username,
		//Grade:     grade,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Animal_Database", // 签发人
		},
	}

	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	//使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString(jwtSecret) // 设置密钥，用于签名验证
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken 校验 Token
func VerifyToken(tokenString string) (claims *Claims, err error) {
	// 解析 Token, 后面的这是个回调函数,使用指定的secret签名并获得完整的编码后的字符串token
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 验证 Token 的有效性
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		//idStr := claims["id"]
		//id = uint(idStr.(float64))
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
