package handler

import (
	fsm "Magaz/services/bots/telegram/v2/utils/stateMngs"
	"github.com/mymmrac/telego"
)

func HandleClientUpdate(sm *fsm.RuleBasedFSM, bot *telego.Bot, update telego.Update) {
	sm.Context["bot"] = bot
	sm.Context["update"] = update

}

func HandleClientCallbackQuery(sm *fsm.RuleBasedFSM, bot *telego.Bot, query telego.CallbackQuery) {

}
