package repo

import (
	"alati/model"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	config    = "config/%s/%s"
	all       = "config/"
	group     = "configGroups/%s/%s"
	allGroups = "configGroups"
	groupMap  = "configGroups/%s/%s/config/"
)

func constructKey(name string, version string) string {
	return fmt.Sprintf(config, name, version)
}

func constructKeyGroup(name string, version string) string {
	return fmt.Sprintf(group, name, version)
}

func constructKeyInGroup(g model.ConfigGroup, config model.Config) string {
	k := fmt.Sprintf(groupMap, g.Name, strconv.Itoa(g.Version))

	keys := make([]string, 0, len(config.Labels))
	for key := range config.Labels {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var builder strings.Builder

	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s:%s", key, config.Labels[key]))
		builder.WriteString("/")
	}

	result := builder.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	k += result
	k += "/" + config.Name + "/" + strconv.Itoa(config.Version)
	return k
}
