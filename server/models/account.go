package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Cprime50/Gopay/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type AccountRepository struct {
// 	DB *gorm.DB
// }

// AccountNumberGenerator generates unique 10-digit account numbers.
type AccountNumberGenerator struct {
	counter int64
	mu      sync.Mutex
}

// // NewAccountRepository initializes new account repository
// func NewAccountRepository() *AccountRepository {
// 	ds, err := InitDS()
// 	if err != nil {
// 		log.Fatal("Failed to initialize DataSources:", err)
// 	}
// 	return &AccountRepository{
// 		DB: ds.DB,
// 	}
// }

type Account struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID `gorm:"primary_key"`
	Email         string    `gorm:"uniqueIndex;not null;type:varchar(250)" json:"email"`
	AccountNumber int64     `gorm:"uniqueIndex;column:account_number;not null"`
	Balance       float64   `gorm:"type:decimal(10,2)"`
	FirstName     string    `gorm:"type:varchar(100);not null"`
	LastName      string    `gorm:"type:varchar(100);not null"`
	Password      string    `gorm:"type:varchar(100);not null"`
	ImageUrl      string    `gorm:"image_url" json:"imageUrl"`
	RoleID        uint      `gorm:"not null;DEFAULT:4" json:"role_id"`
	Role          Role      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	IsActive      bool      `gorm:"type:boolean"`
}

// CreateAccount creates a new account in the database
func (db *DataSources) CreateAccount(ctx context.Context, account *Account) error {
	// Check if an account with the given email already exists
	_, err := db.GetAccountByEmail(ctx, account.Email)
	switch {
	case err == nil:
		log.Printf("Could not create an account with email: %v. Reason: Account already exists\n", account.Email)
		return helper.NewConflict("email", account.Email)

	case errors.Is(err, gorm.ErrRecordNotFound):
		break

	default:
		log.Printf("Error checking account existence: %v\n", err)
		return helper.NewInternal()
	}

	// Generate account number
	var gen *AccountNumberGenerator
	accountNumber := gen.GenerateAccountNumber(ctx, db)

	// initialize account number
	newAccount := &Account{
		AccountNumber: accountNumber,
	}
	if err := db.DB.Create(newAccount).Error; err != nil {
		log.Printf("Error creating account: %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

// GenerateAccountNumber generates a unique 10-digit account number.
func (gen *AccountNumberGenerator) GenerateAccountNumber(ctx context.Context, db *DataSources) int64 {
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

// GetAccountByAccountNum gets a user's account from the database based on the account number
func (ds *DataSources) GetAccountByAccountNum(ctx context.Context, accountNumber int64) (*Account, error) {
	db := ds.GetDB()
	var account Account
	err := db.WithContext(ctx).Where("account_number = ?", accountNumber).First(&account).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %d not found", accountNumber)
			return nil, helper.NewNotFound("account_number", fmt.Sprintf("%d", accountNumber))
		}
		log.Println("Error querying account:", err)
		return nil, helper.NewInternal()
	}
	return &account, nil
}

// GetAccountByEmail gets a user's account from the database based on the email
func (ds *DataSources) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {

	db := ds.GetDB()
	if db == nil {
		log.Println("DB is nil")
		return nil, helper.NewInternal()
	}
	var account Account
	err := db.Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with email %s not found", email)
			return nil, helper.NewNotFound("email", email)
		}
		log.Println("Error querying account:", err)
		return nil, helper.NewInternal()
	}
	return &account, nil
}

// GetAccountByID gets a user's account from the database based on the ID
func (db *DataSources) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	var account Account
	err := db.DB.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %s not found", id)
			return nil, helper.NewNotFound("id", id.String())
		}
		log.Println("Error querying account in db:", err)
		return nil, helper.NewInternal()
	}
	return &account, nil
}

// Update Account details
func (db *DataSources) UpdateAccount(ctx context.Context, account *Account) error {
	// Update only the fields that are filled in the account
	// Omit sensitive fields like password, balance, and role
	err := db.DB.WithContext(ctx).Omit("password", "balance", "role").Updates(&account).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %s not found", account.ID)
			return helper.NewNotFound("id", account.ID.String())
		}
		log.Println("Error querying account:", err)
		return helper.NewInternal()
	}
	return nil
}

func (db *DataSources) ChangeUserStatus(ctx context.Context, account *Account) error {
	// Check if the accounts exists based on the provided accountID
	_, err := db.GetAccountByID(ctx, account.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %s not found", account.ID)
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
		return helper.NewInternal()
	}
	return nil
}

// Resetpassword updates the users password
func (db *DataSources) ResetPassword(ctx context.Context, account *Account) error {
	// Hash the new password
	hashedPassword, err := HashPassword(account.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return helper.NewInternal()
	}

	// Update user password where username, email match and the user is active
	err = db.DB.WithContext(ctx).Model(&Account{}).
		Where("email = ?", account.Email).
		Updates(map[string]interface{}{"password": hashedPassword}).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with email %s not found\n", account.Email)
			return helper.NewNotFound("email", account.Email)
		}
		log.Println("Error reseting password", err)
		return helper.NewInternal()
	}
	return nil
}

// Upload Image by file to cloudinary and save the img url to db
func (db *DataSources) UpdateImgByFile(ctx context.Context, id uuid.UUID, imgFile *File) error {
	media := NewImageRepository()
	imageURL, err := media.FileUpload(imgFile)
	if err != nil {
		log.Println("Error uploaidng image by file", err)
		return helper.NewInternal()
	}
	err = db.DB.WithContext(ctx).Model(&Account{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"image": imageURL}).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %s not found", id)
			return helper.NewNotFound("id", id.String())
		}
		log.Println("Error querying account to update image:", err)
		return helper.NewInternal()
	}
	return nil
}

// Upload Image by url to cloudinary and save the img url to db
func (db *DataSources) UpdateImgByUrl(ctx context.Context, id uuid.UUID, imgUrl *Url) error {
	media := NewImageRepository()
	imageURL, err := media.RemoteUpload(imgUrl)
	if err != nil {
		log.Println("Error uploaidng image by Url", err)
		return helper.NewInternal()
	}
	err = db.DB.WithContext(ctx).Model(&Account{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"image": imageURL}).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %s not found", id)
			return helper.NewNotFound("id", id.String())
		}
		log.Println("Error querying account to updating image:", err)
		return helper.NewInternal()
	}
	return nil
}

// GetAllAccount gets a list of all account in db
func (db *DataSources) GetAllAccount(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	if err := db.DB.WithContext(ctx).Omit("password").Find(&accounts).Error; err != nil {
		log.Println("Error getting users accounts:", err)
		return nil, helper.NewInternal()
	}
	return accounts, nil
}

// DeleteAccount deletes an account based on the provided ID
func (db *DataSources) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	var account *Account
	if err := db.DB.WithContext(ctx).Where("id = ?", id).Delete(&account).Error; err != nil {
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				fmt.Printf("account with ID %s not found", account.ID)
				return helper.NewNotFound("id", account.ID.String())
			}
			log.Println("Error deleting account:", err)
			return helper.NewInternal()
		}
	}
	return nil
}
