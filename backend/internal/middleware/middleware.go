package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

func LogDetailedRequestsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path
		headers := c.Request.Header
		fmt.Printf("Received %s request for %s with headers: %v\n", method, path, headers)
		c.Next()
	}
}

// TODO: Pass Handler instead ?
// SessionMiddleware returns a Gin middleware function that manages sessions
func SessionMiddleware(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create or retrieve the session
		session, err := store.Get(c.Request, "session-name")

		// Set session options for persistent cookies
		expiration := time.Now().Add(24 * time.Hour)
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   int(expiration.Sub(time.Now()).Seconds()),
			HttpOnly: true,
			Secure:   false, //TODO: set to true in production
			SameSite: http.SameSiteStrictMode,
		}
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		// Save the session in the Gin context
		c.Set("session", session)

		// Continue processing the request
		c.Next()
	}
}

// AuthRequired Middleware to check if the user is authenticated
func AuthRequired(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the session
		session, _ := store.Get(c.Request, "session-name")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// If not authenticated, redirect to login page
			c.Redirect(http.StatusFound, "/api/login")
			c.Abort()
			return
		}

		// If authenticated, continue to the next handler
		c.Next()
	}
}
