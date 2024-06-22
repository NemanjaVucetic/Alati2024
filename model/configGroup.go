package model

// swagger:model ConfigGroup
type ConfigGroup struct {
	// Name of the configuration group
	// Required: true
	Name string `json:"name"`

	// Version of the configuration group
	// Required: true
	Version int `json:"version"`

	// Configs in the group
	// Required: true
	Configs map[string]Config `json:"configs"`
}

type ConfigGroupRepository interface {
	Get(id string) (*ConfigGroup, error)
	GetAll() ([]ConfigGroup, error)
	Put(c *ConfigGroup) (*ConfigGroup, error)
	Delete(id string) error
	AddConfigToGroup(group ConfigGroup, config Config) error
	RemoveConfigFromGroup(group ConfigGroup, config Config) error
	GetConfigsByLabels(prefixGroup string, prefixConf string) ([]Config, error)
	DeleteConfigsByLabels(prefixGroup string, prefixConf string) error
}
