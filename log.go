package main

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func runLogger(chanLogs chan UsersLogs, stop chan os.Signal, interval time.Duration) {
	logrus.Info("Message logger has been started")
	inserts := make([]UsersLogs, 0, cap(logs))
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-stop:
			close(chanLogs)
			for log := range chanLogs {
				inserts = append(inserts, log)
			}
			if err := db.Create(inserts).Error; err != nil {
				WarnAdmin(err)
			}
			logrus.Info("Message logger was stopped.")
			break
		case <-ticker.C:
			for i := 0; i < len(chanLogs); i++ {
				inserts = append(inserts, <-chanLogs)
			}
			if err := db.Create(inserts).Error; err != nil {
				WarnAdmin(err)

				for _, insert := range inserts {
					user := NewUser(insert.ID, func(err error) {
						WarnAdmin(err)
					})
					if !user.Exists() {
						if err = db.Create(&Users{
							ID:         insert.ID,
							MyLang:     "en",
							ToLang:     "es",
							Act:        sql.NullString{},
							Mailing:    true,
							Usings:     0,
							Lang:       "en",
							ReferrerID: 0,
							Blocked:    false,
						}).Error; err != nil {
							WarnAdmin(err)
						}
					}
				}
			}
			inserts = nil
			logrus.Info("message logs were saved")
		}
	}
}
