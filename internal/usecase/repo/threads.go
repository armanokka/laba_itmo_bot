package repo

import (
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"sort"
	"strconv"
	"unicode"
)

func (t TranslationRepo) GetThreadsBySubject(subject entity.Subject) ([]entity.Thread, error) {
	threads := make([]entity.Thread, 0, 15)
	err := t.client.Model(&Threads{}).Where("subject = ?", subject).Find(&threads).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(threads, func(i, j int) bool {
		l1, n1b, n1a := extractNumber(threads[i].Name) // letter 1, number 1 before dot, number 1 before dot
		l2, n2b, n2a := extractNumber(threads[j].Name)
		if l1 != l2 {
			return l1 < l2
		}
		if n1b != n2b {
			return n1b < n2b
		}
		return n1a < n2a
	})
	return threads, nil
}

func extractNumber(s string) (letter string, beforeDot int, afterDot int) {
	dot := false
	for _, ch := range s {
		if unicode.IsLetter(ch) {
			letter += string(ch)
			continue
		}
		if ch == '.' || ch == ',' {
			dot = true
			continue
		}
		if unicode.IsDigit(ch) {
			if dot {
				afterDot *= 10
				n, _ := strconv.Atoi(string(ch))
				afterDot += n
				continue
			}
			beforeDot *= 10
			n, _ := strconv.Atoi(string(ch))
			beforeDot += n
		}
	}
	return letter, beforeDot, afterDot
}

func (t TranslationRepo) AddThread(name string, subject entity.Subject) error {
	return t.client.Create(&Threads{
		Subject: subject,
		Name:    name,
	}).Error
}

func (t TranslationRepo) DeleteThread(threadID int) error {
	panic("implement me!")
	// tODO: need to configure foreign keys in db
	return nil
}

func (t TranslationRepo) RenameThread(threadID int, newName string) error {
	query := t.client.Model(&Threads{}).Where("id = ?", threadID).Update("name", newName)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (t TranslationRepo) GetThreadNameByID(threadID int) (thread string, err error) {
	query := t.client.Model(&Threads{}).Where("id = ?", threadID).Select("name").Find(&thread)
	if query.Error != nil {
		return "", query.Error
	}
	if query.RowsAffected == 0 {
		return "", ErrNotFound
	}
	return thread, nil
}
