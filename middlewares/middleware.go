package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// รับ header ตย. Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// ตัดเอา Bearer ออกให้เหลือ แต่ token
		tokenString := strings.Split(authHeader, " ")[1]
		// ตรวจสอบ ว่า token ตรงกัน หรือ หมดอายุใหม return token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// ถ้า Token ถูกต้อง ให้ตั้งค่า Bearer Token ใน Header
			c.Header("Authorization", authHeader)
			//set username ใน claims
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
