package history

import (
	"fmt"
	"time"
)

type HistoryItem struct {
	Index     int       `json:"index"`     // 0-based absolute chronological position
	Command   string    `json:"command"`   // Raw command string
	Timestamp time.Time `json:"timestamp"` // Timestamp if available
	CWD       string    `json:"cwd"`       // Working directory if available
	ExitCode  int       `json:"exit_code"` // Exit code (-1 if unknown)
	Host      string    `json:"host"`      // Hostname if available
	User      string    `json:"user"`      // User if available
	Source    string    `json:"source"`    // "bash", "zsh", "fish"
}

func (h HistoryItem) DisplayTime() string {
	if h.Timestamp.IsZero() {
		return ""
	}
	return h.Timestamp.Format("2006-01-02 15:04:05")
}

func (h HistoryItem) ShortTime() string {
	if h.Timestamp.IsZero() {
		return fmt.Sprintf("#%-5d", h.Index+1)
	}
	now := time.Now()
	if h.Timestamp.Year() == now.Year() {
		return h.Timestamp.Format("01-02 15:04")
	}
	return h.Timestamp.Format("2006-01-02 15:04")
}

func (h HistoryItem) FormatSourceBadge() string {
	switch h.Source {
	case "bash":
		return "[bash]"
	case "zsh":
		return "[zsh]"
	case "fish":
		return "[fish]"
	default:
		return fmt.Sprintf("[%s]", h.Source)
	}
}
