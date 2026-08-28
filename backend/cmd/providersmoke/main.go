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
	fs := external.NewFlashScoreProvider(os.Getenv("FLASHSCORE_URL"))
	tsdb := external.NewTheSportsDBProvider(os.Getenv("THESPORTSDB_KEY"), os.Getenv("THESPORTSDB_URL"))

	reg := external.NewRegistry(api, fs, tsdb)
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
	fixtures, err := api.GetFixtures(ctx, 140, "2025", "")
	if err != nil {
		fmt.Println("GetFixtures error:", err)
	} else if len(fixtures) > 0 {
		fmt.Printf("Fixtures (La Liga) fetched: %d (first: %+v)\n", len(fixtures), fixtures[0])
	} else {
		fmt.Println("Fixtures: no data")
	}

	// 5. Live matches (realtime provider)
	live, err := fs.GetLiveMatches(ctx)
	if err != nil {
		fmt.Println("FlashScore live error:", err)
	} else {
		fmt.Printf("FlashScore live matches: %d\n", len(live))
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
	if p := reg.ByKind(ctx, external.ProviderMedia); p != nil {
		fmt.Println("Media provider:", p.Name())
	} else {
		fmt.Println("Media provider: none healthy")
	}
}
