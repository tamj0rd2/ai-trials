# Pair Programming Statistics CLI Tool

## Overview
This tool analyzes git commit history to determine how many days each developer has pair programmed with others, based on explicit "Co-authored-by" lines in commit messages. It generates an HTML table showing these statistics for the last 2 months.

## Prerequisites
- Go 1.20 or later
- A local git repository (the tool analyzes the current working directory)

## How to Run
1. **Clone or download this repository.**
2. **Build the tool:**
   ```sh
   go build -o pairstairs
   ```
3. **Run the tool:**
   ```sh
   ./pairstairs
   ```
   This will analyze the git history in the current directory and generate `output.html` with the pair programming statistics.

## Filtering by Team
You can limit the output to a specific team by providing a JSON file containing a list of developer names and/or email addresses.

### Example `my_team.json`:
```json
[
  "alice@example.com",
  "Bob",
  "carol@example.com"
]
```

### Run with Team Filter
```sh
./pairstairs --team my_team.json
```
This will restrict the output table to only the developers listed in `my_team.json`.

- The file must be a valid JSON array of strings (names or emails).
- If the `--team` flag is not provided, all developers found in the repository will be included in the output.

## Output
- The tool generates `output.html` in the current directory.
- The file will automatically open in your default browser (on macOS).

## Troubleshooting
- Ensure you are running the tool in a directory with a valid git repository.
- If you provide a team file, make sure it is valid JSON and contains a list of developer names or emails.

## License
MIT

