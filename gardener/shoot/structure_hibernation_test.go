package shoot

import (
	"testing"

	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

func getHibernationTestData() (map[string]interface{}, *corev1beta1.Hibernation) {
	start := "00 17 * * 1"
	end := "00 00 * * 1"
	location := "Pacific/Auckland"
	enabled := true
	hibernationRaw := map[string]interface{}{
		"enabled": enabled,
		"schedules": []interface{}{
			map[string]interface{}{
				"start":    start,
				"end":      end,
				"location": location,
			},
		},
	}
	hibernation := &corev1beta1.Hibernation{
		Enabled: &enabled,
		Schedules: []corev1beta1.HibernationSchedule{corev1beta1.HibernationSchedule{
			Start:    &start,
			End:      &end,
			Location: &location,
		}},
	}

	return hibernationRaw, hibernation
}

func TestExpandHibernation(t *testing.T) {
	hibernationRaw, hibernation := getHibernationTestData()
	data := schema.TestResourceDataRaw(t, HibernationResource().Schema, hibernationRaw)
	out := ExpandHibernation([]interface{}{data.Get("")})
	if diff := cmp.Diff(hibernation, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenHibernation(t *testing.T) {
	hibernationRaw, hibernation := getHibernationTestData()
	out := FlattenHibernation(hibernation)[0].(map[string]interface{})
	data := schema.TestResourceDataRaw(t, HibernationResource().Schema, hibernationRaw).Get("")
	if diff := cmp.Diff(data, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
