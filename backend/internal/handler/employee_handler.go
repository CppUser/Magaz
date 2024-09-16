package handler

import (
	"Magaz/backend/internal/repository"
	crud "Magaz/backend/internal/storage"
	crud2 "Magaz/backend/internal/storage/crud"
	"Magaz/backend/internal/storage/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) HEmployeeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptHeader := c.GetHeader("Accept")
		if acceptHeader == "text/event-stream" {
			clientChan := make(chan interface{})
			h.SSES.Register <- clientChan
			defer func() {
				h.SSES.Unregister <- clientChan
			}()

			c.Writer.Header().Set("Content-Type", "text/event-stream")
			c.Writer.Header().Set("Cache-Control", "no-cache")
			c.Writer.Header().Set("Connection", "keep-alive")

			for {
				select {
				case message := <-clientChan:
					// Type assert to map[string]interface{} before accessing fields
					orderDetails, ok := message.(map[string]interface{})
					if !ok {
						h.Logger.Error("Failed to type assert message")
						continue
					}

					// Extract details from the order
					orderID := orderDetails["ID"].(uint)
					cityName := orderDetails["CityName"].(string)
					productName := orderDetails["ProductName"].(string)
					quantity := orderDetails["Quantity"].(float32)
					due := orderDetails["Due"].(uint)
					username := orderDetails["Username"].(string)
					createdAt := orderDetails["CreatedAt"].(string)

					// Generate HTML fragment
					htmlFragment := fmt.Sprintf(`
                        <tr id="order-%d" hx-swap-oob="true">
                            <td>%d</td>
                            <td>%s</td>
                            <td>%s</td>
                            <td>%.2f</td>
                            <td>%d</td>
                            <td>%s</td>
                            <td>%s</td>
                        </tr>`,
						orderID, orderID, cityName, productName, quantity, due, username, createdAt)

					// Send the data to the client
					fmt.Fprintf(c.Writer, "data: %s\n\n", htmlFragment)
					c.Writer.Flush()
				case <-c.Writer.CloseNotify():
					return
				}
			}
		}

		// Regular page request rendering orders
		tmpl := h.TmplCache["employee.page2.gohtml"]

		// Fetch orders
		orders, err := crud.GetAll[models.Order](h.DB)
		if err != nil {
			h.Logger.Error("Failed to fetch orders from database", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to load orders")
			return
		}

		var orderViews []repository.OrderView
		for _, ords := range orders {
			// Map order models to views
			orderViews = append(orderViews, createOrderView(h.DB, ords))
		}

		// Render the template with order views
		data := gin.H{
			"Orders": orderViews,
		}

		if tmpl != nil {
			if err := tmpl.ExecuteTemplate(c.Writer, "base", data); err != nil {
				h.Logger.Error("Error executing template", zap.Error(err))
				c.String(http.StatusInternalServerError, "Error executing template: %v", err)
				return
			}
		} else {
			h.Logger.Error("Template not found", zap.String("template", "employee.page.gohtml"))
			c.String(http.StatusInternalServerError, "Template not found")
		}
	}
}

