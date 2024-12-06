package logger

import log "github.com/sirupsen/logrus"

func InitLogger() {
	log.SetLevel(log.InfoLevel)
}
