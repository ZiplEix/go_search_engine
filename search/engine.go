package search

import (
	"fmt"
	"search_engine/db"
	"time"
)

func RunEngine() {
	fmt.Println("Starting search engine crawl...")
	defer fmt.Println("Search engine crawl complete.")

	// Get settings from the database
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		fmt.Printf("failed to get search settings: %s", err.Error())
		return
	}

	if !settings.SearchOn {
		fmt.Println("Search engine is disabled")
		return
	}

	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCrawUrls(int(settings.Amount))
	if err != nil {
		fmt.Printf("failed to get next urls: %s", err.Error())
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdateUrl(db.CrawledUrl{
				Url:             next.Url,
				Success:         false,
				CrawDuration:    result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Heading:         result.CrawlData.Heading,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Printf("failed to update a fail url: %s", err.Error())
				return
			}
			continue
		}
		err := next.UpdateUrl(db.CrawledUrl{
			Url:             next.Url,
			Success:         result.Success,
			CrawDuration:    result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageTitle,
			PageDescription: result.CrawlData.PageDescription,
			Heading:         result.CrawlData.Heading,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Printf("failed to update a success url: %s", err.Error())
			return
		}

		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // end of range

	if !settings.AddNew {
		return
	}

	for _, newUrl := range newUrls {
		fmt.Println("Adding new url: ", newUrl.Url)
		err := newUrl.Save()
		if err != nil {
			fmt.Printf("failed to save new url: %s", err.Error())
			return
		}
	}

	fmt.Printf("\nAdded %d new urls in the database\n", len(newUrls))
}
