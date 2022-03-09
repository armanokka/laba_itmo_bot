package bot

import (
	translate2 "github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"strings"
)

func (app App) chosenInlineResult(result tgbotapi.ChosenInlineResult) {
	pp.Println(result)
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
		from, err := translate2.DetectLanguageGoogle(cutStringUTF16(result.Query, 100))
		if err != nil {
			warn(err)
			return
		}

		tr, err := translate2.FlexibleTranslate(from, arr[1], result.Query)
		if err != nil {
			warn(err)
			return
		}

		if _, err = app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          0,
				ChannelUsername: "",
				MessageID:       0,
				InlineMessageID: result.InlineMessageID,
				ReplyMarkup:     nil,
			},
			Text:                  tr,
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			pp.Println(err)
		}
	}
}
