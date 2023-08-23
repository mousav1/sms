# Golang SMS Gateway

This is a Go Package for SMS Gateway

List of supported gateways:

-   [Ghasedak](https://ghasedak.me/)


-  Will be added soon.


## Install

```bash
go get https://github.com/mousav1/sms
```

## Configure

Add the config file to the project

In the configuration file, you can choose the default driver and also set the settings of each driver.

```json
{
    "default_driver": "Ghasedak",
    "drivers": {
      "Ghasedak": {
        "api_key": "api_key",
        "line_number": "line_number",
        "host": "api.ghasedak.me"
      }
    }
  }
```


## Usage

In your code just use it like this.

```go
	config, err := sms.LoadConfig("sms_config.json")
	if err != nil {
		fmt.Println("Failed to load SMS smsGateway configuration:", err)
		return
	}

	smsGateway, err := provider.NewSMSGateway(config)
	if err != nil {
		fmt.Println("Failed to create SMS smsGateway:", err)
		return
	}

	// Use the SMS smsGateway to send messages
	response, err := smsGateway.SendSMS("Number", "message")
	if err != nil {
		fmt.Println("Failed to send SMS:", err)
		return
	}

```

#### Add custom driver

1- Add custom driver to the driver directory

2- Create the code structure for the driver

3- Implement driver methods:

```go

// CreateProvider creates an instance of the mydriver provider.
func (g *mydriver) CreateProvider(config config.DriverConfig) (sms.SMSProvider, error) {}

// SendSMS sends an SMS using mydriver.
func (g *mydriver) SendSMS(to, message string) (sms.Response, error) {}

```

4- Add the driver to the provider.go in package:

```go
func GetProviderFactory(driverName string) (SMSProviderFactory, error) {
	switch driverName {
	case "Ghasedak":
		return &driver.Ghasedak{}, nil
	case "mydriver":
		return &driver.mydriver{}, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}
}
```

