package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// // AdminHandler handles the incoming requests from the admin
func (h *Handler) AdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		data := gin.H{
			"Title": "Admin",
			"Body":  "This is the Admin page",
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

		//tmpl, err := template.ParseFiles("./web/templates/base.layout.gohtml", "./web/templates/admin.page.gohtml")
		//if err != nil {
		//	c.String(http.StatusInternalServerError, "Error parsing template: %v", err)
		//	return
		//}
		//
		//data := gin.H{
		//	"Title": "Admin Page",
		//	"Body":  "This is the admin page content.",
		//}
		//
		//// Ensure correct template execution
		//if err := tmpl.ExecuteTemplate(c.Writer, "base", data); err != nil {
		//	c.String(http.StatusInternalServerError, "Error executing template: %v", err)
		//}

	}
}

func (h *Handler) AdminHandlerTest() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.HTML(200, "admin.html", nil)

	}
}
