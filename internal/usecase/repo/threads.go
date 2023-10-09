package repo

import (
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"sort"
)

func (t TranslationRepo) GetThreadsBySubject(subject entity.Subject) ([]entity.Thread, error) {
	threads := make([]entity.Thread, 0, 15)
	err := t.client.Model(&Threads{}).Where("subject = ?", subject).Find(&threads).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(threads, func(i, j int) bool {
		return threads[i].Name < threads[j].Name
	})
	return threads, nil
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
