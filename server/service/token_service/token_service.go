package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/Cprime50/Gopay/helper"
	models "github.com/Cprime50/Gopay/models/account"
	"github.com/google/uuid"
)

// tokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type TokenService struct {
	TokenRepository models.TokenRepository
	PrivKey         *rsa.PrivateKey
	PubKey          *rsa.PublicKey
	RefreshSecret   string
	TokenExp        int64
	RefreshExp      int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository models.TokenRepository
	PrivKey         *rsa.PrivateKey
	PubKey          *rsa.PublicKey
	RefreshSecret   string
	TokenExp        int64
	RefreshExp      int64
}

// NewPairFromUser generates a new set of ID and refresh tokens for the specified user account.
// If a previous refresh token is provided, it is removed from the token repository.
// The newly generated refresh token is stored in the repository for future validation.
func (s *TokenService) NewPairFromUser(ctx context.Context, account *models.Account, prevTokenID string) (*models.TokenPair, error) {
	// Remove the previous refresh token from the repository if provided
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, account.ID.String(), prevTokenID); err != nil {
			log.Printf("Failed to delete previous refreshToken for UID: %v, TokenID: %v\n", account.ID.String(), prevTokenID)
			return nil, err
		}
	}

	// Generate a new ID token
	idToken, err := generateJWT(account, s.PrivKey, s.TokenExp)
	if err != nil {
		log.Printf("Error generating ID token for Account ID: %v. Error: %v\n", account.ID, err.Error())
		return nil, helper.NewInternal()
	}

	// Generate a new refresh token
	refreshToken, err := generateRefreshToken(account, s.RefreshSecret, s.RefreshExp)
	if err != nil {
		log.Printf("Error generating refresh token for Account ID: %v. Error: %v\n", account.ID, err.Error())
		return nil, helper.NewInternal()
	}

	// Store the newly generated refresh token in the repository
	if err := s.TokenRepository.SetRefreshToken(ctx, account.ID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for UID: %v. Error: %v\n", account.ID, err.Error())
		return nil, helper.NewInternal()
	}

	// Return the newly generated tokens as a TokenPair
	return &models.TokenPair{
		IDToken:      models.IDToken{SignedString: idToken},
		RefreshToken: models.RefreshToken{SignedString: refreshToken.SignedString, ID: refreshToken.ID, AccountID: account.ID},
	}, nil
}

// Signout revokes all valid tokens for a user by reaching out to the repository layer.
func (s *TokenService) Signout(ctx context.Context, id uuid.UUID) error {
	// Delete all valid refresh tokens associated with the user ID from the repository
	return s.TokenRepository.DeleteUserRefreshTokens(ctx, id.String())
}

// ValidateJWT validates the provided ID token JWT string using the public RSA key.
func (s *TokenService) ValidateJWT(tokenString string) (*models.Account, error) {
	// Validate and parse the ID token using the provided public RSA key
	claims, err := validateJWT(tokenString, s.PubKey)
	if err != nil {
		log.Printf("Unable to validate or parse ID token - Error: %v\n", err)
		return nil, helper.NewAuthorization("Unable to verify user from ID token")
	}
	return claims.Account, nil
}

// JWTAuthAdmin validates the id token jwt string
func (s *TokenService) ValidateAdminJWT(tokenString string) (*models.Account, error) {
	claims, err := validateAdminJWT(tokenString, s.PubKey) // uses public RSA key
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, helper.NewAuthorization("Unable to verify admin from idToken")
	}
	return claims.Account, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (s *TokenService) ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	// validate actual JWT with string a secret
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, helper.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.ID)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.ID, err)
		return nil, helper.NewAuthorization("Unable to verify user from refresh token")
	}

	return &models.RefreshToken{
		SignedString: tokenString,
		ID:           tokenUUID,
		AccountID:    claims.AccountID,
	}, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (s *TokenService) ValidateAdminRefreshToken(tokenString string) (*models.RefreshToken, error) {
	// validate actual JWT with string a secret
	claims, err := validateAdminRefreshToken(tokenString, s.RefreshSecret)

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, helper.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.ID)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.ID, err)
		return nil, helper.NewAuthorization("Unable to verify user from refresh token")
	}

	return &models.RefreshToken{
		SignedString: tokenString,
		ID:           tokenUUID,
		AccountID:    claims.AccountID,
	}, nil
}
