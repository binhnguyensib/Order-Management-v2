package handler

import (
	"intern-project-v2/domain"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *authHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

func (ah *authHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var customer domain.CustomerRegiser
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	newCustomer, err := ah.authUsecase.Register(ctx, &customer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, newCustomer)
}

func (ah *authHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var loginReq domain.CustomerLogin
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	customer, token, err := ah.authUsecase.Login(ctx, &loginReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"customer": customer, "token": token})
}
