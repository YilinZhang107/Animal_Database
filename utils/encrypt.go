/*
* @Author: Oatmeal107
* @Date:   2023/6/15 20:28
 */

package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 后端加密密码  bcrypt加密
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword 比较密码
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
