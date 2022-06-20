package pkg

import (
        "context"
        "encoding/json"
        "fmt"
        "log"
        "time"

        cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
        taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
        "google.golang.org/protobuf/types/known/timestamppb"
)

type cloudTasksQueue struct {
        client *cloudtasks.Client
}

func CreateTaskQueue() (TasksQueue, error) {
        client, err := cloudtasks.NewClient(context.Background())
        log.Printf("%v", err)
        if err != nil {
                return cloudTasksQueue{}, err
        }

        return cloudTasksQueue{client: client}, nil
}

func (tasks_queue cloudTasksQueue) CreateHttpTask(queueId string, url string, payload interface{}, scheduleTime time.Time) error {
        queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueId)
        body, jsonErr := json.Marshal(payload)
        if jsonErr != nil {
                return fmt.Errorf("json.Marshal: %v", jsonErr)
        }
        req := &taskspb.CreateTaskRequest{
                Parent: queuePath,
                Task: &taskspb.Task{
                        MessageType: &taskspb.Task_HttpRequest{
                                HttpRequest: &taskspb.HttpRequest{
                                        Url:  url,
                                        Body: body,
                                },
                        },
                        ScheduleTime: timestamppb.New(scheduleTime),
                },
        }

        _, err := tasks_queue.client.CreateTask(context.Background(), req)
        if err != nil {
                return fmt.Errorf("cloudtasks.CreateTask: %v", err)
        }

        return nil
}