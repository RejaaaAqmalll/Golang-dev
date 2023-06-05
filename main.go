package main

import (
	"nyoba/configg"
	"nyoba/controllers/auth"
	"nyoba/controllers/home"
	"nyoba/middleware"
	"nyoba/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db := configg.KoneksiData()
	db.AutoMigrate(&models.Users{})

	r := gin.Default()

	//Routes
	u := r.Group("/user")
	{
		u.POST("/register", auth.Register)
		u.POST("/login", auth.Login)
		u.GET("/validate", middleware.RequireAuth, auth.Validate)
		u.POST("/logout", middleware.RequireAuth, auth.Logout)
		u.POST("/forgotpw", auth.ForgotPassword)
		u.POST("/auth_verify", auth.Authyverify)
		u.POST("/reset-password", auth.ResetPassword)
		u.GET("/dashboard", middleware.RequireAuth, home.Dashboard)
	}

	r.Run(":9000")

}
