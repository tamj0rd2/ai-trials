package main

type Developers map[string]struct{}
type Pair [2]string

type TeamMembers []string

type Commit struct {
	Date       string
	Developers Developers
}
