package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func main() {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "app"
	}

	// Connect to Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379" // Default k8s service name
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Hello from %s (%s)\n", appName, os.Getenv("HOSTNAME"))
	})

	http.HandleFunc("/count", func(w http.ResponseWriter, _ *http.Request) {
		count, err := rdb.Incr(ctx, "counter").Result()
		if err != nil {
			http.Error(w, "Redis error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Count: %d\n", count)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if os.Getenv("APP_READY") != "true" {
			w.WriteHeader(503)
			return
		}
		// Also check Redis connection
		if err := rdb.Ping(ctx).Err(); err != nil {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, _ *http.Request) {
		dbPass := os.Getenv("DB_PASSWORD")
		masked := "not set"
		if len(dbPass) > 0 {
			masked = dbPass[:3] + "***"
		}
		fmt.Fprintf(w, "APP_NAME: %s\nAPP_READY: %s\nDB_PASSWORD: %s\nREDIS_ADDR: %s\n",
			os.Getenv("APP_NAME"),
			os.Getenv("APP_READY"),
			masked,
			redisAddr)
	})

	http.ListenAndServe(":8080", nil)
}
