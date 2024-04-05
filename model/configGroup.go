package model

type ConfigGroup struct {
	Name    string
	Version float64
	Configs map[string]Config
}

type ConfigGroupRepository interface {
	Get(name string, version int) (ConfigGroup, error)
	Add(c ConfigGroup)
	Delete(name string, version int) error
}
