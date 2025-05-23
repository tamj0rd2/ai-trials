package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// main is the entry point of the application. It fetches commits,
// calculates pair programming days, generates HTML, and outputs the results.
func main() {
	commits, err := getCommitsFromGit()
	if err != nil {
		fmt.Println("Error getting commits:", err)
		return
	}

	pairDays, devList := calculatePairDays(commits)

	html := generateHTMLTable(pairDays, devList)

	if err := writeAndOpenHTML(html, "output.html"); err != nil {
		fmt.Println("Error writing or opening HTML:", err)
	}
}

// NewDevelopers creates a new Developers map from a list of developer names
func NewDevelopers(developers ...string) Developers {
	devs := make(Developers)
	for _, dev := range developers {
		devs[dev] = struct{}{}
	}
	return devs
}

// NewCommit creates a new Commit with the specified date and developers
func NewCommit(date string, developers ...string) Commit {
	return Commit{
		Date:       date,
		Developers: NewDevelopers(developers...),
	}
}

// calculatePairDays analyzes commits and returns a map of developer pairs to the number
// of days they paired together, along with a sorted list of all developers.
func calculatePairDays(commits []Commit) (map[Pair]int, []string) {
	// Map to track pairs by day
	pairsByDay := make(map[string]map[Pair]bool)

	// Map to collect all developers
	allDevs := make(map[string]bool)

	// Process each commit
	for _, commit := range commits {
		date := commit.Date

		// Skip if commit has fewer than 2 developers (no pairing)
		if len(commit.Developers) < 2 {
			// Still add individual developers to the list
			for dev := range commit.Developers {
				allDevs[dev] = true
			}
			continue
		}

		// Initialize map for this day if not exists
		if pairsByDay[date] == nil {
			pairsByDay[date] = make(map[Pair]bool)
		}

		// Add all developers to the set
		for dev := range commit.Developers {
			allDevs[dev] = true
		}

		// Create all possible pairs from developers in this commit
		devs := make([]string, 0, len(commit.Developers))
		for dev := range commit.Developers {
			devs = append(devs, dev)
		}

		// Create all possible pairs
		for i := 0; i < len(devs); i++ {
			for j := i + 1; j < len(devs); j++ {
				// Create a pair with names in lexical order
				pair := Pair{devs[i], devs[j]}
				if pair[0] > pair[1] {
					pair[0], pair[1] = pair[1], pair[0]
				}

				// Mark this pair as working together on this day
				pairsByDay[date][pair] = true
			}
		}
	}

	// Count days for each pair
	pairDays := make(map[Pair]int)
	for _, pairsMap := range pairsByDay {
		for pair := range pairsMap {
			pairDays[pair]++
		}
	}

	// Create sorted list of developers
	devList := make([]string, 0, len(allDevs))
	for dev := range allDevs {
		devList = append(devList, dev)
	}
	sort.Strings(devList)

	return pairDays, devList
}

// getCommitsFromGit fetches commit data from git for the last 2 months
func getCommitsFromGit() ([]Commit, error) {
	// Get date 2 months ago
	since := time.Now().AddDate(0, -2, 0).Format("2006-01-02")

	// Setup git log format with custom delimiter
	gitLogFormat := "%an|%ad|%BEND_OF_COMMIT"

	// Run git log command
	cmd := exec.Command("git", "log", "--since="+since, "--pretty=format:"+gitLogFormat, "--date=short")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error getting git log: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting git log: %w", err)
	}

	commits := []Commit{}
	commitLines := []string{}
	scanner := bufio.NewScanner(stdout)

	// Parse output line by line
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "END_OF_COMMIT") {
			// End of commit reached, process it
			commitLines = append(commitLines, line)
			commit := parseCommit(commitLines)
			commits = append(commits, commit)
			commitLines = []string{}
		} else {
			// Add line to current commit
			commitLines = append(commitLines, line)
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("error waiting for git log: %w", err)
	}

	return commits, nil
}

// parseCommit extracts information from commit message lines
func parseCommit(lines []string) Commit {
	if len(lines) == 0 {
		return Commit{}
	}

	// Parse first line to get author and date
	firstLine := lines[0]
	parts := strings.SplitN(firstLine, "|", 3)

	if len(parts) < 2 {
		return Commit{}
	}

	author := parts[0]
	date := parts[1]

	// Create developers map and add author
	developers := NewDevelopers(author)

	// Look for co-authors in the rest of the lines
	for _, line := range lines[1:] {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "Co-authored-by:") {
			// Extract name from co-author line
			coAuthorInfo := strings.TrimPrefix(trimmedLine, "Co-authored-by:")
			coAuthorName := extractNameFromCoAuthorLine(coAuthorInfo)
			if coAuthorName != "" {
				developers[coAuthorName] = struct{}{}
			}
		}
	}

	return Commit{
		Date:       date,
		Developers: developers,
	}
}

// extractNameFromCoAuthorLine extracts the developer name from a co-author line
func extractNameFromCoAuthorLine(coAuthorInfo string) string {
	coAuthorInfo = strings.TrimSpace(coAuthorInfo)

	// Extract name before email
	if idx := strings.Index(coAuthorInfo, "<"); idx > 0 {
		return strings.TrimSpace(coAuthorInfo[:idx])
	}

	return coAuthorInfo
}

// generateHTMLTable creates HTML output with pairing statistics
func generateHTMLTable(pairDays map[Pair]int, devList []string) string {
	html := `<html>
<head>
    <title>Pair Programming Stats</title>
    <style>
        body, table, th, td {
            font-family: sans-serif;
        }
        table, th, td {
            border: 1px solid #ccc;
            border-collapse: collapse;
        }
        th, td {
            padding: 8px;
        }
    </style>
</head>
<body>
    <h1>Pair Programming Days Table (Last 2 Months)</h1>
    <table>
`

	// Header row
	html += "<tr><th></th>"
	for _, dev := range devList {
		html += fmt.Sprintf("<th>%s</th>", dev)
	}
	html += "</tr>\n"

	// Data rows
	for _, rowDev := range devList {
		html += fmt.Sprintf("<tr><th>%s</th>", rowDev)

		for _, colDev := range devList {
			if rowDev == colDev {
				html += "<td>-</td>"
			} else {
				// Create pair (in lexicographically sorted order)
				pair := Pair{rowDev, colDev}
				if pair[0] > pair[1] {
					pair[0], pair[1] = pair[1], pair[0]
				}

				// Get pair days count
				count := 0
				if days, ok := pairDays[pair]; ok {
					count = days
				}

				html += fmt.Sprintf("<td>%d</td>", count)
			}
		}

		html += "</tr>\n"
	}

	html += `    </table>
</body>
</html>`

	return html
}

// writeAndOpenHTML writes HTML content to a file and opens it
func writeAndOpenHTML(html, filename string) error {
	// Write to file
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating HTML file: %w", err)
	}

	// Use a closure to handle the close error properly
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing HTML file: %w", closeErr)
		}
	}()

	_, err = f.WriteString(html)
	if err != nil {
		return fmt.Errorf("error writing HTML: %w", err)
	}

	fmt.Println("Wrote", filename)

	// Open the created file
	if err := exec.Command("open", filename).Start(); err != nil {
		return fmt.Errorf("error opening HTML file: %w", err)
	}

	return nil
}
