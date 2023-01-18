package discover

import (
	"errors"
	"fmt"
	"log"

	"github.com/atikahe/cinefinder-api/pkg/set"

	"github.com/playwright-community/playwright-go"
)

func extractKeywords(page playwright.Page, parent string, child string, keywords *set.Set) error {
	locator, err := page.Locator(parent)
	if err != nil {
		return errors.New("could not locate parent")
	}

	elementHandles, err := locator.ElementHandles()
	if err != nil {
		return errors.New("could not get elementHandles")
	}

	for _, element := range elementHandles {
		keywordElement, err := element.QuerySelector(child)
		if err != nil {
			return errors.New("could not get title element")
		}
		if keywordElement == nil {
			continue
		}

		keyword, err := keywordElement.TextContent()
		if err != nil && keyword == "" {
			return errors.New("could not get text content")
		}

		// Limit keywords length to that of a tweet limit
		if !(*keywords).Contains(keyword) && len(keyword) < 280 && keyword != "" {
			(*keywords).Add(keyword)
		}
	}

	return nil
}

func Run(q string, pw *playwright.Playwright) ([]string, error) {
	// Using playwright is a hacky way to reach the intention of this project.
	// However, it's close enough as a prototype to demonstrate what we want to achieve.

	browser, err := pw.Chromium.LaunchPersistentContext(
		"/tmp/playwright",
		playwright.BrowserTypeLaunchPersistentContextOptions{Headless: playwright.Bool(true)},
	)
	if err != nil {
		log.Fatalf("could not launch browser")
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err = page.Goto(fmt.Sprintf("https://www.google.com/search?q=%s", q)); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	keywords := set.New()

	// Extract keyword from complementary result (re: movie details section)
	if err = extractKeywords(page, `[role="complementary"]`, `[role="heading"]`, &keywords); err != nil {
		return nil, err
	}

	// Extract keyword from featured snippet results
	if err = extractKeywords(page, ".hgKElc", "b", &keywords); err != nil {
		return nil, err
	}

	// Extract keyword from regular search results
	if err = extractKeywords(page, ".VwiC3b", "span > em", &keywords); err != nil {
		return nil, err
	}

	return keywords.Values(), nil
}
