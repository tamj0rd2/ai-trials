package main

import (
	"reflect"
	"testing"
)

func TestPairingLogic_SingleDaySinglePair(t *testing.T) {
	commits := []Commit{
		NewCommit("2025-05-01", "Alice", "Bob"),
	}
	pairDays, devList := calculatePairDays(commits)
	if pairDays[[2]string{"Alice", "Bob"}] != 1 {
		t.Errorf("Expected 1 day for Alice/Bob, got %d", pairDays[[2]string{"Alice", "Bob"}])
	}
	if len(devList) != 2 || devList[0] != "Alice" || devList[1] != "Bob" {
		t.Errorf("Unexpected devList: %v", devList)
	}
}

func TestPairingLogic_MultipleDays(t *testing.T) {
	commits := []Commit{
		NewCommit("2025-05-01", "Alice", "Bob"),
		NewCommit("2025-05-02", "Alice", "Bob"),
	}
	pairDays, _ := calculatePairDays(commits)
	if pairDays[[2]string{"Alice", "Bob"}] != 2 {
		t.Errorf("Expected 2 days for Alice/Bob, got %d", pairDays[[2]string{"Alice", "Bob"}])
	}
}

func TestPairingLogic_MultiplePairsSameDay(t *testing.T) {
	commits := []Commit{
		NewCommit("2025-05-01", "Alice", "Bob", "Carol"),
	}
	pairDays, _ := calculatePairDays(commits)
	expected := map[[2]string]int{
		{"Alice", "Bob"}:   1,
		{"Alice", "Carol"}: 1,
		{"Bob", "Carol"}:   1,
	}
	for pair, want := range expected {
		if got := pairDays[pair]; got != want {
			t.Errorf("Expected %d days for %v, got %d", want, pair, got)
		}
	}
}

func TestPairingLogic_MultipleCommitsSameDay(t *testing.T) {
	commits := []Commit{
		NewCommit("2025-05-01", "Alice", "Bob"),
		NewCommit("2025-05-01", "Alice", "Carol"),
	}
	pairDays, _ := calculatePairDays(commits)
	expected := map[Pair]int{
		Pair{"Alice", "Bob"}:   1,
		Pair{"Alice", "Carol"}: 1,
	}
	for pair, want := range expected {
		if got := pairDays[[2]string{pair[0], pair[1]}]; got != want {
			t.Errorf("Expected %d days for %v, got %d", want, pair, got)
		}
	}
	if _, exists := pairDays[[2]string{"Bob", "Carol"}]; exists {
		t.Errorf("Did not expect Bob/Carol to be counted as a pair")
	}
}

func TestPairingLogic_NoPairs(t *testing.T) {
	commits := []Commit{
		NewCommit("2025-05-01", "Alice"),
	}
	pairDays, devList := calculatePairDays(commits)
	if len(pairDays) != 0 {
		t.Errorf("Expected 0 pairs, got %v", pairDays)
	}
	if !reflect.DeepEqual(devList, []string{"Alice"}) {
		t.Errorf("Unexpected devList: %v", devList)
	}
}

func TestParseCommit_NoCoAuthor(t *testing.T) {
	lines := []string{
		"Alice|2025-05-01|Initial commit",
		"END_OF_COMMIT",
	}
	commit := parseCommit(lines)
	wantDevs := Developers{"Alice": {}}
	if !reflect.DeepEqual(commit.Developers, wantDevs) {
		t.Errorf("Expected only Alice, got %v", commit.Developers)
	}
	if commit.Date != "2025-05-01" {
		t.Errorf("Expected date 2025-05-01, got %s", commit.Date)
	}
}

func TestParseCommit_WithCoAuthor(t *testing.T) {
	lines := []string{
		"Bob|2025-05-02|Add feature",
		"",
		"Co-authored-by: Carol <carol@example.com>",
		"END_OF_COMMIT",
	}
	commit := parseCommit(lines)
	expected := Developers{"Bob": {}, "Carol": {}}
	if !reflect.DeepEqual(commit.Developers, expected) {
		t.Errorf("Expected Bob and Carol, got %v", commit.Developers)
	}
	if commit.Date != "2025-05-02" {
		t.Errorf("Expected date 2025-05-02, got %s", commit.Date)
	}
}

func TestParseCommit_MultipleCoAuthors(t *testing.T) {
	lines := []string{
		"Alice|2025-05-03|Refactor",
		"",
		"Co-authored-by: Bob <bob@example.com>",
		"Co-authored-by: Carol <carol@example.com>",
		"END_OF_COMMIT",
	}
	commit := parseCommit(lines)
	expected := Developers{"Alice": {}, "Bob": {}, "Carol": {}}
	if !reflect.DeepEqual(commit.Developers, expected) {
		t.Errorf("Expected Alice, Bob, and Carol, got %v", commit.Developers)
	}
	if commit.Date != "2025-05-03" {
		t.Errorf("Expected date 2025-05-03, got %s", commit.Date)
	}
}

func TestParseCommit_WeirdSpacing(t *testing.T) {
	lines := []string{
		"Alice|2025-05-04|Fix bug",
		"",
		"  Co-authored-by: Bob <bob@example.com>  ",
		"END_OF_COMMIT",
	}
	commit := parseCommit(lines)
	expected := Developers{"Alice": {}, "Bob": {}}
	if !reflect.DeepEqual(commit.Developers, expected) {
		t.Errorf("Expected Alice and Bob, got %v", commit.Developers)
	}
}

func TestParseCommit_MultiLineRealistic(t *testing.T) {
	lines := []string{
		"Jane Doe|2025-05-22|Implement new feature X",
		"",
		"Co-authored-by: John Smith <john.smith@example.com>",
		"Co-authored-by: Alex Lee <alex.lee@example.com>",
		"END_OF_COMMIT",
	}
	commit := parseCommit(lines)
	expected := Developers{"Jane Doe": {}, "John Smith": {}, "Alex Lee": {}}
	if !reflect.DeepEqual(commit.Developers, expected) {
		t.Errorf("got %v", commit.Developers)
	}
}
