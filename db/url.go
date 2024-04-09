package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	gorm.Model
	Url             string        `json:"url" gorm:"unique;not null"`
	Success         bool          `json:"success" gorm:"not null"`
	CrawDuration    time.Duration `json:"crawDuration"`
	ResponseCode    int           `json:"responseCode"`
	PageTitle       string        `json:"pageTitle"`
	PageDescription string        `json:"pageDescription"`
	Heading         string        `json:"heading"`
	LastTested      *time.Time    `json:"lastTested"`
	Indexed         bool          `json:"indexed" gorm:"default:false"`
}

func (crawled *CrawledUrl) UpdateUrl(input CrawledUrl) error {
	var existingUrl CrawledUrl

	if err := DBConn.Where("url = ?", input.Url).First(&existingUrl).Error; err == nil {
		input.ID = existingUrl.ID
		err = DBConn.Model(&CrawledUrl{}).Where("id = ?", input.ID).Updates(input).Error
		if err != nil {
			return fmt.Errorf("unable to update url: %v", err)
		}

		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error while checking if url exists: %v", err)
	}

	err := DBConn.Create(&input).Error
	if err != nil {
		return fmt.Errorf("unable to create url: %v", err)
	}

	return nil
}

func (crawled *CrawledUrl) GetNextCrawUrls(limit int) ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, fmt.Errorf("unable to get next urls: %v", tx.Error)
	}

	return urls, nil
}

func (crawled *CrawledUrl) Save() error {
	var existingUrl CrawledUrl

	if err := DBConn.Where("url = ?", crawled.Url).First(&existingUrl).Error; err == nil {
		crawled.ID = existingUrl.ID
		err = DBConn.Model(&CrawledUrl{}).Where("id = ?", crawled.ID).Updates(crawled).Error
		if err != nil {
			return fmt.Errorf("unable to update url: %v", err)
		}

		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error while checking if url exists: %v", err)
	}

	err := DBConn.Create(&crawled).Error
	if err != nil {
		return fmt.Errorf("unable to create url: %v", err)
	}

	return nil
}
