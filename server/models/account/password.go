package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// hash password before we save it
func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypting password failed", err)
		return "", err
	}
	hashedPw := string(passwordHash)

	return hashedPw, nil
}
