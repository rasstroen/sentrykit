# Sentry Kit


Sentry Kit - is go-kit-like logger. It can handles your error and shows in sentry.

# Get Started

```go
package main

import (
	"errors"
	"github.com/fr05t1k/sentrykit"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"time"
)

func main() {
	client, _ := sentry.NewClient(sentry.ClientOptions{
		Dsn:         os.Getenv("SENTRY_DSN"),
		Environment: "dev",
	})

	logger := sentrykit.NewSentryLogger(client)

	// msg will be used as Event Title in Sentry
	logger.Log("msg", "User updated", "tries", 1, "properties", map[string]string{"username": "gopher"})
	// if you provide err field it should be error type to get calstack for sentry
	logger.Log("err", errors.New("test error"))

	client.Flush(time.Second * 5)
}

```

If you specify `err` field sentry considers as an exception and shows your callstack etc. Otherwise sentry triggers a message and gets the message from `msg` field.

![image](https://user-images.githubusercontent.com/2131624/68679460-94d5d200-0568-11ea-988e-dbcfee0fcacb.png)
