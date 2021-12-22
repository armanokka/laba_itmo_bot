package bot

import (
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

type User struct {
	tables.Users
	bot *botapi.BotAPI
	error func(error)
	db *gorm.DB
}

func (u User) Exists() bool {
	var exists bool
	if err := u.db.Raw("SELECT EXISTS(SELECT lang FROM users WHERE id=?)", u.ID).Find(&exists).Error; err != nil {
		u.error(err)
	}
	return exists
}
// Create creates user in u.db. Also fills a user, e.g. Fill()
func (u *User) Create(user tables.Users) {
	if err := u.db.Create(&user).Error; err != nil {
		u.error(err)
	} else {
		u.Users = user
	}
}

func (u *User) Fill() {
	if err := u.db.Model(&tables.Users{}).Where("id = ?", u.ID).Find(&u).Error; err != nil {
		u.error(err)
	}
}

func (u *User) Update(user tables.Users) {
	//if err := u.db.Model(&u.Users).Updates(user).Error; err != nil {
	//	u.error(err)
	//}
	if err := u.db.Model(&tables.Users{}).Where("id = ?", u.ID).Updates(user).Error; err != nil {
		u.error(err)
	}
	u.Fill()
	//c.Set(format(u.ID), u.Users, cache.DefaultExpiration)
}

func (u User) Localize(text string, placeholders ...interface{}) string {
	return localize(text, u.Lang, placeholders...)
}

func (u *User) UpdateLastActivity() {
	if err := u.db.Model(&tables.Users{}).Where("id = ?", u.ID).Update("last_activity", time.Now()).Error; err != nil {
		u.error(err)
		return
	}
	u.LastActivity = time.Now()
}

func (u *User) UpdateAct(act string) {
	if err := u.db.Model(&tables.Users{}).Where("id = ?", u.ID).Update("act", act).Error; err != nil {
		u.error(err)
	}
}

func (u *User) ResetAct() {
	u.UpdateAct("")
}

func (u User) StartMessage() Message {
	return Message{
		Text:     u.Localize("/start"),
		Keyboard: tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(u.Localize("My Language")),
				tgbotapi.NewKeyboardButton(u.Localize("Translate Language")),
			),
		),
	}
}

func (u User) IncrUsings() {
	if err := u.db.Model(&tables.Users{}).Exec("UPDATE users SET usings=usings+1 WHERE id = ?", u.ID).Error; err != nil {
		u.error(err)
	}
}

func (u User) SendStart(message tgbotapi.Message) {
	if !u.Exists() {
		if message.From.LanguageCode == "" || !in(config.BotLocalizedLangs, message.From.LanguageCode) {
			message.From.LanguageCode = "en"
		}
		u.Create(tables.Users{
			ID:      message.From.ID,
			Usings:  0,
			Lang:    message.From.LanguageCode,
			Blocked: false,
		})
	} else {
		u.Fill()
	}

	u.bot.Send(tgbotapi.MessageConfig{
		BaseChat:               tgbotapi.BaseChat{
			ChatID:                   message.From.ID,
			ChannelUsername:          "",
			ReplyToMessageID:         0,
			ReplyMarkup:              tgbotapi.NewRemoveKeyboard(true),
			DisableNotification:      true,
			AllowSendingWithoutReply: false,
		},
		Text:                  u.Localize("Просто напиши мне текст, а я его переведу"),
	})
	u.UpdateAct("setup_langs")
	return
}
