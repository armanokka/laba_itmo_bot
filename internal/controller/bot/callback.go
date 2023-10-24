package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
)

func (app App) OnCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	log := app.log.With(zap.Int64("id", callback.From.ID))

	defer func() {
		if err := recover(); err != nil {
			app.notifyAdmin(err)
		}
	}()
	warn := func(err error) {
		_, file, line, _ := runtime.Caller(2)
		log.Error("", zap.Error(err), zap.String("line", file+":"+strconv.Itoa(line)))
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "произошла ошибка"))
		app.notifyAdmin(err)
	}

	user, err := app.repo.GetUserByID(callback.From.ID)
	if err != nil {
		warn(err)
	}

	data := strings.Split(callback.Data, ":")

	log = log.With(zap.String("callback_data", callback.Data))
	log.Debug("new callback")
	switch data[0] {
	case "menu":
		app.bot.Send(app.createMenu(callback.From.ID, callback.Message.MessageID))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "teacher_menu":
		edit, err := app.createTeacherMenu(callback.From.ID, callback.Message.MessageID)
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
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
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "teacher_menu")))
		app.bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, "Лабораторные работы какого потока вы хотели бы проверить?", keyboard))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "enter_queue": // data[1] - subject int, data[2] - labID int,
		subject, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		labID, err := strconv.Atoi(data[2])
		if err != nil {
			warn(err)
			return
		}

		passed, err := app.repo.UserPassedLab(callback.From.ID, labID)
		if err != nil {
			warn(err)
			return
		}
		if passed {
			app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы не можете записаться на лабу №1, потому что вы уже сдали её."))
			return
		}

		in, err := app.repo.UserInAnyQueue(callback.From.ID, entity.Subject(subject))
		if err != nil {
			warn(err)
			return
		}
		if in {
			app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Нельзя записаться на две лабы одновременно. Вы уже записаны на другую лабу по этому предмету"))
			return
		}

		var threadID int
		switch entity.Subject(subject) {
		case entity.IT:
			threadID = *user.ITThreadID
		case entity.OPD:
			threadID = *user.OPDThreadID
		case entity.Programming:
			threadID = *user.ProgrammingThreadID
		}
		if err = app.repo.AddUserToQueue(callback.From.ID, threadID, labID); err != nil {
			warn(err)
			return
		}

		labName, err := app.repo.GetLaboratoryNameByID(labID)
		if err != nil {
			warn(err)
			return
		}

		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID, labID, labName, entity.Subject(subject))
		if err != nil {
			warn(err)
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы встали в очередь на лабу №"+strconv.Itoa(labID)))
	case "leave_queue": // data[1] - subject int, data[2] - lab ID
		subject, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
		}
		labID, err := strconv.Atoi(data[2])
		if err != nil {
			warn(err)
		}
		var thread int
		switch entity.Subject(subject) {
		case entity.IT:
			thread = *user.ITThreadID
		case entity.OPD:
			thread = *user.OPDThreadID
		case entity.Programming:
			thread = *user.ProgrammingThreadID
		}
		if err = app.repo.RemoveUserFromQueue(callback.From.ID, thread, labID); err != nil {
			warn(err)
			return
		}
		labName, err := app.repo.GetLaboratoryNameByID(labID)
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, thread, labID, labName, entity.Subject(subject))
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы вышли из очереди"))
	case "update_queue": // data[1] - subject int, data[2] - lab ID int
		subject, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
		}
		labID, err := strconv.Atoi(data[2])
		if err != nil {
			warn(err)
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

		labName, err := app.repo.GetLaboratoryNameByID(labID)
		if err != nil {
			warn(err)
			return
		}

		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID, labID, labName, entity.Subject(subject))
		if err != nil {
			warn(err)
			return
		}
		edit.Text += "\n<i>" + app.now().Format("15:04:05") + "</i>"
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Информация актуальна на "+app.now().Format("15:04:05")))
	case "show_queue": // data[1] - subject int, data[2] - labID int
		subject, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		labID, err := strconv.Atoi(data[2])
		if err != nil {
			warn(err)
			return
		}

		passed, err := app.repo.UserPassedLab(callback.From.ID, labID)
		if err != nil {
			warn(err)
			return
		}
		if passed {
			app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "Вы не можете записаться на лабу №1, потому что вы уже сдали её."))
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

		labName, err := app.repo.GetLaboratoryNameByID(labID)
		if err != nil {
			warn(err)
			return
		}

		edit, err := app.createQueueMessage(callback.From.ID, callback.Message.MessageID, threadID, labID, labName, entity.Subject(subject))
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
		name, err := app.repo.GetThreadNameByID(*threadID)
		if err != nil {
			warn(err)
			return
		}
		// We made sure the user has thread
		edit, err := app.createLabSelection(callback.From.ID, callback.Message.MessageID, *threadID, name, entity.Subject(subjectID))
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "show_teacher_labs_selection": // data[1] - thread
		labs, err := app.repo.GetLaboratoriesBySubject(*user.TeacherSubject)
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		for i, lab := range labs {
			btn := tgbotapi.NewInlineKeyboardButtonData("ЛР №"+lab.Name, fmt.Sprintf("check_lab:%s:%d", data[1], lab.ID))
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
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "start_checking_labs")))
		app.bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, "Вы выбрали поток "+data[1]+"\nКакую лабораторную работу вы хотели бы проверить?", keyboard))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "update_check_lab": // data[1] - threadID int, data[2] - lab ID int
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Информация обновлена"))
		fallthrough
	case "check_lab": // data[1] - threadID int, data[2] - lab ID int
		labID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		threadName, err := app.repo.GetThreadNameByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, *user.TeacherSubject, threadName, threadID, int(labID))
		if err != nil {
			warn(err)
			return
		}
		if data[0] == "update_check_lab" {
			edit.Text += "\n<i>" + app.now().Format("15:04:05") + "</i>"
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "accept_lab": // data[1] - threadID int, data[2] - lab ID int, data[3] - student ID int64
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		labID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}
		studentID, err := strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		if err = app.repo.MarkLabAsPassed(studentID, int(labID)); err != nil {
			warn(err)
			return
		}

		// Пытаемся записать юзера на следующую лабу
		labs, err := app.repo.GetLaboratoriesBySubject(*user.TeacherSubject)
		if err != nil {
			warn(err)
			return
		}
		nextLab := -1
		for _, lab := range labs {
			if lab.ID > int(labID) {
				nextLab = lab.ID
				break
			}
		}
		text := ""
		if nextLab != -1 {
			if err = app.repo.AddUserToQueue(studentID, threadID, nextLab); err != nil {
				warn(err)
				return
			}
			text = fmt.Sprintf("Вы записаны на сдачу следующей лабы: №%d", nextLab)
		}

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: studentID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Главное меню", "menu"))),
				DisableNotification: false,
			},
			Text: fmt.Sprintf("✅ Поздравляем! Вы сдали лабу №%d\n%s", labID, text),
		})
		threadName, err := app.repo.GetThreadNameByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, *user.TeacherSubject, threadName, threadID, int(labID))
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "✅ ЛР зачтена"))
	case "lab_retake": // data[1] - threadID int, data[2] - lab ID int, data[3] - student ID int64
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		labID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}
		studentID, err := strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		if err = app.repo.MarkLabAsNotPassed(studentID, int(labID)); err != nil {
			warn(err)
			return
		}
		if err = app.repo.AddUserToQueue(studentID, threadID, int(labID)); err != nil {
			warn(err)
			return
		}

		labName, err := app.repo.GetLaboratoryNameByID(int(labID))
		if err != nil {
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
			Text: fmt.Sprintf("Вас отправили на пересдачу лабы №%s\nВы встали в конец очереди.", labName),
		})

		threadName, err := app.repo.GetThreadNameByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, *user.TeacherSubject, threadName, threadID, int(labID))
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(edit)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "✅ Студент отправлен на пересдачу"))
	case "student_missing": // data[1] - thread ID, data[2] - lab ID int, data[3] - student ID int64
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		labID, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			warn(err)
			return
		}
		studentID, err := strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			warn(err)
			return
		}

		if err = app.repo.MarkLabAsNotPassed(studentID, int(labID)); err != nil {
			warn(err)
			return
		}

		labName, err := app.repo.GetLaboratoryNameByID(int(labID))
		if err != nil {
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
			Text: fmt.Sprintf("Вы не явились на сдачу лабы №%s\nВы убраны из очереди. Вы можете встать в неё обратно.", labName),
		})

		threadName, err := app.repo.GetThreadNameByID(threadID)
		if err != nil {
			warn(err)
			return
		}
		edit, err := app.createCheckLabMenu(callback.From.ID, callback.Message.MessageID, *user.TeacherSubject, threadName, threadID, int(labID))
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
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "teacher_menu")))
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
		threadName, err := app.repo.GetThreadNameByID(threadID)
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
			Text:      "Что бы вы хотели сделать с потоком <b>" + threadName + "</b>?",
			ParseMode: tgbotapi.ModeHTML,
		})
	case "delete_thread": // data[1] - threadID int
		threadID, err := strconv.Atoi(data[1])
		if err != nil {
			warn(err)
			return
		}
		threadName, err := app.repo.GetThreadNameByID(threadID)
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
			Text:      "Поток <b>" + threadName + "</b> успешно удалён.",
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
