package distributor

import (
        "encoding/json"
        "errors"
        "io/ioutil"
        "net/http"
        "net/http/httptest"
        "os"
        "smk_alerting_platform/pkg"
        "strings"
        "testing"

        "github.com/stretchr/testify/assert"
)

func TestDistributorWithHTTPClient_NoURlsPerChecker(t *testing.T) {
        rr := httptest.NewRecorder()
        DistributorWithHTTPClient(rr, nil, nil)
        assert.Equal(t, http.StatusInternalServerError, rr.Code)
        assert.Equal(t, "", rr.Body.String())
}

func TestDistributorWithHTTPClient_BadRequest(t *testing.T) {
        os.Setenv("URLS_PER_CHECKER", "12")
        readCloser := ioutil.NopCloser(strings.NewReader("bad request"))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        DistributorWithHTTPClient(rr, request, nil)
        assert.Equal(t, http.StatusBadRequest, rr.Code)
        assert.Equal(t, "json.NewDecoder: invalid character 'b' looking for beginning of value", rr.Body.String())
}

func TestDistributorWithHTTPClient_EmptyUrls(t *testing.T) {
        os.Setenv("URLS_PER_CHECKER", "12")
        payload := pkg.UrlsPayload{Urls: nil}
        b, _ := json.Marshal(payload)
        readCloser := ioutil.NopCloser(strings.NewReader(string(b)))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        client := pkg.NewFakeHttpClient(http.StatusOK, nil)
        DistributorWithHTTPClient(rr, request, client)
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Equal(t, "", rr.Body.String())
        var expectedRequests []http.Request = nil
        assert.Equal(t, expectedRequests, client.Requests)
}

func TestDistributorWithHTTPClient_OkRequest(t *testing.T) {
        os.Setenv("URLS_PER_CHECKER", "2")
        os.Setenv("FUNCTION_BASE_URL", "https://base_url")
        payload := pkg.UrlsPayload{Urls: []string{"url1", "url2", "url3"}}
        b, _ := json.Marshal(payload)
        readCloser := ioutil.NopCloser(strings.NewReader(string(b)))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        client := pkg.NewFakeHttpClient(http.StatusOK, nil)
        DistributorWithHTTPClient(rr, request, client)
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Equal(t, "", rr.Body.String())

        var payload1 pkg.UrlsPayload
        json.NewDecoder(client.Requests[0].Body).Decode(&payload1)
        var payload2 pkg.UrlsPayload
        json.NewDecoder(client.Requests[1].Body).Decode(&payload2)
        assert.Equal(t, pkg.UrlsPayload{Urls: []string{"url3"}}, payload1)
        assert.Equal(t, pkg.UrlsPayload{Urls: []string{"url1", "url2"}}, payload2)
}

func TestDistributorWithHTTPClient_ClientStatusNotOk(t *testing.T) {
        os.Setenv("URLS_PER_CHECKER", "1")
        os.Setenv("FUNCTION_BASE_URL", "https://base_url")
        payload := pkg.UrlsPayload{Urls: []string{"url1"}}
        b, _ := json.Marshal(payload)
        readCloser := ioutil.NopCloser(strings.NewReader(string(b)))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        client := pkg.NewFakeHttpClient(http.StatusInternalServerError, nil)
        DistributorWithHTTPClient(rr, request, client)
        assert.Equal(t, http.StatusInternalServerError, rr.Code)
        assert.Equal(t,
                "Error from SiteChecker, url: https://base_url/site-checker, err: <nil>, response: &{ 500  0 0 map[] <nil> 0 [] false false map[] <nil> <nil>}\n",
                rr.Body.String())
}
func TestDistributorWithHTTPClient_ClientErrorNotNil(t *testing.T) {
        os.Setenv("URLS_PER_CHECKER", "1")
        os.Setenv("FUNCTION_BASE_URL", "https://base_url")
        payload := pkg.UrlsPayload{Urls: []string{"url1"}}
        b, _ := json.Marshal(payload)
        readCloser := ioutil.NopCloser(strings.NewReader(string(b)))
        request := &http.Request{Body: readCloser}
        rr := httptest.NewRecorder()

        client := pkg.NewFakeHttpClient(http.StatusOK, errors.New("some error"))
        DistributorWithHTTPClient(rr, request, client)
        assert.Equal(t, http.StatusInternalServerError, rr.Code)
        assert.Equal(t,
                "Error from SiteChecker, url: https://base_url/site-checker, err: some error, response: &{ 200  0 0 map[] <nil> 0 [] false false map[] <nil> <nil>}\n",
                rr.Body.String())
}