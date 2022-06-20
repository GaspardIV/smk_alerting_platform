package site_checker

import (
        "smk_alerting_platform/pkg"

        "encoding/json"
        "errors"
        "io/ioutil"
        "net/http"
        "net/http/httptest"
        "strconv"
        "strings"
        "testing"
        "time"

        "github.com/stretchr/testify/assert"
)

var (
        siteRunning = pkg.SiteInfo{
                ID:                          "id",
                Url:                         "site",
                PrimaryAdministratorEmail:   "primary_admin",
                SecondaryAdministratorEmail: "secondary_admin",
                State:                       pkg.Running,
                ConfirmationHash:            "confirmation_link",
                ResolvedHash:                "resolved_link",
                Frequency:                   10,
                TimeUntilReporting:          10,
                AllowedResponseTime:         10,
        }
        siteUnavailable = pkg.SiteInfo{
                ID:                          "id",
                Url:                         "site",
                PrimaryAdministratorEmail:   "primary_admin",
                SecondaryAdministratorEmail: "secondary_admin",
                State:                       pkg.Unavailable,
                ConfirmationHash:            "confirmation_link",
                ResolvedHash:                "resolved_link",
                Frequency:                   10,
                TimeUntilReporting:          10,
                AllowedResponseTime:         10,
        }
        siteUnexpected = pkg.SiteInfo{
                Url:   "site",
                State: "some_state",
        }
)

func TestSiteCheckerWithDatabaseHTTPClient_BadRequest(t *testing.T) {
        readCloser := ioutil.NopCloser(strings.NewReader("bad request"))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        SiteCheckerWithDatabaseHTTPClientQueue(rr, request, nil, nil, nil)
        assert.Equal(t, http.StatusBadRequest, rr.Code)
        assert.Equal(t, "json.NewDecoder: invalid character 'b' looking for beginning of value", rr.Body.String())
}

func TestSiteCheckerWithDatabaseHTTPClient_EmptyUrls(t *testing.T) {
        payload := pkg.UrlsPayload{Urls: nil}
        b, _ := json.Marshal(payload)
        readCloser := ioutil.NopCloser(strings.NewReader(string(b)))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        SiteCheckerWithDatabaseHTTPClientQueue(rr, request, nil, nil, nil)
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Equal(t, "", rr.Body.String())
}

func TestSiteCheckerWithDatabaseHTTPClient_UnexpectedStatus(t *testing.T) {
        payload := pkg.UrlsPayload{Urls: []string{"site"}}
        b, _ := json.Marshal(payload)
        readCloser := ioutil.NopCloser(strings.NewReader(string(b)))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected})
        SiteCheckerWithDatabaseHTTPClientQueue(rr, request, database, nil, nil)
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Equal(t, "", rr.Body.String())
}

func TestCheckSite_PreviouslyUnexpectedStatus(t *testing.T) {
        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected})
        shared := &pkg.SharedData{}
        shared.Wg.Add(1)
        checkSite(siteUnexpected.Url, database, nil, nil, shared)
        shared.Wg.Wait()

        siteInfo, err := database.GetSiteInfo("url", siteUnexpected.Url)
        assert.Equal(t, siteUnexpected, siteInfo)
        assert.Nil(t, err)
}

func TestCheckSite_SiteNotFoundInDatabase(t *testing.T) {
        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{siteRunning})
        shared := &pkg.SharedData{}
        shared.Wg.Add(1)
        checkSite("some_url", database, nil, nil, shared)
        shared.Wg.Wait()

        siteInfo, err := database.GetSiteInfo("url", siteRunning.Url)
        assert.Equal(t, siteRunning, siteInfo)
        assert.Nil(t, err)
}

