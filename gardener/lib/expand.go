package lib

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func ExpandSet(set *schema.Set) []string {
	result := make([]string, set.Len())
	for i, k := range set.List() {
		result[i] = k.(string)
	}

	return result
}

func ExpandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if v, ok := v.(string); ok {
			result[k] = v
		}
	}
	return result
}

func ExpandBoolMap(m map[string]interface{}) map[string]bool {
	result := make(map[string]bool)
	for k, v := range m {
		if v, ok := v.(bool); ok {
			result[k] = v
		}
	}
	return result
}
