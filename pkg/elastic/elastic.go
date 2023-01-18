package elastic

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/mitchellh/mapstructure"
)

var (
	INCLUDE_PAGES = []string{
		"person", // tmdb
		"tv",     // tmdb
		"movie",  // tmdb
		"title",  // imdb
		"name",   // imdb
	}

	EXCLUDE_PAGES = []string{
		"discuss", // url_path_dir_3, tmdb
		"keyword", // url_path_dir_1, tmdb
		"video",   // url_path_dir_1, imdb
		"images",  // url_path_dir_3, tmdb
		"changes", // url_path_dir_3, tmdb
	}
)

type Elastic struct {
	Client *elasticsearch.Client
}

func New(cloudID string, apiKey string) (*Elastic, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		CloudID: cloudID,
		APIKey:  apiKey,
	})

	return &Elastic{
		Client: client,
	}, err
}

func (es *Elastic) Index(index string, doc []byte) (*esapi.Response, error) {
	res, err := es.Client.Index(
		index,
		bytes.NewReader(doc),
	)
	return res, err
}

func (es *Elastic) Search(q string, idx string, fields ...string) (*ElasticResponse, error) {
	// Build elasticsearch query
	var buf bytes.Buffer
	query := es.baseQuery(q, fields...)
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// Execute search
	res, err := es.Client.Search(
		// es.Client.Search.WithContext(ctx),
		es.Client.Search.WithIndex(idx),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Prepare result
	var r map[string]interface{}
	var response ElasticResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	mapstructure.Decode(r, &response)

	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(response.Hits.Total.Value),
		int(response.Took),
	)
	return &response, nil
}

func (es *Elastic) baseQuery(q string, fields ...string) map[string]interface{} {
	if len(fields) < 1 {
		fields = []string{"title"}
	}

	baseQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  q,
							"fields": fields,
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"match": map[string]string{
							"url_path_dir_3": EXCLUDE_PAGES[0],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_1": EXCLUDE_PAGES[1],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_1": EXCLUDE_PAGES[2],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_3": EXCLUDE_PAGES[3],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_3": EXCLUDE_PAGES[4],
						},
					},
				},
				"should": []map[string]interface{}{
					{
						"match": map[string]string{
							"url_path_dir_1": INCLUDE_PAGES[0],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_1": INCLUDE_PAGES[1],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_1": INCLUDE_PAGES[2],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_1": INCLUDE_PAGES[3],
						},
					},
					{
						"match": map[string]string{
							"url_path_dir_1": INCLUDE_PAGES[4],
						},
					},
				},
			},
		},
	}

	return baseQuery
}
