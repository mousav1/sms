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
type Ghasedak struct {
	APIKey     string
	LineNumber string
	Host       string
}

type ResultGhasedak struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResultErrorGhasedak struct {
	*ResultGhasedak `json:"result"`
}

type MessageResult struct {
	*ResultGhasedak `json:"result"`
	Items           []int64 `json:"items"`
}

// CreateProvider creates an instance of the Ghasedak provider.
func (g *Ghasedak) CreateProvider(config config.DriverConfig) (sms.SMSProvider, error) {
	return &Ghasedak{
		APIKey:     config.APIKey,
		LineNumber: config.LineNumber,
		Host:       config.Host,
	}, nil
}

// SendSMS sends an SMS using Ghasedak.
func (g *Ghasedak) SendSMS(to, message string) (response.Response, error) {
	apiURL := fmt.Sprintf("http://%s/v2/sms/send/simple?agent=go", g.Host)

	body := url.Values{
		"message":    {message},
		"receptor":   {to},
		"linenumber": {g.LineNumber},
	}

	headers := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Apikey":         g.APIKey,
		"Accept":         "application/json",
		"Accept-Charset": "utf-8",
	}

	request := request.NewRequest(g.APIKey, g.Host)

	resp, err := request.Execute(apiURL, "POST", body, headers)

	defer resp.Body.Close()

	response := response.Response{}
	response.Status = resp.StatusCode

	if http.StatusOK != resp.StatusCode {
		re := new(ResultErrorGhasedak)
		err = json.NewDecoder(resp.Body).Decode(&re)
		if err != nil {
			return response, &errors.Error{
				Status:  resp.StatusCode,
				Message: resp.Status,
				Err:     err,
			}
		}
		return response, &errors.Error{
			Status:  re.ResultGhasedak.Code,
			Message: re.ResultGhasedak.Message,
			Err:     err,
		}
	}

	result := new(MessageResult)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return response, &errors.Error{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Err:     err,
		}
	}

	response.Message = result.Message
	response.MessageID = int64(result.Items[0])

	return response, nil
}
