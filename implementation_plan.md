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

