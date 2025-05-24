# Function Ordering

- Place all public (exported) functions, including `main`, at the top of your Go files, immediately after imports and type definitions.
- This provides a high-level overview of the program's API and flow, making it easier for developers to understand what the code does at a glance.
- Helper and unexported (private) functions should follow the public functions.

## Rationale

- Developers are most interested in the overall structure, entry points, and API of the program.
- Details can be explored as needed by reading further into the file.

# Function Length and Focus

**Instruction:**

- Keep functions short and focused. Each function should do one thing and do it well. If a function is getting long or is handling multiple responsibilities, break it up into smaller, well-named helper functions. This improves readability, maintainability, and testability.

**Example (Anti-pattern):**

The following `main` function is too large and handles too many concerns at once:

```go
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
				html += "<td>-</td>"
			} else {
				pair := Pair{rowDev, colDev}
				if pair[0] > pair[1] {
					pair[0], pair[1] = pair[1], pair[0]
				}
				count := 0
				if v, ok := pairDays[pair]; ok {
					count = v
				}
				html += fmt.Sprintf("<td>%d</td>", count)
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
	// Open the created file
	if err := exec.Command("open", "output.html").Start(); err != nil {
		fmt.Println("Error opening HTML file:", err)
	}
}
```

**Better Alternative:**

Break the `main` function into smaller, focused helpers:

```go
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

func getCommitsFromGit() ([]Commit, error) { /* ... */ }
func generateHTMLTable(pairDays map[Pair]int, devList []string) string { /* ... */ }
func writeAndOpenHTML(html, filename string) error { /* ... */ }
```

This approach makes each function easier to read, test, and maintain.

# Self-explanatory Functions Over Comments

**Instruction:**

Functions should be self-explanatory. If a comment is needed, it's a sign that a new function should probably be introduced.

This isn't always the case, but it mostly is. A case where you might want comments is when writing complicated algorithms that people may find difficult to understand.

## Rationale

- Well-named functions serve as their own documentation
- Breaking complex operations into aptly named smaller functions improves readability
- Code should explain "how", function names should explain "what"
- Reserve comments for particularly complex algorithms or edge cases that cannot be made clear through refactoring

**Example (Anti-pattern):**

```go
func process(data []int) int {
    // Calculate the sum of all even numbers in the array
    sum := 0
    for _, v := range data {
        if v%2 == 0 {
            sum += v
        }
    }
    return sum
}
```

**Better Alternative:**

```go
func sumEvenNumbers(data []int) int {
    sum := 0
    for _, v := range data {
        if v%2 == 0 {
            sum += v
        }
    }
    return sum
}
```

The improved function name makes the comment unnecessary and makes the code's purpose immediately clear.

# Workflow for Implementation Requests

- When the user requests a change, first describe how you would implement it, including which files or areas of the codebase you would touch and why.
- Ask if this approach can be added to the implementation plan.
- Only update the implementation plan after receiving user approval.
- Wait for explicit user confirmation before making any code changes.
- You don't need to follow this entire process for changes to guidelines.md. You just need my approval to edit this file.

