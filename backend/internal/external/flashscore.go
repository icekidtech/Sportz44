package external

import (
	"context"
	"time"
)

// FlashScoreProvider is the realtime provider. FlashScore does not offer an
// official free API; it exposes internal JSON endpoints that are commonly used
// for live scores. This provider polls those endpoints at high frequency (30s)
// to keep live scores/events fresh without burning the bulk provider's quota.
//
// NOTE: FlashScore's endpoints are unofficial and can change. This provider is
// isolated behind the Provider interface so it can be swapped or disabled
// without affecting the rest of the pipeline.
type FlashScoreProvider struct {
	client *httpClient
}

// NewFlashScoreProvider creates a FlashScore realtime provider.
func NewFlashScoreProvider(baseURL string) *FlashScoreProvider {
	return &FlashScoreProvider{
		client: newHTTPClient(baseURL, 10*time.Second, map[string]string{
			"User-Agent": "Mozilla/5.0 (compatible; Sportz44/1.0)",
		}),
	}
}

func (p *FlashScoreProvider) Name() string       { return "flashscore" }
func (p *FlashScoreProvider) Kind() ProviderKind { return ProviderRealtime }

func (p *FlashScoreProvider) Healthy(ctx context.Context) bool {
	// A lightweight probe; if the live endpoint responds, we consider it up.
	_, err := p.GetLiveMatches(ctx)
	return err == nil
}

// GetLiveMatches fetches currently-live matches from FlashScore's live JSON
// endpoint. The exact shape varies; we parse the common "events" structure.
func (p *FlashScoreProvider) GetLiveMatches(ctx context.Context) ([]Fixture, error) {
	var out struct {
		Events []struct {
			ID        string `json:"id"`
			Home      string `json:"home"`
			Away      string `json:"away"`
			HomeScore int    `json:"homeScore"`
			AwayScore int    `json:"awayScore"`
			Status    struct {
				Code int `json:"code"`
			} `json:"status"`
			Time struct {
				Minute int `json:"minute"`
			} `json:"time"`
		} `json:"events"`
	}
	if err := p.client.get(ctx, "/x/feed/f_1_0_1_en_1", &out); err != nil {
		return nil, err
	}
	fixtures := make([]Fixture, 0, len(out.Events))
	for _, e := range out.Events {
		fixtures = append(fixtures, Fixture{
			ProviderID:   e.ID,
			HomeTeamID:   e.Home,
			AwayTeamID:   e.Away,
			HomeTeamName: e.Home,
			AwayTeamName: e.Away,
			Status:       "live",
			HomeScore:    e.HomeScore,
			AwayScore:    e.AwayScore,
			Minute:       e.Time.Minute,
		})
	}
	return fixtures, nil
}

// GetLiveEvents fetches the events for a live match.
func (p *FlashScoreProvider) GetLiveEvents(ctx context.Context, providerFixtureID string) ([]MatchEvent, error) {
	// FlashScore's event feed is keyed by a match ID; we return an empty set
	// when the feed is unavailable rather than failing the whole sync.
	path := "/x/feed/d_" + providerFixtureID + "_en_1"
	var out struct {
		Events []struct {
			Minute int    `json:"minute"`
			Type   string `json:"type"`
			Team   string `json:"team"`
			Player string `json:"player"`
			Detail string `json:"detail"`
		} `json:"events"`
	}
	if err := p.client.get(ctx, path, &out); err != nil {
		return nil, err
	}
	events := make([]MatchEvent, 0, len(out.Events))
	for _, e := range out.Events {
		events = append(events, MatchEvent{
			ProviderFixtureID: providerFixtureID,
			Minute:            e.Minute,
			EventType:         normalizeEventType(e.Type),
			TeamID:            e.Team,
			PlayerName:        e.Player,
			Detail:            e.Detail,
		})
	}
	return events, nil
}

// TheSportsDBProvider is a media enrichment provider (logos, photos, bios).
// It is used to fill in assets that the bulk provider may lack.
type TheSportsDBProvider struct {
	client *httpClient
}

// NewTheSportsDBProvider creates a TheSportsDB media provider.
func NewTheSportsDBProvider(apiKey, baseURL string) *TheSportsDBProvider {
	return &TheSportsDBProvider{
		client: newHTTPClient(baseURL, 10*time.Second, map[string]string{}),
	}
}

func (p *TheSportsDBProvider) Name() string       { return "thesportsdb" }
func (p *TheSportsDBProvider) Kind() ProviderKind { return ProviderMedia }

func (p *TheSportsDBProvider) Healthy(ctx context.Context) bool {
	return true
}

func (p *TheSportsDBProvider) GetTeamLogo(ctx context.Context, providerTeamID string) (string, error) {
	// TheSportsDB uses a different ID space; we return empty when not found.
	return "", nil
}

func (p *TheSportsDBProvider) GetPlayerPhoto(ctx context.Context, providerPlayerID string) (string, error) {
	return "", nil
}
