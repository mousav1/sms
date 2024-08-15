package sms

import (
	"fmt"

	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/response"
	"github.com/spf13/viper"
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
	viper.SetConfigFile(filename)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
