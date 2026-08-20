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

	// Strict Height Math:
	// Total terminal height = m.Height
	// - Search Box: 4 lines
	// - Hits Box borders: 2 lines
	// - Context Box borders: 2 lines
	// - Status bar: 1 line
	// - Newline separators (3 newlines between 4 components): 3 lines
	// Total fixed overhead = 4 + 2 + 2 + 1 + 3 = 12 lines.

	totalAvailInner := m.Height - 12
	if totalAvailInner < 2 {
		totalAvailInner = 2
	}

	topInnerHeight := totalAvailInner / 2
	if topInnerHeight < 1 {
		topInnerHeight = 1
	}
	bottomInnerHeight := totalAvailInner - topInnerHeight
	if bottomInnerHeight < 1 {
		bottomInnerHeight = 1
	}

	var components []string

	// 1. Search Box (4 lines)
	searchBoxStr := m.renderSearchBox(boxWidth)
	components = append(components, searchBoxStr)

	// 2. Search Hits Box (topInnerHeight + 2 lines)
	hitsBoxStr := m.renderSearchHitsBox(boxWidth, topInnerHeight)
	components = append(components, hitsBoxStr)

	// 3. Context Preview Box (bottomInnerHeight + 2 lines)
	contextBoxStr := m.renderContextBox(boxWidth, bottomInnerHeight)
	components = append(components, contextBoxStr)

	// 4. Status Bar (1 line)
	var notification string
	if m.CopiedNotification != "" {
		notification = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")).Bold(true).Render(" " + m.CopiedNotification)
	}

	scrollHint := ""
	if m.CommandOffset > 0 {
		scrollHint = lipgloss.NewStyle().Foreground(lipgloss.Color("#F9E2AF")).Bold(true).Render(fmt.Sprintf(" [Offset +%d] ", m.CommandOffset))
	}

	helpBar := fmt.Sprintf(
		" %s %s  %s %s  %s %s  %s %s  %s %s%s%s",
		m.Styles.KeyHint.Render("↑/↓"), m.Styles.KeyDesc.Render("Browse Hits"),
		m.Styles.KeyHint.Render("←/→"), m.Styles.KeyDesc.Render("Scroll Line"),
		m.Styles.KeyHint.Render("+/-"), m.Styles.KeyDesc.Render(fmt.Sprintf("Context N=%d", m.ContextRadius)),
		m.Styles.KeyHint.Render("Enter"), m.Styles.KeyDesc.Render("Select"),
		m.Styles.KeyHint.Render("Ctrl+R"), m.Styles.KeyDesc.Render("Mode"),
		scrollHint,
		notification,
	)
	components = append(components, m.Styles.StatusBar.Render(helpBar))

	return strings.Join(components, "\n")
}

