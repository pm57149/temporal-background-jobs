package main

import (
	cron "github.com/pm57149/temporal-background-jobs/cronjobs"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "cron", worker.Options{})

	w.RegisterWorkflow(cron.SampleCronWorkflow)
	w.RegisterActivity(cron.DoSomething)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
