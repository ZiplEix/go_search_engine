package utils

import (
	"fmt"
	"search_engine/search"

	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	_ = c.AddFunc("0 * * * *", search.RunEngine) // Run Every Hour
	_ = c.AddFunc("15 * * * *", search.RunIndex) // Run Every Hour at 15
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", cronCount)
}
