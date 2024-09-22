package handler

import (
	"Magaz/backend/internal/repository"
	"Magaz/backend/internal/storage/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// // AdminHandler handles the incoming requests from the admin
func (h *Handler) AdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

func (h *Handler) AdminStatisticsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

// // GetProductsAdminHandler handles the incoming requests from the admin
func (h *Handler) GetProductsAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		citiesWithProducts, err := repository.FetchCityProcdducts(h.DB)
		if err != nil {
			h.Logger.Error("Failed to fetch city products", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to load city products")
			return
		}

		// Prepare the data to be passed to the template
		data := gin.H{
			"CitiesWithProducts": citiesWithProducts,
		}

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "product_section", data)
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

func (h *Handler) PostAdminAddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Parse form data
		var productData struct {
			City               string    `json:"city"`
			Product            string    `json:"product"`
			Quantities         []float32 `json:"quantities"`
			QuantityPricePairs []struct {
				Quantity float32 `json:"quantity"`
				Price    float32 `json:"price"`
			} `json:"quantityPricePairs"`
		}

		if err := c.ShouldBindJSON(&productData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create new product and city instances based on the parsed data
		newCity := models.City{Name: productData.City}
		newProduct := models.Product{Name: productData.Product}

		// Assume city already exists for simplicity; otherwise, find or create
		if err := h.DB.Where("name = ?", newCity.Name).FirstOrCreate(&newCity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find or create city"})
			return
		}

		// Create the new product in the city
		cityProduct := models.CityProduct{
			CityID:        newCity.ID,
			Product:       newProduct,
			TotalQuantity: 0, // Placeholder; adjust as needed
		}

		if err := h.DB.Create(&cityProduct).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
			return
		}

		// Add product prices to the database
		for _, qp := range productData.QuantityPricePairs {
			productPrice := models.QtnPrice{
				CityProductID: cityProduct.ID,
				Quantity:      qp.Quantity,
				Price:         qp.Price,
			}

			if err := h.DB.Create(&productPrice).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product prices"})
				return
			}
		}

		// Respond to the client
		c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
	}
}

func (h *Handler) AdminGetProductAddr() gin.HandlerFunc {
	return func(c *gin.Context) {
		//var items []models.Address

		// Get query parameters
		cityIDStr := c.Query("cityID")
		productIDStr := c.Query("productID")
		quantityIDStr := c.Query("quantityID")

		// Convert query parameters to uint
		cityID, err := strconv.ParseUint(cityIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
			return
		}
		productID, err := strconv.ParseUint(productIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		quantityID, err := strconv.ParseUint(quantityIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity ID"})
			return
		}

		var addresses []models.Address
		// Preload the Employee who added the address and fetch addresses for the given city, product, and quantity
		if err := h.DB.Preload("AddedBy").Where("city_id = ? AND product_id = ? AND qtn_price_id = ?", cityID, productID, quantityID).Find(&addresses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
			return
		}

		// Fetch the actual quantity from QtnPrice table
		var qtnPrice models.QtnPrice
		if err := h.DB.First(&qtnPrice, quantityID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quantity"})
			return
		}

		// Prepare the data to be sent in the response
		var items []map[string]interface{}
		for _, addr := range addresses {
			if !addr.Released {
				items = append(items, map[string]interface{}{
					"ID":          addr.ID,
					"Description": addr.Description,
					"Quantity":    qtnPrice.Quantity, //TODO: need to fetch actual quantity that passed by quantityId
					"AddedAt":     addr.AddedAt,
					"AddedBy":     addr.AddedBy.Username, // Assuming AddedBy is fetched with the Address model
					"Image":       addr.Image,
				})
			}

		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}
func (h *Handler) PostAdminAddProductAddr() gin.HandlerFunc {
	return func(c *gin.Context) {

		var newItem models.Address

		// Bind the form data from the request to the newItem struct
		if err := c.ShouldBind(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle the uploaded image file
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload error: " + err.Error()})
			return
		}

		// Save the file on the server
		imagePath := "./backend/storage/images" + file.Filename
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		newItem.Image = imagePath

		// Save the new item to the database
		if err := h.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item added successfully!"})

	}
}

func (h *Handler) AdminOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

func (h *Handler) AdminDisputesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

func (h *Handler) AdminChatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

func (h *Handler) AdminSettingsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "base", nil)
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
