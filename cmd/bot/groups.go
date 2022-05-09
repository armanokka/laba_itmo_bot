package bot

import (
	"github.com/armanokka/translobot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app App) onGroupMessage(message tgbotapi.Message) {
	if len(message.NewChatMembers) > 0 {
		for _, member := range message.NewChatMembers {
			if member.ID == config.BotID() { // нас пригласили на вечеринку

			}
		}
	}

}
