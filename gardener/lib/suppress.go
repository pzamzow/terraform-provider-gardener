package lib

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// SuppressMissingOptionalConfigurationBlock handles configuration block attributes in the following scenario:
//  * The resource schema includes an optional configuration block with defaults
//  * The API response includes those defaults to refresh into the Terraform state
//  * The operator's configuration omits the optional configuration block
func SuppressMissingOptionalConfigurationBlock(k, old, new string, d *schema.ResourceData) bool {
	return old == "1" && new == "0"
}

func SuppressEmptyNewValue(k, old, new string, d *schema.ResourceData) bool {
	return len(new) == 0
}

func SuppressZeroNewValue(k, old, new string, d *schema.ResourceData) bool {
	return new == "0"
}

func SuppressGardenerLabels(k, old, new string, d *schema.ResourceData) bool {
	return strings.Contains(k, "garden.sapcloud.io") || strings.Contains(k, "gardener.cloud") || strings.HasSuffix(k, "labels.%")
}

func SuppressGardenerAnnotations(k, old, new string, d *schema.ResourceData) bool {
	return strings.Contains(k, "garden.sapcloud.io") || strings.Contains(k, "gardener.cloud") || strings.Contains(k, "kubectl.kubernetes.io") || strings.HasSuffix(k, "annotations.%")
}
