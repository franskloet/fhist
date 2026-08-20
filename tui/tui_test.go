package tui

import (
	"fmt"
	"strings"
	"testing"

	"fhist/history"
)

func TestRenderView(t *testing.T) {
	store := &history.HistoryStore{
		BashItems: []history.HistoryItem{
			{Index: 0, Command: "ls -la", Source: "bash"},
			{Index: 1, Command: "cd /tmp", Source: "bash"},
			{Index: 2, Command: "git status", Source: "bash"},
		},
		ActiveSource: "bash",
	}
	store.RebuildAll()

	m := NewModel(store, "", 5)
	m.Width = 80
	m.Height = 24

	viewOutput := m.View()
	fmt.Println("=== RENDERED VIEW OUTPUT (80x24) ===")
	fmt.Println(viewOutput)
	fmt.Println("=== END VIEW OUTPUT ===")

	lines := strings.Split(viewOutput, "\n")
	fmt.Printf("Total lines outputted: %d (expected <= 24)\n", len(lines))
}
