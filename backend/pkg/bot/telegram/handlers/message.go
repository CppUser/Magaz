package handlers

func IsUserInteractingWithBot(context map[string]interface{}) bool {
	return context["message"] == "/start"
}

// TODO: research on map[string]string insertion order . currently using hack with slice kvp to keep order
type TempMarkup struct {
	Text         string
	CallbackData string
}

//
//type Payload struct {
//	Bot      *telego.Bot
//	Update   interface{}
//	Callback string
//	Cache    *redis.Client
//	DB       *gorm.DB
//}
//
//func StartHandler(payload interface{}) {
//	//TODO: retrieve from db available cities
//	//TODO: do handling if user clicked back (to choose)button with inline keyboard
//	//because right now it handles text only update := payload.(*Payload).Update.(telego.Update)
//	//do if check text do this if callbackquery do update := payload.(*Payload).Update.(telego.CallbackQuery) etc
//
//	//TODO: add reset timer if user didnt reach final state
//
//	update := payload.(*Payload).Update.(telego.Update)
//	bot := payload.(*Payload).Bot
//
//	// Inline keyboard parameters
//	partnerKeyboard := tu.InlineKeyboard(
//		tu.InlineKeyboardRow( // Row 1
//			tu.InlineKeyboardButton("Вакансии"). // Column 1
//								WithCallbackData("partnership:NewYork"),
//		),
//	)
//
//	OrderKeyboard := tu.InlineKeyboard(
//		tu.InlineKeyboardRow( // Row 1
//			tu.InlineKeyboardButton("Оформит"). // Column 1
//								WithCallbackData("order:confirmed"), //TODO: Make order handler ?
//		),
//	)
//
//	// Message parameters
//	greetingsMessage := tu.Message(
//		tu.ID(update.Message.GetChat().ID),
//		"Добро пожаловать в магазин",
//	).WithReplyMarkup(partnerKeyboard)
//
//	orderMessage := tu.Message(
//		tu.ID(update.Message.GetChat().ID),
//		"Желаете оформить заказ?",
//	).WithReplyMarkup(OrderKeyboard)
//
//	_, _ = bot.SendMessage(greetingsMessage)
//	_, _ = bot.SendMessage(orderMessage)
//
//}
//func OrderHandler(payload interface{}) {
//	query := payload.(*Payload).Update.(*telego.Message)
//	bot := payload.(*Payload).Bot
//
//	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
//		ChatID:    tu.ID(query.GetChat().ID),
//		MessageID: query.GetMessageID(),
//		Text:      "Выберите город:",
//		ParseMode: telego.ModeHTML,
//		ReplyMarkup: tu.InlineKeyboard(
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("Краснодар"). // Column 1
//									WithCallbackData("city:Краснодар"),
//			),
//			tu.InlineKeyboardRow( // Row 2
//				tu.InlineKeyboardButton("Армавир"). // Column 1
//									WithCallbackData("city:Армавир"),
//			),
//		),
//	})
//
//}
//
//func CityHandler(payload interface{}) {
//
//	query := payload.(*Payload).Update.(*telego.Message)
//	bot := payload.(*Payload).Bot
//
//	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
//		ChatID:    tu.ID(query.GetChat().ID),
//		MessageID: query.GetMessageID(),
//		Text:      "Выберите желаемый товар:",
//		ParseMode: telego.ModeHTML,
//		ReplyMarkup: tu.InlineKeyboard(
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("Товар 1"). // Column 1
//									WithCallbackData("product:1"),
//			),
//		),
//	})
//}
//
//func ProductHandler(payload interface{}) {
//	query := payload.(*Payload).Update.(*telego.Message)
//	bot := payload.(*Payload).Bot
//
//	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
//		ChatID:    tu.ID(query.GetChat().ID),
//		MessageID: query.GetMessageID(),
//		Text:      "Выберите желаемое количество:",
//		ParseMode: telego.ModeHTML,
//		ReplyMarkup: tu.InlineKeyboard(
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("количество 1 - цена 1"). // Column 1
//											WithCallbackData("quantity:1"),
//			),
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("количество 2 - цена 2"). // Column 1
//											WithCallbackData("quantity:2"),
//			),
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("количество 3 - цена 3"). // Column 1
//											WithCallbackData("quantity:3"),
//			),
//		),
//	})
//
//}
//
//func QuantityHandler(payload interface{}) {
//	query := payload.(*Payload).Update.(*telego.Message)
//	bot := payload.(*Payload).Bot
//
//	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
//		ChatID:    tu.ID(query.GetChat().ID),
//		MessageID: query.GetMessageID(),
//		Text:      "Выберите метод оплаты:",
//		ParseMode: telego.ModeHTML,
//		ReplyMarkup: tu.InlineKeyboard(
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("Перевод на Карту"). // Column 1
//										WithCallbackData("payment:1"),
//			),
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("Оплата Крипто валютой"). // Column 1
//											WithCallbackData("payment:2"),
//			),
//		),
//	})
//
//}
//
//func PaymentHandler(payload interface{}) {
//	query := payload.(*Payload).Update.(*telego.Message)
//	bot := payload.(*Payload).Bot
//
//	_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
//		ChatID:    tu.ID(query.GetChat().ID),
//		MessageID: query.GetMessageID(),
//		Text:      "Завершить заказ?:",
//		ParseMode: telego.ModeHTML,
//		ReplyMarkup: tu.InlineKeyboard(
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("Заказать"). // Column 1
//									WithCallbackData("conformation:true"),
//				tu.InlineKeyboardButton("Отменить"). // Column 1
//									WithCallbackData("conformation:false"),
//			),
//			tu.InlineKeyboardRow( // Row 1
//				tu.InlineKeyboardButton("Добавить в корзину"). // Column 1
//										WithCallbackData("cart:true"),
//			),
//		),
//	})
//
//}
//
//func ConformationHandler(payload interface{}) {
//	// Logic for the end state
//	fmt.Println("Handling conformation state...")
//
//}
//
//func FinalHandler(payload interface{}) {
//	// Logic for the end state
//	fmt.Println("Handling Final state...")
//
//}

