package repo

import (
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"github.com/pkg/errors"
	"time"
)

func (t TranslationRepo) GetQueueBySubject(threadID int, labID int) ([]entity.QueueUser, error) {
	rows, err := t.client.Model(&Queues{}).Raw(`
SELECT queues.user_id, queues.checked, queues.retake, queues.passed, users.first_name, users.last_name, users.patronymic
FROM queues
JOIN users ON users.id=queues.user_id WHERE queues.laboratory_id = ? AND queues.thread_id = ? ORDER BY queues.checked DESC, queues.retake ASC, queues.created_at`, labID, threadID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	queue := make([]entity.QueueUser, 0, 5)
	for rows.Next() {
		var user entity.QueueUser
		if err = rows.Scan(&user.UserID, &user.Checked, &user.Retake, &user.Passed, &user.FirstName, &user.LastName, &user.Patronymic); err != nil {
			return nil, err
		}
		queue = append(queue, user)
	}
	return queue, rows.Err()
}

func (t TranslationRepo) AddUserToQueue(userID int64, userThreadID int, labID int) error {
	took, err := t.UserTookLab(userID, labID)
	if err != nil {
		return err
	}
	return t.client.Create(&Queues{
		UserID:       userID,
		ThreadID:     userThreadID,
		LaboratoryID: labID,
		Checked:      false,
		Retake:       took,
		Passed:       false,
		CreatedAt:    time.Now(),
	}).Error
}

func (t TranslationRepo) RemoveUserFromQueue(userID int64, threadID int, labID int) error {
	return t.client.Where("user_id = ?", userID).Where("thread_id = ?", threadID).Where("laboratory_id = ?", labID).Where("checked = ?", false).Delete(&Queues{}).Error
}

func (t TranslationRepo) MarkLabAsNotPassed(studentID int64, lab int) error {
	query := t.client.Model(&Queues{}).Where("user_id = ?", studentID).Where("laboratory_id = ?", lab).Update("checked", true)
	if query.RowsAffected == 0 {
		return errors.Wrap(ErrNotFound, "MarkLabAsChecked")
	}
	return query.Error
}

func (t TranslationRepo) MarkLabAsPassed(studentID int64, lab int) error {
	query := t.client.Model(&Queues{}).Where("user_id = ?", studentID).Where("laboratory_id = ?", lab).Where("checked = ?", false).Updates(map[string]interface{}{"checked": true, "passed": true})
	if query.RowsAffected == 0 {
		return errors.Wrap(ErrNotFound, "MarkLabAsChecked")
	}
	return query.Error
}

func (t TranslationRepo) UserPassedLab(studentID int64, lab int) (exists bool, err error) {
	return exists, t.client.Model(&Queues{}).Raw(`SELECT EXISTS(SELECT 1 FROM queues WHERE user_id = ? AND laboratory_id = ? AND passed=true)`, studentID, lab).Find(&exists).Error
}
func (t TranslationRepo) UserTookLab(studentID int64, lab int) (exists bool, err error) {
	return exists, t.client.Model(&Queues{}).Raw(`SELECT EXISTS(SELECT 1 FROM queues WHERE user_id = ? AND laboratory_id = ? AND checked=true)`, studentID, lab).Find(&exists).Error
}
