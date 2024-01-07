package controllers

import (
	"fmt"
	"net/http"

	"github.com/CyberTea0X/delta_art/src/backend/models"
	"github.com/CyberTea0X/delta_art/src/backend/utils/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)


func (p *PublicController) CurrentUser(c *gin.Context){
    fmt.Println(c.Request.RequestURI)
    jwt_token := c.MustGet("access_token").(*jwt.Token)

    
    
	user_id, err := token.ExtractUint(jwt_token, "user_id")
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	u,err := models.GetUserByID(p.DB, user_id)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message":"success","data":u})
}
