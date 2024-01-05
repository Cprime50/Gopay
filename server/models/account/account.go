package models

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/Cprime50/Gopay/db"
	"github.com/Cprime50/Gopay/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID `gorm:"primary_key"`
	Email         string    `gorm:"uniqueIndex;not null;type:varchar(250)" json:"email"`
	AccountNumber string    `gorm:"uniqueIndex;column:account_number;not null"`
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
func CreateAccount(ctx context.Context, account *Account) error {
	// Check if an account with the given email already exists
	err := db.DB.WithContext(ctx).Where("email = ?", account.Email).First(&account).Error
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
	accountNumber := GenerateAccountNumber(ctx)

	// initialize account number
	newAccount := &Account{
		AccountNumber: accountNumber,
		Balance:       500,
		RoleID:        2,
		IsActive:      false,
	}
	if err := db.DB.WithContext(ctx).Create(newAccount).Error; err != nil {
		log.Printf("Error creating account: %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

// GenerateAccountNumber generates a unique 10-digit account number.
func GenerateAccountNumber(ctx context.Context) string {
	//Generate a UUID
	uuid := uuid.New()

	//Hash the UUID using MD5
	hasher := md5.New()
	hasher.Write([]byte(uuid.String()))
	hash := hex.EncodeToString(hasher.Sum(nil))

	//Use first 11 digit of hash as account number
	accountNumber := hash[:11]

	// Check if account number is unique
	_, err := GetAccountByAccountNum(ctx, accountNumber)
	if err == nil {
	}
	// If not unique, regenerate and check again
	return accountNumber

}

// GetAccountByAccountNum gets a user's account from the database based on the account number
func GetAccountByAccountNum(ctx context.Context, accountNumber string) (*Account, error) {

	var account Account
	err := db.DB.WithContext(ctx).Where("account_number = ?", accountNumber).First(&account).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			fmt.Printf("account with ID %s not found", accountNumber)
			return nil, helper.NewNotFound("account_number", fmt.Sprintf("%s", accountNumber))
		}
		log.Println("Error querying account:", err)
		return nil, helper.NewInternal()
	}
	return &account, nil
}

// GetAccountByEmail gets a user's account from the database based on the email
func GetAccountByEmail(ctx context.Context, email string) (*Account, error) {

	var account Account
	err := db.DB.WithContext(ctx).Where("email = ?", email).First(&account).Error
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
func GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
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
func UpdateAccount(ctx context.Context, account *Account) error {
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

func ChangeUserStatus(ctx context.Context, account *Account) error {
	// Check if the accounts exists based on the provided accountID
	_, err := GetAccountByID(ctx, account.ID)
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
func ResetPassword(ctx context.Context, account *Account) error {
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
func UpdateImgByFile(ctx context.Context, id uuid.UUID, imgFile *File) error {
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
func UpdateImgByUrl(ctx context.Context, id uuid.UUID, imgUrl *Url) error {
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
func GetAllAccount(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	if err := db.DB.WithContext(ctx).Omit("password").Find(&accounts).Error; err != nil {
		log.Println("Error getting users accounts:", err)
		return nil, helper.NewInternal()
	}
	return accounts, nil
}

// DeleteAccount deletes an account based on the provided ID
func DeleteAccount(ctx context.Context, id uuid.UUID) error {
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
