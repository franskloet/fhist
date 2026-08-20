package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	boxWidth := max(30, m.Width-2)

	// Calculate vertical line budget to fit perfectly inside m.Height
	// Overhead lines:
	// - Search Box (title border + input line + badge line + bottom border) = 4 lines
	// - Search Hits Box (title border + content + bottom border) = topInnerHeight + 2
	// - Context Box (title border + content + bottom border) = bottomInnerHeight + 2
	// - Status bar = 1 line
	// Total fixed overhead = 4 + 2 + 2 + 1 = 9 lines.

	totalAvailInner := m.Height - 9
	if totalAvailInner < 4 {
		totalAvailInner = 4
	}

	topInnerHeight := totalAvailInner / 2
	if topInnerHeight < 2 {
		topInnerHeight = 2
	}
	bottomInnerHeight := totalAvailInner - topInnerHeight
	if bottomInnerHeight < 2 {
		bottomInnerHeight = 2
	}

	var sb strings.Builder

	// 1. Search Box (Top)
	searchBoxStr := m.renderSearchBox(boxWidth)
	sb.WriteString(searchBoxStr + "\n")

	// 2. Search Hits Box (Middle Pane)
	hitsBoxStr := m.renderSearchHitsBox(boxWidth, topInnerHeight)
	sb.WriteString(hitsBoxStr + "\n")

	// 3. Context Preview Box (Bottom Pane)
	contextBoxStr := m.renderContextBox(boxWidth, bottomInnerHeight)
	sb.WriteString(contextBoxStr + "\n")

	// 4. Footer / Status Bar
	var notification string
	if m.CopiedNotification != "" {
		notification = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")).Bold(true).Render(" " + m.CopiedNotification)
	}

	helpBar := fmt.Sprintf(
		" %s %s  %s %s  %s %s  %s %s  %s %s%s",
		m.Styles.KeyHint.Render("↑/↓"), m.Styles.KeyDesc.Render("Browse Hits"),
		m.Styles.KeyHint.Render("+/-"), m.Styles.KeyDesc.Render(fmt.Sprintf("Context N=%d", m.ContextRadius)),
		m.Styles.KeyHint.Render("Enter"), m.Styles.KeyDesc.Render("Select"),
		m.Styles.KeyHint.Render("Ctrl+R"), m.Styles.KeyDesc.Render("Search Mode"),
		m.Styles.KeyHint.Render("Ctrl+S"), m.Styles.KeyDesc.Render("Source"),
		notification,
	)
	sb.WriteString(m.Styles.StatusBar.Render(helpBar))

	return sb.String()
}

func (m Model) renderSearchBox(boxWidth int) string {
	prompt := m.Styles.SearchPrompt.Render("🔍 Type Search Query: ")
	input := m.TextInput.View()
	inputLine := prompt + input

	// Badges
	sourceBadge := m.renderSourceBadge(m.Store.ActiveSource)
	modeBadge := m.Styles.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(lipgloss.Color("#F38BA8")).
		Render(fmt.Sprintf("[%s Mode]", m.SearchMode.String()))

	hitCount := fmt.Sprintf("%d / %d matches", len(m.FilteredIndexes), len(m.CurrentItems))
	counter := m.Styles.CounterText.Render(hitCount)

	badgeLine := fmt.Sprintf("%s %s  %s", sourceBadge, modeBadge, counter)

	boxContent := inputLine + "\n" + badgeLine

	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89B4FA")).
		Padding(0, 1).
		Width(boxWidth)

	return searchStyle.Render(boxContent)
}

func (m Model) renderSourceBadge(source string) string {
	switch source {
	case "bash":
		return m.Styles.BadgeBash.Render("[Source: bash]")
	case "zsh":
		return m.Styles.BadgeZsh.Render("[Source: zsh]")
	case "fish":
		return m.Styles.BadgeFish.Render("[Source: fish]")
	default:
		return m.Styles.BadgeAll.Render(fmt.Sprintf("[Source: %s]", source))
	}
}

