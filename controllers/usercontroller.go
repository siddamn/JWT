package controllers

import (
	"jwt/initialisers"
	"jwt/models"
	"net/http"
	"os"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//hash the password
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//store the user in the database
	var user models.User
	user.Email = body.Email
	user.Password = string(bcryptPassword)

	result := initialisers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var inputEmail string = body.Email
	var inputPassword string = body.Password

	var user models.User
	//find the user in the database by email
	err = initialisers.DB.Where("email = ?", inputEmail).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "User not found",
		})
		return
	}

	//compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid password",
		})
		return
	}

	//generate a JWT token for the user (this part is not implemented in the provided code)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	//sign the token with a secret key (replace "your_secret_key" with
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	//set the token in a cookie
	//cookie is helpful for storing the JWT token on the client side, allowing the client to send it with subsequent requests for authentication and authorization purposes.
	//In this case, it is set to SameSiteLaxMode, which allows the cookie to be sent with top-level navigations and GET requests initiated by third-party websites, but not with other types of cross-site requests.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", signedToken, 3600*24*7, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": "Login successful",
	})
}

func ValidateToken(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	c.JSON(200, gin.H{
		"email": user.(models.User).Email,
	})
}
