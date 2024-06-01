package model

type Config struct {
	Name    string            `json:"name"`
	Version int               `json:"version"`
	Params  map[string]string `json:"params"`
	Labels  map[string]string `json:"labels"`
}

type ConfigRepository interface {
	Get(id string) (*Config, error)
	GetAll() ([]Config, error)
	Put(c *Config) (*Config, error)
	Delete(id string) error
	DeleteAll() error
}