func (m Model) renderSearchHitsBox(boxWidth int, maxLines int) string {
	contentWidth := max(20, boxWidth-4)

	var lines []string

	if len(m.FilteredIndexes) == 0 {
		lines = append(lines, m.Styles.DimmedText.Render("  No matching history commands found."))
	} else {
		startIdx := 0
		if m.SelectedIndex >= maxLines {
			startIdx = m.SelectedIndex - maxLines + 1
		}
		endIdx := min(len(m.FilteredIndexes), startIdx+maxLines)

		for i := startIdx; i < endIdx; i++ {
			itemIdx := m.FilteredIndexes[i]
			item := m.CurrentItems[itemIdx]

			isSelected := (i == m.SelectedIndex)

			idxStr := fmt.Sprintf("#%-5d", item.Index+1)
			srcBadge := fmt.Sprintf("%-6s", item.FormatSourceBadge())

			timeStr := item.ShortTime()
			if timeStr != "" {
				timeStr = fmt.Sprintf("%-12s", timeStr)
			}

			cmdText := sanitizeCmd(item.Command)

			rowStr := fmt.Sprintf("%s %s %s %s", idxStr, srcBadge, timeStr, cmdText)
			rowStr = truncateString(rowStr, contentWidth-3)

			if isSelected {
				selectedLine := m.Styles.SelectedItem.Width(contentWidth).Render("▸ " + rowStr)
				lines = append(lines, selectedLine)
			} else {
				lines = append(lines, "  "+m.Styles.NormalItem.Render(rowStr))
			}
		}
	}

	for len(lines) < maxLines {
		lines = append(lines, "")
	}

	boxContent := strings.Join(lines[:maxLines], "\n")

	hitsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CDD6F4")).
		Padding(0, 1).
		Width(boxWidth)

	// Add header title directly to box
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89B4FA")).Render(" Search Hits (Use ↑ / ↓ to navigate) ")
	_ = title

	return hitsStyle.Render(boxContent)
}

func (m Model) renderContextBox(boxWidth int, maxLines int) string {
	contentWidth := max(20, boxWidth-4)

	if len(m.FilteredIndexes) == 0 || m.SelectedIndex >= len(m.FilteredIndexes) {
		var lines []string
		lines = append(lines, m.Styles.DimmedText.Render("No context available."))
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#CBA6F7")).
			Padding(0, 1).
			Width(boxWidth)
		return boxStyle.Render(strings.Join(lines[:maxLines], "\n"))
	}

	targetMatchIdx := m.FilteredIndexes[m.SelectedIndex]
	targetItem := m.CurrentItems[targetMatchIdx]
	targetHistIndex := targetItem.Index

	titleText := fmt.Sprintf("Context Preview (%d commands BEFORE & AFTER selected hit #%d)", m.ContextRadius, targetHistIndex+1)
	header := m.Styles.ContextTitle.Render(titleText)

	var lines []string
	lines = append(lines, header)

	availLinesForContext := maxLines - 1
	if availLinesForContext < 1 {
		availLinesForContext = 1
	}

	radius := (availLinesForContext - 1) / 2
	if radius > m.ContextRadius {
		radius = m.ContextRadius
	}

	startHistIndex := max(0, targetHistIndex-radius)
	endHistIndex := min(len(m.CurrentItems)-1, targetHistIndex+radius)

	for i := startHistIndex; i <= endHistIndex; i++ {
		item := m.CurrentItems[i]
		offset := i - targetHistIndex

		var tag string
		if offset < 0 {
			tag = fmt.Sprintf("[ -%-2d ]", -offset)
		} else if offset == 0 {
			tag = "[ 🎯 0 ]"
		} else {
			tag = fmt.Sprintf("[ +%-2d ]", offset)
		}

		idxTag := fmt.Sprintf("#%-5d", item.Index+1)
		timeTag := item.ShortTime()
		if timeTag != "" {
			timeTag = "[" + timeTag + "]"
		}

		cmdText := sanitizeCmd(item.Command)

		if offset == 0 {
			lineStr := fmt.Sprintf("%s %s %s %s", tag, idxTag, timeTag, cmdText)
			lineStr = truncateString(lineStr, contentWidth-4)
			highlightedLine := m.Styles.ContextTarget.Width(contentWidth).Render(">>> " + lineStr)
			lines = append(lines, highlightedLine)
		} else {
			lineStr := fmt.Sprintf("%s %s %s %s", tag, idxTag, timeTag, cmdText)
			lineStr = truncateString(lineStr, contentWidth-4)

			tagStyled := m.Styles.ContextTag.Render(tag)
			idxStyled := m.Styles.DimmedText.Render(idxTag)
			timeStyled := m.Styles.DimmedText.Render(timeTag)
			cmdStyled := m.Styles.ContextBefore.Render(cmdText)

			lines = append(lines, fmt.Sprintf("    %s %s %s %s", tagStyled, idxStyled, timeStyled, cmdStyled))
		}
	}

	for len(lines) < maxLines {
		lines = append(lines, "")
	}

	contextStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBA6F7")).
		Padding(0, 1).
		Width(boxWidth)

	return contextStyle.Render(strings.Join(lines[:maxLines], "\n"))
}

func sanitizeCmd(cmd string) string {
	res := strings.ReplaceAll(cmd, "\n", " ↵ ")
	res = strings.ReplaceAll(res, "\t", " ")
	return res
}

func truncateString(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) > maxWidth {
		return string(runes[:maxWidth-3]) + "..."
	}
	return s
}
