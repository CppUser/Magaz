package tg

import (
	"Magaz/backend/pkg/utils/parser"
	"Magaz/backend/pkg/utils/service"
	"Magaz/services/bots/telegram/v2/handler"
	"Magaz/services/bots/telegram/v2/handler/message"
	fsm "Magaz/services/bots/telegram/v2/utils/stateMngs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tgConfig struct {
	BotConfig botsConfig `mapstructure:"tg_bot"`
}

type botsConfig struct {
	WebhookLink     string      `mapstructure:"webhook_link"`
	WebhookBasePath string      `mapstructure:"webhook_base_path"`
	Tokens          []BotTokens `mapstructure:"tokens"`
	Groups          []botGroups `mapstructure:"groups"`
}

type BotTokens struct {
	Type  string `mapstructure:"type"`
	Token string `mapstructure:"id"`
}

type botGroups struct {
	Type    string `mapstructure:"type"`
	GroupID int64  `mapstructure:"id"` // Updated to use `int64` for group IDs
}

type BotService struct {
	running   bool
	bots      []*telego.Bot
	updates   map[string]chan telego.Update
	botTypes  map[string]string
	tgConfig  tgConfig
	fsms       map[string]*fsm.RuleBasedFSM
	waitGroup sync.WaitGroup
}

