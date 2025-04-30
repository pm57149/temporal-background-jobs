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
	ctx := context.Background()
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	customer := []model.Customer{
		{
			AccountNum:       "c1",
			Name:             "Pawan",
			Email:            "pawan@buildcode.com",
			CustomerType:     "new",
			DemoWaitDuration: 10 * time.Second,
		},
	}

	scheduleID := "schedule_567"
	workflowID := "schedule_workflow_8910"

	scheduleHandle, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID:   scheduleID,
		Spec: client.ScheduleSpec{},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  workflows.CustomerWorkflow,
			Args:      []interface{}{customer[0]},
			TaskQueue: "schedule",
		},
	})
	if err != nil {
		log.Fatalln("Unable to create schedule", err)
	}

	//defer func() {
	//	log.Println("Deleting schedule", "ScheduleID", scheduleHandle.GetID())
	//	err = scheduleHandle.Delete(ctx)
	//	if err != nil {
	//		log.Fatalln("Unable to delete schedule", err)
	//	}
	//}()

	//err = scheduleHandle.Trigger(ctx, client.ScheduleTriggerOptions{
	//	Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL,
	//})
	//
	//if err != nil {
	//	log.Fatalln("Unable to trigger schedule", err)
	//}

	log.Println("Updating schedule", "ScheduleID", scheduleHandle.GetID())
	err = scheduleHandle.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(schedule client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			schedule.Description.Schedule.Spec = &client.ScheduleSpec{

				Calendars: []client.ScheduleCalendarSpec{
					{
						Hour: []client.ScheduleRange{
							{
								Start: 17,
							},
						},
						DayOfWeek: []client.ScheduleRange{
							{
								Start: 5,
							},
						},
					},
				},

				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: 5 * time.Second,
					},
				},
			}

			schedule.Description.Schedule.State.Paused = true
			schedule.Description.Schedule.State.LimitedActions = true
			schedule.Description.Schedule.State.RemainingActions = 10

			return &client.ScheduleUpdate{
				Schedule: &schedule.Description.Schedule,
			}, nil
		},
	})
	if err != nil {
		log.Fatalln("Unable to update schedule", err)
	}

}
