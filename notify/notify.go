package notify

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gen2brain/beeep"
)

const (
    envMarkName  = "GOLANG_CLI_REMINDER"
    envMarkValue = "1"
)

// SpawnNotificationProcessNoDB spawns a child process passing reminder time and message
func SpawnNotificationProcessNoDB(reminderTime time.Time, message string) error {
    exePath, err := os.Executable()
    if err != nil {
        return fmt.Errorf("unable to get executable path: %w", err)
    }

    args := []string{envMarkName, reminderTime.Format(time.RFC3339), message}

    cmd := exec.Command(exePath, args...)
    cmd.Env = append(os.Environ(), envMarkName+"="+envMarkValue)

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start notification process: %w", err)
    }
    return nil
}

func RunReminderNoDB(reminderTimeStr string, message string) {
    remindAt, err := time.Parse(time.RFC3339, reminderTimeStr)
    if err != nil {
        fmt.Println("Invalid reminder time format:", err)
        os.Exit(1)
    }

    remaining := time.Until(remindAt)
    
    if remaining <= 0 {
        fmt.Println("Reminder time is in the past or now, sending notification immediately.")
    } else {
        fmt.Printf("Successfully set your reminder\n Time remaining: %s\n", formatDurationHHMMSS(remaining))
        time.Sleep(remaining)
        fmt.Println("Reminder time reached!")
    }

    playBeep()

    iconPath := "assets/information.png"
    if err := beeep.Alert("Reminder", message, iconPath); err != nil {
        fmt.Println("Failed to send notification:", err)
    }
}

// formatDurationHHMMSS formats a time.Duration as HH:MM:SS
func formatDurationHHMMSS(d time.Duration) string {
    d = d.Truncate(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second
    return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}


// playBeep tries to play a simple beep sound cross platform, with fallbacks
func playBeep() {
    // Windows (rundll32 user32.dll,MessageBeep)
    if isWindows() {
        exec.Command("rundll32", "user32.dll,MessageBeep").Run()
        return
    }
    // Mac: use afplay for Ping sound
    if isMac() {
        exec.Command("afplay", "/System/Library/Sounds/Ping.aiff").Run()
        return
    }
    // Linux: try paplay or aplay or terminal bell
    if isLinux() {
        if err := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run(); err == nil {
            return
        }
        if err := exec.Command("aplay", "/usr/share/sounds/alsa/Front_Center.wav").Run(); err == nil {
            return
        }
        // Terminal bell fallback
        fmt.Print("\a")
    }
}

func isWindows() bool {
    return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func isMac() bool {
    // OSTYPE not always set; use runtime.GOOS alternatively if needed
    return os.Getenv("OSTYPE") == "darwin" || os.Getenv("OSTYPE") == "darwin20"
}

func isLinux() bool {
    osname := os.Getenv("OSTYPE")
    return osname == "linux" || osname == "linux-gnu"
}
