package model

import "context"

// swagger:model Config
type Config struct {
	// Name of the configuration
	// Required: true
	Name string `json:"name"`

	// Version of the configuration
	// Required: true
	Version int `json:"version"`

	// Params are key-value pairs for configuration
	// Required: true
	Params map[string]string `json:"params"`

	// Labels are key-value pairs for configuration
	// Required: true
	Labels map[string]string `json:"labels"`
}
type ConfigRepository interface {

	Get(id string, ctx context.Context) (*Config, error)
	GetAll(ctx context.Context) ([]Config, error)
	Put(c *Config, id string, ctx context.Context) (*Config, error)
	Delete(id string, ctx context.Context) error
	DeleteAll(ctx context.Context) error

}
