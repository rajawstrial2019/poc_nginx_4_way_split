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
	OFFICE_AUTH_NAME = "'OFFICE Auth module'"
	OFFICE_AUTH_PORT = "9982"

	TENANT_HEADER = "X-CASB-TENANT"
)

func main() {

	file, err := os.OpenFile("/var/log/casb/office_authserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("ROOT - Got request.")

		nonce := r.Header.Get("nonce")
		currentTime := time.Now()
		fmt.Fprintf(w, "Hit the root for %s\n\n\n", OFFICE_AUTH_NAME)
		fmt.Fprintf(w, "- Response from %s running at port %s.\n", OFFICE_AUTH_NAME, OFFICE_AUTH_PORT)
		fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
		fmt.Fprintf(w, "- Nonce: %s\n", nonce)
		fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
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
			fmt.Fprintf(w, "- Response from %s running at port %s.\n", OFFICE_AUTH_NAME, OFFICE_AUTH_PORT)
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
			fmt.Fprintf(w, "- Response from %s running at port %s.\n", OFFICE_AUTH_NAME, OFFICE_AUTH_PORT)
			fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
			fmt.Fprintf(w, "- Nonce: %s\n", nonce)
			//fmt.Fprintf(w, "- Headers: %s.\n", r.Header)
			fmt.Fprintf(w, "- Response Headers set: %s= %s.\n", TENANT_HEADER, w.Header().Get(TENANT_HEADER))
			fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
			return
		} else {
			if strings.Contains(original_uri_header, "office") {
				log.Print("AUTH - Office is OK.")

				success := GetTenantFromAuthToken(w, r)
				if success {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
				}

				fmt.Fprintf(w, "- Response from %s running at port %s.\n", OFFICE_AUTH_NAME, OFFICE_AUTH_PORT)
				fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
				fmt.Fprintf(w, "- Nonce: %s\n", nonce)
				fmt.Fprintf(w, "- Authorization: %s\n", authorization_header)
				//fmt.Fprintf(w, "- Request Headers: %s.\n", r.Header)
				fmt.Fprintf(w, "- Response Headers set: %s= %s.\n", TENANT_HEADER, w.Header().Get(TENANT_HEADER))
				fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
			} else {
				log.Print("AUTH - Non-Office request is NOT OK.")

				w.WriteHeader(http.StatusUnauthorized)

				fmt.Fprintf(w, " !!! X-Original-URI is not for this Auth module - Unauthorized Error - Status code 401 !!!\n\n\n")
				fmt.Fprintf(w, "- Response from %s running at port %s.\n", OFFICE_AUTH_NAME, OFFICE_AUTH_PORT)
				fmt.Fprintf(w, "- Host: %s URL:%s\n", r.Host, r.URL.Path)
				fmt.Fprintf(w, "- Nonce: %s\n", nonce)
				fmt.Fprintf(w, "- Headers: %s.\n", r.Header)
				fmt.Fprintf(w, "- Current Time is: %s.", currentTime.Format("01-02-2006 15:04:05"))
			}
		}
	})

	runningPort := fmt.Sprintf(":%s", OFFICE_AUTH_PORT)
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
			fmt.Fprintf(w, "!!!Invalid office token, Bearer key word missing - 'Bearer office_tenant=acme-3648'\n\n")
			log.Printf("!!!Invalid office token, Bearer key word missing - 'Bearer office_tenant=acme-3648'\n\n")
		} else {
			tenantSplit := strings.Split(tokenSplit[1], "=")
			if len(tenantSplit) > 1 {
				if tenantSplit[0] == "office_tenant" {
					w.Header().Set(TENANT_HEADER, tenantSplit[1])
					w.WriteHeader(http.StatusOK)
					log.Printf("Validated box token.\n\n")
					return true
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintf(w, "!!!Invalid office token, first item has to be tenant - 'Bearer office_tenant=acme-3648'\n\n")
					log.Printf("!!!Invalid office token, first item has to be tenant - 'Bearer office_tenant=acme-3648'\n\n")
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "!!!Invalid office token, tenant missing - 'Bearer office_tenant=acme-3648'\n\n")
				log.Printf("!!!Invalid office token, tenant missing - 'Bearer office_tenant=acme-3648'\n\n")
			}
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "!!!Invalid office token - 'Bearer office_tenant=acme-3648'\n\n")
		log.Printf("!!!Invalid office token - 'Bearer office_tenant=acme-3648'\n\n")
	}
	return false
}
