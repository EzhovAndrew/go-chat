package config

// Config holds the gateway configuration
type Config struct {
	HTTPPort string
	Services ServiceAddresses
}

// ServiceAddresses contains addresses of backend gRPC services
type ServiceAddresses struct {
	Auth          string
	Users         string
	Chat          string
	Social        string
	Notifications string
}

// New creates a new Config with default values
func New() *Config {
	return &Config{
		HTTPPort: ":8080",
		Services: ServiceAddresses{
			Auth:          "auth:8080",
			Users:         "users:8080",
			Chat:          "chat:8080",
			Social:        "social:8080",
			Notifications: "notifications:8080",
		},
	}
}
