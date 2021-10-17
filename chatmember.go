package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func handleMyChatMember(update tgbotapi.ChatMemberUpdated) {
	warn := func(err error) {
		bot.Send(tgbotapi.NewMessage(update.From.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", update.From.LanguageCode)))
		WarnAdmin(err)
		logrus.Error(err)
	}
	user := NewUser(update.Chat.ID, warn)

	switch update.NewChatMember.Status {
	case "member": // юзер разбанил бота
		analytics.User("{bot_was_UNblocked}", &update.From)
		analytics.Bot(update.From.ID, "{bot_was_unblocked}", "bot_UNblocked")
		db.Model(&Users{}).Where("id = ?", update.From.ID).Update("blocked", false)
		bot.Send(tgbotapi.NewMessage(update.From.ID, user.Localize("Welcome. We are glad that you are with us again. ✋")))

		user.WriteBotLog("bot_was_unblocked", "")
	case "kicked":
		analytics.User("{bot_was_blocked}", &update.From)
		analytics.Bot(update.From.ID, "{bot_was_blocked}", "bot_blocked")
		db.Model(&Users{}).Where("id = ?", update.From.ID).Update("blocked", true)

		user.WriteBotLog("bot_was_blocked", "")
	}
}
