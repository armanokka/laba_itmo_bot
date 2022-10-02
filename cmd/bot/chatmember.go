package bot

import (
	"github.com/armanokka/translobot/pkg/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (app App) onMyChatMember(update tgbotapi.ChatMemberUpdated) {

	user, err := app.db.GetUserByID(update.From.ID)
	isNew := false
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.log.Error("", zap.Error(err))
			return
		}
		isNew = true
	}
	user.SetLang(update.From.LanguageCode)
	switch update.NewChatMember.Status {
	case "member":
		if isNew {
			app.bot.Send(tgbotapi.NewSticker(update.From.ID, tgbotapi.FileID(`CAACAgIAAxkBAAESRh1jObyhH-hcOotE2u08d_mARIZYqwACJQADspiaDg82d3yOk8EtKgQ`)))
			return
		}
		app.bot.Send(tgbotapi.NewSticker(update.From.ID, tgbotapi.FileID(`CAACAgIAAxkBAAESRi1jOb9UZF2V6FqZWt05EJap4JHdMgACKwcAAkb7rAR5-qgf7bN-0CoE`)))
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, user.Localize("Glad to see you again, %s", update.From.FirstName)))
	case "kicked":
		if err := app.analytics.User(tgbotapi.Message{
			MessageID:  0,
			From:       &update.From,
			SenderChat: nil,
			Date:       0,
			Chat: &tgbotapi.Chat{
				ID:        update.From.ID,
				Type:      "private",
				UserName:  update.From.UserName,
				FirstName: update.From.FirstName,
				LastName:  update.From.LastName,
			},
			Text: ":bot_was_blocked",
		}); err != nil {
			app.notifyAdmin(err)
		}
		if err = app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": true}); err != nil {
			app.notifyAdmin(err)
		}
	}
}
