package telegram

import (
	"Magaz/internal/config"
	tgconfig "Magaz/pkg/bot/telegram/config"
	"Magaz/pkg/bot/telegram/handlers"
	"Magaz/pkg/utils/state/fsm"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Bot struct {
	Config           *config.TGBotConfig
	UpdateChanBuffer uint
	UpdatesChan      chan telego.Update
	Cache            *redis.Client
	FSM              *fsm.FSM
}

// TODO: refactor code move some logic to handlers
// InitBot initializes the Telegram bot
func (b *Bot) InitBot() {

	//TODO: need to handle error properly, currently removed do to code complaint
	cfg, _ := tgconfig.LoadConfig("bot_config", "yaml", []string{".", "config/"})
	//if err != nil {
	//	b.Config.Logger.Fatal("Failed to load bot configs", zap.String("error", err.Error()))
	//}

	////////////////////////////////////////////////////////////////////////////////////////////////
	//TODO: move to separate function , possibly add to bot struct
	sm := fsm.NewFSM(fsm.State(cfg.InitialState))
	for state, stateConfig := range cfg.States {
		sm.AddState(fsm.State(state))

		for event, nextState := range stateConfig.Transitions {
			sm.AddTransition(fsm.State(state), fsm.Event{Name: event}, fsm.State(nextState))
		}

		switch stateConfig.Handler {
		case "startHandler":
			sm.AddStateHandler(fsm.State(state), handlers.StartHandler)
		case "cityHandler":
			sm.AddStateHandler(fsm.State(state), handlers.CityHandler)
		case "productHandler":
			sm.AddStateHandler(fsm.State(state), handlers.ProductHandler)
		case "quantityHandler":
			sm.AddStateHandler(fsm.State(state), handlers.QuantityHandler)
		case "paymentHandler":
			sm.AddStateHandler(fsm.State(state), handlers.PaymentHandler)
		case "endHandler":
			sm.AddStateHandler(fsm.State(state), handlers.EndHandler)

		}
	}
	b.FSM = sm
	////////////////////////////////////////////////////////////////////////////////////////////////

	bot, err := telego.NewBot(b.Config.Token)
	if err != nil {
		//TODO: need to handle the error differently without direct call to zap.String
		b.Config.Logger.Fatal("failed create new bot api instance", zap.String("error", err.Error()))
	}
	b.Config.API = bot

	//TODO: find better way to make channel
	//TODO: also need to safly send updates to channel, checking is channel is open or closed
	// Initialize the updates channel
	b.UpdatesChan = make(chan telego.Update, b.UpdateChanBuffer)

	//TODO: refer to SetWebhookParams to setup additional parameters (like certificate, pending updates, etc.)
	_ = b.Config.API.SetWebhook(&telego.SetWebhookParams{
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

	info, _ := b.Config.API.GetWebhookInfo()
	b.Config.Logger.Info("Webhook Info", zap.Any("info", info)) //TODO: in prod it needs to be in JSON format

}

// TODO: name properly
func (b *Bot) ReceiveUpdates() {

	bh, _ := th.NewBotHandler(b.Config.API, b.UpdatesChan)

	bh.Stop()

	// Register new handler with match on command `/start`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		event := fsm.Event{
			Name: "start",
			Payload: handlers.Payload{
				Bot:    bot,
				Update: update,
			},
		}

		if err := handlers.StartHandler(event, b.FSM); err != nil {
			b.Config.Logger.Error("Failed to handle start event", zap.String("error", err.Error()))
		}

	}, th.CommandEqual("start"))

	// Register new handler with match on a call back query with data equal to `go` and non-nil message
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		// Send message
		_, _ = bot.SendMessage(tu.Message(tu.ID(query.Message.GetChat().ID), "GO"))

		// Answer callback query
		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Done"))
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("go"))

	bh.Start()

	//// Define a context key
	//type userID bool
	//var userIDKey userID
	//
	//// Apply middleware that will retrieve user ID from update
	//bh.Use(func(bot *telego.Bot, update telego.Update, next th.Handler) {
	//	// Get initial context
	//	ctx := update.Context()
	//
	//	if update.Message != nil && update.Message.From != nil {
	//		// Set user ID in context
	//		ctx = context.WithValue(ctx, userIDKey, update.Message.From.ID)
	//	}
	//
	//	// Update context
	//	update = update.WithContext(ctx)
	//	next(bot, update)
	//})
	//
	//// Handle messages
	//bh.Handle(func(bot *telego.Bot, update telego.Update) {
	//	ctx := update.Context()
	//
	//	// Retrieve user ID from context
	//	fmt.Println("User ID:", ctx.Value(userIDKey))
	//}, th.AnyMessage())
}
