package handler

import (
	fsm "Magaz/services/bots/telegram/v2/utils/stateMngs"
	"fmt"
	"github.com/mymmrac/telego"
)

func HandleEmplUpdate(sm *fsm.RuleBasedFSM, bot *telego.Bot, update telego.Update) {
	fmt.Printf("Handling employee update: %+v\n", update)
}
func HandleEmplCallbackQuery(sm *fsm.RuleBasedFSM, bot *telego.Bot, query telego.CallbackQuery) {

}
