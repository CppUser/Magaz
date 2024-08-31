package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// EmployeeHandler handles the incoming requests from the employee
func (h *Handler) EmployeeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["employee.page.gohtml"]

		data := gin.H{
			"Title": "Employee",
			"Body":  "This is the Employee page",
		}

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", data)
			if err != nil {
				h.Logger.Error("Error executing template", zap.Error(err))
				c.String(http.StatusInternalServerError, "Error executing template: %v", err)
			}
		} else {
			h.Logger.Error("Template not found", zap.String("template", "admin.page.gohtml"))
			c.String(http.StatusInternalServerError, "Template not found")
		}
	}
}
