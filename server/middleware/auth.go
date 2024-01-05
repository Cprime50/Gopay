package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Cprime50/Gopay/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// AuthUser extracts a user from the Authorization header
// which is of the form "Bearer token"
// It sets the user to the context if the user exists
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// validate token here
		account, err := ValidateJWT(token)

		if err != nil {
			err := helper.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("account", account)

		c.Next()
	}
}

// JWTAuthAdminMiddleware is a middleware function that checks for both regular authentication and admin privileges.
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		token, err := extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Validate token for regular user authentication
		account, err := ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Validate admin role
		accountAdmin, _err := ValidateAdminJWT(token)
		if _err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only Administrator is allowed to perform this action"})
			c.Abort()
			return
		}

		// Set account information to the context
		c.Set("account", account)
		c.Set("accountAdmin", accountAdmin)

		c.Next()
	}
}

// extractTokenFromHeader extracts the token from the Authorization header.
func extractTokenFromHeader(c *gin.Context) (string, error) {
	h := authHeader{}

	// bind Authorization Header to h and check for validation errors
	if err := c.ShouldBindHeader(&h); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			// we used this type in bind_data to extract desired fields from errs
			// you might consider extracting it
			var invalidArgs []invalidArgument
			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}
			err := helper.NewBadRequest("Invalid request parameters. See invalidArgs")
			fmt.Println("invalid request param for auth header", err)
			return "", err
		}

		// otherwise error type is unknown
		err := helper.NewInternal()
		fmt.Println("error", err)
		return "", err
	}

	tokenHeader := strings.Split(h.Token, "Bearer ")
	if len(tokenHeader) < 2 {
		err := helper.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")
		return "", err
	}

	return tokenHeader[1], nil
}
