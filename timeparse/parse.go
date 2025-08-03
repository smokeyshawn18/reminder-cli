package timeparse

import (
	"fmt"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

// ParseTime parses natural language time strings into time.Time
func ParseTime(input string) (time.Time, error) {
    now := time.Now()
    w := when.New(nil)
    w.Add(en.All...)
    w.Add(common.All...)

    result, err := w.Parse(input, now)
    if err != nil {
        return time.Time{}, fmt.Errorf("parse error: %w", err)
    }
    if result == nil {
        return time.Time{}, fmt.Errorf("could not parse time: %q", input)
    }

    return result.Time, nil
}
