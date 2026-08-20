package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.TextInput.Width = max(20, msg.Width-30)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if len(m.FilteredIndexes) > 0 && m.SelectedIndex < len(m.FilteredIndexes) {
				idx := m.FilteredIndexes[m.SelectedIndex]
				m.SelectedCommand = m.CurrentItems[idx].Command
				// Copy to system clipboard
				_ = clipboard.WriteAll(m.SelectedCommand)
			}
			return m, tea.Quit

		case "ctrl+o":
			if len(m.FilteredIndexes) > 0 && m.SelectedIndex < len(m.FilteredIndexes) {
				idx := m.FilteredIndexes[m.SelectedIndex]
				m.SelectedCommand = m.CurrentItems[idx].Command
				m.PrintContextOnExit = true
				_ = clipboard.WriteAll(m.SelectedCommand)
			}
			return m, tea.Quit

		case "up", "ctrl+p":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}

		case "down", "ctrl+n":
			if m.SelectedIndex < len(m.FilteredIndexes)-1 {
				m.SelectedIndex++
			}

		case "pgup":
			m.SelectedIndex = max(0, m.SelectedIndex-10)

		case "pgdown":
			m.SelectedIndex = min(len(m.FilteredIndexes)-1, m.SelectedIndex+10)

		case "home":
			m.SelectedIndex = 0

		case "end":
			if len(m.FilteredIndexes) > 0 {
				m.SelectedIndex = len(m.FilteredIndexes) - 1
			}

		case "+", "=", "]":
			if m.ContextRadius < 50 {
				m.ContextRadius++
				m.CopiedNotification = fmt.Sprintf("Context size N = %d", m.ContextRadius)
			}

		case "-", "_", "[":
			if m.ContextRadius > 1 {
				m.ContextRadius--
				m.CopiedNotification = fmt.Sprintf("Context size N = %d", m.ContextRadius)
			}

		case "ctrl+r":
			m.SearchMode = (m.SearchMode + 1) % 4
			m.ApplyFilter()

		case "ctrl+s":
			// Cycle active source: bash -> zsh -> fish -> all
			switch m.Store.ActiveSource {
			case "bash":
				m.Store.ActiveSource = "zsh"
			case "zsh":
				m.Store.ActiveSource = "fish"
			case "fish":
				m.Store.ActiveSource = "all"
			default:
				m.Store.ActiveSource = "bash"
			}
			m.UpdateCurrentItems()
			m.ApplyFilter()
			m.CopiedNotification = fmt.Sprintf("Source switched to: %s", m.Store.ActiveSource)

		case "ctrl+e":
			// Copy full context window to clipboard
			if len(m.FilteredIndexes) > 0 && m.SelectedIndex < len(m.FilteredIndexes) {
				idx := m.FilteredIndexes[m.SelectedIndex]
				targetHist := m.CurrentItems[idx]

				start := max(0, targetHist.Index-m.ContextRadius)
				end := min(len(m.CurrentItems)-1, targetHist.Index+m.ContextRadius)

				var sb string
				for i := start; i <= end; i++ {
					prefix := "  "
					if i == targetHist.Index {
						prefix = "> "
					}
					sb += fmt.Sprintf("%s%s\n", prefix, m.CurrentItems[i].Command)
				}
				_ = clipboard.WriteAll(sb)
				m.CopiedNotification = "Copied context window to clipboard!"
			}
		}
	}

	// Update text input
	prevValue := m.TextInput.Value()
	m.TextInput, cmd = m.TextInput.Update(msg)
	if m.TextInput.Value() != prevValue {
		m.SelectedIndex = 0
		m.ApplyFilter()
	}

	return m, cmd
}
