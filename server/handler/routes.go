package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupRoutes() {
	// starting route
	h.router.GET("/", func(c *gin.Context) {
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "Welcome Gopay Server")
	})

	// Create a group for base routes
	baseRoutes := h.router.Group("/api")
	{
		baseRoutes.POST("/register", h.Signup)
		baseRoutes.POST("/login", h.Signin)
	}

	// Basic Authenticated routes
	authRoutes := h.router.Group("/api")
	authRoutes.Use(TimeoutMiddleware(h.TimeoutDuration), h.TokenService.AuthUser())
	{

	}

	// Admin routes
	adminRoutes := h.router.Group("/api/admin")
	adminRoutes.Use(TimeoutMiddleware(h.TimeoutDuration), h.TokenService.AuthAdmin())
	{

	}
}

// TimeoutMiddleware converts TimeoutMiddleware to be compatible with gin.HandlerFunc
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
