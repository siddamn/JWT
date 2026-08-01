package middlewares

import (
	"fmt"
	"jwt/initialisers"
	"jwt/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	//fmt.Println("RequireAuth middleware called")
	//get the cookie from the request, validate the token, find user with the token, attach the user to the request context and continue to the next middleware or handler function

	cookie, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	//validate the token
	//ParseWithClaims parses the JWT token from the cookie and validates it using the provided secret key. It checks the signing method and returns an error if the token is invalid or if the signing method is unexpected. If the token is valid, it allows the request to proceed to the next middleware or handler function.
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	//we have a valid token, we can extract the claims and use them to find the user in the database. For now, we will just print the claims to the console.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	//check expiration
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "expires",
		})
		return
	}

	//user from sub
	var userfromtoken models.User
	//returns rows affected, error
	//err = initialisers.DB.First(&userfromtoken, claims["sub"]).Error
	err = initialisers.DB.Where("id = ?", claims["sub"]).First(&userfromtoken).Error
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	//set the user in the context
	c.Set("user", userfromtoken)

	c.Next()
}
