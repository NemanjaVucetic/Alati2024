package model

import "context"

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

	Get(id string, ctx context.Context) (*ConfigGroup, error)
	GetAll(ctx context.Context) ([]ConfigGroup, error)
	Put(c *ConfigGroup, id string, ctx context.Context) (*ConfigGroup, error)
	Delete(id string, ctx context.Context) error
	AddConfigToGroup(group ConfigGroup, config Config, id string, ctx context.Context) error
	RemoveConfigFromGroup(group ConfigGroup, config Config, id string, ctx context.Context) error
	GetConfigsByLabels(prefixGroup string, prefixConf string, ctx context.Context) ([]Config, error)
	DeleteConfigsByLabels(prefixGroup string, prefixConf string, ctx context.Context) error

}
