package shoot

import (
	"testing"

	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/api/core/v1"
)

func getTaintTestData() ([]interface{}, []v1.Taint) {
	taintRaw := []interface{}{
		map[string]interface{}{
			"key":    "key",
			"value":  "value",
			"effect": "NoExecute",
		},
	}
	taint := []v1.Taint{v1.Taint{
		Key:    "key",
		Value:  "value",
		Effect: v1.TaintEffectNoExecute,
	}}

	return taintRaw, taint
}

func TestExpandTaint(t *testing.T) {
	taintRaw, taint := getTaintTestData()
	data := schema.TestResourceDataRaw(t, TaintResource().Schema, taintRaw[0].(map[string]interface{}))
	out := ExpandTaint([]interface{}{data.Get("")})
	if diff := cmp.Diff(taint, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenTaint(t *testing.T) {
	taintRaw, taint := getTaintTestData()
	out := FlattenTaint(taint)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, TaintResource().Schema, taintRaw[0].(map[string]interface{})).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
