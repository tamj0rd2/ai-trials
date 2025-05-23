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

## 8. Use Flat PairOverview Data Structure
- Replace nested or map-based pair statistics with a flat slice of PairOverview structs for simplicity and clarity.
- Define a PairOverview struct with the following fields:
  - Dev1: string (first developer)
  - Dev2: string (second developer)
  - LastPaired: string or time.Time (date they last paired)
  - DaysPaired: int (number of days paired in the last 2 months)
- When processing commits, for each unique pair, record the most recent date they paired and the total number of days paired.
- Use this slice as the main data structure for passing pair statistics to the template and for any further processing or filtering.

## 9. Highlight Pair Recency in HTML Table (Updated)
- Pass the slice of PairOverview to the template for rendering.
- In the template, derive the recency category (green/orange/red) from LastPaired.
- Update the HTML template to iterate over the slice and display the relevant information, applying CSS classes for recency highlighting.
- Add CSS styles for green, orange, and red highlighting.
- Test to ensure correct highlighting and data display for all recency categories.
