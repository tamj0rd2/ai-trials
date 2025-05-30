# Pair Programming Statistics CLI Tool: Idea & Specification

## Idea
We want to build an application that can help us record how many times developers on our team have pair programmed with other developers. We can use git coauthors to figure this out.

---

## Specification

### Overview
This command line application analyzes git commit history to determine how many days each developer has pair programmed with other developers, based on explicit "Co-authored-by" lines in commit messages.

### Features
- **Command Line Tool**: Runs in the terminal and generates an HTML file showing the number of days each developer has paired with every other developer.
- **Commit Analysis Window**: Only analyzes commits from the last 2 months, using a rolling window based on the current date.
- **Pair Detection**: Detects pairs using explicit "Co-authored-by" lines in git commit messages. No inference from other metadata.
- **Developer Equality**: Does not distinguish between authors and co-authors; all listed developers in a commit are treated equally for pairing statistics.
- **Output**: Generates an HTML file with a table showing the number of days each developer has paired with every other developer.
- **Repository Scope**: The tool analyzes the git repository in the current working directory. (Note: In the future, support for aggregating statistics across multiple repositories may be added.)

### Out of Scope
- No visualization beyond the HTML table.
- No aggregation across multiple repositories (future enhancement).
- No privacy or data retention features beyond what is required for local analysis.

### Implementation Notes
- The tool should be implemented in Go.
- The tool should be runnable locally by any developer with access to the repository.
- The tool should only require access to the local git repository in the current working directory.
- The output should be an HTML file formatted for easy human readability.

### Future Improvements
- **Deduplicate Developer Names**: Ensure that developer names are not duplicated in the output, even if they appear with different email addresses or minor variations.
- **Configurable Developer Teams**: Allow configuration of which developers or teams to display in the output. The repository may contain all developer names, but users should be able to filter or group by teams.

