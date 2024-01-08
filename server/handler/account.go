package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Cprime50/Gopay/helper"
	"github.com/Cprime50/Gopay/middleware"
	models "github.com/Cprime50/Gopay/models/account"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
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
	}

	ctx := c.Request.Context()
	hashedPassword, err := models.HashPassword(account.Password)
	if err != nil {
		log.Printf("Unable to hashpassword for account: %v, due to: %v\n", account.Email, err)
		c.AbortWithStatusJSON(helper.Status(err), gin.H{"message": "Internal server error"})
		return
	}
	account.Password = hashedPassword

	//model layer will handle generatingn account number and initilizing user inputed data
	if err := models.CreateAccount(ctx, account); err != nil {
		log.Printf("Error creating account: %v", err)
		c.AbortWithStatusJSON(helper.Status(err), gin.H{"message": "Internal server error"})
		return
	}

	// create token pair as strings
	tokens, err := middleware.NewPairFromUser(ctx, account, "")

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
	accountGotten, err := models.GetAccountByEmail(ctx, account.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error(), "message": "Email not found, create account"})
			return

		}
		log.Println("Failed to sign in user: error getting user form db", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Something went wrong"})
		return

	}
	// verify password - check if matches
	match, err := helper.ComparePassword(accountGotten.Password, account.Password)
	if err != nil {
		log.Println("error comparing password", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Something went wrong"})
		return
	}
	if !match {
		c.AbortWithStatusJSON(helper.Status(err), gin.H{"error": err.Error(), "message": "Invalid email and password combination"})
		return
	}

	tokens, err := middleware.NewPairFromUser(ctx, account, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		c.AbortWithStatusJSON(helper.Status(err), gin.H{
			"message": "Login unsuccessful",
			"account": account.Email,
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens":  tokens,
		"account": account.Email,
	})
}

// Get user

// Get Users

// Update User

// send OTP

// Forgot Password

// To do, fix bug in tokens where admin not generating fresh token
// Get me
func (h *Handler) Me(c *gin.Context) {
	// A *model.User will eventually be added to context in middleware
	account, exists := c.Get("account")

	// This shouldn't happen, as our middleware ought to throw an error.
	// This is an extra safety measure
	// We'll extract this logic later as it will be common to all handler
	// methods which require a valid user
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := helper.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	id := account.(*models.Account).ID

	// use the Request Context
	ctx := c.Request.Context()

	acct, err := models.GetAccountByID(ctx, id)

	if err != nil {
		log.Printf("Unable to find account: %v\n%v", id, err)
		e := helper.NewNotFound("account", id.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account": acct,
	})
}

// Update Me
