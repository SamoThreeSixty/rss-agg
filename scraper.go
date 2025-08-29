package main

import (
	"time"
	"log"
	"github.com/samothreesixty/rss-agg/internal/db"
	"context"
	"sync"
)

func startScraping(
	db *db.Queries,
	concurrency int32,
	timeBetweenScrapes time.Duration,
) {
	log.Println("Starting scraper...")
	log.Printf("Scraping on %v interval with %v concurrency", timeBetweenScrapes, concurrency)

	ticker := time.NewTicker(timeBetweenScrapes)

	for ;  ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), concurrency)
		if err != nil {
			log.Println("Cannot get feeds to fetch:", err)
			continue
		}
		if len(feeds) == 0 {
			log.Println("No feeds to fetch")
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(db * db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()
	
	feed, err := db.UpdateFeedLastFetchedAt(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
	}

	for _, item := range rssFeed.Channel.Items {
		log.Println("Found post", item.Title)
	}
	log.Printf("Feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Items))
}
