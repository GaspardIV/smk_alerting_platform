package site_checker

import (
        "smk_alerting_platform/pkg"

        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "os"
        "time"
)

func checkSite(url string, database pkg.Database, client pkg.HTTPClient, queue pkg.TasksQueue, shared *pkg.SharedData) {
        defer shared.Wg.Done()

        // Retrieve set info from database and handle unexpected status
        siteInfo, err := database.GetSiteInfo("url", url)
        if err != nil {
                message := fmt.Sprintf("Error while getting site %v info %v", url, err)
                log.Println(message)
                shared.AppendErrorMessage(message)
                return
        }
        if siteInfo.State != pkg.Running && siteInfo.State != pkg.Unavailable {
                log.Printf("Unexpected state %v of site %v", siteInfo.State, url)
                return
        }

        // Send HTTP request
        req, err := http.NewRequest(http.MethodGet, url, nil)
        if err != nil {
                message := fmt.Sprintf("Error while getting site creating request to site %v %v", url, err)
                log.Println(message)
                shared.AppendErrorMessage(message)
                return
        }
        resp, err := client.Do(req)

        // Handle response
        var state string
        shouldNotify := false
        if err != nil || resp.StatusCode != http.StatusOK {
                log.Printf("site %v is unavailable", url)
                shouldNotify =
                        siteInfo.State == pkg.Unavailable && time.Now().Sub(siteInfo.StateChangeTimestamp) >= siteInfo.GetTimeUntilReporting()
                state = pkg.Unavailable
        } else {
                log.Printf("site %v available", url)
                state = pkg.Running
        }

        // Update site info
        if state != siteInfo.State {
                siteInfo.State = state
                siteInfo.StateChangeTimestamp = time.Now()
        }
        err = database.UpdateSite(siteInfo)
        if err != nil {
                message := fmt.Sprintf("Error while updating site %v info %v", url, err)
                log.Println(message)
                shared.AppendErrorMessage(message)
                return
        }

        // Trigger notifier if necessary
        if shouldNotify {
                log.Printf("Notifying administrator of site %v", url)
                queue.CreateHttpTask(
                        "notifier-queue",
                        fmt.Sprintf("%v/notifier", os.Getenv("FUNCTION_BASE_URL")),
                        map[string]interface{}{
                                "url": url,
                        },
                        time.Now(),
                )

                queue.CreateHttpTask(
                        "notifier-queue",
                        fmt.Sprintf("%v/notifier", os.Getenv("FUNCTION_BASE_URL")),
                        map[string]interface{}{
                                "url": url,
                        },
                        time.Now().Add(siteInfo.GetAllowedResponseTime()),
                )
        }
}

func SiteCheckerWithDatabaseHTTPClientQueue(w http.ResponseWriter, r *http.Request, database pkg.Database, client pkg.HTTPClient, queue pkg.TasksQueue) {
        var payload pkg.UrlsPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(fmt.Sprintf("json.NewDecoder: %v", err)))
                return
        }

        shared := &pkg.SharedData{}
        shared.Wg.Add(len(payload.Urls))
        for _, url := range payload.Urls {
                go checkSite(url, database, client, queue, shared)
        }

        shared.Wg.Wait()
        if len(shared.ErrorMessages) == 0 {
                w.WriteHeader(http.StatusOK)
        } else {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(shared.ErrorMessages))
        }
}

func SiteChecker(w http.ResponseWriter, r *http.Request) {
        database, dbErr := pkg.CreateDatabase()
        if dbErr != nil {
                w.Write([]byte("Could not connect to database"))
                w.WriteHeader(http.StatusInternalServerError)
        }
        queue, queueErr := pkg.CreateTaskQueue()
        if queueErr != nil {
                w.Write([]byte("Could not connect to cloud tasks queue"))
                w.WriteHeader(http.StatusInternalServerError)
        }
        SiteCheckerWithDatabaseHTTPClientQueue(w, r, database, &http.Client{}, queue)
}