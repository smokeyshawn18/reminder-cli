package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"reminder-cli/db"
	"reminder-cli/device"
	"reminder-cli/notify"
	"reminder-cli/timeparse"

	"github.com/joho/godotenv"
)


const (
    envMarkName  = "GOLANG_CLI_REMINDER"
    envMarkValue = "1"
)

var version = "v1.0.0"

// RunReminderCLI runs the main CLI logic with given args (like os.Args).
// Returns error instead of exiting, so can be tested and benchmarked.
func RunReminderCLI(args []string) error {
    // Load .env but not fatal on error
    if err := godotenv.Load(); err != nil {
        // Warn, but continue
        fmt.Println("Warning: .env file not loaded - ensure DATABASE_URL is set in environment")
    }

    if err := db.Connect(); err != nil {
        return fmt.Errorf("DB connection error: %w", err)
    }
    defer db.Close()

    if err := db.InitSchema(); err != nil {
        return fmt.Errorf("DB schema initialization error: %w", err)
    }

    if os.Getenv(envMarkName) == envMarkValue {
        if len(args) < 3 {
            return errors.New("no reminder ID specified for notification process")
        }
        reminderID := args[2]
        notify.RunReminder(reminderID)
        return nil
    }

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

    hostname, osName, arch, username, localIPs, publicIP, err := device.GetDeviceInfo()
    if err != nil {
        fmt.Println("Warning: could not get device info:", err)
    }

    id, err := db.SaveReminder(reminderTime, message, publicIP, hostname, osName, arch, username, localIPs)
    if err != nil {
        return fmt.Errorf("failed to save reminder: %w", err)
    }

    if err := notify.SpawnNotificationProcess(id); err != nil {
        return fmt.Errorf("failed to start notification process: %w", err)
    }

    fmt.Printf("Reminder set for %s. You will be notified then.\n", reminderTime.Format(time.RFC1123))
    return nil
}

func main() {
	    if len(os.Args) > 1 && os.Args[1] == "version" {
        fmt.Println("Reminder CLI version:", version)
        return
    }
    if err := RunReminderCLI(os.Args); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
