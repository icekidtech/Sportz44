package external

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ESPNProvider is a bulk + realtime provider backed by ESPN's undocumented
// JSON API (site.api.espn.com). It requires no API key and has no documented
// rate limit, and it serves the CURRENT season for ~140 leagues and cups
// worldwide with full match events (goals + assists, cards, subs), lineups,
// and team stats. This makes it the primary source for current-season data,
// complementing API-Sports (which is limited to 2022-2024 on the free tier).
//
// NOTE: This is an unofficial API. The schema can change without notice, so
// this provider is isolated behind the Provider interface and should be
// cached aggressively. Use the /teams endpoint (404 = invalid slug) rather
// than the scoreboard status code to validate a league slug, because the
// scoreboard returns 400 for both invalid slugs and valid-but-idle leagues.
type ESPNProvider struct {
	client *httpClient
}

// NewESPNProvider creates an ESPN provider. baseURL defaults to ESPN's public
// site API if empty.
//
// NOTE: ESPN's API returns HTTP 403 when a custom User-Agent header is set
// (anti-bot protection). We deliberately send no User-Agent so the request
// looks like a normal client and returns 200.
func NewESPNProvider(baseURL string) *ESPNProvider {
	if baseURL == "" {
		baseURL = "https://site.api.espn.com/apis/site/v2/sports/soccer"
	}
	return &ESPNProvider{
		client: newHTTPClient(baseURL, 15*time.Second, map[string]string{}),
	}
}

func (p *ESPNProvider) Name() string       { return "espn" }
func (p *ESPNProvider) Kind() ProviderKind { return ProviderBulk }

func (p *ESPNProvider) Healthy(ctx context.Context) bool {
	// Probe a well-known league (Premier League). The /teams endpoint returns
	// 404 for invalid slugs and real JSON for valid ones, so a 200 here means
	// the API is reachable.
	var out struct {
		Teams []jsonTeam `json:"teams"`
	}
	return p.client.get(ctx, "/eng.1/teams", &out) == nil
}

// SupportsSeason reports whether this provider can serve the given season.
// ESPN serves the current season (and recent completed ones). It does not
// serve arbitrary historical seasons, so we accept the current season and the
// immediately preceding one.
func (p *ESPNProvider) SupportsSeason(season string) bool {
	y, err := strconv.Atoi(season)
	if err != nil {
		return false
	}
	now := time.Now()
	cur := now.Year()
	if now.Month() < 7 {
		cur--
	}
	return y == cur || y == cur-1
}

// ---- Response DTOs (subset of ESPN's site API) ----

