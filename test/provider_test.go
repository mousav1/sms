package test_test

import (
	"testing"

	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/driver"
	"github.com/mousav1/sms/provider"
)

func TestNewSMSGateway(t *testing.T) {
	// Create a test configuration
	testConfig := &config.Config{
		DefaultDriver: "Ghasedak",
		Drivers: map[string]config.DriverConfig{
			"Ghasedak": {
				APIKey:     "your-api-key",
				LineNumber: "your-line-number",
				Host:       "your-host",
			},
		},
	}

	// Create a new SMS gateway
	smsGateway, err := provider.NewSMSGateway(testConfig)
	if err != nil {
		t.Fatalf("Failed to create SMS gateway: %v", err)
	}

	// Check if the provider is of type *driver.Ghasedak
	if _, ok := smsGateway.Provider.(*driver.Ghasedak); !ok {
		t.Error("SMS gateway provider is not of type *driver.Ghasedak")
	}

	// Send a test SMS
	to := "recipient-number"
	message := "test message"
	response, err := smsGateway.SendSMS(to, message)
	if err != nil {
		t.Fatalf("Failed to send SMS: %v", err)
	}

	// Verify the response
	if !response.Success {
		t.Errorf("SMS sending failed: %s", response.Message)
	}
}
