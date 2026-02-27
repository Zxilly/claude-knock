package main

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Zxilly/claude-knock/internal/hook"
	"github.com/Zxilly/claude-knock/internal/notify"
)

func main() {
	// -Embedding mode: Windows relaunches us to handle notification click.
	for _, arg := range os.Args[1:] {
		if strings.EqualFold(arg, "-Embedding") {
			notify.HandleCOMActivation(30 * time.Second)
			return
		}
	}

	// Hook mode: read stdin, send notification, exit immediately.
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read stdin: %v", err)
	}

	input, err := hook.Parse(data)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	title, message := input.FormatNotification()

	if err := notify.Send(title, message); err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}
}
