package repo

import "github.com/armanokka/laba_itmo_bot/internal/usecase/entity"

func (t TranslationRepo) GetLaboratoriesBySubject(subject entity.Subject) ([]int, error) {
	labs := make([]int, 0, 10)
	return labs, t.client.Model(&Laboratories{}).Where("subject = ?", subject).Select("id").Find(&labs).Error
}

func (t TranslationRepo) GetLaboratoryNameByID(labID int) (laboratoryName string, err error) {
	query := t.client.Model(&Threads{}).Where("id = ?", labID).Select("name").Find(&laboratoryName)
	if query.Error != nil {
		return "", query.Error
	}
	if query.RowsAffected == 0 {
		return "", ErrNotFound
	}
	return laboratoryName, nil
}
