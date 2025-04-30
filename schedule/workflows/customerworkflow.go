package workflows

import (
	"go.temporal.io/sdk/workflow"
	"main/activities"
	"main/model"
	"time"
)

func CustomerWorkflow(ctx workflow.Context, customer model.Customer) (model.Account, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("CustomerWorkflow Info workflow started.", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 6,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	bonus := 100

	var account model.Account
	err := workflow.ExecuteActivity(ctx, activities.GetCustomerAccount, customer).Get(ctx, &account)
	if err != nil {
		logger.Error("GetCustomerAccount failed.", "Error", err)
		return model.Account{}, err
	}

	_ = workflow.Sleep(ctx, customer.DemoWaitDuration)

	err = workflow.ExecuteActivity(ctx, activities.UpdateCustomerAccount, customer, bonus).Get(ctx, &account)
	if err != nil {
		logger.Error("UpdateCustomerAccount failed.", "Error", err)
		return model.Account{}, err
	}

	return account, err
}
