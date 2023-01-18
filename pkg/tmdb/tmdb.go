package TMDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/atikahe/cinefinder-api/payload"

	"github.com/go-resty/resty/v2"
)

type TMDB struct {
	Client *resty.Client
}

func NewTMDB() *TMDB {
	return &TMDB{
		Client: resty.New(),
	}
}

func (t *TMDB) SetBaseURL(url string) *TMDB {
	t.Client.SetBaseURL(url)
	return t
}

func (t *TMDB) SetAuth(apikey string) *TMDB {
	t.Client.SetQueryParam("api_key", apikey)
	return t
}

func (t *TMDB) Search(q string) (*payload.Response, error) {
	params := map[string]string{
		"language":      "en-US",
		"page":          "1",
		"include_adult": "false",
		"query":         q,
	}
	url := fmt.Sprintf("%s/search/movie", t.Client.BaseURL)

	rawResponse, err := t.Client.R().EnableTrace().SetQueryParams(params).Get(url)
	if rawResponse.StatusCode() != http.StatusOK || err != nil {
		return nil, errors.New("request error")
	}

	var response payload.Response
	if err := json.Unmarshal(rawResponse.Body(), &response); err != nil {
		return nil, errors.New("error unmarshalling response")
	}

	return &response, nil
}
