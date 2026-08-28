package external

import (
	"context"
	"time"
)

// ProviderKind categorises a data provider by its role in the ingestion
// pipeline. This lets the ingestion service pick the right source for a task
// and fail over when a provider is rate-limited or unavailable.
type ProviderKind string

const (
	// ProviderBulk is used for low-frequency, quota-heavy syncs (fixtures,
	// squads, competitions, historical data). e.g. API-Sports, TheStatsAPI.
	ProviderBulk ProviderKind = "bulk"
	// ProviderRealtime is used for high-frequency, low-quota live updates
	// (scores, events during a match). e.g. FlashScore.
	ProviderRealtime ProviderKind = "realtime"
	// ProviderMedia is used for enrichment assets (logos, photos, bios).
	// e.g. TheSportsDB.
	ProviderMedia ProviderKind = "media"
)

// Provider is the common interface every football data source implements.
// The ingestion service consumes providers through this interface so the app
// is decoupled from any single vendor and can fail over between them.
type Provider interface {
	// Name returns a stable identifier for the provider (e.g. "apisports").
	Name() string
	// Kind returns the provider's role in the pipeline.
	Kind() ProviderKind
	// Healthy reports whether the provider is currently usable.
	Healthy(ctx context.Context) bool
}

// FixtureProvider supplies fixture/match data (bulk or realtime).
type FixtureProvider interface {
	Provider
	// GetFixtures returns fixtures for a competition/season, optionally
	// filtered by status (e.g. "scheduled", "live", "finished"). The
	// competitionID is the provider-specific identifier (e.g. "140" for
	// API-Sports, "PD" for Football-Data).
	GetFixtures(ctx context.Context, competitionID string, season string, status string) ([]Fixture, error)
	// GetFixture returns a single fixture by its provider-specific ID.
	GetFixture(ctx context.Context, providerFixtureID string) (*Fixture, error)
}

// LiveProvider supplies realtime match events (goals, cards, subs).
type LiveProvider interface {
	Provider
	// GetLiveMatches returns all currently-live matches.
	GetLiveMatches(ctx context.Context) ([]Fixture, error)
	// GetLiveEvents returns the events for a live match.
	GetLiveEvents(ctx context.Context, providerFixtureID string) ([]MatchEvent, error)
}

// SquadProvider supplies club and player data.
type SquadProvider interface {
	Provider
	// GetTeam returns a club by its provider-specific ID.
	GetTeam(ctx context.Context, providerTeamID string) (*Club, error)
	// GetSquad returns the players of a club.
	GetSquad(ctx context.Context, providerTeamID string) ([]Player, error)
}

func (s SquadProvider) GetSquad(ctx context.Context, param any) ([]Player, error) {
	panic("unimplemented")
}

func (s SquadProvider) GetSquad(ctx context.Context, param any) ([]Player, error) {
	panic("unimplemented")
}

// CompetitionProvider supplies league/cup metadata.
type CompetitionProvider interface {
	Provider
	// GetCompetitions returns available competitions for a season.
	GetCompetitions(ctx context.Context, season string) ([]Competition, error)
}

// SeasonProvider is implemented by providers that can only serve certain
// seasons (e.g. API-Sports free tier is limited to 2022-2024, Football-Data
// free tier serves the current season). The registry uses this to pick the
// right provider for a requested season.
type SeasonProvider interface {
	// SupportsSeason reports whether the provider can serve the given season.
	SupportsSeason(season string) bool
}

// MediaProvider supplies enrichment assets (logos, photos, bios).
type MediaProvider interface {
	Provider
	// GetTeamLogo returns the logo URL for a club.
	GetTeamLogo(ctx context.Context, providerTeamID string) (string, error)
	// GetPlayerPhoto returns the photo URL for a player.
	GetPlayerPhoto(ctx context.Context, providerPlayerID string) (string, error)
}

// Fixture is a provider-agnostic representation of a match.
type Fixture struct {
	Provider        string // which provider produced this ("apisports", "footballdata")
	ProviderID      string // provider-specific fixture ID
	CompetitionID   int    // provider-specific numeric competition ID (API-Sports)
	CompetitionCode string // provider-specific competition code (Football-Data)
	Season          string
	HomeTeamID      string // provider-specific team ID
	AwayTeamID      string // provider-specific team ID
	HomeTeamName    string
	AwayTeamName    string
	MatchDate       time.Time
	Status          string // scheduled | live | finished | postponed
	HomeScore       int
	AwayScore       int
	Minute          int
	Venue           string
	Referee         string
}

// MatchEvent is a provider-agnostic representation of an in-match event.
type MatchEvent struct {
	ProviderFixtureID string
	Minute            int
	EventType         string // goal | card | substitution | injury | own_goal | penalty
	TeamID            string // provider-specific team ID
	PlayerID          string // provider-specific player ID
	PlayerName        string
	Detail            string // e.g. "Yellow Card", "Normal Goal"
	Comment           string
}

// Club is a provider-agnostic representation of a football club.
type Club struct {
	Provider      string // which provider produced this
	ProviderID    string
	Name          string
	ShortName     string
	Country       string
	CompetitionID int
	LogoURL       string
	Stadium       string
	Colors        string
	Founded       int
}

// Player is a provider-agnostic representation of a footballer.
type Player struct {
	ProviderID   string
	ClubID       string // provider-specific club ID
	Name         string
	Position     string
	JerseyNumber int
	Nationality  string
	BirthDate    *time.Time
	PhotoURL     string
	Rating       float64
}

// MatchStat is a provider-agnostic representation of a per-team match
// statistic (possession, shots, xG, corners, fouls, cards, ...).
type MatchStat struct {
	TeamID   string // provider-specific team ID
	StatType string // possession | shots | shots_on_target | corners | fouls | yellow_cards | red_cards | xg | ...
	Value    string // provider returns values as strings ("58%", "12", "1.34")
}

// StatsProvider supplies per-match statistics.
type StatsProvider interface {
	Provider
	// GetMatchStats returns the statistics for a match.
	GetMatchStats(ctx context.Context, providerFixtureID string) ([]MatchStat, error)
}

// Competition is a provider-agnostic representation of a league/cup.
type Competition struct {
	Provider   string // which provider produced this
	ProviderID string
	Name       string
	Type       string // League | Cup
	Country    string
	LogoURL    string
	Season     string
}
