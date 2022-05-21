package botdb

import (
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type BotDB struct {
	*gorm.DB
}

func New(conn *gorm.DB) BotDB {
	return BotDB{conn}
}

func (db BotDB) GetRandomUser() (users tables.Users, err error) {
	err = db.Model(&tables.Users{}).Order("RAND()").Limit(1).Find(&users).Error
	return users, err
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
	return db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Create(&user).Error
}

func (db BotDB) UpdateUser(id int64, updates tables.Users) error {
	return db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Model(&tables.Users{}).Where("id = ?", id).Updates(updates).Error
}

func (db BotDB) UpdateUserByMap(id int64, updates map[string]interface{}) error {
	return db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Model(&tables.Users{}).Where("id = ?", id).Updates(updates).Error
}

func (db BotDB) GetAllUsers() (users []tables.Users, err error) {
	err = db.Model(&tables.Users{}).Find(&users).Error
	return users, err
}

func (db BotDB) UpdateUserMetrics(id int64, message string) error {
	if err := db.Model(&tables.Users{}).Exec("UPDATE users SET usings=usings+1, last_activity=? WHERE id=?", time.Now(), id).Error; err != nil {
		return err
	}

	return db.LogUserMessage(id, message)
}

func (db BotDB) GetUsersNumber() (num int64, err error) {
	err = db.Model(&tables.Users{}).Raw("SELECT COUNT(*) FROM users").Find(&num).Error
	return
}

func (db BotDB) GetUsersSlice(offset, count int64, slice []int64) (err error) {
	err = db.Model(&tables.Users{}).Raw("SELECT id FROM users OFFSET ? LIMIT ? ORDER BY id DESC", offset, count).Find(&slice).Error
	return
}

var ErrNoRowsAffected = fmt.Errorf("no rows affected")

func (db BotDB) SwapLangs(userID int64) error {
	query := db.Model(&tables.Users{}).Exec("UPDATE users SET my_lang=(@temp:=my_lang), my_lang = to_lang, to_lang = @temp WHERE id = ?", userID)
	if query.RowsAffected != 1 {
		return ErrNoRowsAffected
	}
	return query.Error
}