// EmployeeHandler handles the incoming requests from the employee
func (h *Handler) EmployeeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if this is an SSE connection or a regular request
		acceptHeader := c.GetHeader("Accept")
		if acceptHeader == "text/event-stream" {
			// This is an SSE request
			clientChan := make(chan interface{})
			h.SSES.Register <- clientChan
			defer func() {
				h.SSES.Unregister <- clientChan
			}()

			// Set the headers for SSE
			c.Writer.Header().Set("Content-Type", "text/event-stream")
			c.Writer.Header().Set("Cache-Control", "no-cache")
			c.Writer.Header().Set("Connection", "keep-alive")
			c.Writer.Header().Set("Transfer-Encoding", "chunked")

			// Keep the connection open and send updates as they arrive
			for {
				select {
				case message := <-clientChan:
					_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", dataToJSON(message))
					c.Writer.Flush()
				case <-c.Writer.CloseNotify():
					return
				}
			}
		}

		// This is a regular HTTP request, so render the template with initial data
		tmpl := h.TmplCache["employee.page.gohtml"]

		// Fetch the initial orders
		orders, err := crud.GetAll[models.Order](h.DB)
		if err != nil {
			h.Logger.Error("Failed to fetch orders from database", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to load orders")
			return
		}

		var orderViews []repository.OrderView
		for _, ords := range orders {
			if !ords.Released {
				product, _ := crud.Get[models.Product, uint](h.DB, ords.ProductID)
				city, _ := crud.Get[models.City, uint](h.DB, ords.CityID)
				user, _ := crud.Get[models.User, int64](h.DB, ords.UserID)
				cardPayment, _ := crud.Get[models.Card, uint](h.DB, ords.PaymentMethodID)
				//cryptoPayment,_ := crud.Get[models.Crypto,uint](h.DB,ords.PaymentMethodID)

				var payment repository.PaymentView
				if ords.PaymentMethodType == "card" {
					payment = repository.PaymentView{
						PaymentCategory: "Перевод на карту",
						CardPayment: repository.CardView{
							BankName:   cardPayment.BankName,
							BankUrl:    cardPayment.BankURL,
							CardNumber: cardPayment.CardNumber,
							FirstName:  cardPayment.FirstName,
							LastName:   cardPayment.LastName,
							UserName:   cardPayment.UserID,
							Password:   cardPayment.Password,
						},
					}
				} else if ords.PaymentMethodType == "crypto" { //TODO: Fill crypto payment
					payment = repository.PaymentView{
						PaymentCategory: "Крипто валюта",
						CryptoPayment:   repository.CryptoView{},
					}
				}

				var addrView repository.AddressView

				addr, err := crud.Get[models.Address, *uint](h.DB, ords.ReleasedAddrID)
				if err != nil || addr == nil {
					addrView = repository.AddressView{}
				} else {
					addrView = repository.AddressView{
						ID:          addr.ID,
						City:        city.Name,
						Product:     product.Name,
						Quantity:    ords.Quantity,
						Description: addr.Description,
						Image:       addr.Image,
						AddedAt:     addr.AddedAt,
					}
				}

				orderView := repository.OrderView{
					ID:          ords.ID,
					ProductName: product.Name,
					CityName:    city.Name,
					Quantity:    ords.Quantity,
					Due:         ords.Due,
					CreatedAt:   ords.CreatedAt,
					Client: repository.UserView{
						ID:        user.ID,
						ChatID:    user.ChatID,
						Username:  user.Username,
						FirstName: user.FirstName,
						LastName:  user.LastName,
					},
					PaymentMethod: payment,
					Address:       addrView,
				}

				// Append the orderView to the slice
				orderViews = append(orderViews, orderView)
			}

		}

		// Pass the initial orders to the template
		data := gin.H{
			"Title":  "Employee Orders",
			"Orders": orderViews,
		}

		if tmpl != nil {
			if err := tmpl.ExecuteTemplate(c.Writer, "base", data); err != nil {
				h.Logger.Error("Error executing template", zap.Error(err))
				c.String(http.StatusInternalServerError, "Error executing template: %v", err)
				return
			}
		} else {
			h.Logger.Error("Template not found", zap.String("template", "employee.page.gohtml"))
			c.String(http.StatusInternalServerError, "Template not found")
		}
	}
}

func (h *Handler) GetOrderHandler() gin.HandlerFunc { //TODO: rename BroadcastLatestOrder
	return func(c *gin.Context) {
		// Create the client channel for this connection
		clientChan := make(chan interface{})

		h.SSES.Register <- clientChan
		defer func() {
			h.SSES.Unregister <- clientChan
		}()

		// Set the headers for SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		// Send existing orders immediately after client connects
		orders, err := fetchOrdersFromDB(h.SSES.DB)
		if err != nil {
			h.SSES.Logger.Error("Failed to fetch orders from database", zap.Error(err))
			return
		}
		_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", dataToJSON(orders))
		c.Writer.Flush()

		// Keep the connection open and send updates as they arrive
		for {
			select {
			case message := <-clientChan:
				_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", dataToJSON(message))
				c.Writer.Flush()
			case <-c.Writer.CloseNotify():
				return
			}
		}
	}
}

