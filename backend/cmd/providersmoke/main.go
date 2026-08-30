package main

import (
	"context"
	"fmt"
	"os"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/joho/godotenv"
)

// A quick smoke test for the provider layer. Run with:
//
//	go run ./cmd/providersmoke
func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	apiKey := os.Getenv("API_SPORTS_KEY")
	apiHost := os.Getenv("API_SPORTS_HOST")
	if apiKey == "" || apiHost == "" {
		fmt.Println("API_SPORTS_KEY / API_SPORTS_HOST not set")
		os.Exit(1)
	}

	api := external.NewAPISportsProvider(apiKey, apiHost)
	espn := external.NewESPNProvider(os.Getenv("ESPN_URL"))

	// Football-Data provider (current season free tier)
	fdKey := os.Getenv("FOOTBALL_DATA_KEY")
	fdURL := os.Getenv("FOOTBALL_DATA_URL")
	fd := external.NewFootballDataProvider(fdKey, fdURL)

	reg := external.NewRegistry(espn, api, fd)
	fmt.Println("Providers registered:", len(reg.All()))

	// 1. Competitions (La Liga = 140)
	comps, err := api.GetCompetitions(ctx, "2026")
	if err != nil {
		fmt.Println("GetCompetitions error:", err)
	} else if len(comps) > 0 {
		fmt.Printf("Competitions (2026) fetched: %d (first: %+v)\n", len(comps), comps[0])
	} else {
		// Try current season (2025) if 2026 not available yet
		comps, err = api.GetCompetitions(ctx, "2025")
		if err != nil {
			fmt.Println("GetCompetitions error:", err)
		} else if len(comps) > 0 {
			fmt.Printf("Competitions (2025) fetched: %d (first: %+v)\n", len(comps), comps[0])
		} else {
			fmt.Println("Competitions: no data for 2026 or 2025")
		}
	}

	// 1a. Football-Data current-season competitions
	season := external.CurrentSeason()
	cp := reg.CompetitionsForSeason(ctx, season)
	compsFD, err := cp.GetCompetitions(ctx, season)
	if err != nil {
		fmt.Println("Football-Data GetCompetitions error:", err)
	} else if len(compsFD) > 0 {
		fmt.Printf("Football-Data competitions (%s): %d (first: %+v)\n", season, len(compsFD), compsFD[0])
	} else {
		fmt.Println("Football-Data competitions: no data")
	}

	// 2. Team (Real Madrid = 541)
	team, err := api.GetTeam(ctx, "541")
	if err != nil {
		fmt.Println("GetTeam error:", err)
	} else {
		fmt.Printf("Team: %+v\n", team)
	}

	// 3. Squad
	squad, err := api.GetSquad(ctx, "541")
	if err != nil {
		fmt.Println("GetSquad error:", err)
	} else {
		fmt.Printf("Squad size: %d (first: %+v)\n", len(squad), squad[0])
	}

	// 4. Fixtures for La Liga
	fixtures, err := api.GetFixtures(ctx, "140", "2025", "")
	if err != nil {
		fmt.Println("GetFixtures error:", err)
	} else if len(fixtures) > 0 {
		fmt.Printf("Fixtures (La Liga) fetched: %d (first: %+v)\n", len(fixtures), fixtures[0])
	} else {
		fmt.Println("Fixtures: no data")
	}

	// 5. Live matches (ESPN realtime provider)
	live, err := espn.GetLiveMatches(ctx)
	if err != nil {
		fmt.Println("ESPN live error:", err)
	} else {
		fmt.Printf("ESPN live matches: %d\n", len(live))
	}

	// 5a. ESPN current-season fixtures (Premier League)
	espnFixtures, err := espn.GetFixtures(ctx, "eng.1", external.CurrentSeason(), "")
	if err != nil {
		fmt.Println("ESPN GetFixtures error:", err)
	} else if len(espnFixtures) > 0 {
		fmt.Printf("ESPN fixtures (PL): %d (first: %+v)\n", len(espnFixtures), espnFixtures[0])
	} else {
		fmt.Println("ESPN fixtures: no data")
	}

	// 6. Registry failover checks
	if p := reg.ByKind(ctx, external.ProviderBulk); p != nil {
		fmt.Println("Bulk provider:", p.Name())
	} else {
		fmt.Println("Bulk provider: none healthy")
	}
	if p := reg.ByKind(ctx, external.ProviderRealtime); p != nil {
		fmt.Println("Realtime provider:", p.Name())
	} else {
		fmt.Println("Realtime provider: none healthy")
	}
}
