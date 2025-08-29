package scraper

import (
	"time"
	"log"
	"github.com/samothreesixty/rss-agg/internal/handlers"
	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/samothreesixty/rss-agg/internal/rss"
	"context"
	"sync"
	"database/sql"
	"github.com/google/uuid"
	"strings"
)

func StartScraping(
	apiCfg *handlers.ApiConfig,
	concurrency int32,
	timeBetweenScrapes time.Duration,
) {
	log.Println("Starting scraper...")
	log.Printf("Scraping on %v interval with %v concurrency", timeBetweenScrapes, concurrency)

	ticker := time.NewTicker(timeBetweenScrapes)

	for ;  ; <-ticker.C {
		feeds, err := apiCfg.DB.GetNextFeedsToFetch(context.Background(), concurrency)
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
			go scrapeFeed(apiCfg, &wg, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(apiCfg *handlers.ApiConfig, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()
	
	feed, err := apiCfg.DB.UpdateFeedLastFetchedAt(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	rssFeed, err := rss.UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
	}

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if (item.Description != "") {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("Could not pase date %v with error %v", item.PubDate, err)
			continue
		}

		_, err = apiCfg.DB.CreatePost(context.Background(), db.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error creating post:", err)
		}
	}

	log.Printf("Feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Items))
}
