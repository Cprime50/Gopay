package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Cprime50/Gopay/helper"
	"github.com/Cprime50/Gopay/services/password"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	DB *gorm.DB
}

// AccountNumberGenerator generates unique 10-digit account numbers.
type AccountNumberGenerator struct {
	counter int64
	mu      sync.Mutex
}

// NewAccountRepository initializes new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		DB: db,
	}
}

type Account struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID `gorm:"primary_key"`
	Email         string    `gorm:"uniqueIndex;not null;type:varchar(250)" json:"email"`
	AccountNumber int64     `gorm:"type:varchar(100);uniqueIndex;column:account_number;not null"`
	Balance       float64   `gorm:"type:decimal(10,2)"`
	FirstName     string    `gorm:"type:varchar(100);not null"`
	LastName      string    `gorm:"type:varchar(100);not null"`
	Password      string    `gorm:"type:varchar(100);not null"`
	ImageUrl      string    `gorm:"image_url" json:"imageUrl"`
	RoleID        uint      `gorm:"not null;DEFAULT:4" json:"role_id"`
	Role          Role      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	IsActive      bool      `gorm:"type:boolean;not null"`
}

// Create new acount in DB
func (db *AccountRepository) CreateAccount(ctx context.Context, account *Account) error {
	var gen *AccountNumberGenerator
	// Check if an account for the given email already exists
	_, err := db.GetAccountByEmail(ctx, account.Email)
	if err == nil {
		log.Printf("Could not create an account with email: %v. Reason: %v\n", account.Email, err.Code.Name())
		return helper.NewConflict("email", account.Email)
		// Log any other errors that may occur
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Could not create an account with email: %v. Reason: %v\n", account.Email, err)
		return helper.NewInternal()
	}
	accountNumber := gen.GenerateAccountNumber(ctx, db)
	account = &Account{
		Email:         account.Email,
		AccountNumber: accountNumber,
		FirstName:     account.FirstName,
		LastName:      account.LastName,
		Password:      account.Password,
		Balance:       500,
		RoleID:        2,
		IsActive:      false,
	}

	if err := db.DB.Create(&account).Error; err != nil {
		return err
	}
	return nil
}

// GenerateAccountNumber generates a unique 10-digit account number.
func (gen *AccountNumberGenerator) GenerateAccountNumber(ctx context.Context, db *AccountRepository) int64 {
	gen.mu.Lock()
	defer gen.mu.Unlock()

	for {
		gen.counter++
		accountNumber := gen.counter % 1e10
		// Check if account number is unique
		_, err := db.GetAccountByAccountNum(ctx, accountNumber)
		if err == nil {
		}
		// If not unique, regenerate and check again
	}
}

// Gets a users account from db based on account number
func (db *AccountRepository) GetAccountByAccountNum(ctx context.Context, account_number int64) (*Account, error) {
	var account *Account
	err := db.DB.WithContext(ctx).Where("account_number = ?", account_number).First(&account).Error
	if err != nil {
		return &Account{}, helper.NewNotFound("account_number", account_number.String())
	}
	return account, nil
}

func (db *AccountRepository) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	var account *Account
	err := db.DB.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		return &Account{}, helper.NewNotFound("email", email)
	}
	return account, nil
}

func (db *AccountRepository) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	var account *Account
	err := db.DB.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		return &Account{}, helper.NewNotFound("id", id.String())
	}
	return account, nil
}

// Update Account details
func (db *AccountRepository) UpdateAccount(ctx context.Context, account *Account) error {
	// Update only the fields that are filled in the account
	// Omit sensitive fields like password, balance, and role
	err := db.DB.WithContext(ctx).Omit("password", "balance", "role").Updates(account).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Errorf("account with ID %s not found", account.ID)
			return helper.NewNotFound("id", account.ID.String())
		}
		log.Println("Error querying account:", err)
		return helper.NewInternal()
	}
	return nil
}

func (db *AccountRepository) ChangeUserStatus(ctx context.Context, account *Account) error {
	// Check if the accounts exists based on the provided accountID
	_, err := db.GetAccountByID(ctx, account.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Errorf("account with ID %s not found", account.ID)
			return helper.NewNotFound("id", account.ID.String())
		}
		log.Println("Error querying account:", err)
		return helper.NewInternal()
	}

	// Create a map of columns and their values that you want to update
	updateColumns := map[string]interface{}{
		"IsActive": account.IsActive,
	}

	// Update only the specified columns in the database
	if err := db.DB.WithContext(ctx).Model(&account).Updates(updateColumns).Error; err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

// Update password
func (db *AccountRepository) ResetPassword(ctx context.Context, account *Account) error {
	// Hash the new password
	hashedPassword, err := password.HashPassword(account.Password)
	if err != nil {
		return err
	}

	// Update user password where username, email match and the user is active
	err = db.DB.WithContext(ctx).Model(&Account{}).
		Where("email = ?", account.Email).
		Updates("password", hashedPassword).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Errorf("account with email %s not found", account.Email)
			return helper.NewNotFound("email", account.Email)
		}
		log.Println("Error querying account:", err)
		return helper.NewInternal()
	}
	return nil
}

// Upload Image by file to cloudinary and save the img url to db
func (db *AccountRepository) UpdateImgByFile(ctx context.Context, id uuid.UUID, imgFile *File) error {
	media := NewImageRepository()
	imageURL, err := media.FileUpload(imgFile)
	if err != nil {
		log.Println("Error uploaidng image by file", err)
		return helper.NewInternal()
	}
	err = db.DB.WithContext(ctx).Model(&Account{}).
		Where("id = ?", id).
		Updates("image", imageURL).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Errorf("account with ID %s not found", id)
			return helper.NewNotFound("id", id.String())
		}
		log.Println("Error querying account:", err)
		return helper.NewInternal()
	}
	return nil
}

// Upload Image by url to cloudinary and save the img url to db
func (db *AccountRepository) UpdateImgByUrl(ctx context.Context, id uuid.UUID, imgUrl *Url) error {
	media := NewImageRepository()
	imageURL, err := media.RemoteUpload(imgUrl)
	if err != nil {
		log.Println("Error uploaidng image by Url", err)
		return helper.NewInternal()
	}
	err = db.DB.WithContext(ctx).Model(&Account{}).
		Where("id = ?", id).
		Updates("image", imageURL).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Errorf("account with ID %s not found", id)
			return helper.NewNotFound("id", id.String())
		}
		log.Println("Error querying account:", err)
		return helper.NewInternal()
	}
	return nil
}

// Get all account
func (db *AccountRepository) GetAllUsers(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	if err := db.DB.WithContext(ctx).Omit("password").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// Delete User
func (db *AccountRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	var account *Account
	if err := db.DB.WithContext(ctx).Where("id = ?", id).Delete(&account).Error; err != nil {
		return err
	}
	return nil
}
