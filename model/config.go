package model

// swagger:model Config
type Config struct {
	// Name of the configuration
	// in: string
	Name string `json:"name"`

	// Version of the configuration
	// in: int
	Version int `json:"version"`

	// Parameters of the configuration
	// in: map[string]string
	Params map[string]string `json:"params"`

	// Labels associated with the configuration
	// in: map[string]string
	Labels map[string]string `json:"labels"`
}

type ConfigRepository interface {
	Get(name string, version int) (Config, error)
	Add(c Config)
	Delete(name string, version int) error
}
