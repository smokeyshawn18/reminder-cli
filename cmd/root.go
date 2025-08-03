package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"reminder-cli/notify"
	"reminder-cli/timeparse"

	"github.com/joho/godotenv"
)

const (
    envMarkName  = "GOLANG_CLI_REMINDER"
    envMarkValue = "1"
)

var version = "v1.0.0"

// RunReminderCLI runs the main CLI logic.
// NOTE: Runs countdown & notification inline blocking the CLI process.
func RunReminderCLI(args []string) error {
    // Load .env optionally (not needed if no DB)
    _ = godotenv.Load()

    // If running as child process (the marker is present), run notification and exit
    if os.Getenv(envMarkName) == envMarkValue {
        if len(args) < 4 {
            return errors.New("not enough arguments for notification process (need time and message)")
        }
        reminderTimeStr := args[2]
        message := strings.Join(args[3:], " ")
        notify.RunReminderNoDB(reminderTimeStr, message)
        return nil
    }

    // Normal CLI usage; expect args: <time> <message>
    if len(args) < 3 {
        return fmt.Errorf("usage: %s <time> <message>", args[0])
    }

    reminderTime, err := timeparse.ParseTime(args[1])
    if err != nil {
        return fmt.Errorf("time parsing error: %w", err)
    }

    if time.Now().After(reminderTime) {
        return errors.New("reminder time must be in the future")
    }

    message := strings.Join(args[2:], " ")

    // === OPTION 1: Run inline countdown & notify (blocks CLI) ===

    notify.RunReminderNoDB(reminderTime.Format(time.RFC3339), message)

    // === OPTION 2: Spawn background child process (uncomment if needed) ===
    /*
    if err := notify.SpawnNotificationProcessNoDB(reminderTime, message); err != nil {
        return fmt.Errorf("failed to start notification process: %w", err)
    }

    fmt.Printf("Reminder set for %s. You will be notified then.\n", reminderTime.Format(time.RFC1123))
    */

    return nil
}

func main() {
    if len(os.Args) > 1 && os.Args[1] == "version" {
        fmt.Println("Reminder CLI version:", version)
        return
    }
    if err := RunReminderCLI(os.Args); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}
