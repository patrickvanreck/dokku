package resource

import (
	"fmt"

	"github.com/dokku/dokku/plugins/common"
)

// Resource is a collection of resource constraints for apps
type Resource struct {
	CPU            string `json:"cpu"`
	Memory         string `json:"memory"`
	MemorySwap     string `json:"memory-swap"`
	Network        string `json:"network"`
	NetworkIngress string `json:"network-ingress"`
	NetworkEgress  string `json:"network-egress"`
	NvidiaGPU      string `json:"nvidia-gpu"`
}

// ReportSingleApp is an internal function that displays the resource report for one or more apps
func ReportSingleApp(appName, format, infoFlag string) error {
	if err := common.VerifyAppName(appName); err != nil {
		return err
	}

	resources, err := common.PropertyGetAll("resource", appName)
	if err != nil {
		return nil
	}

	flags := map[string]string{}
	for key, value := range resources {
		flag := fmt.Sprintf("--resource-%v", key)
		flags[flag] = value
	}

	flagKeys := []string{}
	for flagKey := range flags {
		flagKeys = append(flagKeys, flagKey)
	}

	trimPrefix := true
	uppercaseFirstCharacter := false
	return common.ReportSingleApp("resource", appName, infoFlag, flags, flagKeys, format, trimPrefix, uppercaseFirstCharacter)
}

// GetResourceValue fetches a single value for a given app/process/request/key combination
func GetResourceValue(appName string, processType string, resourceType string, key string) (string, error) {
	resources, err := common.PropertyGetAll("resource", appName)
	if err != nil {
		return "", err
	}

	defaultValue := ""
	for k, value := range resources {
		if k == propertyKey("_default_", resourceType, key) {
			defaultValue = value
		}
		if k == propertyKey(processType, resourceType, key) {
			return value, nil
		}
	}

	return defaultValue, nil
}

func propertyKey(processType string, resourceType string, key string) string {
	return fmt.Sprintf("%v.%v.%v", processType, resourceType, key)
}
