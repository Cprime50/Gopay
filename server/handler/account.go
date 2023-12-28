package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Cprime50/Gopay/helper"
	models "github.com/Cprime50/Gopay/models/account"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// signupReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
type signupReq struct {
	Firstname       string `json:"first_name" binding:"required,min=3,max=255,alphanum"` //alphanum is to get only alphanuneric chars, meaning just alpabets and numbers without any special chars
	Lastname        string `json:"last_name" binding:"required,min=3,max=255,alphanum"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=255"`
}

type signinReq struct {
	Email    string `json:"email" binding:"required,email,min=3,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

// Signup handler
func (h *Handler) Signup(c *gin.Context) {
	// define a variable to which we'll bind incoming
	// json body, {email, password}
	var input signupReq

	// Bind form input as JSON
	if err := c.BindJSON(&input); err != nil {
		// check for error in required field
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage, "message": "please correctly provide the relevant fields"})
		return
	}

	// Ensure password provided and confirmedPassword match
	if input.Password != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		c.Abort()
	}

	// update the user table with new data
	account := &models.Account{
		FirstName: input.Firstname,
		LastName:  input.Lastname,
		Email:     input.Email,
		Password:  input.Password,
		Balance:   500,
		RoleID:    2,
		IsActive:  false,
	}

	ctx := c.Request.Context()
	err := h.AccountService.Signup(ctx, account)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.AbortWithStatusJSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// create token pair as strings
	tokens, err := h.TokenService.NewPairFromUser(ctx, account, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to clear/delete the created user in the database

		c.AbortWithStatusJSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "account created successfully",
		"user":    account.Email,
		"tokens":  tokens,
	})
}

// Signin used to authenticate extant user
func (h *Handler) Signin(c *gin.Context) {
	var input signinReq

	// Bind input
	if err := c.BindJSON(&input); err != nil {
		// check for error in required field
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}

		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage, "message": "please correctly provide the relevant fields"})
		return
	}

	account := &models.Account{
		Email:    input.Email,
		Password: input.Password,
	}

	ctx := c.Request.Context()
	err := h.AccountService.Signin(ctx, account)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		c.AbortWithStatusJSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, account, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		c.AbortWithStatusJSON(helper.Status(err), gin.H{
			"message": "Login successful",
			"user":    account.Email,
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
