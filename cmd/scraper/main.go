package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/greg-source/ukrainian-theater-events/internal/logger"
	s3storage "github.com/greg-source/ukrainian-theater-events/internal/storage/s3"
)

const (
	frankoURL = "https://sales.ft.org.ua/events"
	s3Key     = "franko_theater/latest.html"
	bucket    = "theater-bot-scraper"
	region    = "eu-north-1"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	bucketName := os.Getenv("SCRAPER_BUCKET")
	if bucketName == "" {
		bucketName = bucket
	}

	s3Client, err := s3storage.NewClient(ctx, region)
	if err != nil {
		logger.Fatalf("failed to create s3 client: %v", err)
	}
	store := s3storage.New(s3Client, bucketName)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
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
		logger.Infof("scraping %s", pageURL)

		var cards string
		var cardCount int
		if err := chromedp.Run(chromeCtx,
			chromedp.Navigate(pageURL),
			chromedp.WaitReady("body"),
			chromedp.Sleep(5*time.Second),
			chromedp.Evaluate(`document.querySelectorAll('.performanceCard').length`, &cardCount),
			chromedp.Evaluate(`Array.from(document.querySelectorAll('.performanceCard')).map(e => e.outerHTML).join('\n')`, &cards),
		); err != nil {
			logger.Fatalf("failed to scrape page %d: %v", page, err)
		}

		if cardCount == 0 {
			logger.Infof("page %d has no events, stopping", page)
			break
		}

		logger.Infof("page %d: found %d events", page, cardCount)
		allCards.WriteString(cards)
		allCards.WriteString("\n")
		totalCards += cardCount
	}

	if totalCards == 0 {
		logger.Fatalf("no events found across all pages")
	}

	combined := fmt.Sprintf("<html><body>\n%s</body></html>", allCards.String())
	if err := store.Put(ctx, s3Key, []byte(combined)); err != nil {
		logger.Fatalf("failed to write to s3: %v", err)
	}

	logger.Infof("scraping complete: %d events written to %s", totalCards, s3Key)
}
