package msg

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	fsm "tg/pkg/utils/stateMngs"
)

type CallbackQueryMarkup struct {
	Text         string
	CallbackData string
}

func SendMessage(bot *telego.Bot, text string) fsm.ActionFunc {
	return func(context map[string]interface{}) error {
		update := context["update"].(telego.Update)

		_, err := bot.SendMessage(tu.Message(
			tu.ID(update.Message.GetChat().ID),
			text,
		))
		return err
	}
}

// EditMessage edits message.
func EditMessage(bot *telego.Bot, text string) fsm.ActionFunc {
	return func(context map[string]interface{}) error {
		update := context["message"].(*telego.Message)

		// Edit the message
		_, err := bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    tu.ID(update.GetChat().ID),
			MessageID: update.GetMessageID(),
			Text:      text,
		})
		return err
	}
}

// SendMessageWithMarkup sends a message with optional inline keyboard markup.
func SendMessageWithMarkup(bot *telego.Bot, text string, markup map[string]string) fsm.ActionFunc {
	return func(context map[string]interface{}) error {
		update := context["update"].(telego.Update)

		message := tu.Message(
			tu.ID(update.Message.GetChat().ID),
			text,
		)

		// If markup is provided, add inline keyboard
		if len(markup) > 0 {
			var keyboard [][]telego.InlineKeyboardButton
			var row []telego.InlineKeyboardButton

			for buttonText, callbackData := range markup {
				button := telego.InlineKeyboardButton{
					Text:         buttonText,
					CallbackData: callbackData,
				}
				row = append(row, button)
			}

			keyboard = append(keyboard, row)
			inlineKeyboard := telego.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}
			message.WithReplyMarkup(&inlineKeyboard)
		}

		// Send the message
		_, err := bot.SendMessage(message)
		return err
	}
}

// EditMessageWithMarkup edits  message with optional inline keyboard markup.
func EditMessageWithMarkup(bot *telego.Bot, text string, markup []CallbackQueryMarkup) fsm.ActionFunc {
	return func(context map[string]interface{}) error {
		update := context["message"].(*telego.Message)

		// Prepare the edit message parameters
		editParams := &telego.EditMessageTextParams{
			ChatID:    tu.ID(update.GetChat().ID),
			MessageID: update.GetMessageID(),
			Text:      text,
		}

		// Create the inline keyboard
		var keyboard [][]telego.InlineKeyboardButton

		// Add buttons in the order they are provided
		for _, item := range markup {
			button := telego.InlineKeyboardButton{
				Text:         item.Text,
				CallbackData: item.CallbackData,
			}
			// Each button in its own row (adjust this if needed)
			keyboard = append(keyboard, []telego.InlineKeyboardButton{button})
		}

		// Assign the inline keyboard to the edit message parameters
		if len(keyboard) > 0 {
			inlineKeyboard := telego.InlineKeyboardMarkup{
				InlineKeyboard: keyboard,
			}
			editParams.ReplyMarkup = &inlineKeyboard
		}

		// Edit the message
		_, err := bot.EditMessageText(editParams)
		return err
	}
}

func EditMessageWithOrderDetails(bot *telego.Bot, text string, markup []CallbackQueryMarkup) fsm.ActionFunc {
	return func(context map[string]interface{}) error {
		update := context["message"].(*telego.Message)

		//

		// Edit the message
		_, err := bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    tu.ID(update.GetChat().ID),
			MessageID: update.GetMessageID(),
			Text:      text,
		})
		return err
	}
}
