package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// GenerateFromPassword จะรับพาสเวิร์ดเป็นไบต์และคืนเป็นไบต์แฮชและข้อผิดพลาด
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, password string) bool {
	// CompareHashAndPassword จะเปรียบเทียบรหัสผ่านแฮชกับรหัสผ่านที่เป็นสตริงและคืนค่า true หากตรงกัน
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
