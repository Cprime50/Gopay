package migrations

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Cprime50/Gopay/db"
	models "github.com/Cprime50/Gopay/models/account"
	"github.com/google/uuid"
)

func Migrate() {
	log.Printf("Migrations Started")
	startTime := time.Now()
	err := db.DB.AutoMigrate(&models.Role{}, &models.Account{})
	if err != nil {
		log.Fatal(err)
	}

	_err := seedData() // default data being added into the database upon migration
	if _err != nil {
		log.Fatal(_err)
	}
	log.Println("seeding data complete")
	elapsed := time.Since(startTime)
	log.Printf("Migrate completed in %s", elapsed)

}

// adding some default user data and roles into the db
func seedData() error {
	var roles = []models.Role{{ID: 1, Name: "admin", Description: "Administrator role"}, {ID: 2, Name: "user", Description: "user role"}}
	account, err := createAdminAccount()
	if err != nil {
		fmt.Println("error seeding admin data", err)
	}
	db.DB.Save(&roles)
	db.DB.Save(&account)

	return nil

}

//admin

func createAdminAccount() (*models.Account, error) {
	adminID, err := uuid.Parse(os.Getenv("ADMIN_ID"))
	if err != nil {
		return nil, err
	}
	accountNumber, err := strconv.ParseInt(os.Getenv("ADMIN_ACCOUNT_NUMBER"), 10, 64)
	if err != nil {
		return nil, err
	}

	adminBalance, err := strconv.ParseFloat(os.Getenv("ADMIN_ACCOUNT_BALANCE"), 64)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := models.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		return nil, err
	}

	account := &models.Account{
		ID:            adminID,
		Email:         os.Getenv("ADMIN_EMAIL"),
		Password:      hashedPassword,
		RoleID:        1,
		FirstName:     os.Getenv("ADMIN_FIRSTNAME"),
		LastName:      os.Getenv("ADMIN_LASTNAME"),
		AccountNumber: accountNumber,
		Balance:       adminBalance,
	}

	return account, nil
}
