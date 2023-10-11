package repo

import "github.com/armanokka/laba_itmo_bot/internal/usecase/entity"

func (t TranslationRepo) GetLaboratoriesBySubject(subject entity.Subject) ([]entity.Laboratory, error) {
	labs := make([]entity.Laboratory, 0, 8)
	return labs, t.client.Model(&Laboratories{}).Where("subject = ?", subject).Find(&labs).Error
}

func (t TranslationRepo) GetLaboratoryNameByID(labID int) (laboratoryName string, err error) {
	query := t.client.Model(Laboratories{}).Where("id = ?", labID).Select("name").Find(&laboratoryName)
	if query.Error != nil {
		return "", query.Error
	}
	if query.RowsAffected == 0 {
		return "", ErrNotFound
	}
	return laboratoryName, nil
}
