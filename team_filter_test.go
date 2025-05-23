package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseTeamJSON_Valid(t *testing.T) {
	jsonData := `["alice@example.com", "Bob", "carol@example.com"]`
	var team TeamMembers
	err := json.Unmarshal([]byte(jsonData), &team)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	want := TeamMembers{"alice@example.com", "Bob", "carol@example.com"}
	if !reflect.DeepEqual(team, want) {
		t.Errorf("Expected %v, got %v", want, team)
	}
}

func TestParseTeamJSON_Empty(t *testing.T) {
	jsonData := `[]`
	var team TeamMembers
	err := json.Unmarshal([]byte(jsonData), &team)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	if len(team) != 0 {
		t.Errorf("Expected empty team, got %v", team)
	}
}

func TestFilterDevelopersByTeam(t *testing.T) {
	devList := []string{"alice@example.com", "Bob", "Carol", "dave@example.com"}
	team := TeamMembers{"alice@example.com", "Carol"}
	filtered := FilterDevelopersByTeam(devList, team)
	want := []string{"alice@example.com", "Carol"}
	if !reflect.DeepEqual(filtered, want) {
		t.Errorf("Expected %v, got %v", want, filtered)
	}
}

func TestFilterPairDaysByTeam(t *testing.T) {
	pairDays := map[Pair]int{
		{"alice@example.com", "Bob"}:              2,
		{"alice@example.com", "Carol"}:            1,
		{"Bob", "Carol"}:                          3,
		{"alice@example.com", "dave@example.com"}: 1,
	}
	team := TeamMembers{"alice@example.com", "Carol"}
	filtered := FilterPairDaysByTeam(pairDays, team)
	want := map[Pair]int{
		{"alice@example.com", "Carol"}: 1,
	}
	if !reflect.DeepEqual(filtered, want) {
		t.Errorf("Expected %v, got %v", want, filtered)
	}
}
