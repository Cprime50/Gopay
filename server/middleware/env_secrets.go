package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// remember to fix this whenu wake up
func GenerateRSAKeys() error {
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	pubKeyFile := os.Getenv("PUB_KEY_FILE")

	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	// Encode private key to PEM format
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privKeyBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKeyBytes}

	//Create private key
	privKeyFileHandle, err := os.Create(privKeyFile)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privKeyFileHandle.Close()
	if err := pem.Encode(privKeyFileHandle, privKeyBlock); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	// Generate RSA public key
	pubKey := &privateKey.PublicKey

	// Encode public key to PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubKeyBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}

	//Create pubkey
	pubKeyFileHandle, err := os.Create(pubKeyFile)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer pubKeyFileHandle.Close()
	if err := pem.Encode(pubKeyFileHandle, pubKeyBlock); err != nil {
		return fmt.Errorf("failed to write public key to file: %w", err)
	}

	return nil
}

func privKey() (*rsa.PrivateKey, error) {
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	if privKeyFile == "" {
		log.Fatal("privKeyFile not provided in env file")
	}
	priv, err := os.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}
	return privKey, nil
}

func pubKey() (*rsa.PublicKey, error) {
	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	if pubKeyFile == "" {
		log.Println("pubKeyFile  not provided in env file")
	}
	pub, err := os.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}
	return pubKey, nil
}

func tokenExp() (int64, error) {
	TokenExp := os.Getenv("TOKEN_EXP")
	if TokenExp == "" {
		log.Println("TokenExp not provided in env file, will use default 1800")
		TokenExp = "1800"
	}
	tokenExp, err := strconv.ParseInt(TokenExp, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse TOKEN_EXP as int: %w", err)
	}

	return tokenExp, nil
}

func refreshExp() (int64, error) {
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")
	if refreshTokenExp == "" {
		log.Println("RefreshTokenExp not provided in env file, will use default 259200")
		refreshTokenExp = "259200"
	}
	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	return refreshExp, nil
}

func refreshSecret() string {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	// Provide default values if necessary
	if refreshSecret == "" {
		log.Print("refreshSecret  not provided in env file, will generate random secret")
		bytes := make([]byte, 32)
		rand.Read(bytes)
		refreshSecret = base64.StdEncoding.EncodeToString(bytes)
	}
	return refreshSecret
}
