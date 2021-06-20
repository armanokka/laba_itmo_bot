package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const message = "Bot is working"

type Users struct {
	ID     int64 `gorm:"primaryKey;index;not null"`
	MyLang string `gorm:"default:en"`
	ToLang string `gorm:"default:ar;n"`
	Act    string
}


func main()  {
	bot, err := tgbotapi.NewBotAPI("1737819626:AAEoc8WyCq_8rFQcY4q0vtkhqCKro8AudfI")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open("host=ec2-63-34-97-163.eu-west-1.compute.amazonaws.com user=wzlryrrgxbgsbw password=b578bdbc77b5394a60f57660487149ca2238e0cbaf1cdbfb8b931f1168af24c7 dbname=d21k8q9pl6acl4 port=5432 TimeZone=Europe/Moscow"), &gorm.Config{SkipDefaultTransaction: true, PrepareStmt: false})
	if err != nil {
		panic(err)
	}

	var users []Users
	err = db.Model(&Users{}).Find(&users).Error
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ID, message)
		msg.ParseMode = tgbotapi.ModeHTML
		_, err := bot.Send(msg)
		if e, ok := err.(tgbotapi.Error); ok {
			if e.RetryAfter != 0 {
				time.Sleep(time.Duration(e.RetryAfter) * time.Second)
				msg := tgbotapi.NewMessage(user.ID, message)
				msg.ParseMode = tgbotapi.ModeHTML
				bot.Send(msg)
			}
			if e.Code == 403 {
				db.Model(&Users{}).Where("id = ?", user.ID).Limit(1).Delete(&user)
			}
		}
		pp.Println("sent message to", user.ID)
		time.Sleep(33 * time.Millisecond)
	}
	pp.Println("mailing done")
}
