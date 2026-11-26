package diff

import (
	"fmt"
	"os"
	"strings"
)

// Options controls unified diff generation and styling.
type Options struct {
	Color   bool
	Context int
}

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorBold  = "\033[1m"
)

// ShouldColor determines if color output should be enabled based on TTY and environment.
func ShouldColor() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Unified generates a unified diff between two text strings.
func Unified(fromName, toName, a, b string, opts Options) string {
	if a == b {
		return ""
	}

	aLines := splitLines(a)
	bLines := splitLines(b)

	diffs := computeLCSDiff(aLines, bLines)
	if len(diffs) == 0 {
		return ""
	}

	var sb strings.Builder
	fromHeader := fmt.Sprintf("--- %s", fromName)
	toHeader := fmt.Sprintf("+++ %s", toName)

	if opts.Color {
		sb.WriteString(fmt.Sprintf("%s%s%s\n", colorBold, fromHeader, colorReset))
		sb.WriteString(fmt.Sprintf("%s%s%s\n", colorBold, toHeader, colorReset))
	} else {
		sb.WriteString(fromHeader + "\n")
		sb.WriteString(toHeader + "\n")
	}

	hunks := groupHunks(diffs, opts.Context)
	for _, hunk := range hunks {
		hunkHeader := fmt.Sprintf("@@ -%d,%d +%d,%d @@", hunk.aStart, hunk.aCount, hunk.bStart, hunk.bCount)
		if opts.Color {
			sb.WriteString(fmt.Sprintf("%s%s%s\n", colorCyan, hunkHeader, colorReset))
		} else {
			sb.WriteString(hunkHeader + "\n")
		}

		for _, line := range hunk.lines {
			switch line.op {
			case opDelete:
				if opts.Color {
					sb.WriteString(fmt.Sprintf("%s-%s%s\n", colorRed, line.text, colorReset))
				} else {
					sb.WriteString("-" + line.text + "\n")
				}
			case opInsert:
				if opts.Color {
					sb.WriteString(fmt.Sprintf("%s+%s%s\n", colorGreen, line.text, colorReset))
				} else {
					sb.WriteString("+" + line.text + "\n")
				}
			case opEqual:
				sb.WriteString(" " + line.text + "\n")
			}
		}
	}

	return sb.String()
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

type diffOp int

const (
	opEqual diffOp = iota
	opInsert
	opDelete
)

type diffItem struct {
	op   diffOp
	text string
	aIdx int
	bIdx int
}

func computeLCSDiff(a, b []string) []diffItem {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] >= dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// Backtrack to build the diff
	var items []diffItem
	i, j := n, m
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			items = append(items, diffItem{op: opEqual, text: a[i-1], aIdx: i, bIdx: j})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			items = append(items, diffItem{op: opInsert, text: b[j-1], aIdx: i, bIdx: j})
			j--
		} else if i > 0 && (j == 0 || dp[i][j-1] < dp[i-1][j]) {
			items = append(items, diffItem{op: opDelete, text: a[i-1], aIdx: i, bIdx: j})
			i--
		}
	}

	// Reverse to restore original top-to-bottom order
	for k := 0; k < len(items)/2; k++ {
		opp := len(items) - 1 - k
		items[k], items[opp] = items[opp], items[k]
	}

	return items
}

type hunk struct {
	aStart int
	aCount int
	bStart int
	bCount int
	lines  []diffItem
}

func groupHunks(items []diffItem, context int) []hunk {
	if context <= 0 {
		context = 3
	}

	var hunks []hunk
	hasDiff := false
	for _, it := range items {
		if it.op != opEqual {
			hasDiff = true
			break
		}
	}
	if !hasDiff {
		return nil
	}

	// Find changed indices
	type groupRange struct {
		start int
		end   int
	}
	var changeGroups []groupRange
	inChange := false
	groupStart := 0

	for idx, it := range items {
		if it.op != opEqual {
			if !inChange {
				inChange = true
				groupStart = idx - context
				if groupStart < 0 {
					groupStart = 0
				}
			}
		} else if inChange {
			// Check if another change occurs within 2*context
			hasUpcoming := false
			for future := idx + 1; future <= idx+2*context && future < len(items); future++ {
				if items[future].op != opEqual {
					hasUpcoming = true
					break
				}
			}
			if !hasUpcoming {
				inChange = false
				groupEnd := idx + context
				if groupEnd > len(items) {
					groupEnd = len(items)
				}
				changeGroups = append(changeGroups, groupRange{start: groupStart, end: groupEnd})
			}
		}
	}
	if inChange {
		changeGroups = append(changeGroups, groupRange{start: groupStart, end: len(items)})
	}

	for _, grp := range changeGroups {
		sub := items[grp.start:grp.end]
		aStart, aCount := 0, 0
		bStart, bCount := 0, 0

		for _, l := range sub {
			if l.op == opDelete || l.op == opEqual {
				if aStart == 0 && l.aIdx > 0 {
					aStart = l.aIdx
				}
				aCount++
			}
			if l.op == opInsert || l.op == opEqual {
				if bStart == 0 && l.bIdx > 0 {
					bStart = l.bIdx
				}
				bCount++
			}
		}
		if aStart == 0 {
			aStart = 1
		}
		if bStart == 0 {
			bStart = 1
		}

		hunks = append(hunks, hunk{
			aStart: aStart,
			aCount: aCount,
			bStart: bStart,
			bCount: bCount,
			lines:  sub,
		})
	}

	return hunks
}
