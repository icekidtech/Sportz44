package external

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FootballDataProvider is a bulk provider backed by football-data.org (v4).
// Its free tier serves the CURRENT season for all major competitions (La Liga,
// Champions League, etc.) with no historical-season lockout, which complements
// API-Sports (whose free tier is limited to 2022-2024). It uses competition
// codes (e.g. "PD" for La Liga) rather than numeric IDs.
type FootballDataProvider struct {
	client *httpClient
}

// NewFootballDataProvider creates a Football-Data provider using the given API
// token (register at https://www.football-data.org to get a free one).
func NewFootballDataProvider(apiToken, baseURL string) *FootballDataProvider {
	return &FootballDataProvider{
		client: newHTTPClient(baseURL, 15*time.Second, map[string]string{
			"X-Auth-Token": apiToken,
		}),
	}
}

func (p *FootballDataProvider) Name() string       { return "footballdata" }
func (p *FootballDataProvider) Kind() ProviderKind { return ProviderBulk }

func (p *FootballDataProvider) Healthy(ctx context.Context) bool {
	// A lightweight probe: the /competitions endpoint returns 200 with a
	// non-empty list when the token is valid.
	var out struct {
		Competitions []jsonCompetition `json:"competitions"`
	}
	if err := p.client.get(ctx, "/competitions", &out); err != nil {
		return false
	}
	return len(out.Competitions) > 0
}

// SupportsSeason reports whether this provider can serve the given season.
// The free tier only exposes the current season (e.g. "2026" for 2026/27).
func (p *FootballDataProvider) SupportsSeason(season string) bool {
	return season == currentSeason()
}

// currentSeason returns the start year of the current European season
// (seasons run August-June, so the season year is the start year).
func currentSeason() string {
	now := time.Now()
	if now.Month() >= 7 {
		return strconv.Itoa(now.Year())
	}
	return strconv.Itoa(now.Year() - 1)
}

// CurrentSeason returns the current European season start year.
func CurrentSeason() string {
	return currentSeason()
}

// ---- Response DTOs (subset of football-data.org v4) ----

type jsonCompetition struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Type   string `json:"type"`
	Emblem string `json:"emblem"`
	Area   struct {
		Name string `json:"name"`
	} `json:"area"`
	CurrentSeason struct {
		StartDate string `json:"startDate"`
	} `json:"currentSeason"`
}

type jsonTeam struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	TLA       string `json:"tla"`
	Crest     string `json:"crest"`
	Founded   int    `json:"founded"`
	Venue     string `json:"venue"`
	Area      struct {
		Name string `json:"name"`
	} `json:"area"`
	Squad []jsonPlayer `json:"squad"`
}

type jsonPlayer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
	ShirtNumber int    `json:"shirtNumber"`
}

