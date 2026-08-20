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

	var sb strings.Builder

	// 1. Header (Search Bar & Badges)
	title := m.Styles.HeaderTitle.Render("fhist")
	searchInput := m.TextInput.View()

	// Badges
	sourceBadge := m.renderSourceBadge(m.Store.ActiveSource)
	modeBadge := m.Styles.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(lipgloss.Color("#F38BA8")).
		Render(fmt.Sprintf("[%s]", m.SearchMode.String()))

	hitCount := fmt.Sprintf("%d/%d matches", len(m.FilteredIndexes), len(m.CurrentItems))
	counter := m.Styles.CounterText.Render(hitCount)

	headerLine := fmt.Sprintf("%s %s  %s %s  %s", title, searchInput, sourceBadge, modeBadge, counter)
	sb.WriteString(headerLine + "\n")

	// Calculate vertical layout sizes
	availHeight := m.Height - 6 // Header(2), Divider(1), Footer(1), Margins(2)
	if availHeight < 10 {
		availHeight = 10
	}

	// Calculate Top vs Bottom pane heights (approx 45% top, 55% bottom)
	topPaneHeight := availHeight * 45 / 100
	if topPaneHeight < 4 {
		topPaneHeight = 4
	}
	bottomPaneHeight := availHeight - topPaneHeight
	if bottomPaneHeight < 5 {
		bottomPaneHeight = 5
	}

	// 2. Top Pane: Search Match List
	topPaneStr := m.renderTopPane(topPaneHeight)
	sb.WriteString(topPaneStr + "\n")

	// 3. Separator Bar
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#45475A")).
		Render(strings.Repeat("─", max(20, m.Width-2)))
	sb.WriteString(divider + "\n")

	// 4. Bottom Pane: Context Window (N before and N after selected hit)
	bottomPaneStr := m.renderContextPane(bottomPaneHeight)
	sb.WriteString(bottomPaneStr + "\n")

	// 5. Footer / Status Bar
	var notification string
	if m.CopiedNotification != "" {
		notification = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")).Bold(true).Render(" " + m.CopiedNotification)
	}

	helpBar := fmt.Sprintf(
		"%s %s  %s %s  %s %s  %s %s  %s %s%s",
		m.Styles.KeyHint.Render("↑/↓"), m.Styles.KeyDesc.Render("Nav"),
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
	if len(m.FilteredIndexes) == 0 {
		return m.Styles.DimmedText.Render("  No matching history commands found.")
	}

	startIdx := 0
	if m.SelectedIndex >= maxLines {
		startIdx = m.SelectedIndex - maxLines + 1
	}
	endIdx := min(len(m.FilteredIndexes), startIdx+maxLines)

	var lines []string
	contentWidth := max(20, m.Width-4)

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

		rowStr := fmt.Sprintf(" %s %s %s %s", idxStr, srcBadge, timeStr, cmdText)
		rowStr = truncateString(rowStr, contentWidth)

		if isSelected {
			selectedLine := m.Styles.SelectedItem.Width(contentWidth).Render("> " + rowStr)
			lines = append(lines, selectedLine)
		} else {
			lines = append(lines, "  "+m.Styles.NormalItem.Render(rowStr))
		}
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderContextPane(maxLines int) string {
	if len(m.FilteredIndexes) == 0 || m.SelectedIndex >= len(m.FilteredIndexes) {
		return m.Styles.ContextBox.Width(max(20, m.Width-4)).Render("No context available.")
	}

	targetMatchIdx := m.FilteredIndexes[m.SelectedIndex]
	targetItem := m.CurrentItems[targetMatchIdx]
	targetHistIndex := targetItem.Index

	startHistIndex := max(0, targetHistIndex-m.ContextRadius)
	endHistIndex := min(len(m.CurrentItems)-1, targetHistIndex+m.ContextRadius)

	titleText := fmt.Sprintf("Context Window (N = %d commands before/after hit #%d)", m.ContextRadius, targetHistIndex+1)
	header := m.Styles.ContextTitle.Render("│ " + titleText)

	var lines []string
	lines = append(lines, header)

	contentWidth := max(20, m.Width-6)

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
			lineStr = truncateString(lineStr, contentWidth)
			highlightedLine := m.Styles.ContextTarget.Width(contentWidth).Render(">>> " + lineStr)
			lines = append(lines, highlightedLine)
		} else {
			lineStr := fmt.Sprintf("%s %s %s %s", tag, idxTag, timeTag, cmdText)
			lineStr = truncateString(lineStr, contentWidth)

			tagStyled := m.Styles.ContextTag.Render(tag)
			idxStyled := m.Styles.DimmedText.Render(idxTag)
			timeStyled := m.Styles.DimmedText.Render(timeTag)
			cmdStyled := m.Styles.ContextBefore.Render(cmdText)

			lines = append(lines, fmt.Sprintf("    %s %s %s %s", tagStyled, idxStyled, timeStyled, cmdStyled))
		}
	}

	boxContent := strings.Join(lines, "\n")
	return m.Styles.ContextBox.Width(max(20, m.Width-2)).Render(boxContent)
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
