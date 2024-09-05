package telegram

import (
	"Magaz/backend/internal/config"
	"Magaz/backend/internal/storage/crud"
	"Magaz/backend/internal/storage/models"
	tgconfig "Magaz/backend/pkg/bot/telegram/config"
	"Magaz/backend/pkg/bot/telegram/handlers"
	"Magaz/backend/pkg/utils/state/fsm"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// TODO: Remove load from config .
const (
	StateStart   fsm.State = "start"
	StateCity    fsm.State = "city"
	StateProduct fsm.State = "product"
)

type Bot struct {
	Config           *config.TGBotConfig
	API              *telego.Bot
	Logger           *zap.Logger
	UpdateChanBuffer uint
	UpdatesChan      chan telego.Update
	Cache            *redis.Client
	DB               *gorm.DB
	FSM              *fsm.RuleBasedFSM
}

// TODO: refactor code move some logic to handlers

// InitBot initializes the Telegram bot
func (b *Bot) InitBot() {

	//TODO: need to handle error properly, currently removed do to code complaint
	_, _ = tgconfig.LoadConfig("bot_config", "yaml", []string{".", "backend/config/"})
	//if err != nil {
	//	b.Config.Logger.Fatal("Failed to load bot configs", zap.String("error", err.Error()))
	//}

	////////////////////////////////////////////////////////////////////////////////////////////////
	//TODO: handle fsm creation in api .i.e. for dynamic conv creation. Possibly generate handlers with go generate instruction

	//sm := fsm.NewFSM(fsm.State(smcfg.States[0].Name))
	//for _, state := range smcfg.States {
	//	for _, transition := range state.Transitions {
	//		sm.AddTransition(fsm.State(state.Name), fsm.Event(transition.Event), fsm.State(transition.To))
	//	}
	//
	//}
	//
	//// Set up handlers
	//for _, handler := range smcfg.Handlers {
	//	switch handler.Handler {
	//	case "StartHandler":
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.StartHandler)
	//	case "OrderHandler":
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.OrderHandler)
	//	case "CityHandler":
	//
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.CityHandler)
	//	case "ProductHandler":
	//
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.ProductHandler)
	//	case "QuantityHandler":
	//
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.QuantityHandler)
	//	case "PaymentHandler":
	//
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.PaymentHandler)
	//	case "ConformationHandler":
	//
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.ConformationHandler)
	//	case "FinalHandler":
	//
	//		sm.AddHandler(fsm.Event(handler.Event), handlers.FinalHandler)
	//	default:
	//
	//		log.Fatalf("Unknown handler: %s", handler.Handler)
	//	}
	//}
	//
	////sm.AddTransition(StateStart, "start", StateCity)
	////sm.AddTransition(StateCity, "city", StateProduct)
	////sm.AddTransition(StateProduct, "product", StateStart)
	////
	////sm.AddHandler("start", handlers.StartHandler)
	////sm.AddHandler("city", handlers.CityHandler)
	////sm.AddHandler("product", handlers.ProductHandler)
	//b.FSM = sm

	////////////////////////////////////////////////////////////////////////////////////////////////

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
	clientRules := []fsm.Rule{
		{
			Event:      "/start",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				handlers.SendMessageWithMarkup(bot, "Добрый день", map[string]string{
					"Вакансии": "hire",
				}),

				//check if user exist in cache

				handlers.SendMessageWithMarkup(bot, "Желаете оформить заказ?", map[string]string{
					"Оформит": "order",
				}),
			},
		},
		{
			Event:      "hire",
			Conditions: []fsm.ConditionFunc{},
			Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Тут будет сообщение о открытых вакансиях")},
		},
		{
			Event:      "order",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {

					// Fetch cities from the database
					cities, err := crud.GetAllCities(b.DB)
					if err != nil {
						return err
					}

					// Generate city markup
					cityMarkup := make([]handlers.TempMarkup, len(cities))
					for i, city := range cities {
						cityMarkup[i] = handlers.TempMarkup{
							Text:         city.Name,
							CallbackData: "city:" + city.Name, // Pass city name in the callback data
						}
					}

					// Prompt the user to select a city
					return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите ваш город:", cityMarkup)(context)
				},
			},
		},
		{
			Event:      "city",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Extract the city name from the callback data
					callbackData := context["callbackData"].(string)
					cityName := strings.TrimPrefix(callbackData, "city:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "city", cityName)

					// Generate product markup for the selected city
					productMarkup, err := GenerateProductMarkup(b.DB, cityName)
					if err != nil {
						return err
					}

					// Prompt the user to select a product
					return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующий вас товар:", productMarkup)(context)
				},
			},
		},
		{
			Event:      "product",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					productName := strings.TrimPrefix(callbackData, "product:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "product", productName)

					// Generate quantity markup for the selected product
					quantityMarkup, err := GenerateProductPriceMarkup(b.DB, productName)
					if err != nil {
						return err
					}

					// Prompt the user to select a quantity
					return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующее вас количество:", quantityMarkup)(context)
				},
			},
		},
		//TODO: after choosing quantity, offer region where delivery is available
		//TODO: add next rule to add to cart to shop for more items
		{
			Event:      "quantity",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					callbackData := context["callbackData"].(string)
					qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "quantity", qtAmount)

					//TODO: Generate payment method markup

					return handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите метод оплаты:", []handlers.TempMarkup{ //TODO: replace with generated markup
						{Text: "Перевод на карту", CallbackData: "card"},
						{Text: "Оплата Крипто валютой", CallbackData: "crypto"},
					})(context)
				},
			},
		},
		{
			Event:      "card",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "payment", callbackData)

					//TODO: Generate payment method markup

					return handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты перевод на карту", []handlers.TempMarkup{ //TODO: replace with generated markup
						{Text: "Оформить заказ", CallbackData: "confirm"}, //TODO: Add option to add to cart for multiple products
						{Text: "Отменить заказ", CallbackData: "cancel"},
					})(context)
				},
			},
		},
		{
			Event:      "crypto",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Extract the product name from the context
					//callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					//TODO: Generate payment method markup

					return handlers.EditMessageWithMarkup(bot, "Выберите тип крипто валюты", []handlers.TempMarkup{ //TODO: replace with generated markup
						{Text: "Bitcoin", CallbackData: "bitcoin"}, //TODO: Add option to add to cart for multiple products
						{Text: "Etherium", CallbackData: "etherium"},
					})(context)
				},
			},
		},

		{
			Event:      "bitcoin",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "payment", callbackData)

					//TODO: Generate payment method markup

					return handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты bitcoin крипто валютой", []handlers.TempMarkup{ //TODO: replace with generated markup
						{Text: "Оформить заказ", CallbackData: "confirm"}, //TODO: Add option to add to cart for multiple products
						{Text: "Отменить заказ", CallbackData: "cancel"},
					})(context)
				},
			},
		},
		{
			Event:      "etherium",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "payment", callbackData)

					//TODO: Generate payment method markup

					return handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты etherium крипто валютой", []handlers.TempMarkup{ //TODO: replace with generated markup
						{Text: "Оформить заказ", CallbackData: "confirm"}, //TODO: Add option to add to cart for multiple products
						{Text: "Отменить заказ", CallbackData: "cancel"},
					})(context)
				},
			},
		},

		{
			Event:      "confirm",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					user := context["from"].(telego.User)

					choices, err := GetUserChoices(b.Cache, user.ID)
					if err != nil {
						return err
					}
					b.Logger.Info("User choices", zap.Any("choices", choices))

					city, _ := crud.GetCityIDByName(b.DB, choices["city"])
					prID, _ := crud.GetProductIDByCityAndProductName(b.DB, choices["city"], choices["product"])

					qt, _ := strconv.ParseFloat(choices["quantity"], 32)

					var qtnPrice models.QtnPrice
					// Find the price for the given quantity
					if err := b.DB.Where("city_product_id = ? AND quantity = ?", prID, qt).First(&qtnPrice).Error; err != nil { //passing wrong city id
						b.Logger.Error("price not found for the specified quantity", zap.String("error", err.Error()))
					}

					pmt, _ := crud.GetPaymentMethod(b.DB, choices["payment"])

					var message string
					if choices["payment"] == "card" {

						if err := b.DB.Create(&models.Order{
							UserID:            user.ID,
							CityID:            city,
							ProductID:         prID,        //Get product from productCityID
							Quantity:          float32(qt), //TODO: retrieve quantity from cache
							Due:               uint(qtnPrice.Price),
							PaymentMethodType: choices["payment"],
							PaymentMethodID:   pmt.(models.Card).ID, //TODO: retrieve from PaymentMethodType name if card from card if crypto from crypto
							CreatedAt:         time.Now(),
						}).Error; err != nil {
							b.Logger.Error("Failed to create new order in DB", zap.String("error", err.Error()))
						}

						var recentOrder models.Order
						if err := b.DB.Where("user_id = ?", user.ID).Order("created_at DESC").First(&recentOrder).Error; err != nil {
							if errors.Is(err, gorm.ErrRecordNotFound) {
								//return 0, fmt.Errorf("no orders found for user ID: %d", userID)
								b.Logger.Info("no orders found")
							}
							b.Logger.Info("failed to retrieve latest order for user")
						}

						message = fmt.Sprintf(
							"Номер заказа #%d\n"+
								"Город: %s\n"+
								"Товар: %s\n"+
								"Количество: %s\n"+
								"Метод оплаты: %s\n"+
								"Сумма к оплате: %d\n"+
								"\n"+
								"***** Данные Карты *****\n"+
								"Банк: %s\n"+
								"Номер карты: %s\n"+
								"ФИО: %s\n"+
								"*************************\n",
							recentOrder.ID,
							choices["city"],
							choices["product"],
							choices["quantity"],
							"Перевод на Карту",
							uint(qtnPrice.Price),
							pmt.(models.Card).BankName,
							pmt.(models.Card).CardNumber,
							pmt.(models.Card).LastName+" "+pmt.(models.Card).FirstName,
						)
					} else if choices["payment"] == "crypto" {
						//TODO: Implement
					}

					// Prompt the user to select a quantity

					return handlers.EditMessageWithMarkup(bot, message, []handlers.TempMarkup{
						{Text: "Подтеврдить оплату", CallbackData: "payConf"},
					})(context)
				},
			},
		},
		{
			Event:      "cancel",
			Conditions: []fsm.ConditionFunc{},
			Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Ваш заказ успешно отменен")},
		},
		{
			Event:      "payConf",
			Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {

					message := fmt.Sprintf(
						"Пожалуйста ожидайте ответа оператора\n" +
							"Для более быстрой обработки вашего заказа вы можете прикрепите фотографию с оплатой")

					return handlers.EditMessage(bot, message)(context)
				},
			},
		},
		{
			Event:      "status",
			Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
			Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Тут будет сообщение о статусе заказа")},
		},
	}
	// Initialize FSM
	b.FSM = fsm.NewRuleBasedFSM(clientRules)

}

// TODO: name properly
func (b *Bot) ReceiveUpdates() {

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
