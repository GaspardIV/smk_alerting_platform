package scheduler

import (
        "context"
        "fmt"
        "log"
        "os"
        "smk_alerting_platform/pkg"
        "sort"
        "time"
)

type PubSubMessage struct {
        Data []byte `json:"data"`
}

func scheduleDistributors(database pkg.Database, queue pkg.TasksQueue, interval time.Duration, roundingFactor time.Duration) error {
        siteInfos, err := database.GetAllSites()
        if err != nil {
                return fmt.Errorf("database.GetAllSites %v", err)
        }
        start := time.Now()
        end := start.Add(interval).Add(roundingFactor)
        groupedTasks := make(map[int64][]string)

        for _, siteInfo := range siteInfos {
                frequency := siteInfo.GetFrequency()
                timestamp := siteInfo.LastChangeTimestamp.Add(frequency)

                if timestamp.Before(start) {
                        timestamp = timestamp.Add(start.Sub(timestamp).Round(roundingFactor))
                }

                for timestamp.Before(end) {
                        roundedTimestamp := timestamp.Round(roundingFactor).Unix()
                        groupedTasks[roundedTimestamp] = append(groupedTasks[roundedTimestamp], siteInfo.Url)
                        timestamp = timestamp.Add(frequency)
                }
        }

        sortedKeys := make([]int64, len(groupedTasks))
        i := 0
        for k := range groupedTasks {
                sortedKeys[i] = k
                i++
        }
        sort.Slice(sortedKeys, func(i, j int) bool { return sortedKeys[i] < sortedKeys[j] })
        for _, ts := range sortedKeys {
                timestamp := time.Unix(ts, 0)
                log.Printf("Scheduling distributor task at: %v, for urls: %v", timestamp, groupedTasks[ts])

                queue.CreateHttpTask(
                        "scheduler-queue",
                        fmt.Sprintf("%v/distributor", os.Getenv("FUNCTION_BASE_URL")),
                        map[string]interface{}{
                                "urls": groupedTasks[ts],
                        },
                        timestamp,
                )
        }
        return nil
}

func SchedulerWithDatabaseQueue(ctx context.Context, message PubSubMessage, database pkg.Database, queue pkg.TasksQueue) error {
        interval, intervalErr := time.ParseDuration(string(message.Data))
        if intervalErr != nil {
                return fmt.Errorf("interval time.ParseDuration %v", intervalErr)
        }

        roundingFactor, roundingFactorErr := time.ParseDuration(os.Getenv("ROUNDING_FACTOR"))
        if roundingFactorErr != nil {
                return fmt.Errorf("roundingFactor time.ParseDuration %v", roundingFactorErr)
        }

        return scheduleDistributors(database, queue, interval, roundingFactor)
}

func Scheduler(ctx context.Context, message PubSubMessage) error {
        database, dbErr := pkg.CreateDatabase()
        if dbErr != nil {
                return fmt.Errorf("Could not connect to database")
        }
        queue, queueErr := pkg.CreateTaskQueue()
        if queueErr != nil {
                return fmt.Errorf("Could not connect to cloud tasks queue")
        }
        return SchedulerWithDatabaseQueue(ctx, message, database, queue)
}