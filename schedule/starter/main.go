package main

import (
	"context"
	"go.temporal.io/sdk/client"
	"log"
	"main/model"
	"main/workflows"
	"time"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	customers := []model.Customer{
		{
			AccountNum:       "c1",
			Name:             "Pawan",
			Email:            "pawan@buildcode.com",
			CustomerType:     "new",
			DemoWaitDuration: 10 * time.Second,
		},
		{
			AccountNum:       "c2",
			Name:             "Alice",
			Email:            "alice@buildcode.com",
			CustomerType:     "existing",
			DemoWaitDuration: 10 * time.Second,
		},
		{
			AccountNum:       "c3",
			Name:             "Bob",
			Email:            "bob@buildcode.com",
			CustomerType:     "new",
			DemoWaitDuration: 10 * time.Second,
		},
	}

	for _, customer := range customers {
		workflowOptions := client.StartWorkflowOptions{
			ID:        customer.AccountNum,
			TaskQueue: "versioningGoDemoTaskQueue",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CustomerWorkflow, customer)
		if err != nil {
			log.Fatalln("Unable to execute workflow for", customer.AccountNum, err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		var result model.Account
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatalln("Unable get workflow result", err)
		}
		log.Println("Workflow result:", result.Amount)
	}

	//c, err := client.Dial(client.Options{
	//	HostPort: "localhost:7233",
	//})
	//
	//if err != nil {
	//	log.Fatalln("Unable to create client", err)
	//}
	//defer c.Close()
	//
	//// Create 100 customers with varying wait times
	//customers := make([]model.Customer, 100)
	//for i := 0; i < 100; i++ {
	//	waitDuration := time.Duration(rand.Intn(10)+1) * time.Minute // Random wait time between 1 and 10 minutes
	//	if i < 30 {
	//		waitDuration = 15 * time.Minute // The first 30 customers have a wait time of 15 minutes
	//	}
	//
	//	customers[i] = model.Customer{
	//		AccountNum:       "c" + strconv.Itoa(i+1) + "id diff",
	//		Name:             "Customer" + strconv.Itoa(i+1),
	//		Email:            "customer" + strconv.Itoa(i+1) + "@example.com",
	//		CustomerType:     "new",
	//		DemoWaitDuration: waitDuration,
	//	}
	//}
	//
	//// Start all workflows in parallel
	//for _, customer := range customers {
	//	go func(customer model.Customer) {
	//		workflowOptions := client.StartWorkflowOptions{
	//			ID:        customer.AccountNum,
	//			TaskQueue: "versioningGoDemoTaskQueue",
	//		}
	//
	//		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CustomerWorkflow, customer)
	//		if err != nil {
	//			log.Println("Unable to execute workflow for", customer.AccountNum, err)
	//			return
	//		}
	//		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	//
	//		var result model.Account
	//		err = we.Get(context.Background(), &result)
	//		if err != nil {
	//			log.Println("Unable to get workflow result for", customer.AccountNum, err)
	//			return
	//		}
	//		log.Println("Workflow result for", customer.AccountNum, ":", result.Amount)
	//	}(customer)
	//}
	//
	//// Prevent the main function from exiting immediately
	//select {}
}
