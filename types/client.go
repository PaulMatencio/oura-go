package types

import (
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	defaultRetryWaitMin       = 200 * time.Millisecond
	defaultRetryWaitMax       = 3 * time.Second
	defaultRetryMax           = 3
	respReadLimit       int64 = 64000
)

type Client struct {
	HTTPClient    *http.Client
	RetryWaitMin  time.Duration
	RetryWaitMax  time.Duration
	RetryMax      int
	CheckForRetry CheckForRetry
	Backoff       Backoff
}

func NewClient() *Client {
	return &Client{
		HTTPClient:    cleanhttp.DefaultClient(),
		RetryWaitMin:  defaultRetryWaitMin,
		RetryWaitMax:  defaultRetryWaitMax,
		RetryMax:      defaultRetryMax,
		CheckForRetry: DefaultRetryPolicy,
		Backoff:       DefaultBackoff,
	}
}

type Request struct {
	body io.ReadSeeker
	*http.Request
}

func NewRequest(method, url string, body io.ReadSeeker) (*Request, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}
	httpReq, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}
	return &Request{body, httpReq}, nil
}

func (c *Client) Do(req *Request) (*http.Response, error) {
	var i = defaultRetryMax
	for {
		var code int
		if req.body != nil {
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("failed to seek body: %v", err)
			}
		}
		resp, err := c.HTTPClient.Do(req.Request)
		checkOK, checkErr := c.CheckForRetry(resp, err)

		if err != nil {
			log.Error().Msgf("%s %s request failed: %v", req.Method, req.URL, err)
		} else {

			/*


			 */
		}
		if !checkOK {
			if checkErr != nil {
				err = checkErr
			}
			return resp, err
		}

		remain := c.RetryMax - 1
		if remain == 0 {
			break
		}
		i++
		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}
		time.Sleep(wait)

	}
	return nil, fmt.Errorf("%s %s giving up afters %d attempts", req.Method, req.URL, c.RetryMax)

}

func (c *Client) drainBody(body io.ReadCloser) {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
	if err != nil {
		log.Error().Msgf("error reading response body: %v, err")
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url, contentType string, body io.ReadSeeker) (*http.Response, error) {
	req, err := NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}
