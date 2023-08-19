package driver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/mousav1/sms"
	"github.com/mousav1/sms/config"
	"github.com/tidwall/gjson"
)

// Ghasedak represents an SMS gateway provider.
type Ghasedak struct {
	APIKey     string
	LineNumber string
	Host       string
}

// CreateProvider creates an instance of the Ghasedak provider.
func (g *Ghasedak) CreateProvider(config config.DriverConfig) (sms.SMSProvider, error) {
	return &Ghasedak{
		APIKey:     config.APIKey,
		LineNumber: config.LineNumber,
		Host:       config.Host,
	}, nil
}

// SendRequest sends an HTTP request and returns the response.
func (g *Ghasedak) SendRequest(method, url string, headers map[string]string, body url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body.Encode()))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

// SendSMS sends an SMS using Ghasedak.
func (g *Ghasedak) SendSMS(to, message string) (sms.Response, error) {
	apiURL := fmt.Sprintf("http://%s/v2/sms/send/simple?agent=go", g.Host)
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Apikey":       g.APIKey,
	}
	body := url.Values{
		"message":    {message},
		"receptor":   {to},
		"linenumber": {g.LineNumber},
	}

	resp, err := g.SendRequest("POST", apiURL, headers, body)
	if err != nil {
		return sms.Response{Success: false, Message: err.Error()}, err
	}
	defer resp.Body.Close()

	response := sms.Response{}
	response.Code = resp.StatusCode

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return sms.Response{Success: false, Message: err.Error()}, err
		}

		bodyString := string(bodyBytes)
		message := gjson.Get(bodyString, "result.message").String()
		id := gjson.Get(bodyString, "items.0").Int()

		response.Message = message
		response.ID = id
		response.Success = true

		return response, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return sms.Response{Success: false, Message: err.Error()}, err
	}

	bodyString := string(bodyBytes)
	errorMessage := gjson.Get(bodyString, "result.message").String()
	response.Message = errorMessage
	response.Success = false

	return response, fmt.Errorf("failed to send SMS: %s", errorMessage)
}