func (b *BotService) Initialize() error {
	///////TODO: Refactor make call to function simply pass path and struct
	var tgcfg tgConfig
	configPaths := []string{
		".",
		"backend/config/",
	}
	err := parser.Load("bot_config", "yaml", configPaths, &tgcfg)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	b.tgConfig = tgcfg

	var router *gin.Engine

	b.updates = make(map[string]chan telego.Update)
	b.botTypes = make(map[string]string)

	for i, token := range tgcfg.BotConfig.Tokens {
		bot, err := telego.NewBot(token.Token)
		if err != nil {
			return fmt.Errorf("failed to create bot %d: %w", i+1, err)
		}

		whPath := tgcfg.BotConfig.WebhookBasePath + "/" + token.Type
		err = bot.SetWebhook(&telego.SetWebhookParams{
			URL: tgcfg.BotConfig.WebhookLink + whPath + "?token=" + token.Token,
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
		if err != nil {
			return fmt.Errorf("failed to set webhook for bot %d: %w", i+1, err)
		}

		updateChan := make(chan telego.Update, 128)
		b.updates[token.Token] = updateChan
		b.botTypes[token.Token] = token.Type

		router.POST(whPath, func(c *gin.Context) {
			token := c.Query("token")
			var update telego.Update

			if err := c.ShouldBindJSON(&update); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update structure"})
				return
			}

			if updateChan, ok := b.updates[token]; ok {
				updateChan <- update
				c.Status(http.StatusOK)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
			}
		})

		b.bots = append(b.bots, bot)
	}

	clientRules := []fsm.Rule{
		{
			Event:      "/start",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

					//message.SendMessageWithMarkup(bot, "Добрый день", map[string]string{
					//	"Вакансии": "hire",
					//}),


					return message.SendMessageWithMarkup(bot, "Желаете оформить заказ?", map[string]string{
						"Оформит": "order",
					})(context)
				},
			},
		},
		{
			Event:      "hire",
			Conditions: []fsm.ConditionFunc{},
			Actions:    []fsm.ActionFunc{message.EditMessage(bot, "Тут будет сообщение о открытых вакансиях")},
		},
		{
			Event:      "order",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

					// Fetch cities from the database
					cities, err := crud.GetAllCities(b.DB)
					if err != nil {
						return err
					}

					// Generate city markup
					cityMarkup := make([]message.CallbackQueryMarkup, len(cities))
					for i, city := range cities {
						cityMarkup[i] = message.CallbackQueryMarkup{
							Text:         city.Name,
							CallbackData: "city:" + city.Name, // Pass city name in the callback data
						}
					}

					// Prompt the user to select a city
					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите ваш город:", cityMarkup)(context)
				},
			},
		},
		{
			Event:      "city",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

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
					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующий вас товар:", productMarkup)(context)
				},
			},
		},
		{
			Event:      "product",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

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
					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующее вас количество:", quantityMarkup)(context)
				},
			},
		},
		//TODO: add next rule to add to cart to shop for more items
		{
			Event:      "quantity",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

					callbackData := context["callbackData"].(string)
					qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "quantity", qtAmount)

					choices, err := GetUserChoices(b.Cache, user.ID)
					if err != nil {
						return err
					}
					city, _ := crud.GetCityIDByName(b.DB, choices["city"])
					prID, _ := crud.GetProductIDByCityAndProductName(b.DB, choices["city"], choices["product"])
					qt, _ := strconv.ParseFloat(choices["quantity"], 32)

					// Use the GetAvailableAddresses function to check if there are any available addresses
					availableAddresses, err := crud.GetAvailableAddresses(b.DB, city, prID, float32(qt))
					if err != nil || len(availableAddresses) == 0 {
						// If no addresses available, notify the user
						return message.EditMessage(bot, "К сожалению, нет доступных адресов для выбранного количества. Пожалуйста, начните сначала с команды /start.")(context)
					}

					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите метод оплаты:", []message.CallbackQueryMarkup{ //TODO: replace with generated markup
						{Text: "Перевод на карту", CallbackData: "card"},
						//{Text: "Оплата Крипто валютой", CallbackData: "crypto"},
					})(context)
				},
			},
		},
		{
			Event:      "card",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				bot := context["bot"].(*telego.Bot)

				func(context map[string]interface{}) error {
					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "payment", callbackData)

					//TODO: Generate payment method markup

					return message.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты перевод на карту", []message.CallbackQueryMarkup{ //TODO: replace with generated markup
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
					bot := context["bot"].(*telego.Bot)

					// Extract the product name from the context
					//callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					//TODO: Generate payment method markup

					return message.EditMessageWithMarkup(bot, "Выберите тип крипто валюты", []message.CallbackQueryMarkup{ //TODO: replace with generated markup
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
					bot := context["bot"].(*telego.Bot)

					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "payment", callbackData)

					//TODO: Generate payment method markup

					return message.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты bitcoin крипто валютой", []message.CallbackQueryMarkup{ //TODO: replace with generated markup
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
					bot := context["bot"].(*telego.Bot)

					// Extract the product name from the context
					callbackData := context["callbackData"].(string)
					//qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "payment", callbackData)

					//TODO: Generate payment method markup

					return message.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты etherium крипто валютой", []message.CallbackQueryMarkup{ //TODO: replace with generated markup
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
					bot := context["bot"].(*telego.Bot)

					user := context["from"].(telego.User)
					message := context["message"].(*telego.Message)

					//TODO: check if order is in process to avoid flooding
					choices, err := GetUserChoices(b.Cache, user.ID)
					if err != nil {
						return err
					}
					city, _ := crud.GetCityIDByName(b.DB, choices["city"])
					prID, _ := crud.GetProductIDByCityAndProductName(b.DB, choices["city"], choices["product"])
					qt, _ := strconv.ParseFloat(choices["quantity"], 32)

					// Check if the user already has 2 or more active (non-released) orders
					var activeOrders []models.Order
					if err := b.DB.Where("user_id = ? AND released = ?", user.ID, false).Find(&activeOrders).Error; err != nil {
						return fmt.Errorf("failed to fetch active orders: %w", err)
					}

					// If the user already has 2 active (non-released) orders, prevent them from creating more
					if len(activeOrders) >= 2 {
						//TODO: send to operator chat instead right now since there is no option , using old method
						msg := "У вас уже есть 2 активных заказа. Пожалуйста, завершите один из них, прежде чем создавать новый.\n\n" +
							"Если у вас возникли проблемы с заказом свяжитесь с оператором старым методом"

						// Add the details of the active orders to the message
						msg += "Ваши активные заказы:\n"
						for _, order := range activeOrders {
							// Since city and product are not preloaded, use the choices retrieved earlier
							msg += fmt.Sprintf(
								"Заказ #%d\n"+
									"Город: %s\n"+
									"Товар: %s\n"+
									"Количество: %.2f\n"+
									"Сумма к оплате: %d\n\n",
								order.ID,
								choices["city"],
								choices["product"],
								order.Quantity,
								order.Due,
							)
						}

						return message.EditMessage(bot, msg)(context)
					}

					var qtnPrice models.QtnPrice //TODO: Refactor
					// Find the price for the given quantity
					if err := b.DB.Where("city_product_id = ? AND quantity = ?", prID, qt).First(&qtnPrice).Error; err != nil { //passing wrong city id
						b.Logger.Error("price not found for the specified quantity", zap.String("error", err.Error()))
					}

					pmt, _ := crud.GetPaymentMethod(b.DB, choices["payment"])
					address, _ := repository.GetRandomAddress(b.DB, city, prID, float32(qt), user.ID)

					//TODO: Need to figure out how to store custom quantity

					var msg string
					order := models.Order{}
					ordView := repository.OrderView{}

					if choices["payment"] == "card" {

						order = models.Order{
							UserID:            user.ID,
							CityID:            city,
							ProductID:         prID,        //Get product from productCityID
							Quantity:          float32(qt), //TODO: retrieve quantity from cache
							Due:               uint(qtnPrice.Price),
							PaymentMethodType: choices["payment"],
							PaymentMethodID:   pmt.(models.Card).ID, //TODO: retrieve from PaymentMethodType name if card from card if crypto from crypto
							CreatedAt:         time.Now(),
							ReleasedAddrID:    &address.ID,
						}

						var addr models.Address
						if err := b.DB.First(&addr, &address.ID).Error; err == nil {
							if !addr.Assigned {
								addr.Assigned = true
								addr.AssignedUserID = &user.ID
								//TODO add AssignedBy (Bot)
							}
						}

						if err := b.DB.Create(&order).Error; err != nil {
							b.Logger.Error("Failed to create new order in DB", zap.String("error", err.Error()))
						}

						ordView = repository.OrderView{
							ID:          order.ID,
							ProductName: choices["product"],
							CityName:    choices["city"],
							Quantity:    float32(qt),
							Due:         uint(qtnPrice.Price), //TODO: once card implemented need to add all items in cart to due
							CreatedAt:   time.Now(),
							Client: repository.UserView{
								ID:        user.ID,
								ChatID:    message.GetChat().ID,
								Username:  user.Username,
								FirstName: user.FirstName,
								LastName:  user.LastName,
							},
							PaymentMethod: repository.PaymentView{
								PaymentCategory: "Перевод на карту",
								CardPayment: repository.CardView{
									BankName:   pmt.(models.Card).BankName,
									BankUrl:    pmt.(models.Card).BankURL,
									CardNumber: pmt.(models.Card).CardNumber,
									FirstName:  pmt.(models.Card).FirstName,
									LastName:   pmt.(models.Card).LastName,
									UserName:   pmt.(models.Card).UserID,
									Password:   pmt.(models.Card).Password,
									QuickPay:   pmt.(models.Card).QuickPay,
								},
							},
							Address: *address,
						}

						msg = fmt.Sprintf(
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
								"СБП: %s\n"+
								"*************************\n",
							order.ID,
							choices["city"],
							choices["product"],
							choices["quantity"],
							"Перевод на Карту",
							uint(qtnPrice.Price),
							pmt.(models.Card).BankName,
							pmt.(models.Card).CardNumber,
							pmt.(models.Card).LastName+" "+pmt.(models.Card).FirstName,
							pmt.(models.Card).QuickPay,
						)
					} else if choices["payment"] == "crypto" {
						//TODO: Implement
					}

					b.WS.BroadcastOrder(ordView)
					//TODO: Send message to Employee telegram about new order (to personal or to group chat)"4512552536"

					empMessage := fmt.Sprintf("Добавлен новый заказ: %d\n", ordView.ID)

					_, _ = bot.SendMessage(&telego.SendMessageParams{
						ChatID: tu.ID(b.Config.GroupID),
						Text:   empMessage,
					})

					go func(orderID uint, addrID *uint) {
						//TODO: Does not work when there is constant ordering , need to find other way to put timer on order
						//Maybe some database watching system
						<-time.After(15 * time.Minute)

						order.ReleasedAddrID = nil
						if err := b.DB.Save(&order).Error; err != nil {
							b.Logger.Error("Failed to update order status", zap.String("error", err.Error()))
						}

						if addrID != nil {
							var addr models.Address
							if err := b.DB.First(&addr, *addrID).Error; err == nil {
								addr.Released = false
								addr.Assigned = false
								addr.AssignedUserID = nil
								if err := b.DB.Save(&addr).Error; err != nil {
									b.Logger.Error("Failed to unassign address", zap.String("error", err.Error()))
								}
							}
						}
						ordView.Address = repository.AddressView{}
						b.WS.BroadcastOrder(ordView)

					}(order.ID, order.ReleasedAddrID)

					return message.EditMessageWithMarkup(bot, msg, []message.CallbackQueryMarkup{
						{Text: "Подтеврдить оплату", CallbackData: "payConf"},
					})(context)
				},
			},
		},
		{
			Event:      "cancel",
			Conditions: []fsm.ConditionFunc{},
			Actions:    []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)
					return message.EditMessage(bot, "Ваш заказ успешно отменен")(context)
				},
			},
		},
		{
			Event:      "payConf",
			Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

					msg := fmt.Sprintf(
						"Пожалуйста ожидайте ответа оператора\n" +
							"Функция подтвержденя оплаты в данный момент не работает\n" +
							"Если у вас имеется чек на руках отправьте его старым методом\n также укажите номер заказа")

					return message.EditMessage(bot, msg)(context)
				},
			},
		},
		{
			Event:      "status",
			Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
			Actions:    []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)
					return message.EditMessage(bot, "Сообщение о статусе заказа")(context)
				},
			},
		},
	}

	emplRules := []fsm.Rule{
		{
			Event:      "city",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

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
					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующий вас товар:", productMarkup)(context)
				},
			},
		},
		{
			Event:      "product",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

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
					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующее вас количество:", quantityMarkup)(context)
				},
			},
		},
		//TODO: add next rule to add to cart to shop for more items
		{
			Event:      "quantity",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				func(context map[string]interface{}) error {
					bot := context["bot"].(*telego.Bot)

					callbackData := context["callbackData"].(string)
					qtAmount := strings.TrimPrefix(callbackData, "quantity:")

					user := context["from"].(telego.User)
					StoreUserChoice(b.Cache, user.ID, "quantity", qtAmount)

					choices, err := GetUserChoices(b.Cache, user.ID)
					if err != nil {
						return err
					}
					city, _ := crud.GetCityIDByName(b.DB, choices["city"])
					prID, _ := crud.GetProductIDByCityAndProductName(b.DB, choices["city"], choices["product"])
					qt, _ := strconv.ParseFloat(choices["quantity"], 32)

					// Use the GetAvailableAddresses function to check if there are any available addresses
					availableAddresses, err := crud.GetAvailableAddresses(b.DB, city, prID, float32(qt))
					if err != nil || len(availableAddresses) == 0 {
						// If no addresses available, notify the user
						return message.EditMessage(bot, "К сожалению, нет доступных адресов для выбранного количества. Пожалуйста, начните сначала с команды /start.")(context)
					}

					return message.EditMessageWithMarkup(bot, "Пожалуйста выберите метод оплаты:", []message.CallbackQueryMarkup{ //TODO: replace with generated markup
						{Text: "Перевод на карту", CallbackData: "card"},
						//{Text: "Оплата Крипто валютой", CallbackData: "crypto"},
					})(context)
				},
			},
		},
	}

	b.fsms = map[string]*fsm.RuleBasedFSM{
		"client": fsm.NewRuleBasedFSM(clientRules),
		"empl":   fsm.NewRuleBasedFSM(emplRules),
	}



	return nil
}

