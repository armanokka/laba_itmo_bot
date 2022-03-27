package botdb

import (
	"database/sql"
	"github.com/armanokka/translobot/internal/tables"
)

func (db BotDB) MailingExists() (bool, error) {
	var exists bool
	err := db.Model(&tables.Users{}).Raw(`SELECT EXISTS (SELECT 1
FROM information_schema.columns
WHERE table_name = ? AND column_name=?)`, "mailing", "id").Find(&exists).Error
	return exists, err
}

func (db BotDB) DeleteMailuser(id int64) error {
	return db.Model(&tables.Mailing{}).Where("id = ?", id).Delete(&tables.Mailing{}).Error
}

func (db BotDB) DropMailings() error {
	return db.Model(&tables.Users{}).Exec("DROP TABLE IF EXISTS mailing").Error
}

func (db BotDB) CreateMailingTable() error {
	return db.Model(&tables.Mailing{}).Exec(`CREATE TABLE mailing AS (SELECT id FROM users)`).Error
}

func (db BotDB) GetMailersRows() (rows *sql.Rows, err error) {
	return db.Model(&tables.Mailing{}).Select("id").Rows()
}