func TestCheckSite_SiteUnavailable_PreviouslyRunning(t *testing.T) {
        cases := []struct {
                statusCode int
                err        error
        }{
                {
                        http.StatusNotFound,
                        nil,
                },
                {
                        http.StatusOK,
                        errors.New("some error"),
                },
                {
                        http.StatusNotFound,
                        errors.New("some error"),
                },
        }
        for i, c := range cases {
                t.Run(strconv.Itoa(i), func(t *testing.T) {
                        site := siteRunning
                        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{site})
                        shared := &pkg.SharedData{}
                        shared.Wg.Add(1)
                        client := pkg.NewFakeHttpClient(c.statusCode, c.err)
                        checkSite(site.Url, database, client, nil, shared)
                        shared.Wg.Wait()

                        siteInfo, err := database.GetSiteInfo("url", site.Url)
                        assert.Nil(t, err)
                        site = siteUnavailable
                        assert.NotEqual(t, site.StateChangeTimestamp, siteInfo.StateChangeTimestamp)
                        site.StateChangeTimestamp = siteInfo.StateChangeTimestamp
                        pkg.AssertAlmostEqual(t, site, siteInfo)
                })
        }
}

func TestCheckSite_SiteUnavailable_PreviouslyUnavailable(t *testing.T) {
        cases := []struct {
                statusCode int
                err        error
        }{
                {
                        http.StatusNotFound,
                        nil,
                },
                {
                        http.StatusOK,
                        errors.New("some error"),
                },
                {
                        http.StatusNotFound,
                        errors.New("some error"),
                },
        }
        for i, c := range cases {
                t.Run(strconv.Itoa(i), func(t *testing.T) {
                        site := siteUnavailable
                        site.StateChangeTimestamp = time.Now()
                        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{site})

                        client := pkg.NewFakeHttpClient(c.statusCode, c.err)
                        shared := &pkg.SharedData{}
                        shared.Wg.Add(1)
                        checkSite(site.Url, database, client, nil, shared)
                        shared.Wg.Wait()

                        siteInfo, err := database.GetSiteInfo("url", site.Url)
                        assert.Nil(t, err)
                        pkg.AssertAlmostEqual(t, siteInfo, site)
                })
        }
}

func TestCheckSite_SiteUnavailable_PreviouslyUnavailable_NotifierTriggered(t *testing.T) {
        site := siteUnavailable
        site.StateChangeTimestamp = time.Now().Add(-2 * site.GetTimeUntilReporting())
        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{site})

        queue := pkg.NewFakeTasksQueue()
        client := pkg.NewFakeHttpClient(http.StatusInternalServerError, nil)
        shared := &pkg.SharedData{}
        shared.Wg.Add(1)
        checkSite(site.Url, database, client, queue, shared)
        shared.Wg.Wait()

        siteInfo, err := database.GetSiteInfo("url", site.Url)
        assert.Nil(t, err)
        pkg.AssertAlmostEqual(t, site, siteInfo)
        expectedPayload := map[string]interface{}{
                "url": site.Url,
        }
        assert.Equal(t, expectedPayload, queue.LastPayload)
}

func TestCheckSite_SiteRunning_PreviouslyRunning(t *testing.T) {
        site := siteRunning
        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{site})
        client := pkg.NewFakeHttpClient(http.StatusOK, nil)

        shared := &pkg.SharedData{}
        shared.Wg.Add(1)
        checkSite(site.Url, database, client, nil, shared)
        shared.Wg.Wait()

        siteInfo, err := database.GetSiteInfo("url", site.Url)
        assert.Nil(t, err)
        pkg.AssertAlmostEqual(t, site, siteInfo)
}

func TestCheckSite_SiteRunning_PreviouslyUnavailable(t *testing.T) {
        site := siteUnavailable
        database := pkg.CreateLocalDatabase([]pkg.SiteInfo{site})
        client := pkg.NewFakeHttpClient(http.StatusOK, nil)

        shared := &pkg.SharedData{}
        shared.Wg.Add(1)
        checkSite(site.Url, database, client, nil, shared)
        shared.Wg.Wait()

        siteInfo, err := database.GetSiteInfo("url", site.Url)
        assert.Nil(t, err)
        assert.NotEqual(t, siteRunning.StateChangeTimestamp, siteInfo.StateChangeTimestamp)
        siteInfo.StateChangeTimestamp = siteRunning.StateChangeTimestamp
        pkg.AssertAlmostEqual(t, siteRunning, siteInfo)
}