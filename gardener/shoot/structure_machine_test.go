package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getMachineTestData() (map[string]interface{}, *corev1beta1.Machine) {
	version := "1.0.0"
	machineRaw := map[string]interface{}{
		"type": "small",
		"image": []interface{}{
			map[string]interface{}{
				"name":    "foo",
				"version": version,
			},
		},
	}
	machine := &corev1beta1.Machine{
		Type: "small",
		Image: &corev1beta1.ShootMachineImage{
			Name:    "foo",
			Version: &version,
		},
	}

	return machineRaw, machine
}

func TestExpandMachine(t *testing.T) {
	machineRaw, machine := getMachineTestData()
	data := schema.TestResourceDataRaw(t, MachineResource().Schema, machineRaw)
	out := ExpandMachine([]interface{}{data.Get("")})
	if diff := cmp.Diff(machine, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenMachine(t *testing.T) {
	machineRaw, machine := getMachineTestData()
	out := FlattenMachine(machine)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, MachineResource().Schema, machineRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