func (h *Handler) ReleaseOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		orderIDStr := c.Param("orderId")
		//TODO:userIDStr := c.MustGet("userID").(int64) //Fetch employee id

		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		var order models.Order
		//if err := h.DB.Preload("AddrToRelease").Preload("User").First(&order, orderID).Error;
		if err := h.DB.Preload("AddrToRelease").Preload("User").First(&order, orderID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		if order.AddrToRelease == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No address assigned to this order"})
			return
		}

		order.Released = true
		//TODO:order.ReleasedByID = &userID
		order.ReleaseTime = time.Now()

		order.AddrToRelease.Released = true
		order.AddrToRelease.ReleaseDate = time.Now()
		order.AddrToRelease.ReleasedTo = fmt.Sprintf("%d", order.UserID)

		err = h.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&order).Error; err != nil {
				return err
			}

			if err := tx.Save(&order.AddrToRelease).Error; err != nil {
				return err
			}

			return nil
		})

		// Send the message to Telegram with the address details
		if err := h.SendMessage(c, order.User.ChatID, order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Telegram message"})
			return
		}

		// Broadcast the order release status
		message := map[string]interface{}{
			"id":       order.ID,
			"released": true,
		}

		h.SSES.BroadcastMessage(message)

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Order released successfully", "orderId": order.ID})
	}
}

func (h *Handler) GetOrderAddressHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIdStr := c.Query("orderId")

		orderID, err := strconv.Atoi(orderIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		}

		order, err := crud2.GetOrderByID(h.DB, orderID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		}

		addresses, err := crud2.GetAvailableAddresses(h.DB, order.CityID, order.ProductID, order.Quantity)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		}

		// Sorting addresses by ID in ascending order (if not already sorted)
		sort.Slice(addresses, func(i, j int) bool {
			return addresses[i].ID < addresses[j].ID
		})

		addressViews := make([]repository.AddressView, len(addresses))
		for i, addr := range addresses {
			addressViews[i] = repository.AddressView{
				ID:          addr.ID,
				City:        order.City.Name,
				Product:     order.Product.Name,
				Quantity:    order.Quantity,
				Description: addr.Description,
				Image:       addr.Image,
				AddedAt:     addr.AddedAt,
				//TODO:AddedBy:     addr.AddedBy.Username,
			}

		}

		c.JSON(http.StatusOK, gin.H{"addresses": addressViews})

	}
}

