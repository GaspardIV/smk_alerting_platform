package distributor

import (
        "smk_alerting_platform/pkg"

        "encoding/json"
        "fmt"
        "log"
        "math"
        "net/http"
        "os"
        "strconv"
        "strings"
)

func triggerSiteChecker(urls []string, client pkg.HTTPClient, shared *pkg.SharedData) {
        defer shared.Wg.Done()

        body, err := json.Marshal(pkg.UrlsPayload{Urls: urls})
        if err != nil {
                message := fmt.Sprintf("Error while marshaling, url: %v, err: %v", urls, err)
                log.Println(message)
                shared.AppendErrorMessage(message)
                return
        }

        siteCheckerUrl := fmt.Sprintf("%v/site-checker", os.Getenv("FUNCTION_BASE_URL"))
        req, err := http.NewRequest(http.MethodPost, siteCheckerUrl, strings.NewReader(string(body)))
        if err != nil {
                message := fmt.Sprintf("Error while creating request to SiteChecker, url: %v, err: %v", siteCheckerUrl, err)
                log.Println(message)
                shared.AppendErrorMessage(message)
                return
        }
        resp, err := client.Do(req)
        if err != nil || resp.StatusCode != http.StatusOK {
                message := fmt.Sprintf("Error from SiteChecker, url: %v, err: %v, response: %v", siteCheckerUrl, err, resp)
                log.Println(message)
                shared.AppendErrorMessage(message)
                return
        }

        log.Printf("Triggered SiteChecker resp: %v, err: %v", resp, err)
}

func DistributorWithHTTPClient(w http.ResponseWriter, r *http.Request, client pkg.HTTPClient) {
        urlsPerChecker, err := strconv.Atoi(os.Getenv("URLS_PER_CHECKER"))
        if err != nil {
                log.Printf("urlsPerChecker %v", err)
                w.WriteHeader(http.StatusInternalServerError)
                return
        }

        var payload pkg.UrlsPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(fmt.Sprintf("json.NewDecoder: %v", err)))
                return
        }
        urls := payload.Urls

        shared := &pkg.SharedData{}
        shared.Wg.Add(int(math.Ceil(float64(len(urls)) / float64(urlsPerChecker))))
        for index := 0; index < len(urls); index += urlsPerChecker {
                endIndex := index + urlsPerChecker
                if endIndex > len(urls) {
                        endIndex = len(urls)
                }
                go triggerSiteChecker(urls[index:endIndex], client, shared)
        }
        shared.Wg.Wait()
        if len(shared.ErrorMessages) == 0 {
                w.WriteHeader(http.StatusOK)
        } else {
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(shared.ErrorMessages))
        }
}

func Distributor(w http.ResponseWriter, r *http.Request) {
        DistributorWithHTTPClient(w, r, &http.Client{})
}