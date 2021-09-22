package model

import log "github.com/sirupsen/logrus"

type AppContext struct {
	Name string
	Port string
	Logger    *log.Entry     //should be used only for debug logging
}
