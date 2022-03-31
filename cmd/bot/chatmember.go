package bot

import (
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app App) onMyChatMember(update tgbotapi.ChatMemberUpdated) {
	user := tables.Users{Lang: update.From.LanguageCode}

	//warn := func(err error) {
	//	app.bot.Send(tgbotapi.NewMessage(update.From.ID, user.Localize("Произошла ошибка")))
	//	app.notifyAdmin(err)
	//}
	defer func() {
		if err := app.db.UpdateUserMetrics(update.From.ID, "chatmember:"+update.NewChatMember.Status); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	switch update.NewChatMember.Status {
	case "member": // юзер разбанил бота
		app.analytics.User("{bot_was_UNblocked}", &update.From)
		if err := app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": false}); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
		app.bot.Send(tgbotapi.NewMessage(update.From.ID, user.Localize("Добро пожаловать. Мы рады, что вы снова с нами. ✋")))
	case "kicked":
		app.analytics.User("{bot_was_blocked}", &update.From)
		if err := app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": true}); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}
}
