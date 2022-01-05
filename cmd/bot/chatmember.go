package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
)

func (app App) onMyChatMember(update tgbotapi.ChatMemberUpdated) {
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", update.From.LanguageCode)))
		app.notifyAdmin(err)
		logrus.Error(err)
	}
	localizer := i18n.NewLocalizer(app.bundle, update.From.LanguageCode)
	defer func() {
		if err := app.db.UpdateUserLastActivity(update.From.ID); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	switch update.NewChatMember.Status {
	case "member": // юзер разбанил бота
		app.analytics.User("{bot_was_UNblocked}", &update.From)
		if err := app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": false}); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
		locale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Welcome. We are glad that you are with us again. ✋"})
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, locale))
	case "kicked":
		app.analytics.User("{bot_was_blocked}", &update.From)
		if err := app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": true}); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}
}
