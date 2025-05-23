package main

// FilterDevelopersByTeam returns only developers present in the team
func FilterDevelopersByTeam(devList []string, team TeamMembers) []string {
	teamSet := make(map[string]struct{})
	for _, dev := range team {
		teamSet[dev] = struct{}{}
	}
	var filtered []string
	for _, dev := range devList {
		if _, ok := teamSet[dev]; ok {
			filtered = append(filtered, dev)
		}
	}
	return filtered
}

// FilterPairDaysByTeam returns only pairs where both developers are in the team
func FilterPairDaysByTeam(pairDays map[Pair]int, team TeamMembers) map[Pair]int {
	teamSet := make(map[string]struct{})
	for _, dev := range team {
		teamSet[dev] = struct{}{}
	}
	filtered := make(map[Pair]int)
	for pair, days := range pairDays {
		if _, ok1 := teamSet[pair[0]]; ok1 {
			if _, ok2 := teamSet[pair[1]]; ok2 {
				filtered[pair] = days
			}
		}
	}
	return filtered
}
