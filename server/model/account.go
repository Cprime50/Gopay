package models

import (
	"context"
	"errors"
	"log"

	"github.com/Cprime50/Gopay/helper/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	DB *gorm.DB
}

// NewAccountRepository initializes new account repository
func NewAccountRepository(db *gorm.DB) model.AccountRepository {
	return &AccountRepository{
		DB: db,
	}
}

type Account struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID `gorm:"primary_key"`
	Email         string    `gorm:"uniqueIndex;not null;type:varchar(250)" json:"email"`
	AccountNumber int64     `gorm:"type:varchar(100);uniqueIndex;column:account_number"`
	Balance       float64   `gorm:"type:varchar(100)"`
	FirstName     string    ``
	LastName      string    ``
	Password      string    ``
	Image         string    ``
	Role          Role
}

func (db *AccountRepository) CreateAccount(ctx context.Context, account *Account) error {

	// Check if an account for the given email already exists
	var existingAccount Account
	_, err := db.GetAccountByEmail()
	if err == nil {
		log.Printf("Could not create an account with email: %v. Reason: %v\n", account.Email, err.Code.Name())
		return errors.NewConflict("email", account.Email)
		// Log any other errors that may occur
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Could not create an account with email: %v. Reason: %v\n", account.Email, err)
		return errors.NewInternal()
	}

	if err := db.DB.Create(&account).Error; err != nil {
		return err
	}
	return nil
}

func (db *AccountRepository) GetAccountByAccountNum(ctx context.Context, account_number int64) (*Account, error) {
	var account *Account
	err := db.DB.WithContext(ctx).Where("account_number = ?", account_number).First(&account).Error
	if err != nil {
		return &Account{}, err
	}
	return account, nil
}

func (db *AccountRepository) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	var account *Account
	err := db.DB.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		return &Account{}, err
	}
	return account, nil
}
