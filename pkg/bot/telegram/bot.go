package telegram

import (
	"Magaz/internal/config"
	tgconfig "Magaz/pkg/bot/telegram/config"
	"Magaz/pkg/bot/telegram/handlers"
	"Magaz/pkg/utils/state/fsm"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	_, _ = tgconfig.LoadConfig("bot_config", "yaml", []string{".", "config/"})
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
	cityMarkup, err := FetchCitiesFromDB(b.DB)
	rules := []fsm.Rule{
		{
			Event:      "/start",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{
				handlers.SendMessageWithMarkup(bot, "Добрый день", map[string]string{
					"Вакансии": "hire",
				}),
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
			Actions: []fsm.ActionFunc{handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите ваш город:",
				cityMarkup),
			},
		},
		{
			Event:      "city",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующий вас товар:",
				[]handlers.TempMarkup{
					{Text: "Товар 1", CallbackData: "product"},
					{Text: "Товар 2", CallbackData: "product"},
				}),
			},
		},
		{
			Event:      "product",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите интересующее вас количество:",
				[]handlers.TempMarkup{
					{Text: "Количество 1", CallbackData: "quantity"},
					{Text: "Количество 2", CallbackData: "quantity"},
				}),
			},
		},
		//TODO: after choosing quantity, offer region where delivery is available
		//TODO: add next rule to add to cart to shop for more items
		{
			Event:      "quantity",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{handlers.EditMessageWithMarkup(bot, "Пожалуйста выберите метод оплаты:",
				[]handlers.TempMarkup{
					{Text: "Перевод на карту", CallbackData: "card"},
					{Text: "Оплата Крипто валютой", CallbackData: "crypto"},
				}),
			},
		},
		{
			Event:      "card",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты перевод на карту:",
				[]handlers.TempMarkup{
					{Text: "Оформить заказ", CallbackData: "confirm"},
					{Text: "Отменить заказ", CallbackData: "cancel"},
				}),
			},
		},
		{
			Event:      "crypto",
			Conditions: []fsm.ConditionFunc{},
			Actions: []fsm.ActionFunc{handlers.EditMessageWithMarkup(bot, "Вы выбрали метод оплаты крипто валютой:",
				[]handlers.TempMarkup{
					{Text: "Оформить заказ", CallbackData: "confirm"},
					{Text: "Отменить заказ", CallbackData: "cancel"},
				}),
			},
		},
		//TODO: send order details to costumer
		{
			Event:      "confirm",
			Conditions: []fsm.ConditionFunc{},
			Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Ваш заказ успешно оформлен")},
		},
		{
			Event:      "cancel",
			Conditions: []fsm.ConditionFunc{},
			Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Ваш заказ успешно отменен")},
		},
		{
			Event:      "status",
			Conditions: []fsm.ConditionFunc{}, //TODO: wait response from operator to confirm payment
			Actions:    []fsm.ActionFunc{handlers.EditMessage(bot, "Тут будет сообщение о статусе заказа")},
		},
	}
	// Initialize FSM
	b.FSM = fsm.NewRuleBasedFSM(rules)

}

// TODO: name properly
func (b *Bot) ReceiveUpdates() {

	bh, _ := th.NewBotHandler(b.API, b.UpdatesChan)

	bh.Stop()

	//Handling text messages
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		b.FSM.Context["update"] = update

		err := b.FSM.Trigger(fsm.Event(update.Message.Text))
		if err != nil {
			b.Logger.Error("Failed to trigger event", zap.String("error", err.Error()))
		}

	}, th.AnyCommand())

	//Handling callback queries
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		b.FSM.Context["message"] = query.Message
		err := b.FSM.Trigger(fsm.Event(query.Data))
		if err != nil {
			b.Logger.Error("Failed to trigger callback event", zap.String("error", err.Error()))
		}

	}, th.AnyCallbackQuery())

	bh.Start()

}