func (h *Handler) PostOrderAddressHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request struct {
			OrderID   int  `json:"order_id"`
			AddressID uint `json:"address_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		order, err := crud2.GetOrderByID(h.DB, request.OrderID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		}

		address, err := crud2.GetAddressByID(h.DB, request.AddressID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		}

		order.ReleasedAddrID = &request.AddressID

		if err := h.DB.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
			return
		}

		address.Assigned = true
		address.AssignedUserID = &order.UserID
		if err := h.DB.Save(&address).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order"})
			return
		}

		cardPayment, _ := crud.Get[models.Card, uint](h.DB, order.PaymentMethodID)

		var payment repository.PaymentView
		if order.PaymentMethodType == "card" {
			payment = repository.PaymentView{
				PaymentCategory: "Перевод на карту",
				CardPayment: repository.CardView{
					BankName:   cardPayment.BankName,
					BankUrl:    cardPayment.BankURL,
					CardNumber: cardPayment.CardNumber,
					FirstName:  cardPayment.FirstName,
					LastName:   cardPayment.LastName,
					UserName:   cardPayment.UserID,
					Password:   cardPayment.Password,
				},
			}
		} else if order.PaymentMethodType == "crypto" { //TODO: Fill crypto payment
			payment = repository.PaymentView{
				PaymentCategory: "Крипто валюта",
				CryptoPayment:   repository.CryptoView{},
			}
		}

		// Prepare the OrderView structure for broadcasting
		orderView := repository.OrderView{
			ID:          order.ID,
			ProductName: order.Product.Name,
			CityName:    order.City.Name,
			Quantity:    order.Quantity,
			Due:         order.Due,
			CreatedAt:   order.CreatedAt,
			Client: repository.UserView{
				ID:        order.User.ID,
				ChatID:    order.User.ChatID,
				Username:  order.User.Username,
				FirstName: order.User.FirstName,
				LastName:  order.User.LastName,
			},
			PaymentMethod: payment,
			Address: repository.AddressView{
				ID:          address.ID,
				City:        address.City.Name,
				Product:     address.Product.Name,
				Quantity:    order.Quantity,
				Description: address.Description,
				Image:       address.Image,
				AddedAt:     address.AddedAt,
				AddedBy:     "Employee Name", //TODO: fetch and set the actual employee name if needed
			},
		}

		// Broadcast real-time update (SSE, WebSocket, etc.)
		h.SSES.BroadcastMessage(orderView)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"order":   orderView,
		})

	}
}

// SendMessage sends a message to the user via Telegram
func (h *Handler) SendMessage(c *gin.Context, chatID int64, order models.Order) error {
	// Construct the message text
	messageText := fmt.Sprintf(
		"Детали по заказу #%d.\n\n  %s\n",
		order.ID,
		order.AddrToRelease.Description,
	)

	imageURL := h.GetImageURL(c, order.AddrToRelease.Image)

	photo := tu.Photo(
		// Chat ID as String (target username)
		tu.ID(chatID),

		// Send using file from disk
		//tu.File(mustOpen(order.AddrToRelease.Image)),
		tu.FileFromURL(imageURL),
	).WithCaption(messageText)

	// Sending photo
	_, err := h.Bot.API.SendPhoto(photo)
	if err != nil {
		return fmt.Errorf("failed to send image to Telegram: %v", err)
	}

	//msg := tu.Message(
	//	tu.ID(chatID),
	//	messageText,
	//)
	//
	//if _, err := h.Bot.API.SendMessage(msg); err != nil {
	//	//TODO: use zap log
	//	return fmt.Errorf("failed to send message to Telegram: %v", err)
	//}
	return nil
}

// Helper function to convert any data to JSON
func dataToJSON(data interface{}) string {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "[]"
	}
	return string(dataJSON)
}

// Fetch existing orders from the database
func fetchOrdersFromDB(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order
	if err := db.Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func mustOpen(filename string) *os.File {
	cwd, _ := os.Getwd()

	filePath := filepath.Join(cwd, filename)

	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file: %s, error: %v", filePath, err))
	}
	return file
}

func createOrderView(db *gorm.DB, order models.Order) repository.OrderView {
	product, _ := crud.Get[models.Product, uint](db, order.ProductID)
	city, _ := crud.Get[models.City, uint](db, order.CityID)
	user, _ := crud.Get[models.User, int64](db, order.UserID)

	return repository.OrderView{
		ID:          order.ID,
		ProductName: product.Name,
		CityName:    city.Name,
		Quantity:    order.Quantity,
		Due:         order.Due,
		CreatedAt:   order.CreatedAt,
		Client: repository.UserView{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		PaymentMethod: repository.PaymentView{},
		Address:       repository.AddressView{},
	}
}

func (h *Handler) GetImageURL(c *gin.Context, image string) string {
	// Get the request's host and scheme
	scheme := "http"
	if c.Request.TLS != nil { // Check if request is using HTTPS
		scheme = "https"
	}

	// Construct the full URL to the image
	return fmt.Sprintf("%s://%s/api/get/images/%s", scheme, c.Request.Host, image)
}
