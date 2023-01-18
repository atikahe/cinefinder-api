package main

import (
	"log"
	"net/http"

	"github.com/atikahe/cinefinder-api/handler"
	"github.com/atikahe/cinefinder-api/pkg/playwright"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize playwright browser for Discovery feature
	pw, err := playwright.Init()
	if err != nil {
		log.Fatalf("unable to set up playwright")
	}
	defer pw.Stop()

	// Initialize web framework
	e := echo.New()
	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "OK")
	})

	api := e.Group("/api")
	api.GET("/search", handler.Search)
	api.GET("/discover", func(ctx echo.Context) error {
		return handler.Discover(ctx, pw)
	})

	log.Fatal(e.Start(":8009"))
}
