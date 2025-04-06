package config

import "time"

var (
	Config config
)

// config represents the configuration settings for the application.
// It includes the host, port, timeouts, and storage directory for file uploads.
type config struct {
	// Host specifies the hostname or IP address of the server.
	Host string `json:"host"`

	// Port specifies the port number the server listens on.
	Port int `json:"port"`

	// Timeout defines the default duration for operations that do not specify a timeout.
	// This value is represented as a string (e.g., "30s", "2m").
	Timeout string `json:"timeout"`

	// ReadTimeout specifies the maximum duration for reading data from the server.
	// This is a parsed version of the Timeout field.
	ReadTimeout time.Duration `json:"read_timeout"`

	// WriteTimeout specifies the maximum duration for writing data to the server.
	// This is a parsed version of the Timeout field.
	WriteTimeout time.Duration `json:"write_timeout"`

	// Storage specifies the directory for storing uploaded files.
	Storage string `json:"-"`
}
