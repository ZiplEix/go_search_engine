package db

import (
	"fmt"

	"gorm.io/gorm"
)

type SearchIndex struct {
	gorm.Model
	Value string
	Urls  []CrawledUrl `gorm:"many2many:token_urls;"` // many to many relationship
}

func (s *SearchIndex) TableName() string {
	return "search_index"
}

func (s *SearchIndex) Save(index map[string][]uint, crawledUrls []CrawledUrl) error {
	for value, ids := range index {
		newIndex := &SearchIndex{
			Value: value,
		}
		if err := DBConn.Where(SearchIndex{Value: value}).FirstOrCreate(newIndex).Error; err != nil {
			return fmt.Errorf("unable to create search index: %v", err)
		}
		var urlsToAppend []CrawledUrl
		for _, id := range ids {
			for _, url := range crawledUrls {
				if url.ID == id {
					urlsToAppend = append(urlsToAppend, url)
					break
				}
			}
		}
		if err := DBConn.Model(newIndex).Association("Urls").Append(&urlsToAppend); err != nil {
			return fmt.Errorf("unable to append urls to search index: %v", err)
		}
	}

	return nil
}
