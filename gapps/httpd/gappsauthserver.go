package main

import (
	"fmt"
	"os"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/httpd/api"
	"CASB-Securlet/securlets/gapps/notification-be/authserver/model"
	"github.com/sirupsen/logrus"
)

const (
	AUTH_NAME = "'GAPPS Auth module'"
	AUTH_PORT = "9982"
)

var Logger = logrus.StandardLogger()

func main() {
	var appCtx = model.AppContext{
		Name:    AUTH_NAME,
		Port:    AUTH_PORT,
	}

	file, err := os.OpenFile("/var/log/casb/gapps_authserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Fatal(err)
	}

	defer file.Close()

	Logger.SetOutput(file)
	appCtx.Logger = logrus.NewEntry(Logger)


	server, err := api.NewServer(&appCtx)
	if err != nil {
		return
	}

	if err := server.Start(":" + AUTH_PORT); err != nil {
		fmt.Println("Fatal Error - Server can't start.", err)
	}
}
