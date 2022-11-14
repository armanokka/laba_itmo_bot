package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (app App) onChosenInlineResult(update tgbotapi.ChosenInlineResult) {
	if err := app.analytics.InlineChosenInlineResult(update); err != nil {
		app.notifyAdmin(err)
		return
	}
	app.log.Debug("sent chosen_inline_result to the dashbot")
}
