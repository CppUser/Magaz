package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// EmployeeHandler handles the incoming requests from the employee
func (h *Handler) EmployeeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		data := gin.H{
			"Title": "Employee",
			"Body":  "This is the Employee page",
		}
		// Respond with status OK
		//RenderTemplate(c, "employee.page.gohtml", data)

		c.HTML(http.StatusOK, "employee.page.gohtml", data)
	}
}