func (b *BotService) Start() error {
	if b.running {
		return fmt.Errorf("service is already running")
	}

	b.running = true

	for token, updateChan := range b.updates {
		b.waitGroup.Add(1)
		go func(token string, updateChan chan telego.Update) {
			defer b.waitGroup.Done()

			botType := b.botTypes[token]

			var bot *telego.Bot
			for _, b := range b.bots {
				if b.Token() == token {
					bot = b
					break
				}
			}
			if bot == nil {
				fmt.Printf("Bot with token %s not found\n", token)
				return
			}

			bh, err := th.NewBotHandler(bot, updateChan)
			if err != nil {
				fmt.Printf("Failed to create bot handler for bot %s: %v\n", token, err)
				return
			}
			bh.Stop()

			bh.Handle(func(bot *telego.Bot, update telego.Update) {
				switch botType {
				case "client":
					handler.HandleClientUpdate(b.fsms["client"],bot, update)
				case "empl":
					handler.HandleEmplUpdate(b.fsms["empl"],bot, update)
				default:
					fmt.Printf("Unknown bot type for token %s\n", token)
				}
			}, th.AnyCommand())

			bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
				switch botType {
				case "client":
					handler.HandleClientCallbackQuery(b.fsms["client"],bot, query)
				case "empl":
					handler.HandleEmplCallbackQuery(b.fsms["empl"],bot, query)
				default:
					fmt.Printf("Unknown bot type for token %s\n", token)
				}
			}, th.AnyCallbackQuery())
		}(token, updateChan)
	}

	return nil
}

func (b *BotService) Stop() error {
	if !b.running {
		return fmt.Errorf("service is not running")
	}

	b.running = false

	for i, bot := range b.bots {
		log.Printf("Stopping webhook for bot %d with token: %s\n", i+1, bot.Token())
		if err := bot.StopWebhook(); err != nil {
			return fmt.Errorf("failed to stop webhook for bot %d: %w", i+1, err)
		}
	}

	for _, updateChan := range b.updates {
		close(updateChan)
	}

	b.waitGroup.Wait()

	return nil
}

func (b *BotService) Status() string {
	if b.running {
		return "Running"
	}
	return "Stopped"
}

func NewBotService() service.Service {
	return &BotService{}
}
