//TODO: Current setup allows for bot to join group without permission. that can cause a problem with flood and load of large groups to
//infintrate  bot
//Try to imlement something like this in backend or here in ReceiveUpdate
/*

 // Define allowed group chat IDs
    allowedGroups := map[int64]bool{
        -1001234567890: true, // Replace with your actual allowed group chat IDs
        -1009876543210: true, // Add more allowed groups if needed
    }

    // Fetch updates (messages) from the bot
    updates, err := bot.GetUpdates(nil)
    if err != nil {
        log.Fatalf("Failed to get updates: %s", err)
    }

    for _, update := range updates {
        // Check if it's a message and comes from a group chat
        if update.Message != nil && update.Message.Chat.Type == "supergroup" {
            chatID := update.Message.Chat.ID

            // Check if the chat is in the allowed groups list
            if _, ok := allowedGroups[chatID]; !ok {
                // If not allowed, ignore the message
                log.Printf("Ignoring message from unauthorized group chat ID: %d", chatID)
                continue
            }

            // If allowed, process the message
            log.Printf("Message from allowed group: %d, content: %s", chatID, update.Message.Text)

            // Example: reply to the message
            _, err = bot.SendMessage(&telego.SendMessageParams{
                ChatID: telego.ChatID{ID: chatID},
                Text:   "Message received from authorized group!",
            })
            if err != nil {
                log.Printf("Failed to send message: %s", err)
            }
        }
    }


*/

package telegram

