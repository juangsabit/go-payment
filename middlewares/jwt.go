package middlewares

import (
	"fmt"
	"go-payment/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// err := token.TokenValid(c)
		tokenString, err := c.Cookie("Authorization")

		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("API_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
				return
			}

			// find the user with token sub
			var user models.User
			models.DB.Debug().First(&user, claims["user_id"])

			if user.ID == 0 {
				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
				return
			}

			// attach to req
			c.Set("user", user)
			c.Set("userID", user.ID)

			// continue
			c.Next()
		} else {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
	}
}
