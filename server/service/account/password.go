package service

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// hash password before we save it
func hashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypting password failed", err)
		return "", err
	}
	hashedPw := string(passwordHash)

	return hashedPw, nil
}

// confirm user password whether it matches
func comparePassword(password string, comfPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(comfPassword))
	// we can get two types of error here
	if err != nil {
		switch {
		// error when  password doesnt match
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil

		//and error when something unexpected happens
		default:
			return false, err
		}
	}
	return true, nil
}
