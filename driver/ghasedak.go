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

type GhasedakProvider struct {
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
func (g *GhasedakProvider) CreateProvider(config config.DriverConfig) (sms.SMSProvider, error) {
	return &GhasedakProvider{
		APIKey:     config.APIKey,
		LineNumber: config.LineNumber,
		Host:       config.Host,
	}, nil
}

func (g *GhasedakProvider) SendSMS(to, message string) (response.Response, error) {
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

	req := request.NewRequest(g.APIKey, g.Host)

	resp, err := req.Execute(apiURL, "POST", body, headers)
	if err != nil {
		return response.Response{}, err
	}
	defer resp.Body.Close()

	respBody := response.Response{Status: resp.StatusCode}
	if resp.StatusCode != http.StatusOK {
		re := new(ResultErrorGhasedak)
		if err := json.NewDecoder(resp.Body).Decode(&re); err != nil {
			return respBody, &errors.Error{
				Status:  resp.StatusCode,
				Message: resp.Status,
				Err:     err,
			}
		}
		return respBody, &errors.Error{
			Status:  re.ResultGhasedak.Code,
			Message: re.ResultGhasedak.Message,
			Err:     err,
		}
	}

	result := new(MessageResult)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return respBody, &errors.Error{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Err:     err,
		}
	}

	respBody.Message = result.ResultGhasedak.Message
	if len(result.Items) > 0 {
		respBody.MessageID = result.Items[0]
	}
	return respBody, nil
}
