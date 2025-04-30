package main

import (
	"context"
	"github.com/pborman/uuid"
	cron "github.com/pm57149/temporal-background-jobs/cronjobs"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	workflowID := "cron_" + uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    "cron",
		CronSchedule: "*/2 * * * *",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cron.SampleCronWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
