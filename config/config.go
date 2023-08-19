package config

// Config represents the SMS gateway configuration.
type Config struct {
	DefaultDriver string                  `json:"default_driver"`
	Drivers       map[string]DriverConfig `json:"drivers"`
}

// DriverConfig represents the configuration for a specific SMS gateway driver.
type DriverConfig struct {
	// Driver-specific configuration fields.
	APIKey     string `json:"api_key"`
	LineNumber string `json:"LineNumber"`
	Host       string `json:"host"`
}
