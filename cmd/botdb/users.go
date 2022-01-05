package botdb

import (
	"github.com/armanokka/translobot/internal/tables"
	"gorm.io/gorm"
)

type BotDB struct {
	*gorm.DB
}

func New(conn *gorm.DB) BotDB {
	return BotDB{conn}
}

func (db BotDB) GetUserByID(id int64) (tables.Users, error) {
	var user tables.Users
	req := db.Model(&tables.Users{}).Where("id = ?", id).Find(&user)
	if req.RowsAffected == 0 {
		return tables.Users{}, gorm.ErrRecordNotFound
	}
	return user, req.Error
}

func (db BotDB) CreateUser(user tables.Users) (err error) {
	return db.Create(&user).Error
}

func (db BotDB) UpdateUser(id int64, updates tables.Users) (error) {
	return db.Model(&tables.Users{}).Where("id = ?", id).Updates(updates).Error
}

func (db BotDB) UpdateUserByMap(id int64, updates map[string]interface{}) error {
	return db.Model(&tables.Users{}).Where("id = ?", id).Updates(updates).Error
}

func (db BotDB) GetAllUsers() (users []tables.Users, err error) {
	err = db.Model(&tables.Users{}).Find(&users).Error
	return users, err
}