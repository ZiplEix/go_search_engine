package utils

import (
	"fmt"
	"search_engine/search"

	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	_ = c.AddFunc("0 * * * *", search.RunEngine) // Run Every Hour
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", cronCount)
}
