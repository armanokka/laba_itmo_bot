package botdb

import (
	"database/sql"
	"encoding/base64"
	"github.com/armanokka/translobot/internal/tables"
	"time"
)

func (db BotDB) LogUserMessage(id int64, text string) error {
	return db.Create(&tables.UsersLogs{
		ID:      id,
		Intent:  sql.NullString{},
		Text:    base64.StdEncoding.EncodeToString([]byte(text)),
		FromBot: false,
		Date:    time.Now(),
	}).Error
}

func (db BotDB) LogBotMessage(toID int64, intent string, text string) error {
	return db.Create(&tables.UsersLogs{
		ID:      toID,
		Intent:  sql.NullString{
			String: intent,
			Valid:  true,
		},
		Text:    base64.StdEncoding.EncodeToString([]byte(text)),
		FromBot: true,
		Date:    time.Now(),
	}).Error
}

// GetUserLogs returns logs that's field is have base64 encoded text
func (db BotDB) GetUserLogs(id int64, limit int) ([]tables.UsersLogs, error) {
	logs := make([]tables.UsersLogs, 0, limit)
	err := db.Model(&tables.UsersLogs{}).Where("id = ?", id).Order("date ASC").Limit(limit).Find(&logs).Error
	return logs, err
}