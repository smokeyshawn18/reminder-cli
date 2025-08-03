package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql compatibility
)

var db *sql.DB

// Reminder struct models a reminder in the database
type Reminder struct {
    ID        int
    RemindAt  time.Time
    Message   string
    Notified  bool
    UserIP    string
    Hostname  string
    OSName    string
    Arch      string
    Username  string
    LocalIPs  []string
    CreatedAt time.Time
}

// Connect opens a connection to the Postgres DB using DATABASE_URL env var
func Connect() error {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return fmt.Errorf("DATABASE_URL environment variable not set")
    }
    var err error
    db, err = sql.Open("pgx", dbURL)
    if err != nil {
        return fmt.Errorf("unable to open DB connection: %w", err)
    }
    if err = db.Ping(); err != nil {
        return fmt.Errorf("unable to ping DB: %w", err)
    }
    return nil
}

// Close closes the DB connection pool
func Close() {
    if db != nil {
        db.Close()
    }
}

// InitSchema creates the reminders table if it does not exist
func InitSchema() error {
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS reminders (
        id SERIAL PRIMARY KEY,
        remind_at TIMESTAMPTZ NOT NULL,
        message TEXT NOT NULL,
        notified BOOLEAN NOT NULL DEFAULT FALSE,
        user_ip TEXT,
        hostname TEXT,
        os_name TEXT,
        arch TEXT,
        username TEXT,
        local_ips TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );`
    _, err := db.Exec(sqlStmt)
    return err
}

// SaveReminder inserts a new reminder with device info and returns newly created ID
func SaveReminder(remindAt time.Time, message, userIP, hostname, osName, arch, username string, localIPs []string) (int, error) {
    localIPsStr := strings.Join(localIPs, ",")
    var id int

    err := db.QueryRow(
        `INSERT INTO reminders(remind_at, message, user_ip, hostname, os_name, arch, username, local_ips)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
        remindAt, message, userIP, hostname, osName, arch, username, localIPsStr,
    ).Scan(&id)

    if err != nil {
        return 0, fmt.Errorf("failed to insert reminder: %w", err)
    }
    return id, nil
}

// GetReminderByID fetches a Reminder by its ID
func GetReminderByID(id string) (*Reminder, error) {
    r := &Reminder{}
    var localIPsStr string
    err := db.QueryRow(
        `SELECT id, remind_at, message, notified, user_ip, hostname, os_name, arch, username, local_ips, created_at 
         FROM reminders WHERE id=$1`, id).
        Scan(&r.ID, &r.RemindAt, &r.Message, &r.Notified,
            &r.UserIP, &r.Hostname, &r.OSName, &r.Arch, &r.Username, &localIPsStr, &r.CreatedAt)
    if err != nil {
        return nil, err
    }
    if localIPsStr != "" {
        r.LocalIPs = strings.Split(localIPsStr, ",")
    }
    return r, nil
}

// MarkNotified sets the notified field for a reminder so it won’t trigger again
func MarkNotified(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := db.ExecContext(ctx, `UPDATE reminders SET notified = TRUE WHERE id = $1`, id)
    return err
}
