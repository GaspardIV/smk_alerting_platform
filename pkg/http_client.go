package pkg

import "net/http"

type UrlsPayload struct {
        Urls []string `json:"urls"`
}

type HTTPClient interface {
        Do(req *http.Request) (*http.Response, error)
}

type FakeHTTPClient struct {
        StatusCode int
        Err        error
        Requests   []http.Request
}

func NewFakeHttpClient(code int, err error) *FakeHTTPClient {
        return &FakeHTTPClient{
                StatusCode: code,
                Err:        err,
                Requests:   nil,
        }
}

func (c *FakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
        c.Requests = append(c.Requests, *req)
        return &http.Response{StatusCode: c.StatusCode}, c.Err
}