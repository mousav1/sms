package driver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mousav1/sms"
	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/errors"
	"github.com/mousav1/sms/request"
	"github.com/mousav1/sms/response"
)

type KavenegarProvider struct {
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

// CreateProvider creates an instance of the Kavenegar provider.
func (k *KavenegarProvider) CreateProvider(config config.DriverConfig) (sms.SMSProvider, error) {
	return &KavenegarProvider{
		APIKey:     config.APIKey,
		LineNumber: config.LineNumber,
		Host:       config.Host,
	}, nil
}

func (k *KavenegarProvider) SendSMS(to, message string) (response.Response, error) {
	apiURL := fmt.Sprintf("http://%s/v1/%s/sms/send.json", k.Host, k.APIKey)

	body := url.Values{
		"message":    {message},
		"receptor":   {to},
		"linenumber": {k.LineNumber},
	}

	headers := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Accept":         "application/json",
		"Accept-Charset": "utf-8",
	}

	req := request.NewRequest(k.APIKey, k.Host)

	resp, err := req.Execute(apiURL, "POST", body, headers)
	if err != nil {
		return response.Response{}, err
	}
	defer resp.Body.Close()

	respBody := response.Response{Status: resp.StatusCode}
	if resp.StatusCode != http.StatusOK {
		re := new(ResultErrorKavenegar)
		if err := json.NewDecoder(resp.Body).Decode(&re); err != nil {
			return respBody, &errors.Error{
				Status:  resp.StatusCode,
				Message: resp.Status,
				Err:     err,
			}
		}
		return respBody, &errors.Error{
			Status:  re.ResultKavenegar.Status,
			Message: re.ResultKavenegar.Message,
			Err:     err,
		}
	}

	result := new(MessageResultKavenegar)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return respBody, &errors.Error{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Err:     err,
		}
	}

	respBody.Message = result.Message
	respBody.Status = result.Status
	return respBody, nil
}
