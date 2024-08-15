package provider

import (
	"testing"

	"github.com/mousav1/sms"
	"github.com/mousav1/sms/config"
	"github.com/mousav1/sms/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSMSProvider is a mock implementation of the SMSProvider interface
type MockSMSProvider struct {
	mock.Mock
}

// SendSMS is a mock method for sending SMS
func (m *MockSMSProvider) SendSMS(to, message string) (response.Response, error) {
	args := m.Called(to, message)
	return args.Get(0).(response.Response), args.Error(1)
}

// MockSMSProviderFactory is a mock implementation of the SMSProviderFactory interface
type MockSMSProviderFactory struct {
	mock.Mock
}

func (m *MockSMSProviderFactory) CreateProvider(cfg config.DriverConfig) (sms.SMSProvider, error) {
	args := m.Called(cfg)
	return args.Get(0).(sms.SMSProvider), args.Error(1)
}

// TestNewSMSGateway tests the NewSMSGateway function
func TestNewSMSGateway(t *testing.T) {
	mockProvider := new(MockSMSProvider)
	mockFactory := new(MockSMSProviderFactory)

	// Set up the mock provider
	mockProvider.On("SendSMS", "1234567890", "Test message").Return(response.Response{Message: "Success"}, nil)

	// Set up the mock factory to return the mock provider
	mockFactory.On("CreateProvider", config.DriverConfig{
		APIKey:     "mockApiKey",
		LineNumber: "mockLineNumber",
		Host:       "mockHost",
	}).Return(mockProvider, nil)

	// Set mock factory in providerFactories
	SetProviderFactories(map[string]SMSProviderFactory{
		"MockDriver": mockFactory,
	})

	cfg := &config.Config{
		DefaultDriver: "MockDriver",
		Drivers: map[string]config.DriverConfig{
			"MockDriver": {
				APIKey:     "mockApiKey",
				LineNumber: "mockLineNumber",
				Host:       "mockHost",
			},
		},
	}

	gateway, err := NewSMSGateway(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, gateway)
	assert.Equal(t, mockProvider, gateway.Provider)

	// Verify that CreateProvider was called with the correct config
	mockFactory.AssertCalled(t, "CreateProvider", config.DriverConfig{
		APIKey:     "mockApiKey",
		LineNumber: "mockLineNumber",
		Host:       "mockHost",
	})

	// Test the SendSMS method
	response, err := gateway.Provider.SendSMS("1234567890", "Test message")
	assert.NoError(t, err)
	assert.Equal(t, "Success", response.Message)
}

// TestNewSMSGateway_InvalidDriver tests the case where the driver is not found
func TestNewSMSGateway_InvalidDriver(t *testing.T) {
	cfg := &config.Config{
		DefaultDriver: "InvalidDriver",
		Drivers:       map[string]config.DriverConfig{}, // No drivers are defined here
	}

	_, err := NewSMSGateway(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "default driver 'InvalidDriver' not found in the configuration")
}

// TestNewSMSGateway_NoConfig tests the case where no driver config is found
func TestNewSMSGateway_NoConfig(t *testing.T) {
	mockFactory := new(MockSMSProviderFactory)

	SetProviderFactories(map[string]SMSProviderFactory{
		"MockDriver": mockFactory,
	})

	cfg := &config.Config{
		DefaultDriver: "MockDriver",
		Drivers:       map[string]config.DriverConfig{},
	}

	_, err := NewSMSGateway(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "default driver 'MockDriver' not found in the configuration")
}
