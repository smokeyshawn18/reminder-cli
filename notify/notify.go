package notify

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"reminder-cli/db"

	"github.com/gen2brain/beeep"
)

const (
    envMarkName  = "GOLANG_CLI_REMINDER"
    envMarkValue = "1"
)

// SpawnNotificationProcess launches a child process to run RunReminder with the reminder ID
func SpawnNotificationProcess(reminderID int) error {
    exePath, err := os.Executable()
    if err != nil {
        return fmt.Errorf("unable to get executable path: %w", err)
    }

    // Pass the environment variable and reminder ID as argument
    cmd := exec.Command(exePath, envMarkName, strconv.Itoa(reminderID))
    cmd.Env = append(os.Environ(), envMarkName+"="+envMarkValue)

    // Start the process asynchronously
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start notification process: %w", err)
    }
    return nil
}

// RunReminder is run by child process; waits until reminder time and sends notification
func RunReminder(idStr string) {
    // Fetch reminder from DB by ID
    reminder, err := db.GetReminderByID(idStr)
    if err != nil {
        fmt.Println("Failed to fetch reminder:", err)
        os.Exit(1)
    }

    // Sleep until reminder time
    waitDuration := time.Until(reminder.RemindAt)
    if waitDuration > 0 {
        time.Sleep(waitDuration)
    }

    // Send desktop notification using beeep
    iconPath := "assets/information.png" // Optional icon file
    if err := beeep.Alert("Reminder", reminder.Message, iconPath); err != nil {
        fmt.Println("Failed to send notification:", err)
    }

    // Mark reminder as notified in DB to avoid re-notifying
    err = db.MarkNotified(reminder.ID)
    if err != nil {
        fmt.Println("Warning: failed to update reminder as notified:", err)
    }
}
