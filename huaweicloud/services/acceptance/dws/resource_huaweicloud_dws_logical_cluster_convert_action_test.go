package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/go-uuid"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataLogicalClusterConvertAction_basic(t *testing.T) {
	var (
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_dws_logical_cluster_convert_action.test"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataLogicalClusterConvertAction_nonexistentCluster(name),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testAccDataLogicalClusterConvertAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccDataLogicalClusterConvertAction_nonexistentCluster(name string) string {
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dws_logical_cluster_convert_action" "nonexistent_cluster" {
  cluster_id           = "%[1]s"
  logical_cluster_name = "%[2]s"
}
`, randomUUID, name)
}

func testAccDataLogicalClusterConvertAction_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_logical_cluster_convert_action" "test" {
  cluster_id           = "%[1]s"
  logical_cluster_name = "%[2]s"
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
