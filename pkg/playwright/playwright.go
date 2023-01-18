package playwright

import (
	"errors"
	"fmt"
	"log"

	"github.com/playwright-community/playwright-go"
)

func Init() (*playwright.Playwright, error) {
	runOptions := playwright.RunOptions{
		Browsers: []string{"chromium"},
		Verbose:  false,
	}
	err := playwright.Install(&runOptions)
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

	pw, err := playwright.Run(&runOptions)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't start headless browser: %v", err))
	}

	return pw, nil
}
