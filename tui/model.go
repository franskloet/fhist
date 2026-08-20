package tui

import (
	"regexp"
	"strings"

	"fhist/history"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

type SearchMode int

const (
	ModeFuzzy SearchMode = iota
	ModeContains
	ModePrefix
	ModeRegex
)

func (m SearchMode) String() string {
	switch m {
	case ModeFuzzy:
		return "Fuzzy"
	case ModeContains:
		return "Contains"
	case ModePrefix:
		return "Prefix"
	case ModeRegex:
		return "Regex"
	default:
		return "Fuzzy"
	}
}

type Model struct {
	Store              *history.HistoryStore
	CurrentItems       []history.HistoryItem
	FilteredIndexes    []int
	SelectedIndex      int
	ContextRadius      int // N commands before and after hit
	SearchMode         SearchMode
	TextInput          textinput.Model
	Styles             Styles
	Width              int
	Height             int
	SelectedCommand    string
	PrintContextOnExit bool
	CopiedNotification string
	Err                error
}

func NewModel(store *history.HistoryStore, initialQuery string, contextRadius int) Model {
	ti := textinput.New()
	ti.Placeholder = "Type search query..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 35
	if initialQuery != "" {
		ti.SetValue(initialQuery)
	}

	if contextRadius <= 0 {
		contextRadius = 5
	}

	m := Model{
		Store:         store,
		ContextRadius: contextRadius,
		SearchMode:    ModeFuzzy,
		TextInput:     ti,
		Styles:        DefaultStyles(),
		SelectedIndex: 0,
	}

	m.UpdateCurrentItems()
	m.ApplyFilter()

	return m
}

func (m *Model) UpdateCurrentItems() {
	m.CurrentItems = m.Store.GetCurrentItems()
}

func (m *Model) ApplyFilter() {
	query := strings.TrimSpace(m.TextInput.Value())
	if query == "" {
		// Return all items in reverse chronological order (most recent first)
		m.FilteredIndexes = make([]int, len(m.CurrentItems))
		n := len(m.CurrentItems)
		for i := 0; i < n; i++ {
			m.FilteredIndexes[i] = n - 1 - i
		}
		if m.SelectedIndex >= len(m.FilteredIndexes) {
			m.SelectedIndex = max(0, len(m.FilteredIndexes)-1)
		}
		return
	}

	var results []int
	n := len(m.CurrentItems)

	switch m.SearchMode {
	case ModeFuzzy:
		// Convert items to string slice
		strList := make([]string, n)
		for i, item := range m.CurrentItems {
			strList[i] = item.Command
		}
		matches := fuzzy.Find(query, strList)
		// Reverse fuzzy matches to prioritize recent matches
		for i := len(matches) - 1; i >= 0; i-- {
			results = append(results, matches[i].Index)
		}

	case ModeContains:
		q := strings.ToLower(query)
		for i := n - 1; i >= 0; i-- {
			if strings.Contains(strings.ToLower(m.CurrentItems[i].Command), q) {
				results = append(results, i)
			}
		}

	case ModePrefix:
		q := strings.ToLower(query)
		for i := n - 1; i >= 0; i-- {
			if strings.HasPrefix(strings.ToLower(m.CurrentItems[i].Command), q) {
				results = append(results, i)
			}
		}

	case ModeRegex:
		re, err := regexp.Compile("(?i)" + query)
		if err == nil {
			for i := n - 1; i >= 0; i-- {
				if re.MatchString(m.CurrentItems[i].Command) {
					results = append(results, i)
				}
			}
		}
	}

	m.FilteredIndexes = results
	if m.SelectedIndex >= len(m.FilteredIndexes) {
		m.SelectedIndex = max(0, len(m.FilteredIndexes)-1)
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
