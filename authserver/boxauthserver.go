package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	BOX_AUTH_NAME = "'BOX Auth module'"
	BOX_AUTH_PORT = "9981"

	TENANT_HEADER = "X-CASB-TENANT"
)

func main() {

	file, err := os.OpenFile("/var/log/casb/box_authserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("ROOT - Got request.")

		nonce := r.Header.Get("nonce")
		currentTime := time.Now()
		fmt.Fprintf(w, "Hit the root for %s\n\n\n", BOX_AUTH_NAME)
		fmt.Fprintf(w, "- Response from %s running at port %s.\n", BOX_AUTH_NAME, BOX_AUTH_PORT)
		fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
		fmt.Fprintf(w, "- Nonce: %s\n", nonce)
		fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
	})

	//Health Check endpoint
	http.HandleFunc("/api/admin/v1/api-status/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("BOX - HEALTH CHECK - Got request.")
		currentTime := time.Now()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Health check response at %s.", currentTime.Format("01-02-2006 15:04:05"))
	})
	
	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("AUTH - Got request.")

		nonce := r.Header.Get("nonce")
		currentTime := time.Now()
		log.Printf("Header - %s", r.Header)

		var original_uri_header string
		var authorization_header string
		original_uri_header = r.Header.Get("X-Original-URI")
		authorization_header = r.Header.Get("Authorization")
		if len(authorization_header) == 0 {
			log.Print("AUTH - Authorization header is missing.")

			w.WriteHeader(http.StatusUnauthorized)

			fmt.Fprintf(w, "Authorization header is missing")
			fmt.Fprintf(w, "- Response from %s running at port %s.\n", BOX_AUTH_NAME, BOX_AUTH_PORT)
			fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
			fmt.Fprintf(w, "- Nonce: %s\n", nonce)
			fmt.Fprintf(w, "- Headers %s.\n", r.Header)
			fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
			return
		}

		if len(original_uri_header) == 0 {
			log.Print("AUTH - X-Original-URI header is missing.")

			w.WriteHeader(http.StatusUnauthorized)

			fmt.Fprintf(w, " !!! X-Original-URI header is missing - Unauthorized Error - Status code 401 !!!\n\n\n")
			fmt.Fprintf(w, "- Response from %s running at port %s.\n", BOX_AUTH_NAME, BOX_AUTH_PORT)
			fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
			fmt.Fprintf(w, "- Nonce: %s\n", nonce)
			fmt.Fprintf(w, "- Headers: %s.\n", r.Header)
			fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
			return
		} else {
			if strings.Contains(original_uri_header, "box") {
				log.Print("AUTH - Box is OK.")

				GetTenantFromAuthToken(w, r)

				fmt.Fprintf(w, "- Response from %s running at port %s.\n", BOX_AUTH_NAME, BOX_AUTH_PORT)
				fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
				fmt.Fprintf(w, "- Nonce: %s\n", nonce)
				fmt.Fprintf(w, "- Authorization: %s\n", authorization_header)
				//fmt.Fprintf(w, "- Request Headers: %s.\n", r.Header)
				fmt.Fprintf(w, "- Response Headers set: %s= %s.\n", TENANT_HEADER, w.Header().Get(TENANT_HEADER))
				fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
			} else {
				log.Print("AUTH - Non-Box request is NOT OK.")

				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprintf(w, " !!! X-Original-URI is not for this Auth module - Unauthorized Error - Status code 401 !!!\n\n\n")
				fmt.Fprintf(w, "- Response from %s running at port %s.\n", BOX_AUTH_NAME, BOX_AUTH_PORT)
				fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
				fmt.Fprintf(w, "- Nonce: %s\n", nonce)
				fmt.Fprintf(w, "- Headers: %s.\n", r.Header)
				fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))	
			}
		}
	})

	runningPort := fmt.Sprintf(":%s", BOX_AUTH_PORT)
	log.Print("Logging to a file in Go!")
	http.ListenAndServe(runningPort, nil)
}

func GetTenantFromAuthToken(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	log.Printf("Inside Token parsing")
	tokenSplit := strings.Split(auth, " ")
	if len(tokenSplit) > 1 {
		if tokenSplit[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "!!!Invalid box token, Bearer key word missing - 'Bearer box_tenant=gonukona-2343'\n\n")
			log.Printf("!!!Invalid box token, Bearer key word missing - 'Bearer box_tenant=gonukona-2343'\n\n")
		} else {
			tenantSplit := strings.Split(tokenSplit[1], "=")
			if len(tenantSplit) > 1 {
				if tenantSplit[0] == "box_tenant" {
					w.Header().Add(TENANT_HEADER, tenantSplit[1])
					w.WriteHeader(http.StatusOK)
					log.Printf("Validated box token.\n\n")
					return true
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintf(w, "!!!Invalid box token, first item has to be tenant - 'Bearer box_tenant=gonukona-2343'\n\n")
					log.Printf("!!!Invalid box token, first item has to be tenant - 'Bearer box_tenant=gonukona-2343'\n\n")
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "!!!Invalid box token, tenant missing - 'Bearer box_tenant=gonukona-2343'\n\n")
				log.Printf("!!!Invalid box token, tenant missing - 'Bearer box_tenant=gonukona-2343'\n\n")
			}
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "!!!Invalid box token - 'Bearer box_tenant=gonukona-2343'\n\n")
		log.Printf("!!!Invalid box token - 'Bearer box_tenant=gonukona-2343'\n\n")
	}
	return false
}
