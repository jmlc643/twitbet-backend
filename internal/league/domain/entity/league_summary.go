package entity

import "github.com/google/uuid"

type LeagueRole string

const (
	RoleAdmin  LeagueRole = "ADMIN"
	RoleMember LeagueRole = "MIEMBRO"
)

type LeagueSummary struct {
	LeagueID         uuid.UUID
	Slug             string
	Name             string
	Role             LeagueRole
	ParticipantCount int
	Balance          float64
	Status           string
}