func (m Model) renderSearchBox(boxWidth int) string {
	contentWidth := max(20, boxWidth-4)

	prompt := m.Styles.SearchPrompt.Render("🔍 Search: ")
	input := m.TextInput.View()
	inputLine := prompt + input
	inputLine = truncateString(inputLine, contentWidth)

	// Badges
	sourceBadge := m.renderSourceBadge(m.Store.ActiveSource)
	modeBadge := m.Styles.BadgeBox.
		Foreground(lipgloss.Color("#11111B")).
		Background(lipgloss.Color("#F38BA8")).
		Render(fmt.Sprintf("[%s Mode]", m.SearchMode.String()))

	hitCount := fmt.Sprintf("%d / %d matches", len(m.FilteredIndexes), len(m.CurrentItems))
	counter := m.Styles.CounterText.Render(hitCount)

	badgeLine := fmt.Sprintf("%s %s  %s", sourceBadge, modeBadge, counter)
	badgeLine = truncateString(badgeLine, contentWidth)

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
		selIdx := min(m.SelectedIndex, len(m.FilteredIndexes)-1)
		selIdx = max(0, selIdx)

		startIdx := max(0, selIdx-maxLines+1)
		endIdx := min(len(m.FilteredIndexes), startIdx+maxLines)

		for i := startIdx; i < endIdx; i++ {
			if i < 0 || i >= len(m.FilteredIndexes) {
				continue
			}
			itemIdx := m.FilteredIndexes[i]
			if itemIdx < 0 || itemIdx >= len(m.CurrentItems) {
				continue
			}
			item := m.CurrentItems[itemIdx]

			isSelected := (i == selIdx)

			idxStr := fmt.Sprintf("#%-5d", item.Index+1)
			srcBadge := fmt.Sprintf("%-6s", item.FormatSourceBadge())

			timeStr := item.ShortTime()
			if timeStr != "" {
				timeStr = fmt.Sprintf("%-12s", timeStr)
			}

			cmdText := sanitizeCmd(item.Command)
			cmdText = applyOffset(cmdText, m.CommandOffset)

			if isSelected {
				maxCmdLen := max(10, contentWidth-len(idxStr)-len(srcBadge)-len(timeStr)-6)
				cmdTruncated := truncateString(cmdText, maxCmdLen)
				rowStr := fmt.Sprintf("▸ %s %s %s %s", idxStr, srcBadge, timeStr, cmdTruncated)
				selectedLine := m.Styles.SelectedItem.Render(rowStr)
				lines = append(lines, selectedLine)
			} else {
				maxCmdLen := max(10, contentWidth-len(idxStr)-len(srcBadge)-len(timeStr)-6)
				cmdTruncated := truncateString(cmdText, maxCmdLen)
				rowStr := fmt.Sprintf("  %s %s %s %s", idxStr, srcBadge, timeStr, cmdTruncated)
				lines = append(lines, m.Styles.NormalItem.Render(rowStr))
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

	return hitsStyle.Render(boxContent)
}

func (m Model) renderContextBox(boxWidth int, maxLines int) string {
	contentWidth := max(20, boxWidth-4)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#CBA6F7")).
		Padding(0, 1).
		Width(boxWidth)

	selIdx := min(m.SelectedIndex, len(m.FilteredIndexes)-1)
	selIdx = max(0, selIdx)

	if len(m.FilteredIndexes) == 0 || selIdx >= len(m.FilteredIndexes) {
		var lines []string
		lines = append(lines, m.Styles.DimmedText.Render("No context available."))
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		return boxStyle.Render(strings.Join(lines[:maxLines], "\n"))
	}

	targetMatchIdx := m.FilteredIndexes[selIdx]
	if targetMatchIdx < 0 || targetMatchIdx >= len(m.CurrentItems) {
		var lines []string
		lines = append(lines, m.Styles.DimmedText.Render("No context available."))
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		return boxStyle.Render(strings.Join(lines[:maxLines], "\n"))
	}

	targetItem := m.CurrentItems[targetMatchIdx]
	targetHistIndex := targetItem.Index

	titleText := fmt.Sprintf("Context Preview (%d commands BEFORE & AFTER selected hit #%d)", m.ContextRadius, targetHistIndex+1)
	header := m.Styles.ContextTitle.Render(truncateString(titleText, contentWidth))

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
		if i < 0 || i >= len(m.CurrentItems) {
			continue
		}
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
		cmdText = applyOffset(cmdText, m.CommandOffset)

		if offset == 0 {
			maxCmdLen := max(10, contentWidth-len(tag)-len(idxTag)-len(timeTag)-12)
			cmdTruncated := truncateString(cmdText, maxCmdLen)
			lineStr := fmt.Sprintf(">>> %s %s %s %s", tag, idxTag, timeTag, cmdTruncated)
			highlightedLine := m.Styles.ContextTarget.Render(lineStr)
			lines = append(lines, highlightedLine)
		} else {
			maxCmdLen := max(10, contentWidth-len(tag)-len(idxTag)-len(timeTag)-12)
			cmdTruncated := truncateString(cmdText, maxCmdLen)

			tagStyled := m.Styles.ContextTag.Render(tag)
			idxStyled := m.Styles.DimmedText.Render(idxTag)
			timeStyled := m.Styles.DimmedText.Render(timeTag)
			cmdStyled := m.Styles.ContextBefore.Render(cmdTruncated)

			lines = append(lines, fmt.Sprintf("    %s %s %s %s", tagStyled, idxStyled, timeStyled, cmdStyled))
		}
	}

	for len(lines) < maxLines {
		lines = append(lines, "")
	}

	return boxStyle.Render(strings.Join(lines[:maxLines], "\n"))
}

func applyOffset(cmd string, offset int) string {
	if offset <= 0 {
		return cmd
	}
	runes := []rune(cmd)
	if offset >= len(runes) {
		return "« "
	}
	return "« " + string(runes[offset:])
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
