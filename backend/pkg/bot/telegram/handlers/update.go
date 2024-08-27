package handlers

//type user struct {
//	ID    int64  `json:"id"`
//	State string `json:"state"`
//}
//
//func (b *Bot) HandleUpdate(update telego.Update) {
//	// Handle the update
//	b.Config.Logger.Info("Received update", zap.Any("update", update))
//
//	//Check if user exists in cache
//	_, err := b.Cache.Get(context.Background(), fmt.Sprintf("id:%d", update.Message.From.ID)).Result()
//	if err != nil {
//		b.Config.Logger.Info("User not found in cache", zap.Int("id", int(update.Message.From.ID)))
//	}
//
//	cachedUser, err := convert.ToJSON(user{
//		ID:    update.Message.From.ID,
//		State: "start",
//	})
//	if err != nil {
//		b.Config.Logger.Error("Failed to convert user to JSON", zap.Int("id", int(update.Message.From.ID)))
//	}
//
//	//TODO: check if user is in DB
//
//	//add user to cache
//	err = b.Cache.Set(context.Background(), fmt.Sprintf("id:%d", update.Message.From.ID), cachedUser, 0).Err()
//	if err != nil {
//		b.Config.Logger.Error("Failed to add user to cache", zap.Int("id", int(update.Message.From.ID)))
//	} else {
//		b.Config.Logger.Info("User added to cache", zap.Int("id", int(update.Message.From.ID)))
//	}
//
//	val, errr := b.Cache.Get(context.Background(), fmt.Sprintf("id:%d", update.Message.From.ID)).Result()
//	if errr != nil {
//		b.Config.Logger.Error("Failed to get user from cache", zap.Int("id", int(update.Message.From.ID)))
//
//	}
//	fmt.Println(val)
//
//}
