package external

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// APISportsProvider is the primary bulk data provider (API-Sports v3).
// It supplies fixtures, squads, competitions, and match events. Its free tier
// refreshes events every 5-10 minutes, so it is used for bulk/historical sync
// rather than sub-30s realtime.
type APISportsProvider struct {
	client *httpClient
}

// NewAPISportsProvider creates an API-Sports provider using the given API key.
func NewAPISportsProvider(apiKey, host string) *APISportsProvider {
	return &APISportsProvider{
		client: newHTTPClient(host, 15*time.Second, map[string]string{
			"x-apisports-key": apiKey,
		}),
	}
}

func (p *APISportsProvider) Name() string       { return "apisports" }
func (p *APISportsProvider) Kind() ProviderKind { return ProviderBulk }

func (p *APISportsProvider) Healthy(ctx context.Context) bool {
	// A lightweight status probe; the /status endpoint returns 200 when the
	// key is valid and quota remains.
	var out struct {
		Response struct {
			Account struct {
				Requests []struct {
					Current int `json:"current"`
					Limit   int `json:"limit"`
				} `json:"requests"`
			} `json:"account"`
		} `json:"response"`
	}
	if err := p.client.get(ctx, "/status", &out); err != nil {
		return false
	}
	return true
}

// SupportsSeason reports whether this provider can serve the given season.
// The API-Sports free plan only grants access to seasons 2022-2024; paid
// plans unlock the full range. Update this if the account plan changes.
func (p *APISportsProvider) SupportsSeason(season string) bool {
	y, err := strconv.Atoi(season)
	if err != nil {
		return false
	}
	return y >= 2022 && y <= 2024
}

// ---- Response DTOs (subset of API-Sports v3) ----

type apiFixturesResponse struct {
	Response []apiFixture `json:"response"`
}

type apiFixture struct {
	ID     int    `json:"id"`
	Date   string `json:"date"`
	League struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"league"`
	Teams struct {
		Home apiTeamRef `json:"home"`
		Away apiTeamRef `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	} `json:"goals"`
	Score struct {
		Fulltime struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"fulltime"`
	} `json:"score"`
	Fixture struct {
		Status struct {
			Short string `json:"short"`
		} `json:"status"`
		Venue struct {
			Name string `json:"name"`
		} `json:"venue"`
		Referee string `json:"referee"`
	} `json:"fixture"`
}

type apiTeamRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ---- FixtureProvider ----

func (p *APISportsProvider) GetFixtures(ctx context.Context, competitionID string, season, status string) ([]Fixture, error) {
	leagueID, _ := strconv.Atoi(competitionID)
	path := fmt.Sprintf("/fixtures?league=%d&season=%s", leagueID, season)
	if status != "" {
		path += "&status=" + status
	}
	var out apiFixturesResponse
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	fixtures := make([]Fixture, 0, len(out.Response))
	for _, f := range out.Response {
		fixtures = append(fixtures, p.mapFixture(f))
	}
	return fixtures, nil
}

func (p *APISportsProvider) GetFixture(ctx context.Context, providerFixtureID string) (*Fixture, error) {
	path := "/fixtures?id=" + providerFixtureID
	var out apiFixturesResponse
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if len(out.Response) == 0 {
		return nil, fmt.Errorf("fixture %s not found", providerFixtureID)
	}
	f := p.mapFixture(out.Response[0])
	return &f, nil
}

func (p *APISportsProvider) mapFixture(f apiFixture) Fixture {
	homeScore, awayScore := 0, 0
	if f.Goals.Home != nil {
		homeScore = *f.Goals.Home
	}
	if f.Goals.Away != nil {
		awayScore = *f.Goals.Away
	}
	date, _ := time.Parse(time.RFC3339, f.Date)
	return Fixture{
		ProviderID:    strconv.Itoa(f.ID),
		CompetitionID: f.League.ID,
		Season:        "",
		HomeTeamID:    strconv.Itoa(f.Teams.Home.ID),
		AwayTeamID:    strconv.Itoa(f.Teams.Away.ID),
		HomeTeamName:  f.Teams.Home.Name,
		AwayTeamName:  f.Teams.Away.Name,
		MatchDate:     date,
		Status:        normalizeStatus(f.Fixture.Status.Short),
		HomeScore:     homeScore,
		AwayScore:     awayScore,
		Venue:         f.Fixture.Venue.Name,
		Referee:       f.Fixture.Referee,
	}
}

// ---- LiveProvider ----

func (p *APISportsProvider) GetLiveMatches(ctx context.Context) ([]Fixture, error) {
	return p.GetFixtures(ctx, "", "", "live")
}

func (p *APISportsProvider) GetLiveEvents(ctx context.Context, providerFixtureID string) ([]MatchEvent, error) {
	path := "/fixtures/events?fixture=" + providerFixtureID
	var out struct {
		Response []struct {
			Time struct {
				Elapsed int `json:"elapsed"`
			} `json:"time"`
			Team struct {
				ID int `json:"id"`
			} `json:"team"`
			Player struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"player"`
			Type     string `json:"type"`
			Detail   string `json:"detail"`
			Comments string `json:"comments"`
		} `json:"response"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	events := make([]MatchEvent, 0, len(out.Response))
	for _, e := range out.Response {
		events = append(events, MatchEvent{
			ProviderFixtureID: providerFixtureID,
			Minute:            e.Time.Elapsed,
			EventType:         normalizeEventType(e.Type),
			TeamID:            strconv.Itoa(e.Team.ID),
			PlayerID:          strconv.Itoa(e.Player.ID),
			PlayerName:        e.Player.Name,
			Detail:            e.Detail,
			Comment:           e.Comments,
		})
	}
	return events, nil
}

// ---- SquadProvider ----

func (p *APISportsProvider) GetTeam(ctx context.Context, providerTeamID string) (*Club, error) {
	path := "/teams?id=" + providerTeamID
	var out struct {
		Response []struct {
			Team struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Code    string `json:"code"`
				Country string `json:"country"`
				Logo    string `json:"logo"`
				Founded int    `json:"founded"`
			} `json:"team"`
			Venue struct {
				Name string `json:"name"`
			} `json:"venue"`
		} `json:"response"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if len(out.Response) == 0 {
		return nil, fmt.Errorf("team %s not found", providerTeamID)
	}
	t := out.Response[0]
	return &Club{
		ProviderID: strconv.Itoa(t.Team.ID),
		Name:       t.Team.Name,
		ShortName:  t.Team.Code,
		Country:    t.Team.Country,
		LogoURL:    t.Team.Logo,
		Stadium:    t.Venue.Name,
		Founded:    t.Team.Founded,
	}, nil
}

