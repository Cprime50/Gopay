package handler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	account_service "github.com/Cprime50/Gopay/service/account_service"
	token_service "github.com/Cprime50/Gopay/service/token_service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type Handler struct {
	router          *gin.Engine
	AccountService  *account_service.AccountService
	TokenService    *token_service.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

// Initilizes and retuens new handler
func (h *Handler) NewHandler(router *gin.Engine) (*Handler, error) {
	log.Print("setting up handler")

	// read in ACCOUNT_API_URL
	baseURL := os.Getenv("ACCOUNT_API_URL")

	// read in HANDLER_TIMEOUT
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	hTimeout, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	maxbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	//set timeout middleware
	timeoutDuration := time.Duration(time.Duration(hTimeout) * time.Second)

	// Initialize AccountService
	accountService := account_service.AccountService{} // Replace with your actual initialization

	// Initialize TokenService
	tokenService := token_service.TokenService{}

	handler := &Handler{
		router:          router,
		AccountService:  &accountService,
		TokenService:    &tokenService,
		BaseURL:         baseURL,
		TimeoutDuration: timeoutDuration,
		MaxBodyBytes:    maxbb,
	}

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*, http://*, *"},
		AllowMethods:     []string{"GET, POST, PUT, DELETE, OPTIONS"},
		AllowHeaders:     []string{"Origin, Accept, Authorization, Content-Type, X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == baseURL
		},
		MaxAge: 12 * time.Hour,
	}))

	// //routes and middleware setup
	// handler.SetupRoutes() //to set up routes

	return handler, nil

}
