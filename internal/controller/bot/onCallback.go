package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
)

func (app *App) OnCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	log := app.log.With(zap.Int64("id", callback.From.ID), zap.String("callback_data", callback.Data))
	//log.Debug("new callback")

	defer func() {
		if err := recover(); err != nil {
			app.notifyAdmin(err)
		}
	}()
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "Что-то пошло не так... Попробуйте написать боту /start"))
		app.notifyAdmin(err)

		_, file, line, _ := runtime.Caller(2)
		log.Error("", zap.Error(err), zap.String("line", file+":"+strconv.Itoa(line)))
	}

	user, err := app.repo.GetUserByID(callback.From.ID)
	if err != nil {
		warn(err)
		return
	}

	data := strings.Split(callback.Data, ":")

	switch data[0] {
	case "menu":
		if user.TeacherSubject != nil {
			edit, err := app.createTeacherMainMenu(callback.From.ID, callback.Message.MessageID)
			if err != nil {
				warn(err)
				return
			}
			app.bot.Send(edit)
			app.SetCurrentPassingStudent(0)
		} else {
			app.bot.Send(app.createMainMenu(callback.From.ID, callback.Message.MessageID))
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "start_checking_labs":
		threads, err := app.repo.GetThreadsBySubject(*user.TeacherSubject)
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		for i, thread := range threads {
			btn := tgbotapi.NewInlineKeyboardButtonData(thread.Name, fmt.Sprintf("show_teacher_labs_selection:%d", thread.ID))
			if i%4 == 0 || len(keyboard.InlineKeyboard) == 0 {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
				continue
			}
			l := len(keyboard.InlineKeyboard) - 1
			if l < 0 {
				l = 0
			}
			keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], btn)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "menu")))
		app.SetCurrentPassingStudent(0)
		app.bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, "Лабораторные работы какого потока вы хотели бы проверить?", keyboard))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "enter_queue": // data[1] - threadID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}

		lab, err := app.repo.GetNextLab(callback.From.ID, thread.Subject)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы уже сдали все лабы ;)"))
				return
			}
			warn(err)
			return
		}

		if err = app.repo.AddUserToQueue(callback.From.ID, threadID, lab.ID); err != nil {
			warn(err)
			return
		}

		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы встали в очередь"))
	case "leave_queue": // data[1] - thread ID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("❌ Выйти", "confirm_leaving_queue:"+data[1]),
				tgbotapi.NewInlineKeyboardButtonData("✅ Остаться", "show_queue:"+strconv.Itoa(int(thread.Subject))),
			))
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.From.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: &keyboard,
			},
			Text: "Вы уверены, что хотите выйти из очереди? Отменить это действие будет невозможно",
		})
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "open_change_my_lab_menu": // data[1] - thread id
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createLabSelection(callback.From.ID, callback.Message.MessageID, threadID, thread.Subject)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "change_my_lab": // data[1] - thread ID, data[2] - new lab ID
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		newLabID, err := strconv.Atoi(data[2])
		if err != nil {
			warn(err)
			return
		}
		if err = app.repo.ChangeUserLab(callback.From.ID, threadID, newLabID); err != nil {
			warn(err)
			return
		}
		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, fmt.Sprintf("Теперь вы сдаёте лабу №%d", newLabID)))
	case "confirm_leaving_queue":
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		if err = app.repo.RemoveUserFromQueue(callback.From.ID, threadID); err != nil {
			warn(err)
			return
		}
		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы вышли из очереди"))
	case "update_queue": // data[1] - thread ID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}

		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		edit.Text += "\n<i>" + app.now().Format("15:04:05") + "</i>"
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Информация актуальна на "+app.now().Format("15:04:05")))
	case "show_queue": // data[1] - subject int
		subject, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}

		var threadID int
		switch entity.Subject(subject) {
		case entity.IT:
			threadID = *user.ITThreadID
		case entity.Programming:
			threadID = *user.ProgrammingThreadID
		case entity.OPD:
			threadID = *user.OPDThreadID
		}

		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "set_thread": // data[1] - subject int, data[2] - threadID int
		subject, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		threadID, err := strconv.Atoi(data[2])
		if err != nil {
			warn(err)
			return
		}
		column := ""
		switch entity.Subject(subject) {
		case entity.IT:
			column = "it_thread_id"
			user.ITThreadID = &threadID
		case entity.Programming:
			column = "programming_thread_id"
			user.ProgrammingThreadID = &threadID
		case entity.OPD:
			column = "opd_thread_id"
			user.OPDThreadID = &threadID
		}
		if err = app.repo.UpdateUserByID(callback.From.ID, column, threadID); err != nil {
			warn(err)
			return
		}
		fallthrough
	case "show_labs_selection": // data[1] - subject int, data[2] - threadID int
		subjectID, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		// Before parsing thread we check if user has it for this subject
		var threadID *int
		switch entity.Subject(subjectID) {
		case entity.IT:
			threadID = user.ITThreadID
		case entity.Programming:
			threadID = user.ProgrammingThreadID
		case entity.OPD:
			threadID = user.OPDThreadID
		}
		if threadID == nil {
			threads, err := app.repo.GetThreadsBySubject(entity.Subject(subjectID))
			if err != nil {
				warn(err)
				return
			}

			if len(threads) == 0 {
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "menu")))
				app.bot.Send(tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      callback.From.ID,
						MessageID:   callback.Message.MessageID,
						ReplyMarkup: &keyboard,
					},
					Text: entity.Subject(subjectID).Name() + "\n\n" + "Список потоков пуст. Обратись к создателю бота @armanokka за помощью",
				})
				return
			}
			keyboard := tgbotapi.NewInlineKeyboardMarkup()
			for i, thread := range threads {
				btn := tgbotapi.NewInlineKeyboardButtonData(thread.Name, fmt.Sprintf("set_thread:%s:%d", strconv.Itoa(int(subjectID)), thread.ID))
				if i%4 == 0 || len(keyboard.InlineKeyboard) == 0 {
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
					continue
				}
				l := len(keyboard.InlineKeyboard) - 1
				if l < 0 {
					l = 0
				}
				keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], btn)
			}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "menu")))
			app.bot.Send(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      callback.From.ID,
					MessageID:   callback.Message.MessageID,
					ReplyMarkup: &keyboard,
				},
				Text:                  "<b>" + entity.Subject(subjectID).Name() + "</b>\n\n" + "Выбери свой поток:",
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			})
			app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
			return
		}
		// We made sure the user has thread
		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, *threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "show_teacher_labs_selection": // data[1] - thread
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "update_check_lab": // data[1] - threadID int
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Информация актуальна на "+app.now().Format("15:04:05")))
		fallthrough
	case "check_lab": // data[1] - threadID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		if data[0] == "update_check_lab" {
			edit.Text += "\n<i>" + app.now().Format("15:04:05") + "</i>"
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "accept_lab": // data[1] - threadID int data[2] - student ID int64
		// TODO сделать общую очередь для всех лаб
		// TODO отмечать того сдающего, который открыт у учителя
		// TODO считать количество пересдач
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		studentID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		if _, err = app.repo.GradeLab(studentID, threadID, true); err != nil {
			warn(err)
			return
		}

		nextLab, err := app.repo.GetNextLab(studentID, thread.Subject)
		if err != nil && !errors.Is(err, repo.ErrNotFound) {
			warn(err)
			return
		}

		// Пытаемся записать юзера на следующую лабу

		text := ""
		if nextLab.ID != 0 {
			if err = app.repo.AddUserToQueue(studentID, threadID, nextLab.ID); err != nil {
				warn(err)
				return
			}
			text = fmt.Sprintf("Вы записаны на сдачу следующей лабы: №%s", nextLab.Name)
		}

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: studentID,
			},
			Text: fmt.Sprintf("✅ Поздравляем! Вы сдали лабу №%s\n", nextLab.Name) + text,
		})
		app.bot.Send(tgbotapi.StickerConfig{tgbotapi.BaseFile{ //nolint:govet
			BaseChat: tgbotapi.BaseChat{
				ChatID: studentID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Главное меню", "menu"))),
			},
			File: tgbotapi.FileID(`CAACAgIAAxkBAAEXUKplTNfITDQ1wwXlwzx4U87NahYcUQAC5BQAAqt86UviEhEqhf3MYjME`),
		}})
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "✅ ЛР зачтена"))
	case "lab_retake": // data[1] - threadID int, data[2] - student ID int64
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		studentID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		labID, err := app.repo.GradeLab(studentID, threadID, false)
		if err != nil {
			warn(err)
			return
		}
		if err = app.repo.AddUserToQueue(studentID, threadID, labID); err != nil {
			warn(err)
			return
		}

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: studentID,
			},
			Text: fmt.Sprintf("Вас отправили на пересдачу лабы по %s.\nВы встали в конец очереди.", thread.Subject.NameGenitiveCase()),
		})
		app.bot.Send(tgbotapi.StickerConfig{tgbotapi.BaseFile{ //nolint:govet
			BaseChat: tgbotapi.BaseChat{
				ChatID: studentID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Главное меню", "menu"))),
			},
			File: tgbotapi.FileID(`CAACAgIAAxkBAAEXTJNlS6NRr3e-fv4vgcQnbjPTCluxcAACaBoAAkYKqUjgFcYYloWvkzME`),
		}})

		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "✅ Студент отправлен на пересдачу"))
	case "student_missing": // s1] - thread ID, data[2] - student ID int64
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		studentID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		if _, err = app.repo.GradeLab(studentID, threadID, false); err != nil {
			warn(err)
			return
		}

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: studentID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Главное меню", "menu"))),
				DisableNotification: true,
			},
			Text: fmt.Sprintf("Вы не явились на сдачу лабы по %s\nВы убраны из очереди. Вы можете встать в неё обратно.", thread.Subject.NameGenitiveCase()),
		})
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, threadID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "✅ Студент отправлен на пересдачу"))
	case "manage_threads":
		threads, err := app.repo.GetThreadsBySubject(*user.TeacherSubject)
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		for i, thread := range threads {
			btn := tgbotapi.NewInlineKeyboardButtonData(thread.Name, fmt.Sprintf("manage_thread:%d", thread.ID))
			if i%4 == 0 || len(keyboard.InlineKeyboard) == 0 {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
				continue
			}
			l := len(keyboard.InlineKeyboard) - 1
			if l < 0 {
				l = 0
			}
			keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], btn)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить поток", "add_thread")))
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "menu")))
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.From.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: &keyboard,
			},
			Text: "Управление потоками",
		})
	case "manage_thread": // data[1] - thread ID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить поток", fmt.Sprintf("delete_thread:%d", threadID))),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✏️ Переименовать поток", fmt.Sprintf("rename_thread:%d", threadID))),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "manage_threads")))
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.From.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: &keyboard,
			},
			Text:      "Что бы вы хотели сделать с потоком <b>" + thread.Name + "</b>?",
			ParseMode: tgbotapi.ModeHTML,
		})
	case "delete_thread": // data[1] - threadID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		thread, err := app.repo.GetThreadByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		if err = app.repo.DeleteThread(threadID); err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "manage_threads")))
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      callback.From.ID,
				MessageID:   callback.Message.MessageID,
				ReplyMarkup: &keyboard,
			},
			Text:      "Поток <b>" + thread.Name + "</b> успешно удалён.",
			ParseMode: tgbotapi.ModeHTML,
		})
	case "rename_thread": // data[1] - threadID int
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Фича пока не реализована!"))
		//threadID, err := strconv.Atoi(data[1])
		//if err != nil {
		//	warn(err)
		//	return
		//}
		//threadName, err := app.repo.GetThreadNameByID(threadID)
		//if err != nil {
		//	warn(err)
		//	return
		//}

	}
}
