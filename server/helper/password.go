package helper

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// confirm user password checks if password matches hashed one in db
func ComparePassword(hashedPassword string, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
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
