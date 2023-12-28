package service

import (
	"context"
	"errors"
	"log"

	"github.com/Cprime50/Gopay/helper"
	models "github.com/Cprime50/Gopay/models/account"
	"gorm.io/gorm"
)

// AccountService acts as a struct for injecting an implementation of AccountRepository
// for use in service methods
type AccountService struct {
	AccountRepository models.AccountRepository
	ImageRepository   models.ImageRepository
}

// hold repositories that will eventually be injected into Account service layer
type AccountConfig struct {
	AccountRepository models.AccountRepository
	ImageRepository   models.ImageRepository
}

// NewAccountService is a factory function for
// initializing an AccpuntService with its repository layer dependencies
func NewAccountService(c *AccountConfig) *AccountService {
	return &AccountService{
		AccountRepository: c.AccountRepository,
		ImageRepository:   c.ImageRepository,
	}
}

// Hashes users password and creates a new user account
func (s *AccountService) Signup(ctx context.Context, account *models.Account) error {
	//hash password
	hashedPassword, err := models.HashPassword(account.Password)
	if err != nil {
		log.Printf("Unable to hashpassword for account: %v, due to: %v\n", account.Email, err)
		return helper.NewInternal()
	}
	account.Password = hashedPassword

	//model layer will handle generatingn account number and initilizing user inputed data
	if err := s.AccountRepository.CreateAccount(ctx, account); err != nil {
		log.Printf("Error creating account: %v", err)
		return helper.NewInternal()
	}
	return nil
}

// Signin reaches our to a AccountRepository check if the user exists
// and then compares the supplied password with the provided password
// if a valid email/password combo is provided, u will hold all
// available account fields
func (s *AccountService) Signin(ctx context.Context, account *models.Account) error {
	accountGotten, err := s.AccountRepository.GetAccountByEmail(ctx, account.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NewNotFound("Email not found, create account", err.Error())
		}
		return helper.NewInternal()
	}
	// verify password - check if matches
	match, err := ComparePassword(accountGotten.Password, account.Password)
	if err != nil {
		return helper.NewInternal()
	}
	if !match {
		return helper.NewAuthorization("Invalid email and password combination")
	}

	*account = *accountGotten
	return nil
}
