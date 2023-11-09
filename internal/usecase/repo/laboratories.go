package repo

import "github.com/armanokka/laba_itmo_bot/internal/usecase/entity"

func (t TranslationRepo) GetLaboratoriesBySubject(subject entity.Subject) ([]entity.Laboratory, error) {
	labs := make([]entity.Laboratory, 0, 8) // should be of type Laboratories
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

func (t TranslationRepo) GetNextLab(userID int64, subject entity.Subject) (entity.Laboratory, error) {
	var record Laboratories
	query := t.client.Model(Laboratories{}).Raw(`SELECT * FROM laboratories WHERE id NOT IN (SELECT laboratory_id FROM queues WHERE passed=true AND user_id=?) AND subject = ? ORDER BY id LIMIT 1`, userID, int(subject)).Find(&record)
	if query.Error != nil {
		return entity.Laboratory{}, query.Error
	}
	if query.RowsAffected == 0 {
		return entity.Laboratory{}, ErrNotFound
	}
	return entity.Laboratory{
		ID:      record.ID,
		Subject: record.Subject,
		Name:    record.Name,
	}, nil
}
