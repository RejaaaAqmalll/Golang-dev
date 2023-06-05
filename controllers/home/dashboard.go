package home

import (
	"net/http"
	"nyoba/models"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {
	//  Get name and email
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"email": user.(models.Users).Email,
		"name":  user.(models.Users).Name,
	})

}
