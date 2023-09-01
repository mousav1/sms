package sms

import (
	"encoding/json"
	"os"

	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/response"
)

// SMSProvider defines the interface that SMS gateway providers should implement.
type SMSProvider interface {
	SendSMS(to, message string) (response.Response, error)
}

// SMSGateway represents an SMS gateway.
type SMSGateway struct {
	Provider SMSProvider
}

// SendSMS sends an SMS using the gateway's provider.
func (g *SMSGateway) SendSMS(to, message string) (response.Response, error) {
	return g.Provider.SendSMS(to, message)
}

func LoadConfig(filename string) (*config.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config config.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
