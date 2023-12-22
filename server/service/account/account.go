package service

import (
	"context"
	"log"

	"github.com/Cprime50/Gopay/helper/errors"
)

// AccountService acts as a struct for injecting an implementation of AccountRepository
// for use in service methods
type AccountService struct {
	AccountRepository model.AccountRepository
	ImageRepository   model.ImageRepository
}

// hold repositories that will eventually be injected into Account service layer
type AccountConfig struct {
	AccountRepository model.AccountRepository
	ImageRepository   model.ImageRepository
}

// NewAccountService is a factory function for
// initializing an AccpuntService with its repository layer dependencies
func NewAccountService(c *AccountConfig) model.UserService {
	return &AccountService{
		AccountRepository: c.AccountRepository,
		ImageRepository:   c.ImageRepository,
	}
}

// GenerateAccountNumber generates a unique 10-digit account number
func (s *AccountService) Signup(ctx context.Context, account *model.Account) error {
	//hash password
	password, err := hashPassword(account.Password)
	if err != nil {
		log.Printf("Unable to hashpassword for account: %v, due to: %v\n", account.Email, err)
		return errors.NewInternal()
	}
	account.Password = hashedPassword

	// // Generate unique accountNumber
	// accountNumber, err := GenerateAccountNumber()
	// if err != nil {
	// 	log.Printf("Unable to hashpassword for account: %v, due to: %v\n", account.Email, err)
	// 	return errors.NewInternal()
	// 	return err
	// }
	// _,err := db.GetAccountByAccountNum()
	// if err == nil {
	// 	log.Printf("Could not create an account with email: %v. Reason: %v\n", account.Email, err.Code.Name())
	// 	return errors.NewConflict("email", account.Email)
	// 	// Log any other errors that may occur
	// } else if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	log.Printf("Could not create an account with email: %v. Reason: %v\n", account.Email, err)
	// 	return errors.NewInternal()
	// }
	if err := s.AccountRepository.CreateAccount(ctx, account); err != nil {
		log.Printf("Error creating account: %v", err)
		return errors.NewInternal()
	}
	return nil
}
