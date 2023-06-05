package middleware

import (
	"fmt"
	"net/http"
	"nyoba/configg"
	"nyoba/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// middleware akan di eksekusi dulu sebelum validate
func RequireAuth(c *gin.Context) {

	// Ambil cookie request
	tokenString, err := c.Cookie("Authorizion")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Can't get cookie",
		})
	}

	// Decode cookie yang didapat
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Check expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "You are exp, try to login again",
			})
		}

		// Temukan user yang sedang login
		var user models.Users

		id := claims["id_user"].(float64)
		configg.KoneksiData().Debug().First(&user, int(id))

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error user can't found",
			})
		}

		// Melampirkan ke request
		c.Set("user", user)
		// Lanjutkan
		c.Next()
	} else {
		fmt.Println(err)
	}

}
