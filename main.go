package main

import (
	"github.com/nubesFilius/go-webscrape/sites"
	"log"
	"time"
)

func main() {
	start := time.Now()

	//Get GO TIOBE Stats
	goStats, err := sites.GetTIOBELanguageStats("Go")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("TIOBELanguage: %s, Rank: %d, Rating: %.2f, Change: %.2f", goStats.Name, goStats.Rank, goStats.Rating, goStats.Change)
	log.Printf("Time it took: %dms", time.Since(start).Milliseconds())
	//end of TIOBE stats
}
