package lib

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IdParts(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		err := fmt.Errorf("Unexpected ID format (%q), expected %q", id, "namespace/name")
		return "", "", err
	}

	return parts[0], parts[1], nil
}

func BuildID(meta metav1.ObjectMeta) string {
	return meta.Namespace + "/" + meta.Name
}

func NewStringSet(f schema.SchemaSetFunc, in []string) *schema.Set {
	var out = make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return schema.NewSet(f, out)
}

func FlattenStringMap(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}
