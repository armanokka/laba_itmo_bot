package main

//func handleVoice(message tgbotapi.Message) {
//	warn := func(err error) {
//		bot.Send(tgbotapi.NewMessage(message.Chat.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", message.From.LanguageCode)))
//		WarnAdmin(err)
//		logrus.Error(err)
//	}
//	user := NewUser(message.From.ID, warn)
//
//	analytics.User("voice", message.From)
//
//	//url, err := bot.GetFileDirectURL(message.Voice.FileID)
//	//if err != nil {
//	//	warn(err)
//	//}
//	//
//	//res, err := http.Get(url)
//	//if err != nil {
//	//	warn(err)
//	//}
//
//	bot.Send(tgbotapi.MessageConfig{
//		BaseChat:              tgbotapi.BaseChat{
//			ChatID:                   message.From.ID,
//			ReplyToMessageID:         message.MessageID,
//			ReplyMarkup:              BuildSupportedLanguagesKeyboard("transcribe"),
//		},
//		Text: user.Localize(),
//		// TODO: Handle transcribe:lang
//	})
//}
