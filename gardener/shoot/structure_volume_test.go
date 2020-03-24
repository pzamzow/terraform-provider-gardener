package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getVolumeTestData() (map[string]interface{}, *corev1beta1.Volume) {
	volumeType := "foo"
	volumeRaw := map[string]interface{}{
		"type": volumeType,
		"size": "1Gb",
	}
	volume := &corev1beta1.Volume{
		Type:       &volumeType,
		VolumeSize: "1Gb",
	}

	return volumeRaw, volume
}

func TestExpandVolume(t *testing.T) {
	volumeRaw, volume := getVolumeTestData()
	data := schema.TestResourceDataRaw(t, VolumeResource().Schema, volumeRaw)
	out := ExpandVolume([]interface{}{data.Get("")})
	if diff := cmp.Diff(volume, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenVolume(t *testing.T) {
	volumeRaw, volume := getVolumeTestData()
	out := FlattenVolume(volume)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, VolumeResource().Schema, volumeRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
