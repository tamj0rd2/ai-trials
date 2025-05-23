package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Developers map[string]struct{}
type Pair [2]string

type Commit struct {
	Date       string
	Developers Developers
}

func main() {
	// Get date 2 months ago
	since := time.Now().AddDate(0, -2, 0).Format("2006-01-02")
	gitLogFormat := "%an|%ad|%BEND_OF_COMMIT"
	cmd := exec.Command("git", "log", "--since="+since, "--pretty=format:"+gitLogFormat, "--date=short")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting git log:", err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting git log:", err)
		return
	}
	commits := []Commit{}
	commitLines := []string{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "END_OF_COMMIT") {
			commitLines = append(commitLines, line)
			commit := parseCommit(commitLines)
			commits = append(commits, commit)
			commitLines = []string{}
		} else {
			commitLines = append(commitLines, line)
		}
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for git log:", err)
		return
	}
	// Group commits by date and calculate pair days
	pairDays, devList := calculatePairDays(commits)

	// Generate HTML output
	html := "<html><head><title>Pair Programming Stats</title><style>body,table,th,td{font-family:sans-serif;} table,th,td{border:1px solid #ccc;border-collapse:collapse;}th,td{padding:8px;}</style></head><body>"
	html += "<h1>Pair Programming Days Table (Last 2 Months)</h1>"
	html += "<table>"

	// Header row
	html += "<tr><th></th>"
	for _, d := range devList {
		html += fmt.Sprintf("<th>%s</th>", d)
	}
	html += "</tr>"

	// Data rows
	for _, rowDev := range devList {
		html += fmt.Sprintf("<tr><th>%s</th>", rowDev)
		for _, colDev := range devList {
			if rowDev == colDev {
				html += "<td>-</td"
			} else {
				pair := Pair{rowDev, colDev}
				if pair[0] > pair[1] {
					pair[0], pair[1] = pair[1], pair[0]
				}
				html += fmt.Sprintf("<td>%d</td>", pairDays[pair])
			}
		}
		html += "</tr>"
	}
	html += "</table></body></html>"

	// Write to file
	f, err := os.Create("output.html")
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(html)
	if err != nil {
		fmt.Println("Error writing HTML:", err)
		return
	}
	fmt.Println("Wrote output.html")
}

func parseCommit(lines []string) Commit {
	header := lines[0]
	header = strings.TrimSuffix(header, "END_OF_COMMIT")
	parts := strings.SplitN(header, "|", 3)
	author := strings.TrimSpace(parts[0])
	date := strings.TrimSpace(parts[1])
	message := ""
	if len(parts) > 2 {
		message = parts[2]
	}
	coAuthorRe := regexp.MustCompile(`Co-authored-by: (.+?) <.+?>`)
	coAuthors := map[string]struct{}{}
	for _, match := range coAuthorRe.FindAllStringSubmatch(message, -1) {
		if len(match) > 1 {
			coAuthors[match[1]] = struct{}{}
		}
	}
	devs := map[string]struct{}{author: {}}
	for d := range coAuthors {
		devs[d] = struct{}{}
	}
	return Commit{Date: date, Developers: devs}
}

// NewCommit constructs a Commit from a date and a list of developer names.
func NewCommit(date string, devs ...string) Commit {
	ps := Developers{}
	for _, d := range devs {
		ps[d] = struct{}{}
	}
	return Commit{Date: date, Developers: ps}
}

// calculatePairDays groups commits by date and returns pair day counts and sorted developer list
func calculatePairDays(commits []Commit) (map[Pair]int, []string) {
	commitsByDate := map[string][]Commit{}
	for _, c := range commits {
		commitsByDate[c.Date] = append(commitsByDate[c.Date], c)
	}
	pairDays := map[Pair]map[string]struct{}{} // pair -> set of dates
	allDevs := map[string]struct{}{}
	for date, dailyCommits := range commitsByDate {
		pairsToday := map[Pair]struct{}{}
		for _, c := range dailyCommits {
			devList := []string{}
			for d := range c.Developers {
				devList = append(devList, d)
				allDevs[d] = struct{}{}
			}
			sort.Strings(devList)
			for i := 0; i < len(devList); i++ {
				for j := i + 1; j < len(devList); j++ {
					pair := Pair{devList[i], devList[j]}
					pairsToday[pair] = struct{}{}
				}
			}
		}
		for pair := range pairsToday {
			if pairDays[pair] == nil {
				pairDays[pair] = map[string]struct{}{}
			}
			pairDays[pair][date] = struct{}{}
		}
	}
	result := map[Pair]int{}
	for pair, dates := range pairDays {
		result[pair] = len(dates)
	}
	devList := []string{}
	for d := range allDevs {
		devList = append(devList, d)
	}
	sort.Strings(devList)
	return result, devList
}
