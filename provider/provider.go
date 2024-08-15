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

var providerFactories = map[string]SMSProviderFactory{
	"Kavenegar": &driver.KavenegarProvider{},
	"Ghasedak":  &driver.GhasedakProvider{},
}

func NewSMSGateway(config *config.Config) (*sms.SMSGateway, error) {

	driverName := config.DefaultDriver
	driverConfig, exists := config.Drivers[driverName]
	if !exists {
		return nil, fmt.Errorf("default driver '%s' not found in the configuration", driverName)
	}

	providerFactory, exists := providerFactories[driverName]
	if !exists {
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}

	provider, err := providerFactory.CreateProvider(driverConfig)
	if err != nil {
		return nil, err
	}

	return &sms.SMSGateway{Provider: provider}, nil
}

// SetProviderFactories allows setting provider factories for testing.
func SetProviderFactories(factories map[string]SMSProviderFactory) {
	providerFactories = factories
}
