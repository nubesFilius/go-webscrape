package sites

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type TIOBEStats struct {
	Languages []TIOBELanguage
}

type TIOBELanguage struct {
	Rank   int
	Name   string
	Rating float64
	Change float64
}

// GetTIOBELanguageStats returns the stats for a specific language.
func GetTIOBELanguageStats(languageName string) (*TIOBELanguage, error) {
	stats, err := getAllTIOBEStats()
	if err != nil {
		return nil, err
	}
	for _, lang := range stats.Languages {
		if strings.EqualFold(lang.Name, languageName) {
			return &lang, nil
		}
	}
	return nil, fmt.Errorf("error getting stats for language %s", languageName)
}

// getAllTIOBEStats fetches and parses the TIOBE Index page.
func getAllTIOBEStats() (*TIOBEStats, error) {
	// Request the HTML page.
	res, err := http.Get("https://www.tiobe.com/tiobe-index/")
	if err != nil {
		return nil, fmt.Errorf("failed to GET request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read document from reader: %s", err)
	}

	var stats TIOBEStats
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Find the table rows
	doc.Find("table#top20 tbody tr").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection) {
			defer wg.Done()

			rank, err := strconv.Atoi(s.Find("td").Eq(0).Text())
			if err != nil {
				log.Printf("failed to parse rank: %s", err)
				return
			}

			language := s.Find("td").Eq(4).Text()

			ratingStr := strings.TrimSuffix(s.Find("td").Eq(5).Text(), "%")
			rating, err := strconv.ParseFloat(ratingStr, 64)
			if err != nil {
				log.Printf("failed to parse rating: %s", err)
				return
			}

			changeStr := strings.TrimSuffix(s.Find("td").Eq(6).Text(), "%")
			change, err := strconv.ParseFloat(changeStr, 64)
			if err != nil {
				log.Printf("failed to parse change: %s", err)
				return
			}

			mu.Lock()
			stats.Languages = append(stats.Languages, TIOBELanguage{
				Rank:   rank,
				Name:   language,
				Rating: rating,
				Change: change,
			})
			mu.Unlock()
		}(s)
	})

	wg.Wait()

	return &stats, nil
}
