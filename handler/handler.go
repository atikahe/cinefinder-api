package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/atikahe/cinefinder-api/payload"
	"github.com/atikahe/cinefinder-api/usecase"

	"github.com/labstack/echo/v4"
	"github.com/playwright-community/playwright-go"
)

const (
	TMDB_INDEX    string = "search-movie"
	IMDB_INDEX    string = "search-imdb"
	CONTENT_TYPE  string = "Content-Type"
	CACHE_CONTROL string = "Cache-Control"
	CONNECTION    string = "Connection"
)

func Search(ctx echo.Context) error {
	// Process query
	param := payload.SearchParam{}
	if err := ctx.Bind(&param); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var indices []string
	if len(param.Index) < 1 {
		indices = []string{TMDB_INDEX, IMDB_INDEX}
	}

	// Set response header
	ctx.Response().Header().Set(CONTENT_TYPE, "text/event-stream")
	ctx.Response().Header().Set(CACHE_CONTROL, "no-cache")
	ctx.Response().Header().Set(CONNECTION, "keep-alive")

	// Prepare channel
	dataCh := createChannel(param, indices...)

	// Send response as stream
	flusher, ok := ctx.Response().Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	for data := range dataCh {
		ctx.Response().WriteHeader(data.StatusCode)
		response, _ := json.Marshal(data)
		_, err := ctx.Response().Write(response)
		if err != nil {
			return fmt.Errorf("unable to send data")
		}
		flusher.Flush()
	}

	return nil
}

func createChannel(param payload.SearchParam, indices ...string) chan *payload.HTTPResponse {
	// Init wait group and channel
	var wg sync.WaitGroup
	dataCh := make(chan *payload.HTTPResponse)

	for _, index := range indices {
		wg.Add(1)
		go func(index string) {
			var httpres *payload.HTTPResponse

			// Perform search to elastic
			r, err := usecase.Search(param.Query, index)

			if err != nil {
				httpres = &payload.HTTPResponse{StatusCode: int(http.StatusInternalServerError), Message: err.Error()}
			} else {
				httpres = &payload.HTTPResponse{
					StatusCode: int(http.StatusOK),
					Message:    http.StatusText(http.StatusOK),
					Data:       r,
				}
			}

			dataCh <- httpres
			wg.Done()
		}(index)
	}

	go func() {
		// Close channel when process finished
		wg.Wait()
		close(dataCh)
	}()

	return dataCh
}

func Discover(ctx echo.Context, pw *playwright.Playwright) error {
	param := payload.SearchParam{}
	if err := ctx.Bind(&param); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	texts, err := usecase.Discover(param.Query, pw)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, texts)
}
