package handler

import (
	"github.com/mymmrac/telego"
	msg "tg/internal/bot/handler/message"
	fsm "tg/pkg/utils/stateMngs"
)

func HandleEmplUpdate(sm *fsm.RuleBasedFSM, bot *telego.Bot, update telego.Update) {
	if update.Message.Text == "/start" {
		msg.SendMessage(bot, "Hi client")
	}
}
func HandleEmplCallbackQuery(sm *fsm.RuleBasedFSM, bot *telego.Bot, query telego.CallbackQuery) {

}
