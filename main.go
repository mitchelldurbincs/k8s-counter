package main

import (
	"fmt"
	"net/http"
	"os"
)

var counter int

func main() {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "app"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Hello from %s (%s)\n", appName, os.Getenv("HOSTNAME"))
	})

	http.HandleFunc("/count", func(w http.ResponseWriter, _ *http.Request) {
		counter++
		fmt.Fprintf(w, "Count: %d\n", counter)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if os.Getenv("APP_READY") != "true" {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, _ *http.Request) {
		dbPass := os.Getenv("DB_PASSWORD")
		masked := "not set"
		if len(dbPass) > 0 {
			masked = dbPass[:3] + "***" // Show first 3 chars only
		}
		fmt.Fprintf(w, "APP_NAME: %s\nAPP_READY: %s\nDB_PASSWORD: %s\n",
			os.Getenv("APP_NAME"),
			os.Getenv("APP_READY"),
			masked)
	})

	http.ListenAndServe(":8080", nil)
}
