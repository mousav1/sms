package request

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/mousav1/sms/errors"
)

type Request struct {
	apikey  string
	BaseURL *url.URL
}

func NewRequest(apikey, host string) *Request {
	baseURL, _ := url.Parse(host)
	c := &Request{
		BaseURL: baseURL,
		apikey:  apikey,
	}
	return c
}

func (r *Request) Execute(urlStr, method string, b url.Values, headers map[string]string) (*http.Response, error) {
	body := strings.NewReader(b.Encode())
	ul, _ := url.Parse(urlStr)
	u := r.BaseURL.ResolveReference(ul)
	req, _ := http.NewRequest(method, u.String(), body)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if err, ok := err.(net.Error); ok {
			return nil, err
		}
		if resp == nil {
			return nil, &errors.Error{
				Status:  http.StatusInternalServerError,
				Message: "nil response",
				Err:     err,
			}
		}
		return nil, &errors.Error{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Err:     err,
		}
	}

	defer resp.Body.Close()

	return resp, nil
}
