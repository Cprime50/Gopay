package models

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (Db *DataSources) Migrate() {
	log.Printf("Migrations Started")
	startTime := time.Now()
	err := Db.DB.AutoMigrate(&Role{}, &Account{})
	if err != nil {
		log.Fatal(err)
	}

	_err := Db.seedData() // default data being added into the database upon migration
	if _err != nil {
		log.Fatal(_err)
	}
	log.Println("seeding data complete")
	elapsed := time.Since(startTime)
	log.Printf("Migrate completed in %s", elapsed)

}

// adding some default user data and roles into the db
func (Db *DataSources) seedData() error {

	var roles = []Role{{ID: 1, Name: "admin", Description: "Administrator role"}, {ID: 2, Name: "user", Description: "user role"}}
	account, err := createAdminAccount()
	if err != nil {
		fmt.Println("error seeding admin data", err)
	}
	Db.DB.Save(&roles)
	Db.DB.Save(&account)

	return nil

}

//admin

func createAdminAccount() (*Account, error) {
	adminID, err := uuid.Parse(os.Getenv("ADMIN_ID"))
	if err != nil {
		return nil, err
	}

	adminAccountNumber, err := strconv.ParseInt(os.Getenv("ADMIN_ACCOUNT_NUMBER"), 10, 64)
	if err != nil {
		return nil, err
	}

	adminBalance, err := strconv.ParseFloat(os.Getenv("ADMIN_ACCOUNT_BALANCE"), 64)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		return nil, err
	}

	account := &Account{
		ID:            adminID,
		Email:         os.Getenv("ADMIN_EMAIL"),
		Password:      hashedPassword,
		RoleID:        1,
		FirstName:     os.Getenv("ADMIN_FIRSTNAME"),
		LastName:      os.Getenv("ADMIN_LASTNAME"),
		AccountNumber: adminAccountNumber,
		Balance:       adminBalance,
	}

	return account, nil
}
