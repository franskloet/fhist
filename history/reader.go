package history

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type HistoryStore struct {
	AllItems     []HistoryItem
	BashItems    []HistoryItem
	ZshItems     []HistoryItem
	FishItems    []HistoryItem
	ActiveSource string // "bash", "zsh", "fish", "all"
}

func LoadAllHistory(customPath string, preferredSource string) (*HistoryStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}

	store := &HistoryStore{}

	if customPath != "" {
		items, err := LoadBashHistory(customPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read specified history file: %w", err)
		}
		store.AllItems = items
		store.ActiveSource = "custom"
		return store, nil
	}

	// Detect default shell
	userShell := os.Getenv("SHELL")
	defaultSource := "bash"
	if strings.Contains(userShell, "zsh") {
		defaultSource = "zsh"
	} else if strings.Contains(userShell, "fish") {
		defaultSource = "fish"
	}

	// Load standard history files
	bashPath := filepath.Join(home, ".bash_history")
	if bItems, err := LoadBashHistory(bashPath); err == nil && len(bItems) > 0 {
		store.BashItems = bItems
	}

	zshPath := filepath.Join(home, ".zsh_history")
	if zItems, err := LoadZshHistory(zshPath); err == nil && len(zItems) > 0 {
		store.ZshItems = zItems
	}

	fishPath := filepath.Join(home, ".local", "share", "fish", "fish_history")
	if fItems, err := LoadFishHistory(fishPath); err == nil && len(fItems) > 0 {
		store.FishItems = fItems
	}

	if preferredSource != "" && preferredSource != "auto" {
		store.ActiveSource = preferredSource
	} else {
		store.ActiveSource = defaultSource
	}

	store.RebuildAll()
	return store, nil
}

func (s *HistoryStore) GetCurrentItems() []HistoryItem {
	switch s.ActiveSource {
	case "bash":
		if len(s.BashItems) > 0 {
			return s.BashItems
		}
	case "zsh":
		if len(s.ZshItems) > 0 {
			return s.ZshItems
		}
	case "fish":
		if len(s.FishItems) > 0 {
			return s.FishItems
		}
	}

	// Fallback if specific source is empty
	if len(s.BashItems) > 0 {
		return s.BashItems
	}
	if len(s.ZshItems) > 0 {
		return s.ZshItems
	}
	if len(s.FishItems) > 0 {
		return s.FishItems
	}
	return s.AllItems
}

func (s *HistoryStore) RebuildAll() {
	var merged []HistoryItem
	merged = append(merged, s.BashItems...)
	merged = append(merged, s.ZshItems...)
	merged = append(merged, s.FishItems...)

	// Assign 0-based sequential indexes
	for i := range merged {
		merged[i].Index = i
	}

	s.AllItems = merged
}