func (p *APISportsProvider) GetSquad(ctx context.Context, providerTeamID string) ([]Player, error) {
	path := "/players/squads?team=" + providerTeamID
	var out struct {
		Response []struct {
			Team struct {
				ID int `json:"id"`
			} `json:"team"`
			Players []struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Age      int    `json:"age"`
				Number   int    `json:"number"`
				Position string `json:"position"`
				Photo    string `json:"photo"`
			} `json:"players"`
		} `json:"response"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if len(out.Response) == 0 {
		return nil, fmt.Errorf("squad for team %s not found", providerTeamID)
	}
	players := make([]Player, 0, len(out.Response[0].Players))
	for _, pl := range out.Response[0].Players {
		players = append(players, Player{
			ProviderID:   strconv.Itoa(pl.ID),
			ClubID:       providerTeamID,
			Name:         pl.Name,
			Position:     pl.Position,
			JerseyNumber: pl.Number,
			PhotoURL:     pl.Photo,
		})
	}
	return players, nil
}

// ---- CompetitionProvider ----

func (p *APISportsProvider) GetCompetitions(ctx context.Context, season string) ([]Competition, error) {
	path := "/leagues?season=" + season
	var out struct {
		Response []struct {
			League struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				Logo string `json:"logo"`
			} `json:"league"`
			Country struct {
				Name string `json:"name"`
			} `json:"country"`
		} `json:"response"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	comps := make([]Competition, 0, len(out.Response))
	for _, c := range out.Response {
		comps = append(comps, Competition{
			ProviderID: strconv.Itoa(c.League.ID),
			Name:       c.League.Name,
			Type:       c.League.Type,
			Country:    c.Country.Name,
			LogoURL:    c.League.Logo,
			Season:     season,
		})
	}
	return comps, nil
}

// normalizeStatus maps API-Sports status short codes to our canonical set.
func normalizeStatus(short string) string {
	switch short {
	case "NS", "TBD", "PST":
		return "scheduled"
	case "1H", "2H", "HT", "ET", "P", "BT", "LIVE":
		return "live"
	case "FT", "AET", "PEN":
		return "finished"
	case "CANC", "ABD", "SUSP", "INT", "AWD", "WO":
		return "postponed"
	default:
		return "scheduled"
	}
}

// normalizeEventType maps API-Sports event types to our canonical set.
func normalizeEventType(t string) string {
	switch t {
	case "Goal":
		return "goal"
	case "Card":
		return "card"
	case "subst":
		return "substitution"
	case "Var":
		return "var"
	default:
		return "event"
	}
}
