package main

import (
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func runLogger(chanLogs chan UsersLogs, stop chan os.Signal, interval time.Duration) {
	logrus.Info("Message logger has been started")
	inserts := make([]UsersLogs, 0, cap(logs))
	ticker := time.NewTicker(interval)
	for range ticker.C {
		if len(stop) > 0 {
			close(chanLogs)
			for log := range chanLogs {
				inserts = append(inserts, log)
			}
			if err := db.Create(inserts).Error; err != nil {
				WarnAdmin(err)
			}
			logrus.Info("Message logger was stopped.")
			return
		}
		for i := 0; i < len(chanLogs); i++ {
			inserts = append(inserts, <-chanLogs)
		}
		if err := db.Create(inserts).Error; err != nil {
			WarnAdmin(err)
			break
		}
		inserts = nil
		pp.Println("logs were saved")
	}
}