//
//type Payload struct {
//	Bot    *telego.Bot
//	Update telego.Update
//}
//
//// StartHandler handles the start state
//func StartHandler(event fsm.Event, f *fsm.FSM) error {
//	log.Printf("StartHandler - Current State: %s, Event: %s", f.CurrentState, event.Name)
//
//	// Type assert the payload to the expected type
//	payload, ok := event.Payload.(Payload)
//	if !ok {
//		return fmt.Errorf("invalid payload type, expected *telego.Update")
//	}
//
//	nextState, exists := f.Transitions[f.CurrentState][event]
//	if exists {
//		log.Printf("Transition found: Current State: %s, Event: %s, Next State: %s", f.CurrentState, event.Name, nextState)
//	} else {
//		log.Printf("No transition found for Current State: %s, Event: %s", f.CurrentState, event.Name)
//	}
//
//	// Inline keyboard parameters
//	inlineKeyboard := tu.InlineKeyboard(
//		tu.InlineKeyboardRow( // Row 1
//			tu.InlineKeyboardButton("Go"). // Column 1
//							WithCallbackData("go"),
//			tu.InlineKeyboardButton("Callback data button 2"). // Column 2
//										WithCallbackData("callback_2"),
//		),
//		tu.InlineKeyboardRow( // Row 2
//			tu.InlineKeyboardButton("URL button").WithURL("https://example.com"), // Column 1
//		),
//	)
//
//	// Message parameters
//	message := tu.Message(
//		tu.ID(payload.Update.Message.Chat.ID),
//		"My message",
//	).WithReplyMarkup(inlineKeyboard)
//
//	// Sending message
//	_, _ = payload.Bot.SendMessage(message)
//
//	//return f.Trigger(event)
//	return nil
//}
//
//// CityHandler handles the city state
//func CityHandler(event fsm.Event, f *fsm.FSM) error {
//	fmt.Println("Please enter the state you live in.")
//	return f.Trigger(event)
//}
//
//// ProductHandler handles the product state
//func ProductHandler(event fsm.Event, f *fsm.FSM) error {
//	fmt.Println("Please enter the state you live in.")
//	return f.Trigger(event)
//}
//
//// QuantityHandler handles the quantity state
//func QuantityHandler(event fsm.Event, f *fsm.FSM) error {
//	fmt.Println("Please enter the state you live in.")
//	return f.Trigger(event)
//}
//
//// PaymentHandler handles the payment state
//func PaymentHandler(event fsm.Event, f *fsm.FSM) error {
//	fmt.Println("Please enter the state you live in.")
//	return f.Trigger(event)
//}
//
//// EndHandler handles the end state
//func EndHandler(event fsm.Event, f *fsm.FSM) error {
//	fmt.Println("Please enter the state you live in.")
//	return f.Trigger(event)
//}
