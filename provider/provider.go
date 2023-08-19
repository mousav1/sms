package provider

import (
	"fmt"
	"github.com/mousav1/sms"

	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/driver"
)

type SMSProviderFactory interface {
	CreateProvider(config config.DriverConfig) (sms.SMSProvider, error)
}

func NewSMSGateway(config *config.Config) (*sms.SMSGateway, error) {
	driverName := config.DefaultDriver
	driverConfig, exists := config.Drivers[driverName]
	if !exists {
		return nil, fmt.Errorf("default driver '%s' not found in the configuration", driverName)
	}

	providerFactory, err := GetProviderFactory(driverName)
	if err != nil {
		return nil, err
	}

	provider, err := providerFactory.CreateProvider(driverConfig)
	if err != nil {
		return nil, err
	}

	return &sms.SMSGateway{
		Provider: provider,
	}, nil
}

func GetProviderFactory(driverName string) (SMSProviderFactory, error) {
	switch driverName {
	case "Ghasedak":
		return &driver.Ghasedak{}, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}
}
