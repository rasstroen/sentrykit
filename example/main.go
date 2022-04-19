package main

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/rasstroen/sentrykit"
	"log"
	"os"
	"time"
)

func main() {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         os.Getenv("SENTRY_DSN"),
		Environment: "dev",
	})
	if err != nil {
		panic(err)
	}

	logger := sentrykit.NewSentryLogger(client)

	logger.Log("msg", "User updated", "tries", 1, "properties", map[string]string{"username": "gopher"})
	logger.Log("err", errors.New("test error"))

	if !client.Flush(time.Second * 5) {
		log.Fatal("flush timeout")
	}
}
