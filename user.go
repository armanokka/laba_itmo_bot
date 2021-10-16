package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"time"
)

type User struct {
	Users
	error func(error)
}

// NewUser return User with such id
func NewUser(id int64, errfunc func(error)) User {
	return User{
		Users: Users{ID: id},
		error: errfunc,
	}
}

func (u User) Exists() bool {
	var exists bool
	if err := db.Raw("SELECT EXISTS(SELECT lang FROM users WHERE id=?)", u.ID).Find(&exists).Error; err != nil {
		u.error(err)
	}
	return exists
}
// Create creates user in db. Also fills a user, e.g. Fill()
func (u User) Create(user Users) {
	if err := db.Create(&user).Error; err != nil {
		u.error(err)
	} else {
		u.Users = user
	}
}

func (u *User) Fill() {
	if v, ok := c.Get(format(u.ID)); ok {
		u.Users = v.(Users)
		return
	}
	if err := db.Model(&Users{}).Where("id = ?", u.ID).Find(&u.Users).Error; err != nil {
		u.error(err)
	}
	c.Set(format(u.ID), u.Users, cache.DefaultExpiration)
}

func (u *User) Update(user Users) {
	if err := db.Model(&u.Users).Updates(user).Error; err != nil {
		u.error(err)
	}
	if err := db.Model(&Users{}).Where("id = ?", u.ID).Find(&u.Users).Error; err != nil {
		u.error(err)
	}
	c.Set(format(u.ID), u.Users, cache.DefaultExpiration)
}

func (u User) Localize(text string, placeholders ...interface{}) string {
	return localize(text, u.Lang, placeholders...)
}

//func (u User) WriteLog(intent string) {
//	var exists bool
//	if err := db.Model(&UsersLogs{}).Raw("SELECT EXISTS(SELECT uid FROM users_logs WHERE intent=? AND id=? ORDER BY uid DESC LIMIT 1)", intent, u.ID).Find(&exists).Error; err != nil {
//		u.error(err)
//	}
//	if !exists {
//		err := db.Model(&UsersLogs{}).Raw("UPDATE users_logs SET times=times+1 WHERE id=? AND intent=? ORDER BY uid DESC LIMIT 1", u.ID, intent).Error
//		if err != nil {
//			u.error(err)
//		}
//		return
//	}
//}

func (u User) UpdateLastActivity() {
	if err := db.Model(&Users{}).Where("id = ?", u.ID).Update("last_activity", time.Now()).Error; err != nil {
		u.error(err)
	}
}

func (u User) SendStart() {
	msg := tgbotapi.NewMessage(u.ID, u.Localize("Just send me a text and I will translate it"))
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(u.Localize("My Language")),
			tgbotapi.NewKeyboardButton(u.Localize("Translate Language")),
		),
	)

	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)

	analytics.Bot(u.ID, msg.Text, "Start")
}