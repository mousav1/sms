package config

// Config represents the SMS gateway configuration.
type Config struct {
	DefaultDriver string                  `mapstructure:"default_driver"`
	Drivers       map[string]DriverConfig `mapstructure:"drivers"`
}

// DriverConfig represents the configuration for a specific SMS gateway driver.
type DriverConfig struct {
	APIKey     string `mapstructure:"api_key"`
	LineNumber string `mapstructure:"line_number"`
	Host       string `mapstructure:"host"`
}
