package bot

import (
	"errors"
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (app App) onMyChatMember(update tgbotapi.ChatMemberUpdated) {
	user := tables.Users{Lang: update.From.LanguageCode}
	log := app.log.With()
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, user.Localize("Произошла ошибка")))
		app.notifyAdmin(err)
	}

	var err error
	user, err = app.db.GetUserByID(update.From.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			tolang := ""
			if update.From.LanguageCode == "" || update.From.LanguageCode == "en" {
				update.From.LanguageCode = "en"
				tolang = "ru"
			} else if update.From.LanguageCode == "ru" {
				tolang = "en"
			}
			if err = app.db.CreateUser(tables.Users{
				ID:           update.From.ID,
				MyLang:       update.From.LanguageCode,
				ToLang:       tolang,
				Act:          "",
				Usings:       0,
				Blocked:      false,
				LastActivity: time.Now(),
			}); err != nil {
				warn(err)
				return
			}
			user.MyLang = update.From.LanguageCode
			user.ToLang = tolang

		} else {
			warn(err)
		}
	}
	user.SetLang(update.From.LanguageCode)
	log = log.With(zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang))

	switch update.NewChatMember.Status {
	case "member": // юзер разбанил бота
		app.analytics.User("{bot_was_UNblocked}", &update.From)
		// not change user.Blocked since we change it in app.onMessage
	case "kicked":
		app.analytics.User("{bot_was_blocked}", &update.From)
		if err := app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": true}); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}
}
