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
	ModeContains SearchMode = iota // hstr-style tokenized substring search (Default)
	ModeFuzzy
	ModePrefix
	ModeRegex
)

func (m SearchMode) String() string {
	switch m {
	case ModeContains:
		return "Substring"
	case ModeFuzzy:
		return "Fuzzy"
	case ModePrefix:
		return "Prefix"
	case ModeRegex:
		return "Regex"
	default:
		return "Substring"
	}
}

type Model struct {
	Store              *history.HistoryStore
	CurrentItems       []history.HistoryItem
	FilteredIndexes    []int
	SelectedIndex      int
	ContextRadius      int // N commands before and after hit
	CommandOffset      int // Horizontal scroll offset for long command lines
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
		CommandOffset: 0,
		SearchMode:    ModeContains,
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
	case ModeContains:
		queryLower := strings.ToLower(query)
		tokens := strings.Fields(queryLower)

		for i := n - 1; i >= 0; i-- {
			cmdLower := strings.ToLower(m.CurrentItems[i].Command)
			matchAll := true
			for _, token := range tokens {
				if !strings.Contains(cmdLower, token) {
					matchAll = false
					break
				}
			}
			if matchAll {
				results = append(results, i)
			}
		}

	case ModeFuzzy:
		strList := make([]string, n)
		for i, item := range m.CurrentItems {
			strList[i] = item.Command
		}
		matches := fuzzy.Find(query, strList)
		for i := len(matches) - 1; i >= 0; i-- {
			results = append(results, matches[i].Index)
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
