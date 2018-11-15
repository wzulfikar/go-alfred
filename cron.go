package main

import (
	"fmt"
	"log"
	"time"

	cron "gopkg.in/robfig/cron.v2"
)

// gopkg.in/robfig/cron.v2/doc.go
const cronSpec = "TZ=Asia/Kuala_Lumpur 0 0 8 * * *" // daily at 8am
var loc = localizedTime()

func localizedTime() *time.Location {
	tz := "Asia/Kuala_Lumpur"
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Failed to load time zone %s: %+v", tz, err)
	}
	return loc
}

func StartCron() {
	j := &CronJob{}

	// uncomment to run cron job right away
	// j.Run()

	c := cron.New()
	c.AddJob(cronSpec, j)
	c.Start()

	cronInfo(c)

	// simulate a running system with channel for
	// robfig/cron to operate.
	done := make(chan bool)
	<-done
}

func cronInfo(c *cron.Cron) {
	entry := c.Entries()[0]
	now := time.Now()

	fmt.Println("Cron job info:")
	fmt.Printf("- Prev time       : %s\n", entry.Prev)
	fmt.Printf("- Current time    : %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("- Next time       : %s\n", entry.Next)
	fmt.Printf("- Current to next : %s", diffTime(entry.Next, now))
	fmt.Println()
}

type CronJob struct{}

func (j *CronJob) Run() {
	// fetch knowledge from youtrack
	issues, err := FetchKnowledge()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(issues)
	for i := 0; i < len(*issues); i++ {
		issue := (*issues)[i]

		msg := fmt.Sprintf("*%s*\n%s\n[→ Read more](%s)",
			issue.Summary,
			shorten(issue.Description),
			GetLink(issue.ID))

		SendMsg(msg)
	}
}
