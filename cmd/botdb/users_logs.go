package botdb

import (
	"database/sql"
	"github.com/armanokka/translobot/internal/tables"
	"time"
)

func (db BotDB) LogUserMessage(id int64, text string) error {
	return db.Create(&tables.UsersLogs{
		ID:      id,
		Intent:  sql.NullString{},
		Text:    text,
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
		Text:    text,
		FromBot: true,
		Date:    time.Now(),
	}).Error
}