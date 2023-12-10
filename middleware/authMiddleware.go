package middleware

import (
	"gin-api/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// auth function
func EnsureAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		token = strings.Split(token, " ")[1]
		claims, err := helpers.DecodeToken(token)

		roleType, ok := claims["roleType"].(string)
		if !ok || roleType != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("userLogin", claims)

		// Optionally set additional flags based on the roleType
		if roleType == "admin" {
			c.Set("isAdmin", true)
		}

		c.Next()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func EnsureCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		token = strings.Split(token, " ")[1]
		claims, err := helpers.DecodeToken(token)

		roleType, ok := claims["roleType"].(string)
		if !ok || roleType != "customer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("userLogin", claims)

		// Optionally set additional flags based on the roleType
		if roleType == "customer" {
			c.Set("isAdmin", true)
		}

		c.Next()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
