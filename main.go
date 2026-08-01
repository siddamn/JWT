package main

import (
	"fmt"
	"jwt/controllers"
	"jwt/initialisers"
	"jwt/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initialisers.LoadEnv()
	initialisers.ConnectToDb()
	initialisers.SyncDB()
}
func main() {
	fmt.Println("Hello, World!")

	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	//checks jwt token in the cookie and returns the user if valid, otherwise returns 401
	r.GET("/validate", middlewares.RequireAuth, controllers.ValidateToken)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
