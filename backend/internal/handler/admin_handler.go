package handler

import (
	"Magaz/backend/internal/repository"
	"Magaz/backend/internal/storage/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ProductItem struct {
	Quantity  float32
	Price     float32
	Available int // This can be calculated if needed
}

type ProductView struct {
	City     string
	Products []ProductDetail
}

type ProductDetail struct {
	Name  string
	Total int
	Items []ProductItem
}

// // AdminHandler handles the incoming requests from the admin
func (h *Handler) AdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpl := h.TmplCache["admin.page.gohtml"]

		citiesWithProducts, err := repository.FetchCityProducts(h.DB)
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
			productPrice := models.ProductPrice{
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

func (h *Handler) AdminGetProductItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.Item

		// Fetch items from the database
		if err := h.DB.Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}
func (h *Handler) PostAdminAddProductItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		var newItem models.Item

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
