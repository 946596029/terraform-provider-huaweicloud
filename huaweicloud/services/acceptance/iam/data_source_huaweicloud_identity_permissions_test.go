package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePermissions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_identity_permissions.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		dcNameByName = "data.huaweicloud_identity_permissions.filter_by_name"
		dcByName     = acceptance.InitDataSourceCheck(dcNameByName)

		dcNameByCatalog = "data.huaweicloud_identity_permissions.filter_by_catalog"
		dcByCatalog     = acceptance.InitDataSourceCheck(dcNameByCatalog)	

		dcNameByType = "data.huaweicloud_identity_permissions.filter_by_type"
		dcByType     = acceptance.InitDataSourceCheck(dcNameByType)

		dcNameByCustom = "data.huaweicloud_identity_permissions.filter_by_custom"
		dcByCustom     = acceptance.InitDataSourceCheck(dcNameByCustom)

		dcNameByScopeType = "data.huaweicloud_identity_permissions.filter_by_scope_type"
		dcByScopeType     = acceptance.InitDataSourceCheck(dcNameByScopeType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePermissions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "permissions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByCatalog.CheckResourceExists(),
					resource.TestCheckOutput("catalog_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					dcByCustom.CheckResourceExists(),
					resource.TestCheckOutput("custom_filter_is_useful", "true"),
					dcByScopeType.CheckResourceExists(),
					resource.TestCheckOutput("scope_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourcePermissions_basic(name string) string {
	return fmt.Sprintf(`
# All
data "huaweicloud_identity_permissions" "all" {
}

# Filter by name
data "huaweicloud_identity_permissions" "filter_by_name" {
  name = "%[1]s"
}

output "name_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.filter_by_name.permissions[*].name : v == "KMS Administrator"])
}

# Filter by catalog
data "huaweicloud_identity_permissions" "filter_by_catalog" {
  catalog = "ELB"
}

output "catalog_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.filter_by_catalog.permissions[*].catalog : v == "ELB"])
}

# Filter by type
data "huaweicloud_identity_permissions" "filter_by_type" {
  type    = "system-policy"
  catalog = "OBS"
}

output "type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.filter_by_type.permissions[*].catalog : v == "OBS"])
}

# Filter by custom
data "huaweicloud_identity_permissions" "filter_by_custom" {
  type = "custom"
}

output "custom_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.filter_by_custom.permissions[*].catalog : v == "CUSTOMED"])
}

# Filter by scope type
data "huaweicloud_identity_permissions" "filter_by_scope_type" {
  scope_type = "project"
  name       = "CCE FullAccess"
}

output "scope_type_filter_is_useful" {
  value = alltrue([length(data.huaweicloud_identity_permissions.filter_by_scope_type.permissions) == 1])
}
`, name)
}
