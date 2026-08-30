package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/icekidtech/Sportz44/backend/internal/external"
)

func main() {
	espn := external.NewESPNProvider(os.Getenv("ESPN_URL"))
	ctx := context.Background()

	// Use slug-encoded ID for the PL match
	providerID := "eng.1:401879294"
	fmt.Println("=== Testing GetLiveEvents with slug-encoded ID:", providerID, "===")

	events, err := espn.GetLiveEvents(ctx, providerID)
	if err != nil {
		log.Fatal("GetLiveEvents:", err)
	}
	fmt.Printf("Got %d events\n", len(events))
	for _, ev := range events {
		if ev.EventType == "goal" || ev.EventType == "yellow" || ev.EventType == "red" {
			fmt.Printf("  %s min=%d team=%s player=%s (detail: %s)\n", ev.EventType, ev.Minute, ev.TeamID, ev.PlayerName, ev.Detail)
			if ev.AssistPlayerName != "" {
				fmt.Printf("    assist: %s\n", ev.AssistPlayerName)
			}
		}
	}

	fmt.Println("\n=== Testing GetMatchStats with slug-encoded ID ===")
	stats, err := espn.GetMatchStats(ctx, providerID)
	if err != nil {
		log.Fatal("GetMatchStats:", err)
	}
	fmt.Printf("Got %d stats\n", len(stats))
	for _, s := range stats {
		fmt.Printf("  team=%s %s=%s\n", s.TeamID, s.StatType, s.Value)
	}

	fmt.Println("\n=== Testing GetFixtures (PL today) ===")
	fixtures, err := espn.GetFixtures(ctx, "eng.1", "", "")
	if err != nil {
		log.Fatal("GetFixtures:", err)
	}
	fmt.Printf("Got %d fixtures\n", len(fixtures))
	for _, f := range fixtures {
		fmt.Printf("  %s vs %s (ID: %s, status: %s)\n", f.HomeTeamName, f.AwayTeamName, f.ProviderID, f.Status)
	}

	// Find a match to verify slug-encoded ID works for GetFixture
	for _, f := range fixtures {
		if strings.HasPrefix(f.ProviderID, "eng.1:") {
			fmt.Printf("\n=== Testing GetFixture with slug-encoded ID: %s ===\n", f.ProviderID)
			got, err := espn.GetFixture(ctx, f.ProviderID)
			if err != nil {
				log.Fatal("GetFixture:", err)
			}
			fmt.Printf("  %s %d-%d %s\n", got.Status, got.HomeScore, got.AwayScore, got.Venue)
			break
		}
	}
}
