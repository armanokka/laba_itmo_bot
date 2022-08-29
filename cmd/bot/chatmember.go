package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app App) onMyChatMember(update tgbotapi.ChatMemberUpdated) {
	switch update.NewChatMember.Status {
	case "member": // юзер разбанил бота
		if err := app.analytics.UserStartedBot(*update.NewChatMember.User); err != nil {
			app.notifyAdmin(err)
		}
		// not change user.Blocked since we change it in app.onMessage
	case "kicked":
		if err := app.analytics.User(tgbotapi.Message{
			MessageID:  0,
			From:       update.NewChatMember.User,
			SenderChat: nil,
			Date:       0,
			Chat: &tgbotapi.Chat{
				ID:        update.NewChatMember.User.ID,
				Type:      "private",
				UserName:  update.NewChatMember.User.UserName,
				FirstName: update.NewChatMember.User.FirstName,
				LastName:  update.NewChatMember.User.LastName,
			},
			Text: ":bot_was_blocked",
		}); err != nil {
			app.notifyAdmin(err)
		}
		if err := app.db.UpdateUserByMap(update.From.ID, map[string]interface{}{"blocked": true}); err != nil {
			app.notifyAdmin(err)
		}
	}
}
