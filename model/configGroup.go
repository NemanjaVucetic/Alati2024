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
	Get(name string, version int) (ConfigGroup, error)
	Add(c ConfigGroup) error
	Delete(name string, version int) error
	AddConfigToGroup(group ConfigGroup, config Config) error
	RemoveConfigFromGroup(group ConfigGroup, key string) error
	GetConfigsByLabels(group ConfigGroup, labels *map[string]string) ([]Config, error)
	DeleteConfigsByLabels(group ConfigGroup, labels *map[string]string) error
}
