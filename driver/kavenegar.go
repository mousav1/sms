package driver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mousav1/sms"
	request "github.com/mousav1/sms/Request"
	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/errors"
	"github.com/mousav1/sms/response"
)

// Ghasedak represents an SMS gateway provider.
type kavenegar struct {
	APIKey     string
	LineNumber string
	Host       string
}

type ResultKavenegar struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResultErrorKavenegar struct {
	*ResultKavenegar `json:"result"`
}

type MessageResultKavenegar struct {
	*ResultKavenegar `json:"result"`
	Entries          []Message `json:"entries"`
}

type Message struct {
	MessageID  int    `json:"messageid"`
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statustext"`
	Sender     string `json:"sender"`
	Receptor   string `json:"receptor"`
	Date       int    `json:"date"`
	Cost       int    `json:"cost"`
}

// CreateProvider creates an instance of the Ghasedak provider.
func (g *kavenegar) CreateProvider(config config.DriverConfig) (sms.SMSProvider, error) {
	return &kavenegar{
		APIKey:     config.APIKey,
		LineNumber: config.LineNumber,
		Host:       config.Host,
	}, nil
}

// SendSMS sends an SMS using Ghasedak.
func (g *kavenegar) SendSMS(to, message string) (response.Response, error) {
	apiURL := fmt.Sprintf("http://%s/v1/%s/sms/send.json", g.Host, g.APIKey)

	body := url.Values{
		"message":    {message},
		"receptor":   {to},
		"linenumber": {g.LineNumber},
	}

	headers := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Accept":         "application/json",
		"Accept-Charset": "utf-8",
	}

	request := request.NewRequest(g.APIKey, g.Host)

	resp, err := request.Execute(apiURL, "POST", body, headers)

	defer resp.Body.Close()

	response := response.Response{}
	response.Status = resp.StatusCode

	if http.StatusOK != resp.StatusCode {
		re := new(ResultErrorKavenegar)
		err = json.NewDecoder(resp.Body).Decode(&re)
		if err != nil {
			return response, &errors.Error{
				Status:  resp.StatusCode,
				Message: resp.Status,
				Err:     err,
			}
		}
		return response, &errors.Error{
			Status:  re.ResultKavenegar.Status,
			Message: re.ResultKavenegar.Message,
			Err:     err,
		}
	}

	result := new(MessageResultKavenegar)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return response, &errors.Error{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Err:     err,
		}
	}

	response.Message = result.Message
	response.Status = result.Status

	return response, nil
}
