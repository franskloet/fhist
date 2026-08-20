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
		m.TextInput.Width = max(15, msg.Width-35)

	case tea.KeyMsg:
		keyStr := msg.String()

		switch keyStr {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			selIdx := min(m.SelectedIndex, len(m.FilteredIndexes)-1)
			if len(m.FilteredIndexes) > 0 && selIdx >= 0 {
				idx := m.FilteredIndexes[selIdx]
				m.SelectedCommand = m.CurrentItems[idx].Command
				_ = clipboard.WriteAll(m.SelectedCommand)
			}
			return m, tea.Quit

		case "ctrl+o":
			selIdx := min(m.SelectedIndex, len(m.FilteredIndexes)-1)
			if len(m.FilteredIndexes) > 0 && selIdx >= 0 {
				idx := m.FilteredIndexes[selIdx]
				m.SelectedCommand = m.CurrentItems[idx].Command
				m.PrintContextOnExit = true
				_ = clipboard.WriteAll(m.SelectedCommand)
			}
			return m, tea.Quit

		case "up", "ctrl+p":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
			return m, nil

		case "down", "ctrl+n":
			if m.SelectedIndex < len(m.FilteredIndexes)-1 {
				m.SelectedIndex++
			}
			return m, nil

		case "right":
			m.CommandOffset += 5
			m.CopiedNotification = fmt.Sprintf("Scrolled right (+%d chars)", m.CommandOffset)
			return m, nil

		case "left":
			if m.CommandOffset > 0 {
				m.CommandOffset = max(0, m.CommandOffset-5)
				if m.CommandOffset == 0 {
					m.CopiedNotification = "Scrolled to start"
				} else {
					m.CopiedNotification = fmt.Sprintf("Scrolled (+%d chars)", m.CommandOffset)
				}
			}
			return m, nil

		case "+", "=", "shift+up", "alt+up":
			if m.ContextRadius < 50 {
				m.ContextRadius++
				m.CopiedNotification = fmt.Sprintf("Context size N = %d", m.ContextRadius)
			}
			return m, nil

		case "-", "_", "shift+down", "alt+down":
			if m.ContextRadius > 1 {
				m.ContextRadius--
				m.CopiedNotification = fmt.Sprintf("Context size N = %d", m.ContextRadius)
			}
			return m, nil

		case "pgup":
			m.SelectedIndex = max(0, m.SelectedIndex-10)
			return m, nil

		case "pgdown":
			m.SelectedIndex = min(len(m.FilteredIndexes)-1, m.SelectedIndex+10)
			return m, nil

		case "home":
			m.SelectedIndex = 0
			return m, nil

		case "end":
			if len(m.FilteredIndexes) > 0 {
				m.SelectedIndex = len(m.FilteredIndexes) - 1
			}
			return m, nil

		case "ctrl+r":
			m.SearchMode = (m.SearchMode + 1) % 4
			m.ApplyFilter()
			return m, nil

		case "ctrl+s":
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
			return m, nil

		case "ctrl+e":
			selIdx := min(m.SelectedIndex, len(m.FilteredIndexes)-1)
			if len(m.FilteredIndexes) > 0 && selIdx >= 0 {
				idx := m.FilteredIndexes[selIdx]
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
			return m, nil
		}
	}

	// Route text typing to text input
	prevValue := m.TextInput.Value()
	m.TextInput, cmd = m.TextInput.Update(msg)
	if m.TextInput.Value() != prevValue {
		m.SelectedIndex = 0
		m.ApplyFilter()
	}

	return m, cmd
}
