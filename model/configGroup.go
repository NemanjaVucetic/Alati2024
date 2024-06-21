package model

// swagger:model ConfigGroup
type ConfigGroup struct {
	// Name of the configuration group
	// in: string
	Name string `json:"name"`

	// Version of the configuration group
	// in: int
	Version int `json:"version"`

	// Configurations in the group
	// in: map[string]Config
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
