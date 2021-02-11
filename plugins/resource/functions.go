package resource

import (
	"fmt"

	"github.com/dokku/dokku/plugins/common"
)

func clearByResourceType(appName string, processType string, resourceType string) {
	noun := "limits"
	if resourceType == "reserve" {
		noun = "reservation"
	}

	message := fmt.Sprintf("clearing %v %v", appName, noun)
	if processType != "_default_" && processType != "" {
		message = fmt.Sprintf("%v (%v)", message, processType)
	}
	common.LogInfo2Quiet(message)

	if processType == "" {
		resources, err := common.PropertyGetAll("resource", appName)
		if err != nil {
			return
		}
		for key := range resources {
			common.PropertyDelete("resource", appName, key)
		}
	} else {
		resources := []string{
			"cpu",
			"memory",
			"memory-swap",
			"network",
			"network-ingress",
			"network-egress",
			"nvidia-gpu",
		}

		for _, key := range resources {
			property := propertyKey(processType, resourceType, key)
			common.PropertyDelete("resource", appName, property)
		}
	}
}

func setResourceType(appName string, processType string, r Resource, resourceType string) error {
	if len(processType) == 0 {
		processType = "_default_"
	}

	resources := map[string]string{
		"cpu":             r.CPU,
		"memory":          r.Memory,
		"memory-swap":     r.MemorySwap,
		"network":         r.Network,
		"network-ingress": r.NetworkIngress,
		"network-egress":  r.NetworkEgress,
		"nvidia-gpu":      r.NvidiaGPU,
	}

	hasValues := false
	for _, value := range resources {
		if value != "" {
			hasValues = true
		}
	}

	if !hasValues {
		reportResourceType(appName, processType, resourceType)
		return nil
	}

	noun := "limits"
	if resourceType == "reserve" {
		noun = "reservation"
	}
	message := fmt.Sprintf("Setting resource %v for %v", noun, appName)
	if processType != "_default_" {
		message = fmt.Sprintf("%v (%v)", message, processType)
	}
	common.LogInfo2Quiet(message)

	for key, value := range resources {
		if value == "" {
			continue
		}

		property := propertyKey(processType, resourceType, key)
		if value == "clear" {
			common.LogVerbose(fmt.Sprintf("%v: cleared", key))
			if err := common.PropertyDelete("resource", appName, property); err != nil {
				return err
			}
			continue
		}

		common.LogVerbose(fmt.Sprintf("%v: %v", key, value))
		if err := common.PropertyWrite("resource", appName, property, value); err != nil {
			return err
		}
	}

	return nil
}

func reportResourceType(appName string, processType string, resourceType string) {
	noun := "limits"
	if resourceType == "reserve" {
		noun = "reservation"
	}

	message := fmt.Sprintf("resource %v %v information", noun, appName)
	if processType == "_default_" {
		message = fmt.Sprintf("%v [defaults]", message)
	} else {
		message = fmt.Sprintf("%v (%v)", message, processType)
	}
	common.LogInfo2Quiet(message)

	resources := []string{
		"cpu",
		"memory",
		"memory-swap",
		"network",
		"network-ingress",
		"network-egress",
		"nvidia-gpu",
	}

	for _, key := range resources {
		property := propertyKey(processType, resourceType, key)
		value := common.PropertyGet("resource", appName, property)
		common.LogVerbose(fmt.Sprintf("%v: %v", key, value))
	}
}
