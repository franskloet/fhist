package history

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadZshHistory(path string) ([]HistoryItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []HistoryItem
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var currentTs time.Time
	var currentCmd strings.Builder
	inExtended := false

	flushCmd := func() {
		cmdStr := strings.TrimSpace(currentCmd.String())
		if cmdStr != "" {
			items = append(items, HistoryItem{
				Index:     len(items),
				Command:   cmdStr,
				Timestamp: currentTs,
				ExitCode:  -1,
				Source:    "zsh",
			})
		}
		currentCmd.Reset()
		currentTs = time.Time{}
		inExtended = false
	}

	for scanner.Scan() {
		line := scanner.Text()

		// Check for zsh extended format: `: 1724128500:0;command`
		if strings.HasPrefix(line, ": ") && strings.Contains(line, ";") {
			if inExtended || currentCmd.Len() > 0 {
				flushCmd()
			}

			parts := strings.SplitN(line[2:], ";", 2)
			if len(parts) == 2 {
				meta := parts[0]
				cmdText := parts[1]

				// meta is "timestamp:duration"
				metaParts := strings.Split(meta, ":")
				if len(metaParts) >= 1 {
					if tsVal, err := strconv.ParseInt(metaParts[0], 10, 64); err == nil {
						currentTs = time.Unix(tsVal, 0)
					}
				}

				currentCmd.WriteString(cmdText)
				inExtended = true

				// Check if ends with multiline escape `\`
				if !strings.HasSuffix(cmdText, "\\") {
					flushCmd()
				} else {
					// Remove trailing backslash for multiline continuation
					str := currentCmd.String()
					currentCmd.Reset()
					currentCmd.WriteString(strings.TrimSuffix(str, "\\"))
				}
				continue
			}
		}

		if inExtended {
			currentCmd.WriteString("\n")
			if strings.HasSuffix(line, "\\") {
				currentCmd.WriteString(strings.TrimSuffix(line, "\\"))
			} else {
				currentCmd.WriteString(line)
				flushCmd()
			}
		} else {
			if currentCmd.Len() > 0 {
				currentCmd.WriteString("\n")
			}
			currentCmd.WriteString(line)
			flushCmd()
		}
	}

	if currentCmd.Len() > 0 {
		flushCmd()
	}

	for i := range items {
		items[i].Index = i
	}

	return items, scanner.Err()
}
