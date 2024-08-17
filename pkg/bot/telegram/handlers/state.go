package handlers

import (
	"Magaz/pkg/utils/state/fsm"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Payload struct {
	Bot    *telego.Bot
	Update telego.Update
}

// StartHandler handles the start state
func StartHandler(event fsm.Event, f *fsm.FSM) error {
	// Type assert the payload to the expected type
	payload, ok := event.Payload.(Payload)
	if !ok {
		return fmt.Errorf("invalid payload type, expected *telego.Update")
	}

	// Inline keyboard parameters
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow( // Row 1
			tu.InlineKeyboardButton("Go"). // Column 1
							WithCallbackData("go"),
			tu.InlineKeyboardButton("Callback data button 2"). // Column 2
										WithCallbackData("callback_2"),
		),
		tu.InlineKeyboardRow( // Row 2
			tu.InlineKeyboardButton("URL button").WithURL("https://example.com"), // Column 1
		),
	)

	// Message parameters
	message := tu.Message(
		tu.ID(payload.Update.Message.Chat.ID),
		"My message",
	).WithReplyMarkup(inlineKeyboard)

	// Sending message
	_, _ = payload.Bot.SendMessage(message)

	//return f.Trigger(event)
	return nil
}

// CityHandler handles the city state
func CityHandler(event fsm.Event, f *fsm.FSM) error {
	fmt.Println("Please enter the state you live in.")
	return f.Trigger(event)
}

// ProductHandler handles the product state
func ProductHandler(event fsm.Event, f *fsm.FSM) error {
	fmt.Println("Please enter the state you live in.")
	return f.Trigger(event)
}

// QuantityHandler handles the quantity state
func QuantityHandler(event fsm.Event, f *fsm.FSM) error {
	fmt.Println("Please enter the state you live in.")
	return f.Trigger(event)
}

// PaymentHandler handles the payment state
func PaymentHandler(event fsm.Event, f *fsm.FSM) error {
	fmt.Println("Please enter the state you live in.")
	return f.Trigger(event)
}

// EndHandler handles the end state
func EndHandler(event fsm.Event, f *fsm.FSM) error {
	fmt.Println("Please enter the state you live in.")
	return f.Trigger(event)
}