import (
	"Magaz/backend/internal/config"
	ws "Magaz/backend/internal/system/websocket"
	tgconfig "Magaz/backend/pkg/bot/telegram/config"
	"Magaz/backend/pkg/utils/state/fsm"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type Bot struct {
	Config           *config.TGBotConfig
	API              *telego.Bot
	Logger           *zap.Logger
	UpdateChanBuffer uint
	UpdatesChan      chan telego.Update
	Cache            *redis.Client
	DB               *gorm.DB
	//FSM              *fsm.RuleBasedFSM
	WS *ws.Manager
	//Hub              *sse.SSEHub

}

// TODO: refactor code move some logic to handlers

// InitBot initializes the Telegram bot
func (b *Bot) InitBot() {

	//TODO: need to handle error properly, currently removed do to code complaint
	_, _ = tgconfig.LoadConfig("bot_config", "yaml", []string{".", "backend/config/"})
	//if err != nil {
	//	b.Config.Logger.Fatal("Failed to load bot configs", zap.String("error", err.Error()))
	//}

	bot, err := telego.NewBot(b.Config.Token)
	if err != nil {
		//TODO: need to handle the error differently without direct call to zap.String
		b.Logger.Fatal("failed create new bot api instance", zap.String("error", err.Error()))

	}
	b.API = bot

	//TODO: find better way to make channel
	//TODO: also need to safly send updates to channel, checking is channel is open or closed
	// Initialize the updates channel
	b.UpdatesChan = make(chan telego.Update, b.UpdateChanBuffer)

	//TODO: refer to SetWebhookParams to setup additional parameters (like certificate, pending updates, etc.)
	_ = bot.SetWebhook(&telego.SetWebhookParams{
		URL: b.Config.WebhookLink + b.Config.WebhookPath,
		AllowedUpdates: []string{
			"message",
			"edited_message",
			"callback_query",
			"inline_query",
			"chosen_inline_result",
			"poll",
			"poll_answer",
			"shipping_query",
			"pre_checkout_query",
			"my_chat_member",
			"chat_member",
		},
	})

	info, _ := bot.GetWebhookInfo()
	b.Logger.Info("Webhook Info", zap.Any("info", info)) //TODO: in prod it needs to be in JSON format

	//TODO: TEMPORARY refactor code
	//clientRules := []fsm.Rule{
	//	{
	//		Event:      "/start",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			//TODO: add it when admin has option to post it in settings
	//			//handlers.SendMessageWithMarkup(bot, "Добрый день", map[string]string{
	//			//	"Вакансии": "hire",
	//			//}),
	//
	//			//check if user exist in cache
	//
	//			handlers.SendMessageWithMarkup(bot, "Желаете оформить заказ?", map[string]string{
	//				"Оформит": "order",
	//			}),
	//		},
	//	},
	//	{
	//		Event:      "hire",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Тут будет сообщение о открытых вакансиях")},
	//	},
	//	{
	//		Event:      "order",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//
	//				// Fetch cities from the database
	//				cities, err := crud.GetAllCities(b.DB)
	//				if err != nil {
	//					return err
	//				}
	//
	//				// Generate city markup
	//				cityMarkup := make([]handlers.TempMarkup, len(cities))
	//				for i, city := range cities {
	//					cityMarkup[i] = handlers.TempMarkup{
	//						Text:         city.Name,
	//						CallbackData: "city:" + city.Name, // Pass city name in the callback data
	//					}
	//				}
	//
	//				// Prompt the user to select a city
	//				return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите ваш город:", cityMarkup)(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "city",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				// Extract the city name from the callback data
	//				callbackData := context["callbackData"].(string)
	//				cityName := strings.TrimPrefix(callbackData, "city:")
	//
	//				user := context["from"].(telego.User)
	//				StoreUserChoice(b.Cache, user.ID, "city", cityName)
	//
	//				// Generate product markup for the selected city
	//				productMarkup, err := GenerateProductMarkup(b.DB, cityName)
	//				if err != nil {
	//					return err
	//				}
	//
	//				// Prompt the user to select a product
	//				return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующий вас товар:", productMarkup)(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "product",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				// Extract the product name from the context
	//				callbackData := context["callbackData"].(string)
	//				productName := strings.TrimPrefix(callbackData, "product:")
	//
	//				user := context["from"].(telego.User)
	//				StoreUserChoice(b.Cache, user.ID, "product", productName)
	//
	//				// Generate quantity markup for the selected product
	//				quantityMarkup, err := GenerateProductPriceMarkup(b.DB, productName)
	//				if err != nil {
	//					return err
	//				}
	//
	//				// Prompt the user to select a quantity
	//				return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующее вас количество:", quantityMarkup)(context)
	//			},
	//		},
	//	},
	//	//TODO: add next rule to add to cart to shop for more items
	//	{
	//		Event:      "quantity",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				callbackData := context["callbackData"].(string)
	//				qtAmount := strings.TrimPrefix(callbackData, "quantity:")
	//
	//				user := context["from"].(telego.User)
	//				StoreUserChoice(b.Cache, user.ID, "quantity", qtAmount)
	//
	//				choices, err := GetUserChoices(b.Cache, user.ID)
	//				if err != nil {
	//					return err
	//				}
	//				city, _ := crud.GetCityIDByName(b.DB, choices["city"])
	//				prID, _ := crud.GetProductIDByCityAndProductName(b.DB, choices["city"], choices["product"])
	//				qt, _ := strconv.ParseFloat(choices["quantity"], 32)
	//
	//				// Use the GetAvailableAddresses function to check if there are any available addresses
	//				availableAddresses, err := crud.GetAvailableAddresses(b.DB, city, prID, float32(qt))
	//				if err != nil || len(availableAddresses) == 0 {
	//					// If no addresses available, notify the user
	//					return handlers.EditMessage(bot, "К сожалению, нет доступных адресов для выбранного количества. Пожалуйста, начните сначала с команды /start.")(context)
	//				}
	//
	//				return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите метод оплаты:", []handlers.TempMarkup{ //TODO: replace with generated markup
	//					{Text: "Перевод на карту", CallbackData: "card"},
	//					//{Text: "Оплата Крипто валютой", CallbackData: "crypto"},
	//				})(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "card",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				// Extract the product name from the context
	//				callbackData := context["callbackData"].(string)
	//				//qtAmount := strings.TrimPrefix(callbackData, "quantity:")
	//
	//				user := context["from"].(telego.User)
	//				StoreUserChoice(b.Cache, user.ID, "payment", callbackData)
	//
	//				//TODO: Generate payment method markup
	//
	//				return handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты перевод на карту", []handlers.TempMarkup{ //TODO: replace with generated markup
	//					{Text: "Оформить заказ", CallbackData: "confirm"}, //TODO: Add option to add to cart for multiple products
	//					{Text: "Отменить заказ", CallbackData: "cancel"},
	//				})(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "crypto",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				// Extract the product name from the context
	//				//callbackData := context["callbackData"].(string)
	//				//qtAmount := strings.TrimPrefix(callbackData, "quantity:")
	//
	//				//TODO: Generate payment method markup
	//
	//				return handlers.EditMessageWithMarkup(bot, "Выберите тип крипто валюты", []handlers.TempMarkup{ //TODO: replace with generated markup
	//					{Text: "Bitcoin", CallbackData: "bitcoin"}, //TODO: Add option to add to cart for multiple products
	//					{Text: "Etherium", CallbackData: "etherium"},
	//				})(context)
	//			},
	//		},
	//	},
	//
	//	{
	//		Event:      "bitcoin",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				// Extract the product name from the context
	//				callbackData := context["callbackData"].(string)
	//				//qtAmount := strings.TrimPrefix(callbackData, "quantity:")
	//
	//				user := context["from"].(telego.User)
	//				StoreUserChoice(b.Cache, user.ID, "payment", callbackData)
	//
	//				//TODO: Generate payment method markup
	//
	//				return handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты bitcoin крипто валютой", []handlers.TempMarkup{ //TODO: replace with generated markup
	//					{Text: "Оформить заказ", CallbackData: "confirm"}, //TODO: Add option to add to cart for multiple products
	//					{Text: "Отменить заказ", CallbackData: "cancel"},
	//				})(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "etherium",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				// Extract the product name from the context
	//				callbackData := context["callbackData"].(string)
	//				//qtAmount := strings.TrimPrefix(callbackData, "quantity:")
	//
	//				user := context["from"].(telego.User)
	//				StoreUserChoice(b.Cache, user.ID, "payment", callbackData)
	//
	//				//TODO: Generate payment method markup
	//
	//				return handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты etherium крипто валютой", []handlers.TempMarkup{ //TODO: replace with generated markup
	//					{Text: "Оформить заказ", CallbackData: "confirm"}, //TODO: Add option to add to cart for multiple products
	//					{Text: "Отменить заказ", CallbackData: "cancel"},
	//				})(context)
	//			},
	//		},
	//	},
	//
	//	{
	//		Event:      "confirm",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//				user := context["from"].(telego.User)
	//				message := context["message"].(*telego.Message)
	//
	//				//TODO: check if order is in process to avoid flooding
	//				choices, err := GetUserChoices(b.Cache, user.ID)
	//				if err != nil {
	//					return err
	//				}
	//				city, _ := crud.GetCityIDByName(b.DB, choices["city"])
	//				prID, _ := crud.GetProductIDByCityAndProductName(b.DB, choices["city"], choices["product"])
	//				qt, _ := strconv.ParseFloat(choices["quantity"], 32)
	//
	//				// Check if the user already has 2 or more active (non-released) orders
	//				var activeOrders []models.Order
	//				if err := b.DB.Where("user_id = ? AND released = ?", user.ID, false).Find(&activeOrders).Error; err != nil {
	//					return fmt.Errorf("failed to fetch active orders: %w", err)
	//				}
	//
	//				// If the user already has 2 active (non-released) orders, prevent them from creating more
	//				if len(activeOrders) >= 2 {
	//					//TODO: send to operator chat instead right now since there is no option , using old method
	//					msg := "У вас уже есть 2 активных заказа. Пожалуйста, завершите один из них, прежде чем создавать новый.\n\n" +
	//						"Если у вас возникли проблемы с заказом свяжитесь с оператором старым методом"
	//
	//					// Add the details of the active orders to the message
	//					msg += "Ваши активные заказы:\n"
	//					for _, order := range activeOrders {
	//						// Since city and product are not preloaded, use the choices retrieved earlier
	//						msg += fmt.Sprintf(
	//							"Заказ #%d\n"+
	//								"Город: %s\n"+
	//								"Товар: %s\n"+
	//								"Количество: %.2f\n"+
	//								"Сумма к оплате: %d\n\n",
	//							order.ID,
	//							choices["city"],
	//							choices["product"],
	//							order.Quantity,
	//							order.Due,
	//						)
	//					}
	//
	//					return handlers.EditMessage(bot, msg)(context)
	//				}
	//
	//				var qtnPrice models.QtnPrice //TODO: Refactor
	//				// Find the price for the given quantity
	//				if err := b.DB.Where("city_product_id = ? AND quantity = ?", prID, qt).First(&qtnPrice).Error; err != nil { //passing wrong city id
	//					b.Logger.Error("price not found for the specified quantity", zap.String("error", err.Error()))
	//				}
	//
	//				pmt, _ := crud.GetPaymentMethod(b.DB, choices["payment"])
	//				address, _ := repository.GetRandomAddress(b.DB, city, prID, float32(qt), user.ID)
	//
	//				//TODO: Need to figure out how to store custom quantity
	//
	//				var msg string
	//				order := models.Order{}
	//				ordView := repository.OrderView{}
	//
	//				if choices["payment"] == "card" {
	//
	//					order = models.Order{
	//						UserID:            user.ID,
	//						CityID:            city,
	//						ProductID:         prID,        //Get product from productCityID
	//						Quantity:          float32(qt), //TODO: retrieve quantity from cache
	//						Due:               uint(qtnPrice.Price),
	//						PaymentMethodType: choices["payment"],
	//						PaymentMethodID:   pmt.(models.Card).ID, //TODO: retrieve from PaymentMethodType name if card from card if crypto from crypto
	//						CreatedAt:         time.Now(),
	//						ReleasedAddrID:    &address.ID,
	//					}
	//
	//					var addr models.Address
	//					if err := b.DB.First(&addr, &address.ID).Error; err == nil {
	//						if !addr.Assigned {
	//							addr.Assigned = true
	//							addr.AssignedUserID = &user.ID
	//							//TODO add AssignedBy (Bot)
	//						}
	//					}
	//
	//					if err := b.DB.Create(&order).Error; err != nil {
	//						b.Logger.Error("Failed to create new order in DB", zap.String("error", err.Error()))
	//					}
	//
	//					ordView = repository.OrderView{
	//						ID:          order.ID,
	//						ProductName: choices["product"],
	//						CityName:    choices["city"],
	//						Quantity:    float32(qt),
	//						Due:         uint(qtnPrice.Price), //TODO: once card implemented need to add all items in cart to due
	//						CreatedAt:   time.Now(),
	//						Client: repository.UserView{
	//							ID:        user.ID,
	//							ChatID:    message.GetChat().ID,
	//							Username:  user.Username,
	//							FirstName: user.FirstName,
	//							LastName:  user.LastName,
	//						},
	//						PaymentMethod: repository.PaymentView{
	//							PaymentCategory: "Перевод на карту",
	//							CardPayment: repository.CardView{
	//								BankName:   pmt.(models.Card).BankName,
	//								BankUrl:    pmt.(models.Card).BankURL,
	//								CardNumber: pmt.(models.Card).CardNumber,
	//								FirstName:  pmt.(models.Card).FirstName,
	//								LastName:   pmt.(models.Card).LastName,
	//								UserName:   pmt.(models.Card).UserID,
	//								Password:   pmt.(models.Card).Password,
	//								QuickPay:   pmt.(models.Card).QuickPay,
	//							},
	//						},
	//						Address: *address,
	//					}
	//
	//					msg = fmt.Sprintf(
	//						"Номер заказа #%d\n"+
	//							"Город: %s\n"+
	//							"Товар: %s\n"+
	//							"Количество: %s\n"+
	//							"Метод оплаты: %s\n"+
	//							"Сумма к оплате: %d\n"+
	//							"\n"+
	//							"***** Данные Карты *****\n"+
	//							"Банк: %s\n"+
	//							"Номер карты: %s\n"+
	//							"ФИО: %s\n"+
	//							"СБП: %s\n"+
	//							"*************************\n",
	//						order.ID,
	//						choices["city"],
	//						choices["product"],
	//						choices["quantity"],
	//						"Перевод на Карту",
	//						uint(qtnPrice.Price),
	//						pmt.(models.Card).BankName,
	//						pmt.(models.Card).CardNumber,
	//						pmt.(models.Card).LastName+" "+pmt.(models.Card).FirstName,
	//						pmt.(models.Card).QuickPay,
	//					)
	//				} else if choices["payment"] == "crypto" {
	//					//TODO: Implement
	//				}
	//
	//				b.WS.BroadcastOrder(ordView)
	//				//TODO: Send message to Employee telegram about new order (to personal or to group chat)"4512552536"
	//
	//				empMessage := fmt.Sprintf("Добавлен новый заказ: %d\n", ordView.ID)
	//
	//				_, _ = bot.SendMessage(&telego.SendMessageParams{
	//					ChatID: tu.ID(b.Config.GroupID),
	//					Text:   empMessage,
	//				})
	//
	//				go func(orderID uint, addrID *uint) {
	//					//TODO: Does not work when there is constant ordering , need to find other way to put timer on order
	//					//Maybe some database watching system
	//					<-time.After(15 * time.Minute)
	//
	//					order.ReleasedAddrID = nil
	//					if err := b.DB.Save(&order).Error; err != nil {
	//						b.Logger.Error("Failed to update order status", zap.String("error", err.Error()))
	//					}
	//
	//					if addrID != nil {
	//						var addr models.Address
	//						if err := b.DB.First(&addr, *addrID).Error; err == nil {
	//							addr.Released = false
	//							addr.Assigned = false
	//							addr.AssignedUserID = nil
	//							if err := b.DB.Save(&addr).Error; err != nil {
	//								b.Logger.Error("Failed to unassign address", zap.String("error", err.Error()))
	//							}
	//						}
	//					}
	//					ordView.Address = repository.AddressView{}
	//					b.WS.BroadcastOrder(ordView)
	//
	//				}(order.ID, order.ReleasedAddrID)
	//
	//				return handlers.EditMessageWithMarkup(bot, msg, []handlers.TempMarkup{
	//					{Text: "Подтеврдить оплату", CallbackData: "payConf"},
	//				})(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "cancel",
	//		Conditions: []fsm.ConditionFunc{},
	//		Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Ваш заказ успешно отменен")},
	//	},
	//	{
	//		Event:      "payConf",
	//		Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
	//		Actions: []fsm.ActionFunc{
	//			func(context map[string]interface{}) error {
	//
	//				message := fmt.Sprintf(
	//					"Пожалуйста ожидайте ответа оператора\n" +
	//						"Функция подтвержденя оплаты в данный момент не работает\n" +
	//						"Если у вас имеется чек на руках отправьте его старым методом\n также укажите номер заказа")
	//
	//				return handlers.EditMessage(bot, message)(context)
	//			},
	//		},
	//	},
	//	{
	//		Event:      "status",
	//		Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
	//		Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Тут будет сообщение о статусе заказа")},
	//	},
	//}
	//// Initialize FSM
	//b.FSM = fsm.NewRuleBasedFSM(clientRules)

}

// ReceiveUpdates handle receiving messages
func (b *Bot) ReceiveUpdates() { // TODO: name properly

	bh, _ := th.NewBotHandler(b.API, b.UpdatesChan)

	bh.Stop()

	//Handling text messages
	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		b.FSM.Context["update"] = update

		userStatus, err := CheckUserDetails(b, update)
		if err == nil {
			if userStatus.State == "/start" {

				err := b.FSM.Trigger("/start")
				if err != nil {
					b.Logger.Error("Failed to trigger event", zap.String("error", err.Error()))
				}
			} else {
				err := b.FSM.Trigger(fsm.Event(update.Message.Text))
				if err != nil {
					b.Logger.Error("Failed to trigger event", zap.String("error", err.Error()))
				}
			}

		} else {
			b.Logger.Error("Failed to handle user cache", zap.String("error", err.Error()))
		}

	}, th.AnyCommand())

	//Handling callback queries
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		b.FSM.Context["message"] = query.Message
		b.FSM.Context["from"] = query.From
		b.FSM.Context["callbackData"] = query.Data

		dataParts := strings.Split(query.Data, ":")
		if len(dataParts) > 1 {
			err := b.FSM.Trigger(fsm.Event(dataParts[0]))
			if err != nil {
				b.Logger.Error("Failed to trigger callback event", zap.String("error", err.Error()))
			}

		} else {
			err := b.FSM.Trigger(fsm.Event(query.Data))
			if err != nil {
				b.Logger.Error("Failed to trigger callback event", zap.String("error", err.Error()))
			}
		}

	}, th.AnyCallbackQuery())

	bh.Start()

}
