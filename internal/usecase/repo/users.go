package repo

import (
	"errors"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"gorm.io/gorm"
)

func (t TranslationRepo) GetUserByID(id int64) (entity.User, error) {
	var user Users
	query := t.client.Where("id = ?", id).Find(&user)
	if query.Error != nil {
		return entity.User{}, query.Error
	}
	if query.RowsAffected == 0 {
		return entity.User{}, gorm.ErrRecordNotFound
	}
	return entity.User{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Patronymic: user.Patronymic,
		//Group:             user.Group,
		OPDThreadID:         user.OPDThreadID,
		ProgrammingThreadID: user.ProgrammingThreadID,
		ITThreadID:          user.ITThreadID,
		Action:              user.Action,
		TeacherSubject:      user.TeacherSubject,
	}, nil
}

func (t TranslationRepo) CreateUser(user entity.User) error {
	return t.client.Create(&Users{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Patronymic: user.Patronymic,
		//Group:             user.Group,
		OPDThreadID:         user.OPDThreadID,
		ProgrammingThreadID: user.ProgrammingThreadID,
		ITThreadID:          user.ITThreadID,
		Action:              user.Action,
	}).Error
}

func (t TranslationRepo) UpdateUserByID(id int64, columnValue ...interface{}) error {
	if len(columnValue)%2 != 0 {
		return errors.New("UpdateUserByID: not even number of columns")
	}
	if len(columnValue) == 0 {
		return errors.New("UpdateUserByID: empty columnValue")
	}
	updates := make(map[string]interface{}, len(columnValue)/2)
	for i := 0; i < len(columnValue); i += 2 {
		column, ok := columnValue[i].(string)
		if !ok {
			return errors.New("UpdateUserByID: column must be a string")
		}
		updates[column] = columnValue[i+1]
	}
	return t.client.Where("id = ?", id).Model(&Users{}).Updates(updates).Error
}

func (t TranslationRepo) GetAllUsersIDs() ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (t TranslationRepo) GetUserPassedLaboratoriesIDs(userID int64, userThread int) ([]int, error) {
	labs := make([]int, 0, 3)
	return labs, t.client.Model(&Queues{}).Where("user_id = ?", userID).Where("passed = ?", true).Select("laboratory_id").Find(&labs).Error
}
