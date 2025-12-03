package password

import "golang.org/x/crypto/bcrypt"

// Encrypt 加密密码
func Encrypt(password string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashPassword)
}

// Verify 验证密码
func Verify(password string, hashPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)) == nil
}
