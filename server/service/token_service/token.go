package service

import (
	"crypto/rsa"
	"fmt"
	"log"
	"time"

	models "github.com/Cprime50/Gopay/models"
	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// idTokenCustomClaims holds structure of jwt claims of idToken
type idTokenCustomClaims struct {
	Account *models.Account `json:"account"`
	jwt.RegisteredClaims
}

// generateJWT generates an IDToken which is a jwt with myCustomClaims
// Could call this GenerateJWT, but the signature makes this fairly clear
func generateJWT(account *models.Account, key *rsa.PrivateKey, exp int64) (string, error) {
	//tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	issuedAt := jwt.NewNumericDate(time.Now().UTC())
	expiresAt := jwt.NewNumericDate(issuedAt.Add(time.Hour * time.Duration(exp)))

	claims := idTokenCustomClaims{
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return signedString, nil
}

// refreshTokenData holds the actual signed jwt string along with the ID
// We return the id so it can be used without re-parsing the JWT from signed string
type refreshTokenData struct {
	SignedString string
	ID           uuid.UUID
	ExpiresIn    time.Duration
}

// refreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type refreshTokenCustomClaims struct {
	AccountID uuid.UUID `json:"account_id"`
	RoleID    uint      `json:"role_id"`
	jwt.RegisteredClaims
}

// generateRefreshToken creates a refresh token
// The refresh token stores only the account's ID, role and a string
func generateRefreshToken(account *models.Account, key string, exp int64) (*refreshTokenData, error) {
	issuedAt := jwt.NewNumericDate(time.Now().UTC())
	expiresAt := jwt.NewNumericDate(issuedAt.Add(time.Hour * time.Duration(exp)))
	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
		AccountID: account.ID,
		RoleID:    account.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &refreshTokenData{
		SignedString: signedString,
		ID:           tokenID,
		ExpiresIn:    expiresAt.Sub(time.Now()),
	}, nil
}

// validateIDToken returns the token's claims if the token is valid
func validateJWT(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}
	claims, ok := token.Claims.(*idTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}
	return claims, nil
}

// Validate admin
func validateAdminJWT(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}
	//check role for admin
	claims, ok := token.Claims.(*idTokenCustomClaims)
	accountRole := claims.Account.RoleID
	if !token.Valid && accountRole != 1 {
		return nil, fmt.Errorf("invalid admin token")
	}

	if !ok {
		return nil, fmt.Errorf("token valid but Invalid admin couldn't parse claims")
	}
	return claims, nil
}

// validateRefreshToken uses the secret key to validate a refresh token
func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}
	//check role for admin and validate
	claims, ok := token.Claims.(*refreshTokenCustomClaims)
	accountRole := claims.RoleID
	if !token.Valid && accountRole != 1 {
		return nil, fmt.Errorf("invalid admin refresh token")
	}

	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
	}
	return claims, nil
}

// Validate admin refresh token
func validateAdminRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, fmt.Errorf("error parsing claims to jwt %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("Refresh token is invalid")
	}
	claims, ok := token.Claims.(*refreshTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
	}
	return claims, nil
}
