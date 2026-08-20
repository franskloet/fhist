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

	// Calculate strict pixel-perfect line budget
	// Total height budget = m.Height
	// Overhead lines:
	// - Header: 1 line
	// - Status Bar: 1 line
	// - Top Box Borders: 2 lines
	// - Bottom Box Borders: 2 lines
	// Total fixed overhead = 6 lines.

	totalAvailInner := m.Height - 6
	if totalAvailInner < 6 {
		totalAvailInner = 6
	}

	// Allocate 50% to Top Pane (Search Hits) and 50% to Bottom Pane (Context)
	topInnerHeight := totalAvailInner / 2
	if topInnerHeight < 3 {
		topInnerHeight = 3
	}
	bottomInnerHeight := totalAvailInner - topInnerHeight
	if bottomInnerHeight < 3 {
		bottomInnerHeight = 3
	}

	var sb strings.Builder

	// 1. Header (Search Bar & Badges)
	title := m.Styles.HeaderTitle.Render("fhist")
	searchInput := m.TextInput.View()

	sourceBadge := m.renderSourceBadge(m.Store.ActiveSource)
	modeBadge := m.Styles.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(lipgloss.Color("#F38BA8")).
		Render(fmt.Sprintf("[%s]", m.SearchMode.String()))

	hitCount := fmt.Sprintf("%d/%d matches", len(m.FilteredIndexes), len(m.CurrentItems))
	counter := m.Styles.CounterText.Render(hitCount)

	headerLine := fmt.Sprintf("%s %s  %s %s  %s", title, searchInput, sourceBadge, modeBadge, counter)
	sb.WriteString(headerLine + "\n")

	// 2. Top Pane (Search Matches List Box)
	topPaneContent := m.renderTopPane(topInnerHeight)
	topBox := m.Styles.TableBorder.
		Width(max(20, m.Width-2)).
		Height(topInnerHeight + 2).
		Render(topPaneContent)
	sb.WriteString(topBox + "\n")

	// 3. Bottom Pane (Context Window Box)
	bottomPaneContent := m.renderContextPane(bottomInnerHeight)
	bottomBox := m.Styles.ContextBox.
		Width(max(20, m.Width-2)).
		Height(bottomInnerHeight + 2).
		Render(bottomPaneContent)
	sb.WriteString(bottomBox + "\n")

	// 4. Status Bar
	var notification string
	if m.CopiedNotification != "" {
		notification = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")).Bold(true).Render(" " + m.CopiedNotification)
	}

	helpBar := fmt.Sprintf(
		"%s %s  %s %s  %s %s  %s %s  %s %s%s",
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

func (m Model) renderSourceBadge(source string) string {
	switch source {
	case "bash":
		return m.Styles.BadgeBash.Render("[bash]")
	case "zsh":
		return m.Styles.BadgeZsh.Render("[zsh]")
	case "fish":
		return m.Styles.BadgeFish.Render("[fish]")
	default:
		return m.Styles.BadgeAll.Render(fmt.Sprintf("[%s]", source))
	}
}

func (m Model) renderTopPane(maxLines int) string {
	contentWidth := max(20, m.Width-6)

	if len(m.FilteredIndexes) == 0 {
		var lines []string
		lines = append(lines, m.Styles.DimmedText.Render("  No matching history commands found."))
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		return strings.Join(lines[:maxLines], "\n")
	}

	// Calculate scrolling window for top pane list
	startIdx := 0
	if m.SelectedIndex >= maxLines {
		startIdx = m.SelectedIndex - maxLines + 1
	}
	endIdx := min(len(m.FilteredIndexes), startIdx+maxLines)

	var lines []string

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

	// Fill remaining height with blank lines so height stays fixed
	for len(lines) < maxLines {
		lines = append(lines, "")
	}

	return strings.Join(lines[:maxLines], "\n")
}

func (m Model) renderContextPane(maxLines int) string {
	contentWidth := max(20, m.Width-6)

	if len(m.FilteredIndexes) == 0 || m.SelectedIndex >= len(m.FilteredIndexes) {
		var lines []string
		lines = append(lines, m.Styles.DimmedText.Render("No context available."))
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		return strings.Join(lines[:maxLines], "\n")
	}

	targetMatchIdx := m.FilteredIndexes[m.SelectedIndex]
	targetItem := m.CurrentItems[targetMatchIdx]
	targetHistIndex := targetItem.Index

	// Title header takes 1 line
	titleText := fmt.Sprintf("Context Window (N = %d commands before/after hit #%d)", m.ContextRadius, targetHistIndex+1)
	header := m.Styles.ContextTitle.Render(titleText)

	var lines []string
	lines = append(lines, header)

	availLinesForContext := maxLines - 1
	if availLinesForContext < 1 {
		availLinesForContext = 1
	}

	// Determine how many lines before and after we can fit
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

	// Fill remaining height with blank lines so height stays fixed
	for len(lines) < maxLines {
		lines = append(lines, "")
	}

	return strings.Join(lines[:maxLines], "\n")
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