type espnScoreboard struct {
	Leagues []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"leagues"`
	Events []espnEvent `json:"events"`
}

type espnEvent struct {
	ID           string `json:"id"`
	Date         string `json:"date"`
	Name         string `json:"name"`
	Competitions []struct {
		ID     string     `json:"id"`
		Date   string     `json:"date"`
		Status espnStatus `json:"status"`
		Venue  struct {
			FullName string `json:"fullName"`
		} `json:"venue"`
		Competitors []espnCompetitor `json:"competitors"`
		Details     []espnDetail     `json:"details"`
	} `json:"competitions"`
}

type espnStatus struct {
	Type struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		State       string `json:"state"`
		Completed   bool   `json:"completed"`
		Description string `json:"description"`
		Detail      string `json:"detail"`
	} `json:"type"`
	DisplayClock string `json:"displayClock"`
	Period       int    `json:"period"`
}

type espnCompetitor struct {
	ID       string `json:"id"`
	HomeAway string `json:"homeAway"`
	Winner   bool   `json:"winner"`
	Score    string `json:"score"`
	Team     struct {
		ID           string `json:"id"`
		DisplayName  string `json:"displayName"`
		ShortName    string `json:"shortName"`
		Abbreviation string `json:"abbreviation"`
		Location     string `json:"location"`
	} `json:"team"`
	Statistics []espnStat `json:"statistics"`
}

type espnStat struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	DisplayValue string `json:"displayValue"`
}

type espnDetail struct {
	Type struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"type"`
	Clock struct {
		DisplayValue string `json:"displayValue"`
	} `json:"clock"`
	Team struct {
		ID string `json:"id"`
	} `json:"team"`
	ScoreValue       int  `json:"scoreValue"`
	ScoringPlay      bool `json:"scoringPlay"`
	RedCard          bool `json:"redCard"`
	YellowCard       bool `json:"yellowCard"`
	PenaltyKick      bool `json:"penaltyKick"`
	OwnGoal          bool `json:"ownGoal"`
	Shootout         bool `json:"shootout"`
	AthletesInvolved []struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"athletesInvolved"`
}

// ---- FixtureProvider ----

// GetFixtures returns fixtures for a competition (league slug) and season.
// ESPN's scoreboard endpoint is date-based; we fetch the current day's games
// for the slug. The season parameter is accepted for interface compatibility
// but ESPN serves the current season only.
func (p *ESPNProvider) GetFixtures(ctx context.Context, competitionID string, season, status string) ([]Fixture, error) {
	path := "/" + competitionID + "/scoreboard"
	var out espnScoreboard
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	fixtures := make([]Fixture, 0, len(out.Events))
	for _, e := range out.Events {
		if len(e.Competitions) == 0 {
			continue
		}
		c := e.Competitions[0]
		f := p.mapFixture(e.ID, c.Date, c.Status, c.Venue.FullName, c.Competitors)
		if status != "" && f.Status != status {
			continue
		}
		fixtures = append(fixtures, f)
	}
	return fixtures, nil
}

// GetFixture returns a single fixture by its ESPN event ID.
func (p *ESPNProvider) GetFixture(ctx context.Context, providerFixtureID string) (*Fixture, error) {
	// The summary endpoint returns a single event's full detail.
	path := "/summary?event=" + providerFixtureID
	var out struct {
		Header struct {
			ID   string `json:"id"`
			Date string `json:"date"`
		} `json:"header"`
		Competitions []struct {
			ID     string     `json:"id"`
			Date   string     `json:"date"`
			Status espnStatus `json:"status"`
			Venue  struct {
				FullName string `json:"fullName"`
			} `json:"venue"`
			Competitors []espnCompetitor `json:"competitors"`
		} `json:"competitions"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if len(out.Competitions) == 0 {
		return nil, fmt.Errorf("fixture %s not found", providerFixtureID)
	}
	c := out.Competitions[0]
	f := p.mapFixture(providerFixtureID, c.Date, c.Status, c.Venue.FullName, c.Competitors)
	return &f, nil
}

func (p *ESPNProvider) mapFixture(id, date string, status espnStatus, venue string, competitors []espnCompetitor) Fixture {
	var home, away espnCompetitor
	for _, c := range competitors {
		switch c.HomeAway {
		case "home":
			home = c
		case "away":
			away = c
		}
	}
	homeScore, _ := strconv.Atoi(home.Score)
	awayScore, _ := strconv.Atoi(away.Score)
	parsed, _ := time.Parse(espnTimeLayout, date)
	return Fixture{
		Provider:     p.Name(),
		ProviderID:   id,
		HomeTeamID:   home.Team.ID,
		AwayTeamID:   away.Team.ID,
		HomeTeamName: home.Team.DisplayName,
		AwayTeamName: away.Team.DisplayName,
		MatchDate:    parsed,
		Status:       espnStatusToCanonical(status),
		HomeScore:    homeScore,
		AwayScore:    awayScore,
		Minute:       espnMinute(status),
		Venue:        venue,
	}
}

// espnTimeLayout matches ESPN's date format (e.g. "2026-08-28T19:00Z"), which
// omits seconds and is not a standard RFC3339 layout.
const espnTimeLayout = "2006-01-02T15:04Z"

// ---- LiveProvider ----

