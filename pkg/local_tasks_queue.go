package pkg

import (
        "bytes"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "time"
)

type localTasksQueue struct{}

func CreateLocalTaskQueue() *localTasksQueue {
        return &localTasksQueue{}
}

func (tasks_queue localTasksQueue) CreateHttpTask(queueId string, url string, payload interface{}, scheduleTime time.Time) error {
        body, jsonErr := json.Marshal(payload)
        if jsonErr != nil {
                return fmt.Errorf("json.Marshal: %v", jsonErr)
        }

        time.AfterFunc(time.Until(scheduleTime), func() {
                log.Printf(url)
                resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
                if err != nil {
                        log.Fatalf("http.Post: %v", err)
                }
                resp.Body.Close()
        })
        return nil
}