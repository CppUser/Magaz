package ws

import (
	"Magaz/backend/internal/repository"
	crud "Magaz/backend/internal/storage"
	crud2 "Magaz/backend/internal/storage/crud"
	"Magaz/backend/internal/storage/models"
	"fmt"
)

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type EventHandler func(e Event, c *Client) error

const (
	EventSendMessage  = "send_message"
	EventOrderRelease = "order_release"
	//EventOutAssignAddress response to EventInUpdateAddress
	EventOutAssignAddress = "assign_address"
	//EventInUpdateAddress request to update address
	EventInUpdateAddress = "update_address"
)

type SendMessageEvent struct {
	Type    string `json:"type"`
	Payload string `json:"message"`
}

func SendMessageHandler(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

type OrderReleaseEvent struct {
}

func OrderReleaseHandler(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

type UpdateAddressEvent struct {
}

// TODO: refactor move database manipulation to its own function
func UpdateAddressHandler(event Event, c *Client) error {

	payloadMap, ok := event.Payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("payload is not a valid map: %v", event.Payload)
	}

	addressId, addressOk := payloadMap["addressId"].(float64)
	orderId, orderOk := payloadMap["orderId"].(float64)

	if !addressOk || !orderOk {
		//TODO: Log properly
		return fmt.Errorf("missing or invalid fields in the payload: %v", payloadMap)
	}
	/////////////////////////////////////////////////////////////////////////////////////////////
	//          TODO: Move to AssignAddress handler and call it here
	/////////////////////////////////////////////////////////////////////////////////////////////
	order, err := crud2.GetOrderByID(c.mng.DB, int(orderId))
	if err != nil {
		return err
	}

	address, err := crud2.GetAddressByID(c.mng.DB, uint(addressId))
	if err != nil {
		return err
	}

	// Start a transaction
	tx := c.mng.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//TODO: Refactor modal to use better approach (converting hack)
	order.ReleasedAddrID = func(i uint) *uint { return &i }(uint(addressId))
	if err := c.mng.DB.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	address.Assigned = true
	address.AssignedUserID = &order.UserID

	if err := c.mng.DB.Save(&address).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	//Prep outgoing data to front end

	cardPayment, _ := crud.Get[models.Card, uint](c.mng.DB, order.PaymentMethodID)

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
			AddedBy:     address.AddedBy.Username, //TODO: Use crud.GetWithAsssociation to pull employee (its not preloaded with GetOrderByID)
		},
	}

	//data, err := json.Marshal(orderView)
	//if err != nil {
	//	return fmt.Errorf("failed to marshal broadcast message: %v", err)
	//}

	/////////////////////////////////////////////////////////////////////////////////////////////

	/////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: Use SendMessageHandler ?
	/////////////////////////////////////////////////////////////////////////////////////////////

	//Sending respond to updated data
	var out Event
	out.Type = EventOutAssignAddress
	out.Payload = orderView

	//TODO: Encode message ?
	for client := range c.mng.clients { //TODO: Fix later to send to correct person
		client.egress <- out
	}

	return nil
}

type AssignAddressEvent struct {
}

func AssignAddressHandler(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}
