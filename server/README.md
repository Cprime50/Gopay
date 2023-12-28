
## To-Do

### JWT Admin User Authentication

1. **Implement JWT Admin Authentication Middleware:**
    - [ ] Create a middleware function for JWT admin authentication.
    - [ ] Use the middleware to validate JWT tokens and enforce admin privileges.
    - [ ] Handle unauthorized access appropriately.                                                                  

    ```go
    adminRoutes := router.Group("/api/admin")
    adminRoutes.Use(util.JWTAuthAdminMiddleware())
    {
        // Your admin-specific routes
    }
    ```

2. **JWT Admin Authentication Handlers:**
    - [ ] Define route handlers for admin-specific actions.
    - [ ] Use the `JWTAuthAdminMiddleware` to protect admin routes.

    ```go
    // Example admin-specific handler
    adminRoutes.GET("/dashboard", func(c *gin.Context) {
        // Your admin-specific logic here
        c.JSON(http.StatusOK, gin.H{"message": "Admin Dashboard"})
    })
    ```

3. **Test JWT Admin Authentication:**
    - [ ] Write unit tests to ensure the correct functioning of JWT admin authentication.
    - [ ] Cover both positive and negative scenarios, including unauthorized access attempts.

    ```go
    // Example unit test for admin authentication
    func TestAdminAuthentication(t *testing.T) {
        // Your test logic here
    }
    ```

4. **Update README Documentation:**
    - [ ] Update the README with usage instructions for JWT admin authentication.
    - [ ] Provide examples of how to apply the middleware and use admin-specific routes.

    ```markdown
    ## Usage

    ### J
