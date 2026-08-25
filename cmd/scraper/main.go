package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abrahamjinston-png/theater-events/internal/logger"
	s3storage "github.com/abrahamjinston-png/theater-events/internal/storage/s3"
	"github.com/chromedp/chromedp"
)

const (
	frankoURL    = "https://sales.ft.org.ua/events"
	pageRetries  = 3
	retryBackoff = 15 * time.Second
	pageLoadWait = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	region := mustEnv("AWS_REGION")
	bucketName := mustEnv("S3_BUCKET")
	s3Key := mustEnv("SCRAPER_S3_KEY")

	s3Client, err := s3storage.NewClient(ctx, region)
	if err != nil {
		logger.Fatal("creating s3 client failed", err)
	}
	store := s3storage.New(s3Client, bucketName)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx)
	defer chromeCancel()

	var allCards strings.Builder
	totalCards := 0

	for page := 1; ; page++ {
		pageURL := fmt.Sprintf("%s?page=%d", frankoURL, page)
		logger.Info(fmt.Sprintf("scraping %s", pageURL))

		cards, cardCount, err := scrapePage(chromeCtx, pageURL)
		if err != nil {
			logger.Fatal(fmt.Sprintf("scraping page %d failed after %d attempts", page, pageRetries), err)
		}

		if cardCount == 0 {
			logger.Info(fmt.Sprintf("page %d has no events, stopping", page))
			break
		}

		logger.Info(fmt.Sprintf("page %d: found %d events", page, cardCount))
		allCards.WriteString(cards)
		allCards.WriteString("\n")
		totalCards += cardCount
	}

	if totalCards == 0 {
		logger.Fatal("no events found across all pages", nil)
	}

	combined := fmt.Sprintf("<html><body>\n%s</body></html>", allCards.String())
	if err := store.Put(ctx, s3Key, []byte(combined)); err != nil {
		logger.Fatal("writing to s3 failed", err)
	}

	logger.Info(fmt.Sprintf("scraping complete: %d events written to %s", totalCards, s3Key))
}

func scrapePage(ctx context.Context, pageURL string) (cards string, cardCount int, err error) {
	for attempt := 1; attempt <= pageRetries; attempt++ {
		var c string
		var n int
		err = chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.WaitReady("body"),
			chromedp.Sleep(pageLoadWait),
			chromedp.Evaluate(`document.querySelectorAll('.performanceCard').length`, &n),
			chromedp.Evaluate(`Array.from(document.querySelectorAll('.performanceCard')).map(e => e.outerHTML).join('\n')`, &c),
		)
		if err == nil {
			return c, n, nil
		}
		logger.Error(fmt.Sprintf("scrape attempt %d/%d failed for %s", attempt, pageRetries, pageURL), err)
		if attempt < pageRetries {
			time.Sleep(retryBackoff)
		}
	}
	return "", 0, err
}

func mustEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logger.Fatal(fmt.Sprintf("%s is required", name), nil)
	}
	return v
}