type jsonMatch struct {
	ID       int    `json:"id"`
	UTCDate  string `json:"utcDate"`
	Status   string `json:"status"`
	Venue    string `json:"venue"`
	HomeTeam struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"homeTeam"`
	AwayTeam struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"awayTeam"`
	Score struct {
		FullTime struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"fullTime"`
	} `json:"score"`
	Referees []struct {
		Name string `json:"name"`
	} `json:"referees"`
}

type jsonMatchesResponse struct {
	Competition struct {
		Code string `json:"code"`
	} `json:"competition"`
	Matches []jsonMatch `json:"matches"`
}

// ---- CompetitionProvider ----

func (p *FootballDataProvider) GetCompetitions(ctx context.Context, season string) ([]Competition, error) {
	var out struct {
		Competitions []jsonCompetition `json:"competitions"`
	}
	if err := p.client.get(ctx, "/competitions", &out); err != nil {
		return nil, err
	}
	comps := make([]Competition, 0, len(out.Competitions))
	for _, c := range out.Competitions {
		compSeason := season
		if compSeason == "" {
			compSeason = seasonFromStartDate(c.CurrentSeason.StartDate)
		}
		comps = append(comps, Competition{
			Provider:   p.Name(),
			ProviderID: c.Code,
			Name:       c.Name,
			Type:       normalizeCompetitionType(c.Type),
			Country:    c.Area.Name,
			LogoURL:    c.Emblem,
			Season:     compSeason,
		})
	}
	return comps, nil
}

// ---- FixtureProvider ----

func (p *FootballDataProvider) GetFixtures(ctx context.Context, competitionID string, season, status string) ([]Fixture, error) {
	path := fmt.Sprintf("/competitions/%s/matches", competitionID)
	if season != "" {
		path += "?season=" + season
	}
	if s := mapFDStatusFilter(status); s != "" {
		if strings.Contains(path, "?") {
			path += "&status=" + s
		} else {
			path += "?status=" + s
		}
	}
	var out jsonMatchesResponse
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return p.mapMatches(out, competitionID, season), nil
}

func (p *FootballDataProvider) GetFixture(ctx context.Context, providerFixtureID string) (*Fixture, error) {
	var out struct {
		jsonMatch
		Competition struct {
			Code string `json:"code"`
		} `json:"competition"`
	}
	if err := p.client.get(ctx, "/matches/"+providerFixtureID, &out); err != nil {
		return nil, err
	}
	f := p.mapMatch(out.jsonMatch, out.Competition.Code, "")
	return &f, nil
}

func (p *FootballDataProvider) mapMatches(out jsonMatchesResponse, code, season string) []Fixture {
	fixtures := make([]Fixture, 0, len(out.Matches))
	for _, m := range out.Matches {
		fixtures = append(fixtures, p.mapMatch(m, code, season))
	}
	return fixtures
}

func (p *FootballDataProvider) mapMatch(m jsonMatch, code, season string) Fixture {
	homeScore, awayScore := 0, 0
	if m.Score.FullTime.Home != nil {
		homeScore = *m.Score.FullTime.Home
	}
	if m.Score.FullTime.Away != nil {
		awayScore = *m.Score.FullTime.Away
	}
	referee := ""
	if len(m.Referees) > 0 {
		referee = m.Referees[0].Name
	}
	matchDate, _ := time.Parse(time.RFC3339, m.UTCDate)
	return Fixture{
		Provider:        p.Name(),
		ProviderID:      strconv.Itoa(m.ID),
		CompetitionCode: code,
		Season:          season,
		HomeTeamID:      strconv.Itoa(m.HomeTeam.ID),
		AwayTeamID:      strconv.Itoa(m.AwayTeam.ID),
		HomeTeamName:    m.HomeTeam.Name,
		AwayTeamName:    m.AwayTeam.Name,
		MatchDate:       matchDate,
		Status:          normalizeFDStatus(m.Status),
		HomeScore:       homeScore,
		AwayScore:       awayScore,
		Venue:           m.Venue,
		Referee:         referee,
	}
}

// ---- LiveProvider ----

// GetLiveEvents fetches the events for a match. Football-Data exposes events
// for any match (live or finished) via /matches/{id}/events, including the
// assist credited for each goal.
func (p *FootballDataProvider) GetLiveEvents(ctx context.Context, providerFixtureID string) ([]MatchEvent, error) {
	var out struct {
		Events []struct {
			Type   string `json:"type"`
			Detail string `json:"detail"`
			Time   struct {
				Minute    int  `json:"minute"`
				ExtraTime *int `json:"extraTime"`
			} `json:"time"`
			Team struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"team"`
			Player struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"player"`
			Assist struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"assist"`
		} `json:"events"`
	}
	if err := p.client.get(ctx, "/matches/"+providerFixtureID+"/events", &out); err != nil {
		return nil, err
	}
	events := make([]MatchEvent, 0, len(out.Events))
	for _, e := range out.Events {
		events = append(events, MatchEvent{
			ProviderFixtureID: providerFixtureID,
			Minute:            e.Time.Minute,
			EventType:         normalizeFDEventType(e.Type),
			TeamID:            strconv.Itoa(e.Team.ID),
			PlayerID:          strconv.Itoa(e.Player.ID),
			PlayerName:        e.Player.Name,
			AssistPlayerID:    strconv.Itoa(e.Assist.ID),
			AssistPlayerName:  e.Assist.Name,
			Detail:            e.Detail,
		})
	}
	return events, nil
}

