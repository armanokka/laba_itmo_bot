package bot

import (
	"context"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"github.com/armanokka/laba_itmo_bot/pkg/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

func (app *App) onMessage(ctx context.Context, message tgbotapi.Message) {
	log := app.log.With(zap.Int64("id", message.From.ID), zap.String("message", message.Text))
	defer func() {
		if err := recover(); err != nil {
			app.notifyAdmin(err)
		}
	}()

	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Что-то пошло не так... Попробуйте написать боту /start"))
		app.notifyAdmin(err)

		_, file, line, _ := runtime.Caller(2)
		log.Error("", zap.Error(err), zap.String("line", file+":"+strconv.Itoa(line)))
	}

	//log.Debug("new message")

	if message.Chat.ID < 0 {
		return
	}

	var err error
	user, err := app.repo.GetUserByID(message.From.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			warn(err)
			return
		}
		action := "setup_fio"
		if err = app.repo.CreateUser(entity.User{
			ID:         message.From.ID,
			FirstName:  nil,
			LastName:   nil,
			Patronymic: nil,
			Action:     &action,
		}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, "⭐️ Добро пожаловать в бота по записи на сдачу лаб!"))
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, "Введите свои ФИО:\nПример: Иванов Иван Иванович"))
		return
	}

	if user.Action != nil {
		switch *user.Action {
		case "setup_fio":
			chunks := strings.Fields(message.Text)
			if len(chunks) < 2 {
				app.bot.Send(tgbotapi.NewMessage(message.From.ID, "❌ Ошибка. ФИО должно состоять из двух и более слов. Напиши заново"))
				return
			}
			for i, word := range chunks {
				for _, r := range word {
					if !unicode.Is(unicode.Cyrillic, r) {
						app.bot.Send(tgbotapi.NewMessage(message.From.ID, "❌ Ошибка. ФИО должно состоять из русских букв. Напиши заново"))
						return
					}
				}
				if i <= 2 {
					chunks[i] = toupperfirst(word)
				}
			}
			if len(chunks) >= 3 {
				chunks[2] = strings.Join(chunks[2:], " ")
				chunks = chunks[:3]
				if err = app.repo.UpdateUserByID(message.From.ID, "patronymic", chunks[2]); err != nil {
					warn(err)
					return
				}
			}
			if err = app.repo.UpdateUserByID(message.From.ID, "first_name", chunks[0], "last_name", chunks[1], "action", nil); err != nil {
				warn(err)
				return
			}
		}
	}
	// todo: -добавить фичу подсчета среднего времени сдачи?

	if user.TeacherSubject != nil {
		msg, err := app.createTeacherMainMenu(message.From.ID, 0)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(msg)
		return
	}

	switch message.Command() {
	case "id":
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10)))
	default:
		app.bot.Send(app.createMainMenu(message.From.ID, 0))
	}
	// no code here!
}
