# Implementation Plan: Pair Programming Statistics CLI Tool

## 1. Parse Git Commit History
- Use a process call to `git log` to retrieve all commits from the last 2 months.
- Extract the commit author, all "Co-authored-by" lines, and the commit date from each commit message.

## 2. Build Pairing Data
- For each commit, create a set of all developers involved (author + co-authors).
- Group commits by day (using the commit date).
- For each day, create a set of all unique developer pairs who paired on that day (from all commits on that day).
- For each unique unordered pair, increment a counter representing the number of days they paired (not the number of commits).

## 3. Aggregate Results
- Build a data structure (e.g., a map of developer pairs to counts) to store the number of days each pair has worked together.

## 4. Format Output
- Collect all unique developers.
- Generate an HTML file with a table with developers as both rows and columns, and the cell values as the number of days each pair has worked together.
- Write the HTML output to a file (e.g., `output.html`).

## 5. CLI Integration
- Implement the above logic in a Go main function.
- Ensure the tool runs in the current working directory and only requires access to the local git repository.

## 6. Future-Proofing
- Structure the code to allow for easy extension to multiple repositories in the future.

## 7. Configurable Developer Teams (JSON)
- Allow the user to specify a `my_team.json` file containing a list of developer names and/or email addresses that define the team of interest.
- The file should be a simple JSON array, e.g. `["alice@example.com", "Bob"]`.
- Add a CLI flag (e.g., `--team my_team.json`) to accept this file as input.
- When the flag is provided, filter the output so that only developers listed in `my_team.json` appear in the HTML table (both rows and columns).
- Ensure that all pairing statistics are limited to pairs where both developers are in the specified team.
- If the flag is not provided, default to showing all developers as before.
- Validate the contents of `my_team.json` and handle errors gracefully (e.g., file not found, invalid format).

## 8. Highlight Pair Recency in HTML Table
- Track the most recent date each pair worked together while processing commits.
- For each pair, determine the recency category based on the most recent pairing date:
  - Green: paired in the last 2 weeks
  - Orange: paired in the last month (but not last 2 weeks)
  - Red: paired over a month ago
- Extend the template data structure to include a map of recency classes for each pair (e.g., map[string]map[string]string).
- Update the HTML template to apply a CSS class to each cell based on the recency category.
- Add CSS styles for green, orange, and red highlighting.
- Test to ensure correct highlighting for all recency categories.

