/*
Helper functions
*/
package bot

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"github.com/armanokka/laba_itmo_bot/pkg/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
)

func (app App) createLabSelection(userID int64, messageID int, threadID int, subject entity.Subject) (tgbotapi.EditMessageTextConfig, error) {
	labs, err := app.repo.GetLaboratoriesBySubject(subject)
	if err != nil {
		return tgbotapi.EditMessageTextConfig{}, err
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for i, lab := range labs {
		btn := tgbotapi.NewInlineKeyboardButtonData("Лаба №"+lab.Name, fmt.Sprintf("change_my_lab:%s:%s", strconv.Itoa(threadID), strconv.Itoa(lab.ID)))
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
		tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "show_queue:"+strconv.Itoa(int(subject)))))

	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      userID,
			MessageID:   messageID,
			ReplyMarkup: &keyboard,
		},
		Text: fmt.Sprintf("Какую лабу вы сдаёте? После выбора вы не потеряете свое место в очереди."),
	}, nil
}

func (app *App) createQueueMessage(userID int64, messageID, threadID int) (tgbotapi.EditMessageTextConfig, error) {
	queue, err := app.repo.GetQueueByThreadID(threadID)
	if err != nil {
		return tgbotapi.EditMessageTextConfig{}, err
	}

	people := ""
	before := 0
	in := false

	for i, booking := range queue {
		people += "\n"
		if booking.Patronymic == nil {
			s := ""
			booking.Patronymic = &s
		}
		fio := fmt.Sprintf(`%d. <a href="tg://user?id=%d">%s %s %s</a> / ЛР №%s`, i+1, booking.UserID, booking.FirstName, booking.LastName, *booking.Patronymic, booking.LabName)
		if !booking.Checked && booking.UserID != userID && !in {
			before++
		}
		if booking.Passed {
			fio += " (сдал)"
		}
		if booking.UserID == app.currentPassingStudent {
			fio += "  ⬅️ (сдает сейчас. " + app.now().Format("15:04:05 2/1") + ")"
		} else if app.currentPassingStudent == 0 && i == 0 {
			fio += "  ⬅️ (проверка начнется во время пары)"
		}
		if booking.UserID == userID {
			fio += " (вы)"
			in = true
		}
		people += fio
	}
	inQueueText := "<b>Вы в очереди</b> ✅"
	if !in {
		inQueueText = "<b>Вы не в очереди</b> ❌"
	}
	beforeYouText := ""
	if in {
		beforeYouText = "\n<i>до вас <b>" + strconv.Itoa(before) + "</b> " + declOfNum(before, []string{"человек", "человека", "человек"}) + "</i>"
		if before == 0 {
			beforeYouText = "\nвы сдаете первым\\первой"
		}
	}
	if len(queue) == 0 {
		people = "<i>очередь пуста</i>"
	}

	thread, err := app.repo.GetThreadByID(threadID)
	if err != nil {
		return tgbotapi.EditMessageTextConfig{}, err
	}

	// Creating keyboard
	btn := tgbotapi.NewInlineKeyboardButtonData("❌ Выйти из очереди", "leave_queue:"+strconv.Itoa(threadID))
	if !in {
		btn = tgbotapi.NewInlineKeyboardButtonData("✅ Встать в очередь", fmt.Sprintf("enter_queue:%s", strconv.Itoa(threadID)))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			btn,
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить очередь", fmt.Sprintf("update_queue:%s", strconv.Itoa(threadID))),
		),
	)
	if in {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✍️ Сдаю другую лабу", fmt.Sprintf("open_change_my_lab_menu:%s", strconv.Itoa(threadID))),
		))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "menu"),
	))

	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      userID,
			ReplyMarkup: &keyboard,
			MessageID:   messageID,
		},
		Text: fmt.Sprintf(`Очередь потока <b>%s</b>. %s.  

%s %s

очередь:%s`, thread.Name, thread.Subject.Name(), inQueueText, beforeYouText, people),
		ParseMode: tgbotapi.ModeHTML,
	}, nil
}

func (app *App) createTeacherMainMenu(userID int64, messageID int) (tgbotapi.EditMessageTextConfig, error) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Проверять ЛР", "start_checking_labs"),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("Список потоков", "manage_threads"),
		//	tgbotapi.NewInlineKeyboardButtonData("Список ЛР", "manage_labs"),
		//),
	)
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      userID,
			MessageID:   messageID,
			ReplyMarkup: &keyboard,
		},
		Text:      "Панель управления очередями на сдачу лабораторных работ",
		ParseMode: tgbotapi.ModeHTML,
	}, nil
}

