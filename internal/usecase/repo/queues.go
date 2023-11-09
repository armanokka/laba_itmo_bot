package repo

import (
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
	"time"
)

func (t TranslationRepo) GetQueueByThreadID(threadID int) ([]entity.QueueUser, error) {
	rows, err := t.client.Model(&Queues{}).Raw(`
SELECT queues.user_id, laboratories.name, queues.checked, queues.retake, queues.passed, users.first_name, users.last_name, users.patronymic
FROM queues
         JOIN laboratories ON  laboratories.id=queues.laboratory_id
         JOIN users ON users.id=queues.user_id WHERE queues.thread_id = ? AND checked=false ORDER BY queues.checked DESC, queues.retake ASC, queues.created_at`, threadID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	queue := make([]entity.QueueUser, 0, 5)
	for rows.Next() {
		var user entity.QueueUser
		if err = rows.Scan(&user.UserID, &user.LabName, &user.Checked, &user.Retake, &user.Passed, &user.FirstName, &user.LastName, &user.Patronymic); err != nil {
			return nil, err
		}
		queue = append(queue, user)
	}
	return queue, rows.Err()
}

func (t TranslationRepo) AddUserToQueue(userID int64, threadID int, labID int) error {
	retakes, err := t.UserRetakesLab(userID, labID)
	if err != nil {
		return err
	}
	return t.client.Create(&Queues{
		UserID:       userID,
		ThreadID:     threadID,
		LaboratoryID: labID,
		Checked:      false,
		Retake:       retakes,
		Passed:       false,
		CreatedAt:    time.Now(),
	}).Error
}

func (t TranslationRepo) RemoveUserFromQueue(userID int64, threadID int) error {
	return t.client.Where("user_id = ?", userID).Where("thread_id = ?", threadID).Where("checked = ?", false).Delete(&Queues{}).Error
}

func (t TranslationRepo) GradeLab(studentID int64, threadID int, passed bool) (labID int, err error) {
	var record Queues
	query := t.client.Model(&record).Where("user_id = ?", studentID).Where("thread_id = ?", threadID).Where("checked = ?", false).Clauses(clause.Returning{}).Select("laboratory_id").Updates(map[string]interface{}{"checked": true, "passed": passed})
	if query.RowsAffected == 0 {
		return 0, errors.Wrap(ErrNotFound, "MarkLabAsChecked")
	}
	return record.LaboratoryID, query.Error
}

func (t TranslationRepo) UserInQueue(studentID int64, threadID int) (in bool, err error) { // whether user is in any queue of the subject
	query := t.client.Model(&Queues{}).Raw("SELECT EXISTS(SELECT 1 FROM queues WHERE user_id = ? AND checked=false AND thread_id = ?)", studentID, threadID).Find(&in)
	if query.Error != nil {
		return false, err
	}
	if query.RowsAffected == 0 {
		return false, ErrNotFound
	}
	return in, nil
}

func (t TranslationRepo) UserRetakesLab(studentID int64, threadID int) (in bool, err error) { // whether user is in any queue of the subject
	query := t.client.Model(&Queues{}).Raw("SELECT EXISTS(SELECT 1 FROM queues WHERE user_id = ? AND checked=true AND thread_id = ?)", studentID, threadID).Find(&in)
	if query.Error != nil {
		return false, err
	}
	if query.RowsAffected == 0 {
		return false, ErrNotFound
	}
	return in, nil
}
