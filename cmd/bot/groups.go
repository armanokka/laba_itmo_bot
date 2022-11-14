package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: translate in groups
func (app App) onGroupMessage(message tgbotapi.Message) {
	if len(message.NewChatMembers) > 0 {
		for _, member := range message.NewChatMembers {
			if member.ID == app.bot.Self.ID { // нас пригласили на вечеринку

			}
		}
	}

}
