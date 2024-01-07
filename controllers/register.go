package controllers

import (
	"net/http"
	"github.com/CyberTea0X/delta_art/src/backend/models"
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (p *PublicController) Register(c *gin.Context) {

    var input RegisterInput

    if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    user := models.User{}

	user.Username = input.Username
	user.Password = input.Password
    user.Email = input.Email

    if !models.IsValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
        return
    }

    if _, err := user.SaveUser(p.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
    }

    c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

