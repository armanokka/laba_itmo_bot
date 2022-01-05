package botdb

import "time"

func (db BotDB) UpdateUserLastActivity(id int64) error {
	return db.UpdateUserByMap(id, map[string]interface{}{"last_activity": time.Now()})
}

func (db BotDB) IncreaseUserUsings(id int64) error {
	return db.UpdateUserByMap(id, map[string]interface{}{"usings": "usings+1"})
}