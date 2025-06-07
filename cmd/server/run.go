package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func startWebServer(bind string, port int) {
	addr := fmt.Sprintf("%s:%d", bind, port)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	log.Printf("[HTTP] Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[HTTP] Failed to start server: %v", err)
	}
}

func main() {
	bind := flag.String("bind", getEnv("BIND", "0.0.0.0"), "bind address")
	port := flag.Int("port", getEnvInt("PORT", 8080), "port number")
	redisAddr := flag.String("redis", getEnv("REDIS_ADDR", "localhost:6379"), "redis address")
	flag.Parse()

	redisClient = redis.NewClient(&redis.Options{Addr: *redisAddr})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("[Redis] Connection failed: %v", err)
	}

	go startWebServer(*bind, *port)

	// Placeholder for daemon loop
	log.Println("[main] MCPD CLI Daemon running...")
	for {
		time.Sleep(1 * time.Hour)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		var out int
		_, err := fmt.Sscanf(val, "%d", &out)
		if err == nil {
			return out
		}
	}
	return fallback
}
