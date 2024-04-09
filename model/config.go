package model

type Config struct {
	Name    string            `json:"name"`
	Version int               `json:"version"`
	Params  map[string]string `json:"params"`
}

type ConfigRepository interface {
	Get(name string, version int) (Config, error)
	Add(c Config)
	Delete(name string, version int) error
}
