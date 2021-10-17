package main

import (
	"database/sql"
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

// WriteLog writes messages log of user and bot. For user intent mustn't be passed.
func (u User) WriteBotLog(intent, text string) {
	go func() {
		logs <- UsersLogs{
			ID:      u.ID,
			Intent:  sql.NullString{
				String: cutStringUTF16(intent, 25),
				Valid:  true,
			},
			Text:    cutStringUTF16(text, 518),
			FromBot: true,
			Date:    time.Now(),
		}
	}()
}

func (u User) WriteUserLog(text string) {
	go func() {
		logs <- UsersLogs{
			ID:      u.ID,
			Intent:  sql.NullString{},
			Text:    cutStringUTF16(text, 518),
			FromBot: false,
			Date:    time.Now(),
		}
	}()
}

func (u User) UpdateLastActivity() {
	if err := db.Model(&Users{}).Where("id = ?", u.ID).Update("last_activity", time.Now()).Error; err != nil {
		u.error(err)
	}
}

