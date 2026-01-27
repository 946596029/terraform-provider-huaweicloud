package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProjectQuota_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_identity_project_quota.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjectQuota_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(dcName, "resources.0.type", "project"),
				),
			},
		},
	})
}

func testAccDataSourceProjectQuota_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%[1]s_%[2]s"
  description = "Terraform test project"
}

data "huaweicloud_identity_project_quota" "test" {
  project_id = huaweicloud_identity_project.test.id
}
`, acceptance.HW_REGION_NAME, name)
}
