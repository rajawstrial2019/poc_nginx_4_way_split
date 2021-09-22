package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	BOX_AUTH_NAME = "'BOX Auth module'"
	BOX_AUTH_PORT = "9981"

	BOX_APP_NAME = "'BOX App module'"
	BOX_APP_PORT = "9991"

	GAPPS_AUTH_NAME = "'GAPPS Auth module'"
	GAPPS_AUTH_PORT = "9982"

	GAPPS_APP_NAME = "'GAPPS App module'"
	GAPPS_APP_PORT = "9992"

	TENANT_HEADER = "X-CASB-TENANT"
)

func main() {

	file, err := os.OpenFile("/var/log/casb/gapps_appserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("ROOT - Got request.")

		nonce := r.Header.Get("nonce")
		currentTime := time.Now()
		fmt.Fprintf(w, "Hit the root for %s\n\n\n", GAPPS_AUTH_NAME)
		fmt.Fprintf(w, "- Response from %s running at port %s.\n", GAPPS_AUTH_NAME, GAPPS_APP_PORT)
		fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
		fmt.Fprintf(w, "- Nonce: %s\n", nonce)
		fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
	})

	//Health Check endpoint
	http.HandleFunc("/v2/api/v1/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("GAAPS - HEALTH CHECK - Got request.")
		currentTime := time.Now()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Health check response at %s.", currentTime.Format("01-02-2006 15:04:05"))
	})

	http.HandleFunc("/v2/gapps/notification/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("GAPPS NOTIFICATION - Got request.")

		nonce := r.Header.Get("nonce")
		tenant := r.Header.Get(TENANT_HEADER)
		currentTime := time.Now()

		fmt.Fprintf(w, "Response from %s!! Current Time is: %s\n", GAPPS_APP_NAME, currentTime.Format("01-02-2006 15:04:05"))

		fmt.Fprintf(w, "\n\nSetup:\n")
		fmt.Fprintf(w, "- Nginx running inside container.\n")
		fmt.Fprintf(w, "- Nginx downloaded, compiled and installed with Nginx Auth Module.\n")
		fmt.Fprintf(w, "- Go Lang downloaded and installed.\n")
		fmt.Fprintf(w, "- Nginx Running at port 80.\n")
		fmt.Fprintf(w, "- Go Web apps (Four) compiled and installed.\n")
		fmt.Fprintf(w, "   1. port %s running %s.\n", BOX_AUTH_PORT, BOX_AUTH_NAME)
		fmt.Fprintf(w, "   2. port %s running %s.\n", BOX_APP_PORT, BOX_APP_NAME)
		fmt.Fprintf(w, "   3. port %s running %s.\n", GAPPS_AUTH_PORT, GAPPS_AUTH_NAME)
		fmt.Fprintf(w, "   4. port %s running %s.\n", GAPPS_APP_PORT, GAPPS_APP_NAME)

		fmt.Fprintf(w, "\n\nNginx configure to :\n")
		fmt.Fprintf(w, "- reject calls without 'Authorization' Header. Bad request does not reach Auth Server.\n")
		fmt.Fprintf(w, "- forward incoming request to an Auth Module, based on URI (can be done based on Domain too). \n")
		fmt.Fprintf(w, "   - URI containing '/box/notification' will be routed to Box Auth Server. (https://sedcasb-feature-eoe-gcp-box-notify.casb-sp1-sed-saasdev.elastica-inc.com/box/notification/)\n")
		fmt.Fprintf(w, "   - URI containing '/gapps/notification' will be routed to Gapps Auth Server.\n")
		fmt.Fprintf(w, "- forward Authenticated request to Corresponding App Module.\n")
		fmt.Fprintf(w, "   - Box Auth module is configured to forward Valid request to Box App server.\n")
		fmt.Fprintf(w, "   - Gapps Auth module is configured to forward Valid request to Gapps App server.\n")
		fmt.Fprintf(w, "- forward Custom Header from Auth Module to App module. If Auth module has processed JWT token and extracted information, then it can provide it to App module.\n")

		fmt.Fprintf(w, "\n\nRequest/Response:\n")
		fmt.Fprintf(w, "- Tenant Extracted from Authorization header in Auth Module: %s\n\n", tenant)
		fmt.Fprintf(w, "- X-Original-URI:%s\n", r.Header.Get("X-Original-URI"))
		fmt.Fprintf(w, "- X-Real-IP:%s\n", r.Header.Get("X-Real-IP"))
		fmt.Fprintf(w, "- Host: %s%s\n", r.Host, r.URL.Path)
		fmt.Fprintf(w, "- Nonce: %s\n\n\n", nonce)
		//fmt.Fprintf(w, "- Headers %s.\n", r.Header)
	})

	runningPort := fmt.Sprintf(":%s", GAPPS_APP_PORT)
	log.Print("Logging to a file in Go!")
	http.ListenAndServe(runningPort, nil)
}
