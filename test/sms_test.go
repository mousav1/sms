package test

import (
	"github.com/mousav1/sms"
	"testing"
)

func TestSMSGateway_SendSMS(t *testing.T) {
	// Create a mock SMS provider for testing
	mockProvider := &MockSMSProvider{
		SendSMSFunc: func(to, message string) (sms.Response, error) {
			// Simulate a successful response
			response := sms.Response{
				Success: true,
				Code:    200,
				Message: "SMS sent successfully",
				ID:      12345,
			}
			return response, nil
		},
	}

	// Create an SMS gateway with the mock provider
	smsGateway := &sms.SMSGateway{
		Provider: mockProvider,
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

// MockSMSProvider is a mock implementation of the SMSProvider interface
type MockSMSProvider struct {
	SendSMSFunc func(to, message string) (sms.Response, error)
}

// SendSMS calls the SendSMSFunc function of the mock provider
func (m *MockSMSProvider) SendSMS(to, message string) (sms.Response, error) {
	return m.SendSMSFunc(to, message)
}