func (p *ESPNProvider) GetLiveMatches(ctx context.Context) ([]Fixture, error) {
	// ESPN has no single "all live" endpoint; we poll the major leagues'
	// scoreboards and filter to live matches. This covers the primary
	// European leagues + MLS.
	slugs := []string{"eng.1", "esp.1", "ger.1", "ita.1", "fra.1", "por.1", "ned.1", "sco.1", "usa.1", "bra.1", "arg.1"}
	var live []Fixture
	for _, slug := range slugs {
		fixtures, err := p.GetFixtures(ctx, slug, "", "live")
		if err != nil {
			continue
		}
		live = append(live, fixtures...)
	}
	return live, nil
}

// GetLiveEvents returns the events for a match by its ESPN event ID. ESPN
// embeds the events in the scoreboard's details array; we fetch the summary
// endpoint for the specific event to get its full event list.
func (p *ESPNProvider) GetLiveEvents(ctx context.Context, providerFixtureID string) ([]MatchEvent, error) {
	path := "/summary?event=" + providerFixtureID
	var out struct {
		Competitions []struct {
			Details []espnDetail `json:"details"`
		} `json:"competitions"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if len(out.Competitions) == 0 {
		return nil, fmt.Errorf("fixture %s not found", providerFixtureID)
	}
	events := make([]MatchEvent, 0, len(out.Competitions[0].Details))
	for _, d := range out.Competitions[0].Details {
		ev := MatchEvent{
			ProviderFixtureID: providerFixtureID,
			Minute:            espnMinuteFromClock(d.Clock.DisplayValue),
			EventType:         espnEventType(d),
			TeamID:            d.Team.ID,
			Detail:            d.Type.Text,
		}
		// athletesInvolved[0] is the primary player (scorer for goals, the
		// carded player for cards). For goals, athletesInvolved[1] (if any) is
		// the assister.
		if len(d.AthletesInvolved) > 0 {
			ev.PlayerID = d.AthletesInvolved[0].ID
			ev.PlayerName = d.AthletesInvolved[0].DisplayName
		}
		if ev.EventType == "goal" && len(d.AthletesInvolved) > 1 {
			ev.AssistPlayerID = d.AthletesInvolved[1].ID
			ev.AssistPlayerName = d.AthletesInvolved[1].DisplayName
		}
		events = append(events, ev)
	}
	return events, nil
}

// ---- StatsProvider ----

func (p *ESPNProvider) GetMatchStats(ctx context.Context, providerFixtureID string) ([]MatchStat, error) {
	path := "/summary?event=" + providerFixtureID
	var out struct {
		Competitions []struct {
			Competitors []espnCompetitor `json:"competitors"`
		} `json:"competitions"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	if len(out.Competitions) == 0 {
		return nil, fmt.Errorf("fixture %s not found", providerFixtureID)
	}
	stats := make([]MatchStat, 0, len(out.Competitions[0].Competitors)*8)
	for _, c := range out.Competitions[0].Competitors {
		for _, s := range c.Statistics {
			stats = append(stats, MatchStat{
				TeamID:   c.Team.ID,
				StatType: espnStatType(s.Name),
				Value:    s.DisplayValue,
			})
		}
	}
	return stats, nil
}

// ---- SquadProvider ----

func (p *ESPNProvider) GetTeam(ctx context.Context, providerTeamID string) (*Club, error) {
	// ESPN's team endpoint requires a league slug; we can't resolve a bare
	// team ID to a league without extra calls. Return a minimal club so the
	// interface is satisfied; club metadata is enriched by API-Sports.
	return &Club{
		Provider:   p.Name(),
		ProviderID: providerTeamID,
	}, nil
}

func (p *ESPNProvider) GetSquad(ctx context.Context, providerTeamID string) ([]Player, error) {
	// ESPN's roster endpoint is league-scoped; return empty for now. Squads
	// are enriched by API-Sports.
	return []Player{}, nil
}

// ---- CompetitionProvider ----

func (p *ESPNProvider) GetCompetitions(ctx context.Context, season string) ([]Competition, error) {
	// ESPN has no single "list all leagues" endpoint. We return the well-known
	// set of league slugs we support so the ingestion service can enumerate
	// them. Each is validated lazily on first use.
	known := []struct {
		slug, name, country, typ string
	}{
		{"eng.1", "Premier League", "England", "League"},
		{"eng.2", "Championship", "England", "League"},
		{"eng.fa", "FA Cup", "England", "Cup"},
		{"eng.league_cup", "Carabao Cup", "England", "Cup"},
		{"esp.1", "La Liga", "Spain", "League"},
		{"esp.2", "La Liga 2", "Spain", "League"},
		{"esp.copa_del_rey", "Copa del Rey", "Spain", "Cup"},
		{"ger.1", "Bundesliga", "Germany", "League"},
		{"ger.2", "2. Bundesliga", "Germany", "League"},
		{"ger.dfb_pokal", "DFB-Pokal", "Germany", "Cup"},
		{"ita.1", "Serie A", "Italy", "League"},
		{"ita.2", "Serie B", "Italy", "League"},
		{"ita.coppa_italia", "Coppa Italia", "Italy", "Cup"},
		{"fra.1", "Ligue 1", "France", "League"},
		{"fra.2", "Ligue 2", "France", "League"},
		{"fra.coupe_de_france", "Coupe de France", "France", "Cup"},
		{"por.1", "Primeira Liga", "Portugal", "League"},
		{"por.taca.portugal", "Taça de Portugal", "Portugal", "Cup"},
		{"ned.1", "Eredivisie", "Netherlands", "League"},
		{"ned.cup", "KNVB Beker", "Netherlands", "Cup"},
		{"sco.1", "Premiership", "Scotland", "League"},
		{"sco.tennents", "Scottish Cup", "Scotland", "Cup"},
		{"bel.1", "Belgian Pro League", "Belgium", "League"},
		{"tur.1", "Süper Lig", "Turkey", "League"},
		{"sui.1", "Swiss Super League", "Switzerland", "League"},
		{"aut.1", "Austrian Bundesliga", "Austria", "League"},
		{"gre.1", "Greek Super League", "Greece", "League"},
		{"den.1", "Superliga", "Denmark", "League"},
		{"swe.1", "Allsvenskan", "Sweden", "League"},
		{"nor.1", "Eliteserien", "Norway", "League"},
		{"cze.1", "Czech First League", "Czechia", "League"},
		{"rou.1", "Liga 1", "Romania", "League"},
		{"rus.1", "Russian Premier League", "Russia", "League"},
		{"bra.1", "Série A", "Brazil", "League"},
		{"bra.2", "Série B", "Brazil", "League"},
		{"bra.copa_do_brazil", "Copa do Brasil", "Brazil", "Cup"},
		{"arg.1", "Liga Profesional", "Argentina", "League"},
		{"arg.copa", "Copa Argentina", "Argentina", "Cup"},
		{"uru.1", "Liga AUF", "Uruguay", "League"},
		{"par.1", "Primera División", "Paraguay", "League"},
		{"chi.1", "Chilean Primera", "Chile", "League"},
		{"col.1", "Primera A", "Colombia", "League"},
		{"per.1", "Liga 1", "Peru", "League"},
		{"ecu.1", "LigaPro", "Ecuador", "League"},
		{"usa.1", "MLS", "USA", "League"},
		{"usa.open", "US Open Cup", "USA", "Cup"},
		{"mex.1", "Liga MX", "Mexico", "League"},
		{"jpn.1", "J.League", "Japan", "League"},
		{"chn.1", "Chinese Super League", "China", "League"},
		{"ksa.1", "Saudi Pro League", "Saudi Arabia", "League"},
		{"ind.1", "Indian Super League", "India", "League"},
		{"aus.1", "A-League Men", "Australia", "League"},
		{"tha.1", "Thai League 1", "Thailand", "League"},
		{"rsa.1", "South African Premiership", "South Africa", "League"},
		{"nga.1", "Nigerian Professional League", "Nigeria", "League"},
		{"gha.1", "Ghanaian Premier League", "Ghana", "League"},
		{"uefa.champions", "UEFA Champions League", "Europe", "Cup"},
		{"uefa.europa", "UEFA Europa League", "Europe", "Cup"},
		{"uefa.europa.conf", "UEFA Conference League", "Europe", "Cup"},
		{"uefa.euro", "European Championship", "Europe", "Cup"},
		{"uefa.nations", "UEFA Nations League", "Europe", "Cup"},
		{"fifa.world", "FIFA World Cup", "World", "Cup"},
		{"fifa.cwc", "FIFA Club World Cup", "World", "Cup"},
		{"conmebol.america", "Copa América", "South America", "Cup"},
		{"conmebol.libertadores", "Copa Libertadores", "South America", "Cup"},
		{"conmebol.sudamericana", "Copa Sudamericana", "South America", "Cup"},
		{"concacaf.gold", "Gold Cup", "North America", "Cup"},
		{"caf.nations", "Africa Cup of Nations", "Africa", "Cup"},
		{"afc.asian.cup", "Asian Cup", "Asia", "Cup"},
	}
	comps := make([]Competition, 0, len(known))
	for _, k := range known {
		comps = append(comps, Competition{
			Provider:   p.Name(),
			ProviderID: k.slug,
			Name:       k.name,
			Type:       k.typ,
			Country:    k.country,
			Season:     season,
		})
	}
	return comps, nil
}

// ---- helpers ----

// espnStatusToCanonical maps ESPN status names to our canonical set.
func espnStatusToCanonical(s espnStatus) string {
	name := s.Type.Name
	switch {
	case strings.Contains(name, "FULL_TIME"), strings.Contains(name, "FINAL"), strings.Contains(name, "FULL"):
		return "finished"
	case strings.Contains(name, "LIVE"), strings.Contains(name, "IN_PROGRESS"), strings.Contains(name, "HALFTIME"):
		return "live"
	case strings.Contains(name, "PRE"), strings.Contains(name, "SCHEDULED"), strings.Contains(name, "POSTPONED"):
		return "scheduled"
	default:
		return "scheduled"
	}
}

// espnMinute extracts the current match minute from an ESPN status.
func espnMinute(s espnStatus) int {
	return espnMinuteFromClock(s.DisplayClock)
}

// espnMinuteFromClock parses ESPN's display clock ("45'", "90'+3'") into a
// minute integer.
func espnMinuteFromClock(clock string) int {
	clock = strings.TrimSuffix(clock, "'")
	clock = strings.Split(clock, "+")[0]
	n, _ := strconv.Atoi(strings.TrimSpace(clock))
	return n
}

// espnEventType maps an ESPN detail to our canonical event type.
func espnEventType(d espnDetail) string {
	text := strings.ToLower(d.Type.Text)
	switch {
	case d.OwnGoal:
		return "own_goal"
	case d.PenaltyKick && d.ScoringPlay:
		return "penalty"
	case strings.Contains(text, "goal"):
		return "goal"
	case d.RedCard:
		return "card"
	case d.YellowCard:
		return "card"
	case strings.Contains(text, "substitution"):
		return "substitution"
	default:
		return "event"
	}
}

// espnStatType maps ESPN statistic names to our canonical set.
func espnStatType(name string) string {
	switch name {
	case "possessionPct":
		return "possession"
	case "totalShots":
		return "shots"
	case "shotsOnTarget":
		return "shots_on_target"
	case "wonCorners":
		return "corners"
	case "foulsCommitted":
		return "fouls"
	case "totalGoals":
		return "goals"
	case "goalAssists":
		return "assists"
	case "yellowCards":
		return "yellow_cards"
	case "redCards":
		return "red_cards"
	default:
		return strings.ToLower(name)
	}
}
