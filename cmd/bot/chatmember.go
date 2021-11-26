package bot

import (
	"github.com/armanokka/translobot/internal/tables"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (app app) onMyChatMember(update tgbotapi.ChatMemberUpdated) {
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", update.From.LanguageCode)))
		app.notifyAdmin(err)
		logrus.Error(err)
	}
	user := app.loadUser(update.Chat.ID, warn)
	defer user.UpdateLastActivity()
	switch update.NewChatMember.Status {
	case "member": // юзер разбанил бота
		app.analytics.User("{bot_was_UNblocked}", &update.From)
		app.analytics.Bot(update.From.ID, "{bot_was_unblocked}", "bot_UNblocked")
		app.db.Model(&tables.Users{}).Where("id = ?", update.From.ID).Update("blocked", false)
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, user.Localize("Welcome. We are glad that you are with us again. ✋")))

		app.writeBotLog(update.From.ID, "bot_was_unblocked", "")
	case "kicked":
		app.analytics.User("{bot_was_blocked}", &update.From)
		app.analytics.Bot(update.From.ID, "{bot_was_blocked}", "bot_blocked")
		app.db.Model(&tables.Users{}).Where("id = ?", update.From.ID).Update("blocked", true)

		app.writeBotLog(update.From.ID, "bot_was_blocked", "")
	}
}