// normalizeFDEventType maps Football-Data event types to our canonical set.
func normalizeFDEventType(t string) string {
	switch t {
	case "GOAL":
		return "goal"
	case "OWN_GOAL":
		return "own_goal"
	case "PENALTY_SHOOTOUT_GOAL":
		return "goal"
	case "YELLOW_CARD":
		return "card"
	case "RED_CARD":
		return "card"
	case "SUBSTITUTION":
		return "substitution"
	default:
		return "event"
	}
}

// ---- SquadProvider ----

func (p *FootballDataProvider) GetTeam(ctx context.Context, providerTeamID string) (*Club, error) {
	var out jsonTeam
	if err := p.client.get(ctx, "/teams/"+providerTeamID, &out); err != nil {
		return nil, err
	}
	return &Club{
		Provider:   p.Name(),
		ProviderID: strconv.Itoa(out.ID),
		Name:       out.Name,
		ShortName:  out.TLA,
		Country:    out.Area.Name,
		LogoURL:    out.Crest,
		Stadium:    out.Venue,
		Founded:    out.Founded,
	}, nil
}

func (p *FootballDataProvider) GetSquad(ctx context.Context, providerTeamID string) ([]Player, error) {
	var out jsonTeam
	if err := p.client.get(ctx, "/teams/"+providerTeamID, &out); err != nil {
		return nil, err
	}
	players := make([]Player, 0, len(out.Squad))
	for _, pl := range out.Squad {
		var birth *time.Time
		if t, err := time.Parse("2006-01-02", pl.DateOfBirth); err == nil {
			birth = &t
		}
		players = append(players, Player{
			ProviderID:   strconv.Itoa(pl.ID),
			ClubID:       strconv.Itoa(out.ID),
			Name:         pl.Name,
			Position:     normalizeFDPosition(pl.Position),
			JerseyNumber: pl.ShirtNumber,
			Nationality:  pl.Nationality,
			BirthDate:    birth,
		})
	}
	return players, nil
}

// ---- helpers ----

func seasonFromStartDate(startDate string) string {
	if t, err := time.Parse("2006-01-02", startDate); err == nil {
		return strconv.Itoa(t.Year())
	}
	return currentSeason()
}

func normalizeCompetitionType(t string) string {
	switch strings.ToUpper(t) {
	case "CUP":
		return "Cup"
	default:
		return "League"
	}
}

func normalizeFDStatus(s string) string {
	switch strings.ToUpper(s) {
	case "FINISHED", "AWARDED":
		return "finished"
	case "LIVE", "IN_PLAY", "PAUSED":
		return "live"
	case "POSTPONED":
		return "postponed"
	case "CANCELLED":
		return "cancelled"
	case "SUSPENDED":
		return "suspended"
	default: // SCHEDULED, TIMED
		return "scheduled"
	}
}

// mapFDStatusFilter converts our status filter to a football-data status.
func mapFDStatusFilter(status string) string {
	switch status {
	case "scheduled":
		return "SCHEDULED"
	case "live":
		return "LIVE"
	case "finished":
		return "FINISHED"
	case "postponed":
		return "POSTPONED"
	default:
		return ""
	}
}

func normalizeFDPosition(pos string) string {
	switch strings.ToLower(pos) {
	case "goalkeeper":
		return "Goalkeeper"
	case "defence", "defender":
		return "Defender"
	case "midfield", "midfielder":
		return "Midfielder"
	case "offence", "attacker", "forward":
		return "Forward"
	default:
		return pos
	}
}
