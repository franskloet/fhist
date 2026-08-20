package history

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadFishHistory(path string) ([]HistoryItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var items []HistoryItem
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var currentCmd string
	var currentTs time.Time

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "- cmd:") {
			if currentCmd != "" {
				items = append(items, HistoryItem{
					Index:     len(items),
					Command:   currentCmd,
					Timestamp: currentTs,
					ExitCode:  -1,
					Source:    "fish",
				})
				currentCmd = ""
				currentTs = time.Time{}
			}
			currentCmd = strings.TrimSpace(strings.TrimPrefix(trimmed, "- cmd:"))
		} else if strings.HasPrefix(trimmed, "when:") {
			tsStr := strings.TrimSpace(strings.TrimPrefix(trimmed, "when:"))
			if tsVal, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
				currentTs = time.Unix(tsVal, 0)
			}
		}
	}

	if currentCmd != "" {
		items = append(items, HistoryItem{
			Index:     len(items),
			Command:   currentCmd,
			Timestamp: currentTs,
			ExitCode:  -1,
			Source:    "fish",
		})
	}

	for i := range items {
		items[i].Index = i
	}

	return items, scanner.Err()
}
