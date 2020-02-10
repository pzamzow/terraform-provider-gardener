package gardener

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func MetadataResource(objectName string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:             schema.TypeMap,
				Description:      fmt.Sprintf("An unstructured key value map stored with the %s that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations", objectName),
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ValidateFunc:     ValidateAnnotations,
				DiffSuppressFunc: suppressCreatedByAnnotation,
			},
			"generation": {
				Type:        schema.TypeInt,
				Description: "A sequence number representing a specific generation of the desired state.",
				Computed:    true,
			},
			"labels": {
				Type:             schema.TypeMap,
				Description:      fmt.Sprintf("Map of string keys and values that can be used to organize and categorize (scope and select) the %s. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels", objectName),
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ValidateFunc:     ValidateLabels,
				DiffSuppressFunc: suppressStatusLabel,
			},
			"name": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Name of the %s, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names", objectName),
				Optional:    true,

				Computed:     true,
				ValidateFunc: ValidateName,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("An opaque value that represents the internal version of this %s that can be used by clients to determine when %s has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#concurrency-control-and-consistency", objectName, objectName),
				Computed:    true,
			},
			"self_link": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("A URL representing this %s.", objectName),
				Computed:    true,
			},
			"uid": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The unique in time and space value for this %s. More info: http://kubernetes.io/docs/user-guide/identifiers#uids", objectName),
				Computed:    true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Namespace defines the space within which name of the %s must be unique.", objectName),
				Optional:    true,
				Default:     "default",
			},
		},
	}
}

func ExpandMetadata(in []interface{}) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{}
	if len(in) < 1 {
		return meta
	}
	m := in[0].(map[string]interface{})

	if v, ok := m["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Annotations = expandStringMap(m["annotations"].(map[string]interface{}))
	}
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	meta.Annotations["confirmation.garden.sapcloud.io/deletion"] = "true"
	if v, ok := m["labels"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Labels = expandStringMap(m["labels"].(map[string]interface{}))
	}

	if v, ok := m["generate_name"]; ok {
		meta.GenerateName = v.(string)
	}
	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}
	if v, ok := m["namespace"]; ok {
		meta.Namespace = v.(string)
	}

	return meta
}

func FlattenMetadata(meta metav1.ObjectMeta, d *schema.ResourceData, metaPrefix ...string) []interface{} {
	m := make(map[string]interface{})
	prefix := ""
	if len(metaPrefix) > 0 {
		prefix = metaPrefix[0]
	}
	configAnnotations := d.Get(prefix + "metadata.0.annotations").(map[string]interface{})
	m["annotations"] = RemoveInternalKeysMapMeta(meta.Annotations, configAnnotations)
	if meta.GenerateName != "" {
		m["generate_name"] = meta.GenerateName
	}
	configLabels := d.Get(prefix + "metadata.0.labels").(map[string]interface{})
	m["labels"] = RemoveInternalKeysMapMeta(meta.Labels, configLabels)
	m["name"] = meta.Name
	m["resource_version"] = meta.ResourceVersion
	m["self_link"] = meta.SelfLink
	m["uid"] = fmt.Sprintf("%v", meta.UID)
	m["generation"] = meta.Generation

	if meta.Namespace != "" {
		m["namespace"] = meta.Namespace
	}

	return []interface{}{m}
}
