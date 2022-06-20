package pkg

import "time"

type TasksQueue interface {
        CreateHttpTask(queueId string, url string, payload interface{}, date time.Time) error
}

type FakeTasksQueue struct {
        LastPayload interface{}
}

func NewFakeTasksQueue() *FakeTasksQueue {
        return &FakeTasksQueue{}
}

func (q *FakeTasksQueue) CreateHttpTask(_, _ string, payload interface{}, _ time.Time) error {
        q.LastPayload = payload
        return nil
}