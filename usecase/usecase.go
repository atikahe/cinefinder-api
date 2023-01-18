package usecase

import (
	"fmt"
	"log"

	"github.com/atikahe/cinefinder-api/config"
	"github.com/atikahe/cinefinder-api/pkg/discover"
	"github.com/atikahe/cinefinder-api/pkg/elastic"

	"github.com/playwright-community/playwright-go"
)

type SearchResult struct {
	Source  string                `json:"source"`
	Elapsed string                `json:"elapsed"`
	Data    []elastic.ElasticData `json:"data"`
}

func Search(query string, index string) (*SearchResult, error) {
	// Load env
	env, err := config.LoadEnvConfig(".env")
	if err != nil {
		log.Fatalf("Couldn't load env config: %s", err)
		return nil, err
	}
	if err := env.Validate(); err != nil {
		log.Fatalf("Env config invalid: %s", err)
		return nil, err
	}

	// Connect to Elastic
	es, err := elastic.New(env.ElasticCloudID, env.ElasticAPIKey)
	if err != nil {
		log.Fatalf("init elastic client failed: %s", err)
		return nil, err
	}

	// Search in elastic
	fields := []string{"title", "overview", "meta_description"}
	hits, err := es.Search(query, index, fields...)
	if err != nil {
		log.Fatalf("search elastic failed: %s", err)
		return nil, err
	}

	// Clean search result
	response := SearchResult{
		Source:  index,
		Elapsed: fmt.Sprintf("%fms", hits.Took),
		Data:    []elastic.ElasticData{},
	}
	for _, hit := range hits.Hits.Hits {
		response.Data = append(response.Data, hit.Source)
	}

	return &response, nil
}

func Discover(q string, pw *playwright.Playwright) ([]string, error) {
	texts, err := discover.Run(q, pw)
	if err != nil {
		return nil, err
	}
	return texts, nil
}
