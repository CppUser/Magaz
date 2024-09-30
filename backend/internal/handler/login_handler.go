package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) GETLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		data := gin.H{
			"Title": "Login",
			"Error": "Invalid username or password",
		}
		tmpl, _ := h.TmplCache.GetTemplate("login.page.gohtml")
		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", data)
			if err != nil {
				h.Logger.Error("Error executing template", zap.Error(err))
				c.String(http.StatusInternalServerError, "Error executing template: %v", err)
			}
		} else {
			h.Logger.Error("Template not found", zap.String("template", "login.page.gohtml"))
			c.String(http.StatusInternalServerError, "Template not found")
		}
	}
}

func (h *Handler) POSTLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Process the login form
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Simple authentication logic (replace with actual logic)
		if username == "admin" && password == "password" {
			// Set the session values
			session, _ := h.Session.Get(c.Request, "session-name")
			session.Values["authenticated"] = true
			session.Save(c.Request, c.Writer)

			//TODO: Redirect based on role
			// Redirect to the admin page after successful login
			c.Redirect(http.StatusFound, "/api/admin")
			return
		}

		// If login fails, render the login page with an error message
		data := gin.H{
			"Title": "Login",
			"Error": "Invalid username or password",
		}
		tmpl, _ := h.TmplCache.GetTemplate("login.page.gohtml")
		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", data)
			if err != nil {
				h.Logger.Error("Error executing template", zap.Error(err))
				c.String(http.StatusInternalServerError, "Error executing template: %v", err)
			}
		} else {
			h.Logger.Error("Template not found", zap.String("template", "login.page.gohtml"))
			c.String(http.StatusInternalServerError, "Template not found")
		}
	}
}
