package handler

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	fsm "tg/pkg/utils/stateMngs"
)

func HandleClientUpdate(sm *fsm.RuleBasedFSM, bot *telego.Bot, update telego.Update) {
	sm.Context["bot"] = bot
	sm.Context["update"] = update

	if update.Message != nil {
		//message := update.Message

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.GetChat().ID),
			"Hi",
		))
	} else if update.CallbackQuery != nil {

	}

}

func HandleClientCallbackQuery(sm *fsm.RuleBasedFSM, bot *telego.Bot, query telego.CallbackQuery) {

}
