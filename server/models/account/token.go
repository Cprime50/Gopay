package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Cprime50/Gopay/db"
	"github.com/Cprime50/Gopay/helper"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID           uuid.UUID `json:"-"`
	AccountID    uuid.UUID `json:"-"`
	SignedString string    `json:"refreshToken"`
}

type Token struct {
	SignedString string `json:"token"`
}

// TokenPair used for returning pairs of id and refresh tokens
type TokenPair struct {
	Token
	RefreshToken
}

// SetRefreshToken stores a refresh token with an expiry time
func SetRefreshToken(ctx context.Context, accountID string, tokenID string, expiresIn time.Duration) error {
	// We'll store accountID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", accountID, tokenID)
	if err := db.RedisClient.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for accountID/tokenID: %s/%s: %v\n", accountID, tokenID, err)
		return helper.NewInternal()
	}
	return nil
}

// Deletes a specific refresh token from Redis
func DeleteRefreshToken(ctx context.Context, accountID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", accountID, tokenID)
	result := db.RedisClient.Del(ctx, key)
	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for accountID/tokenID: %s/%s: %v\n", accountID, tokenID, err)
		return helper.NewInternal()
	}

	// Val returns count of deleted keys.
	// If no key was deleted, the refresh token is invalid
	if result.Val() < 1 {
		log.Printf("Refresh token to redis for accountID/tokenID: %s/%s does not exist\n", accountID, tokenID)
		return helper.NewAuthorization("Invalid refresh token")
	}

	return nil
}

// DeleteUserRefreshTokens looks for all tokens beginning with
// accountID and scans to delete them in a non-blocking fashion
func DeleteUserRefreshTokens(ctx context.Context, accountID string) error {
	pattern := fmt.Sprintf("%s*", accountID)

	iter := db.RedisClient.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := db.RedisClient.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}

	// check last value
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iter.Val())
	}

	if failCount > 0 {
		return helper.NewInternal()
	}

	return nil
}
