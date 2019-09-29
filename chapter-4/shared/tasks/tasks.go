package tasks

import (
	"chapter-4/shared/config"
	"chapter-4/shared/run"
	"context"
	"fmt"
	"log"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

// Tasks is the client object that contains
// pointers to the backend clients
type Tasks struct {
	tasks *cloudtasks.Client
	run   *run.Client
}

// NewClient initializes the client with
// the backend services: Cloud Tasks and our
// custom Cloud Run api client
func NewClient() (*Tasks, error) {
	tasks, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	run, err := run.NewClient()
	if err != nil {
		return nil, err
	}

	return &Tasks{tasks, run}, nil
}

func (client *Tasks) AddTask(ctx context.Context,
	queue, receiverFunc, ID string,
	message []byte,
) error {
	// Talk to the Cloud Run API to get the target HTTP url.
	// Example: https://background-rwrmxiaqmq-ew.a.run.app
	serviceUrl, err := client.run.GetServiceUrl(
		config.Region, receiverFunc)
	if err != nil {
		return err
	}
	log.Printf("submitting task to %s", serviceUrl)

	// Build the target HTTP request
	targetRequest := &taskspb.HttpRequest{
		// The URL to send the request to
		Url:  fmt.Sprintf("%v/%v", serviceUrl, ID),
		Body: message,
		AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
			// This tells Cloud Tasks to add an Authorization
			// header with the identity of the calling function
			OidcToken: &taskspb.OidcToken{
				ServiceAccountEmail: config.FunctionServiceAccount,
			},
		},
	}

	err = client.sendTask(ctx, ID, queue, targetRequest)
	if err != nil {
		return err
	}
	return nil
}

func (client *Tasks) sendTask(
	ctx context.Context,
	ID, queue string,
	targetRequest *taskspb.HttpRequest) error {

	// Build the queue name
	queueName := fmt.Sprintf("projects/%v/locations/%v/queues/%v",
		config.ProjectID,
		config.Region,
		queue)

	// Build a unique task name
	taskName := fmt.Sprintf("%v/tasks/%v", queueName, ID)

	// Construct the request to the Cloud Tasks API
	req := &taskspb.CreateTaskRequest{
		Parent: queueName,
		Task: &taskspb.Task{
			Name: taskName,
			PayloadType: &taskspb.Task_HttpRequest{
				HttpRequest: targetRequest,
			},
		},
	}

	// Send the request to the Cloud Tasks API
	_, err := client.tasks.CreateTask(ctx, req)
	return err
}