func (app *App) createCheckLabMenu(userID int64, messageID int, threadID int) (tgbotapi.EditMessageTextConfig, error) {
	queue, err := app.repo.GetQueueByThreadID(threadID)
	if err != nil {
		return tgbotapi.EditMessageTextConfig{}, err
	}

	// Choosing first not checked student
	var currentStudent entity.QueueUser
	fullQueue := ""
	afterCurrentStudentCount := 0
	for i, booking := range queue {
		fio := fmt.Sprintf("\n"+`%d. <a href="tg://user?id=%d">%s %s</a> / ЛР №%s`, i+1, booking.UserID, booking.FirstName, booking.LastName, booking.LabName)

		if booking.Checked {
			fio = "<s>" + fio + "</s>"
		} else {
			if currentStudent.UserID == 0 {
				currentStudent = booking
				fio += "  ⬅️ (сейчас)"
			} else {
				afterCurrentStudentCount++
			}
		}
		fullQueue += fio
	}

	thread, err := app.repo.GetThreadByID(threadID)
	if err != nil {
		return tgbotapi.EditMessageTextConfig{}, err
	}

	if currentStudent.UserID == 0 {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить очередь", fmt.Sprintf("update_check_lab:%d", threadID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "start_checking_labs"),
			))
		return tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      userID,
				MessageID:   messageID,
				ReplyMarkup: &keyboard,
			},
			Text: fmt.Sprintf(`Очередь потока <b>%s</b>. %s.

<b>сейчас никто не сдаёт.</b> <i>0/0</i>
%s

<i>очередь пуста</i>`, thread.Name, thread.Subject.Name(), fullQueue),
			ParseMode: tgbotapi.ModeHTML,
		}, nil
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Зачесть ЛР", fmt.Sprintf("accept_lab:%d:%d", threadID, currentStudent.UserID)),
			tgbotapi.NewInlineKeyboardButtonData("🚫 Пересдача", fmt.Sprintf("lab_retake:%d:%d", threadID, currentStudent.UserID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚷 Студент отсутствует", fmt.Sprintf("student_missing:%d:%d", threadID, currentStudent.UserID)),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("⏭ Пропустить", fmt.Sprintf("student_missing:%d:%d:%d", threadID, labID, currentStudent.UserID)),
		//),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить очередь", fmt.Sprintf("update_check_lab:%d", threadID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "start_checking_labs"),
		),
	)
	app.SetCurrentPassingStudent(currentStudent.UserID)
	if currentStudent.Patronymic == nil {
		s := ""
		currentStudent.Patronymic = &s
	}
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      userID,
			MessageID:   messageID,
			ReplyMarkup: &keyboard,
		},
		Text: fmt.Sprintf(`%s %s

Сейчас сдаёт:  <b><a href="tg://user?id=%d">%s %s %s</a></b> (ЛР <b>№%s</b>)

очередь:%s

<i>Пожалуйста, проверьте лабораторную работу студента, а затем нажмите на одну из кнопок.</i>`, thread.Subject.Name(), thread.Name, currentStudent.UserID, currentStudent.FirstName, currentStudent.LastName, *currentStudent.Patronymic, currentStudent.LabName, fullQueue),
		ParseMode: tgbotapi.ModeHTML,
	}, nil
}

func (app *App) createMainMenu(userID int64, messageID int) tgbotapi.EditMessageTextConfig {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Программирование", fmt.Sprintf("show_labs_selection:%d", int(entity.Programming))),
			tgbotapi.NewInlineKeyboardButtonData("Информатика", fmt.Sprintf("show_labs_selection:%d", int(entity.IT))),
			tgbotapi.NewInlineKeyboardButtonData("ОПД", fmt.Sprintf("show_labs_selection:%d", int(entity.OPD))),
		))
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      userID,
			MessageID:   messageID,
			ReplyMarkup: &keyboard,
		},
		Text:      "<b>Главное меню</b>\nВыбери предмет, чтобы записаться на сдачу лабы по нему.",
		ParseMode: tgbotapi.ModeHTML,
	}
}

// cutString cut string using runes by limit
func cutStringUTF16(text string, limit int) string {
	points := utf16.Encode([]rune(text))
	if len(points) > limit {
		return string(utf16.Decode(points[:limit]))
	}
	return text
}

func toupperfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

