package handler

import (
	"Magaz/backend/internal/repository"
	"Magaz/backend/internal/storage/models"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
			err := tmpl.ExecuteTemplate(c.Writer, "products", data)
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

func (h *Handler) GetAddProductFormHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		//TODO: Pass data with cities and available products
		//citiesWithProducts, err := repository.FetchCityProcdducts(h.DB)
		//if err != nil {
		//	h.Logger.Error("Failed to fetch city products", zap.Error(err))
		//	c.String(http.StatusInternalServerError, "Failed to load city products")
		//	return
		//}
		//
		//// Prepare the data to be passed to the template
		//data := gin.H{
		//	"CitiesWithProducts": citiesWithProducts,
		//}

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "add_product_section", nil)
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
			City       string    `form:"city"`
			Product    string    `form:"product"`
			Quantities []float32 `form:"quantities[]"`
			Prices     []float32 `form:"prices[]"`
		}

		if err := c.ShouldBind(&productData); err != nil {
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

		// Add product prices to the database using the parsed Quantities and Prices
		for i := 0; i < len(productData.Quantities); i++ {
			productPrice := models.QtnPrice{
				CityProductID: cityProduct.ID,
				Quantity:      productData.Quantities[i],
				Price:         productData.Prices[i],
			}

			if err := h.DB.Create(&productPrice).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product prices"})
				return
			}
		}

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
func (h *Handler) GetProductAddrForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmpl := h.TmplCache["admin.page.gohtml"]

		cityIDStr := c.Query("cityID")
		productIDStr := c.Query("productID")
		quantityIDStr := c.Query("quantityID")

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
		if err := h.DB.Preload("AddedBy").Where("city_id = ? AND product_id = ? AND qtn_price_id = ?", cityID, productID, quantityID).Find(&addresses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
			return
		}

		var qtnPrice models.QtnPrice
		if err := h.DB.First(&qtnPrice, quantityID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quantity"})
			return
		}

		var items []map[string]interface{}
		for _, addr := range addresses {
			if !addr.Released {
				items = append(items, map[string]interface{}{
					"ID":          addr.ID,
					"Description": addr.Description,
					"Quantity":    qtnPrice.Quantity,
					"AddedAt":     addr.AddedAt,
					"AddedBy":     addr.AddedBy.Username,
					"Image":       addr.Image,
				})
			}
		}

		data := gin.H{
			"Items":      items,
			"CityID":     cityID,
			"ProductID":  productID,
			"QuantityID": quantityID,
		}

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "prd_address_form", data)
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
func (h *Handler) GetAddAddrForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		tmpl := h.TmplCache["admin.page.gohtml"]

		cityIDStr := c.Query("cityID")
		productIDStr := c.Query("productID")
		quantityIDStr := c.Query("quantityID")

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

		data := gin.H{
			"CityID":     cityID,
			"ProductID":  productID,
			"QuantityID": quantityID,
		}

		//TODO: Pass data with cities and available products
		//citiesWithProducts, err := repository.FetchCityProcdducts(h.DB)
		//if err != nil {
		//	h.Logger.Error("Failed to fetch city products", zap.Error(err))
		//	c.String(http.StatusInternalServerError, "Failed to load city products")
		//	return
		//}
		//
		//// Prepare the data to be passed to the template
		//data := gin.H{
		//	"CitiesWithProducts": citiesWithProducts,
		//}

		if tmpl != nil {
			err := tmpl.ExecuteTemplate(c.Writer, "add_address_form", data)
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
func (h *Handler) PostAdminAddProductAddr() gin.HandlerFunc {
	return func(c *gin.Context) {

		cityIDStr := c.Query("cityID")
		productIDStr := c.Query("productID")
		quantityIDStr := c.Query("quantityID")

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

		err = c.Request.ParseMultipartForm(10 << 20) // Limit your file size to 10 MB
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
			return
		}

		// Extract descriptions and files
		descriptions := c.Request.MultipartForm.Value["description[]"]
		files := c.Request.MultipartForm.File["image[]"]

		if len(descriptions) != len(files) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatched descriptions and image files count"})
			return
		}

		// Iterate over descriptions and images
		for i, description := range descriptions {
			// Open the file for each image
			fileHeader := files[i]
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open image file"})
				return
			}
			defer file.Close()

			// Save the file to the server
			cwd, _ := os.Getwd()
			imagePath := filepath.Join(cwd, "backend", "storage", "images")
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				os.MkdirAll(imagePath, os.ModePerm) // Create directory if it doesn't exist
			}

			// Generate a unique file name
			uuidFilename := uuid.New().String()          // Generate a UUID for unique file name
			fileExt := filepath.Ext(fileHeader.Filename) // Get the original file extension
			fileName := uuidFilename                     // Use UUID as filename (without extension)
			filePath := filepath.Join(imagePath, fileName+fileExt)

			out, err := os.Create(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image file"})
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image to disk"})
				return
			}

			// Store the new address in the database
			newAddress := models.Address{
				CityID:      uint(cityID),
				QtnPriceID:  uint(quantityID),
				ProductID:   uint(productID),
				Description: description,
				Image:       fileName,
				AddedAt:     time.Now(),
				EmployeeID:  1, // Replace with actual employee ID
			}

			if err := h.DB.Create(&newAddress).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save address to database"})
				return
			}
		}

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "All addresses added successfully"})

		//var newItem models.Address
		//
		//// Bind the form data from the request to the newItem struct
		//if err := c.ShouldBind(&newItem); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}
		//
		//// Handle the uploaded image file
		//file, err := c.FormFile("image")
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload error: " + err.Error()})
		//	return
		//}
		//
		//// Save the file on the server
		//imagePath := "./backend/storage/images" + file.Filename
		//if err := c.SaveUploadedFile(file, imagePath); err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		//	return
		//}
		//
		//newItem.Image = imagePath
		//
		//// Save the new item to the database
		//if err := h.DB.Create(&newItem).Error; err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		//	return
		//}
		//
		//c.JSON(http.StatusOK, gin.H{"message": "Item added successfully!"})

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
