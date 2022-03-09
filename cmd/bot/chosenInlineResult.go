package bot

import (
	translate2 "github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (app App) chosenInlineResult(result tgbotapi.ChosenInlineResult) {
	warn := func(err error) {
		app.notifyAdmin(err)
	}
	arr := strings.Split(result.ResultID, ":")
	if len(arr) == 0 {
		arr[0] = result.ResultID
	}
	if len(arr) == 2 {
		_, ok := langs[arr[1]]
		if !ok {
			return
		}
		tr, err := translate2.MicrosoftTranslate("", arr[1], result.Query)
		if err != nil {
			warn(err)
			return
		}

		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          0,
				ChannelUsername: "",
				MessageID:       0,
				InlineMessageID: result.InlineMessageID,
				ReplyMarkup:     nil,
			},
			Text:                  tr.TranslatedText,
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})
	}
}
