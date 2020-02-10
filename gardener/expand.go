package gardener

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func expandSet(set *schema.Set) []string {
	result := make([]string, set.Len())
	for i, k := range set.List() {
		result[i] = k.(string)
	}

	return result
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if v, ok := v.(string); ok {
			result[k] = v
		}
	}
	return result
}

func expandBoolMap(m map[string]interface{}) map[string]bool {
	result := make(map[string]bool)
	for k, v := range m {
		if v, ok := v.(bool); ok {
			result[k] = v
		}
	}
	return result
}

func RemoveInternalKeysMapMeta(aMap map[string]string, bMap map[string]interface{}) map[string]string {
	for key := range aMap {
		if _, ok := bMap[key]; !ok {
			delete(aMap, key)
		}
	}
	return aMap
}

func AddMissingDataForUpdate(old *corev1beta1.Shoot, new *corev1beta1.Shoot) {
	if new.Spec.DNS == nil {
		new.Spec.DNS = &corev1beta1.DNS{}
	}
	if new.Spec.DNS.Domain == nil || *new.Spec.DNS.Domain == "" {
		new.Spec.DNS.Domain = old.Spec.DNS.Domain
	}
	new.Spec.SeedName = old.Spec.SeedName
	new.ObjectMeta.ResourceVersion = old.ObjectMeta.ResourceVersion
	new.ObjectMeta.Finalizers = old.ObjectMeta.Finalizers
}

func RemoveInternalKeysMapSpec(aMap map[string]interface{}, bMap map[string]interface{}) map[string]interface{} {
	for key, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			if val2, ok := bMap[key]; !ok {
				delete(aMap, key)
			} else {
				aMap[key] = RemoveInternalKeysMapSpec(val.(map[string]interface{}), val2.(map[string]interface{}))
			}
		case []interface{}:
			if val2, ok := bMap[key]; !ok {
				delete(aMap, key)
			} else {
				aMap[key] = RemoveInternalKeysArraySpec(val.([]interface{}), val2.([]interface{}))
			}
		default:
			if val2, ok := bMap[key]; !ok || val2 == "" {
				delete(aMap, key)
			}
		}
	}
	return aMap
}

func RemoveInternalKeysArraySpec(ArrayA, ArrayB []interface{}) []interface{} {
	for i, val := range ArrayA {
		switch val.(type) {
		case map[string]interface{}:
			if i >= len(ArrayB) || ArrayB[i] == nil {
				ArrayA = remove(ArrayA, i)
			} else {
				ArrayA[i] = RemoveInternalKeysMapSpec(val.(map[string]interface{}), ArrayB[i].(map[string]interface{}))
			}
		case []interface{}:
			if i >= len(ArrayB) || ArrayB[i] == nil {
				ArrayA = remove(ArrayA, i)
			} else {
				ArrayA[i] = RemoveInternalKeysArraySpec(val.([]interface{}), ArrayB[i].([]interface{}))
			}
		default:
			if i >= len(ArrayB) || ArrayB[i] == nil || ArrayB[i] == "" {
				ArrayA = remove(ArrayA, i)
			}
		}
	}
	return ArrayA
}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