func in[T comparable](arr []T, keys ...T) bool {
	for _, k := range keys {
		exists := false
		for _, v := range arr {
			if k == v {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
}

// inI is in but returning index
func inSlice(arr []string, k string) (int, bool) {
	for i, v := range arr {
		if k == v {
			return i, true
		}
	}
	return 0, false
}
func inFuzzy(arr []string, keys ...string) bool {
	for _, k := range keys {
		exists := false
		for _, v := range arr {
			if k == v || diff(k, v) == 1 {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
}

var SupportedFormattingTags = []string{"b", "strong", "i", "em", "u", "ins", "s", "strike", "del", "span", "tg-spoiler", "a", "code", "pre"}

// validHtml doesn't deal with order of tags, e.g. "</b>hey<b>" is valid
func validHtml(s string) bool {
	d := xml.NewDecoder(strings.NewReader(s))
	tags := make(map[string]int, 2)
	for {
		token, err := d.Token()
		if err != nil && err != io.EOF {
			return false
		}
		if token == nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			tags[t.Name.Local]++
		case xml.EndElement:
			tags[t.Name.Local]--
		}
	}
	for _, count := range tags {
		if count != 0 {
			return false
		}
	}
	return true
}

func inMapValues(m map[string]string, values ...string) bool {
	for _, v := range values {
		var ok bool
		for _, v1 := range m {
			if v == v1 {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func hasPrefix(s, prefix string, maxCharsDifference int) bool {
	if maxCharsDifference < 1 {
		maxCharsDifference = 0
	}
	runes := []rune(s)
	prefixRunes := []rune(prefix)
	if len(runes) > len(prefixRunes) {
		runes = runes[:len(prefixRunes)]
	}
	for i, r := range runes {
		if r != prefixRunes[i] {
			maxCharsDifference--
			if maxCharsDifference < 0 {
				return false
			}
		}
	}
	return maxCharsDifference > -1
}

func maxDiff(source string, arrs [][]string) []string {
	i, maxDifference := 0, 0
	for idx, arr := range arrs {
		difference := 0
		for _, v := range arr {
			if v == "" {
				continue
			}
			difference += diff(source, v)
		}
		if difference > maxDifference {
			i, maxDifference = idx, difference
		}
	}
	return arrs[i]
}

func randid(seed int64) string {
	rand.Seed(seed)
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func concatNonEmpty(separator string, ss ...string) string {
	b := new(bytes.Buffer)
	for i, s := range ss {
		if strings.TrimSpace(s) == "" {
			continue
		}
		if i > 0 {
			b.WriteString(separator)
		}
		b.WriteString(s)
	}
	return b.String()
}

// diff counts difference between s1 and s2 by comparing their characters
func diff(s1, s2 string) (n int) {
	reader := strings.NewReader(s1)
	for _, ch := range s2 {
		r, _, err := reader.ReadRune()
		if err == nil {
			if r != ch {
				n++
			}
			continue
		}
		break
	}
	if l1, l2 := len(s1), len(s2); l1 != l2 {
		if l1 > l2 {
			n += l1 - l2
		} else {
			n += l2 - l1
		}
	}
	return n
}

type Keyboard struct {
	Dictionary          []string            `json:"dictionary"`
	Paraphrase          []string            `json:"paraphrase"`
	Examples            map[string]string   `json:"examples"`
	ReverseTranslations map[string][]string `json:"reverse_translations"`
}

func buildKeyboard(from, to string, ret Keyboard) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔉", fmt.Sprintf("speech:%s:%s", from, to))))
	if len(ret.Examples) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬", fmt.Sprintf("examples:%s:%s", from, to)))
	}
	if len(ret.Dictionary) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("📖", fmt.Sprintf("dictionary:%s:%s", from, to)))
	}
	if len(ret.Paraphrase) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("✨", fmt.Sprintf("paraphrase:%s:%s", from, to)))
	}
	if len(ret.ReverseTranslations) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("📚", fmt.Sprintf("reverse_translations:%s:%s", from, to)))
	}
	return keyboard
}

func IsCtxError(err error) bool {
	if e, ok := err.(errors.Error); ok {
		return IsCtxError(e.Err)
	}
	return errors.Is(err, context.Canceled)
}

func parseKeyboard(text string) (keyboard interface{}, valid bool) {
	kb := tgbotapi.NewInlineKeyboardMarkup()
	if text != "/empty" {
		scanner := bufio.NewScanner(strings.NewReader(text))
		for scanner.Scan() {
			if scanner.Err() != nil {
				return
			}
			btns := strings.Fields(scanner.Text())
			row := tgbotapi.NewInlineKeyboardRow()
			for _, btn := range btns {
				parts := strings.Split(btn, "|") // parts[0] - text on button, parts[1] - link for button
				if len(parts) != 2 {
					return
				}
				row = append(row, tgbotapi.NewInlineKeyboardButtonURL(parts[0], parts[1]))
			}
			kb.InlineKeyboard = append(kb.InlineKeyboard, row)
		}
	}
	if len(kb.InlineKeyboard) > 0 {
		return kb, true
	}
	return
}

func inSliceNotUsed(arr []string, usedIndexes []int, k string) (int, bool) {
	for i, v := range arr {
		if k == v {
			for _, usedIdx := range usedIndexes {
				if i != usedIdx {
					return i, true
				}
			}
		}
	}
	return 0, false
}

func remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}
func removeHtml(s string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	return doc.Text(), nil
}

func Tick(callbackData string, inlineKeyboard [][]tgbotapi.InlineKeyboardButton) {
	for i1, row := range inlineKeyboard {
		for i2, button := range row {
			if button.CallbackData == nil {
				continue
			}
			if *button.CallbackData == callbackData && !strings.HasPrefix(*button.CallbackData, "✅ ") {
				inlineKeyboard[i1][i2].Text = "✅ " + button.Text
				break
			}
		}
	}
}

func UntickAll(inlineKeyboard [][]tgbotapi.InlineKeyboardButton) {
	for i1, row := range inlineKeyboard {
		for i2, button := range row {
			inlineKeyboard[i1][i2].Text = strings.TrimPrefix(button.Text, "✅ ")
		}
	}
}

func declOfNum(number int, titles []string) string {
	if number < 0 {
		number *= -1
	}

	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return titles[currentCase]
}
