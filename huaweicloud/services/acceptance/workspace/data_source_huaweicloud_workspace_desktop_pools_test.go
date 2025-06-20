package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopPoolsDataSource_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_desktop_pools.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByName   = "data.huaweicloud_workspace_desktop_pools.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByType   = "data.huaweicloud_workspace_desktop_pools.filter_by_type"
		dcFilterByType = acceptance.InitDataSourceCheck(filterByType)

		filterByEnterpriseProjectId   = "data.huaweicloud_workspace_desktop_pools.filter_by_enterprise_project_id"
		dcFilterByEnterpriseProjectId = acceptance.InitDataSourceCheck(filterByEnterpriseProjectId)

		filterByMaintenanceMode   = "data.huaweicloud_workspace_desktop_pools.filter_by_maintenance_mode"
		dcFilterByMaintenanceMode = acceptance.InitDataSourceCheck(filterByMaintenanceMode)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDesktopPools_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "pools.#"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.name"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.type"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.created_time"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.desktop_count"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.desktop_used"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.product.#"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.product.0.product_id"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.product.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.product.0.type"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.product.0.cpu"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.product.0.memory"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.image_id"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.image_name"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.root_volume.#"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.root_volume.0.type"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.root_volume.0.size"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.security_groups.#"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.security_groups.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "pools.0.status"),
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcFilterByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcFilterByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					dcFilterByMaintenanceMode.CheckResourceExists(),
					resource.TestCheckOutput("is_in_maintanence_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDesktopPools_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  name_regex = "WORKSPACE"
  visibility = "market"
}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name              = "%[1]s"
  type              = "STATIC"
  description       = "Created by terraform test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  subnet_id         = data.huaweicloud_workspace_service.test.network_ids[0]
  security_groups   = [
    data.huaweicloud_workspace_service.test.desktop_security_group[0].id,
    data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id,
  ]

  product {
    flavor_id = data.huaweicloud_workspace_flavors.test.flavors[0].id
  }

  image {
    image_type = "market"
    image_id   = try(data.huaweicloud_images_images.test.images[0].id, "")
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }
}
`, rName)
}

func testAccDataSourceDesktopPools_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  name = "%[2]s"
  type = "STATIC"
  enterprise_project_id = "0"
  in_maintenance_mode = true
}

# all
data "huaweicloud_workspace_desktop_pools" "all" {
  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

# filter by name
data "huaweicloud_workspace_desktop_pools" "filter_by_name" {
  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]

  name = local.name
}

locals {
  name_filter_result = [
     for v in data.huaweicloud_workspace_desktop_pools.filter_by_name.pools[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# filter by type
data "huaweicloud_workspace_desktop_pools" "filter_by_type" {
  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]

  type = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_type.pools[*].type : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# filter by enterprise project id
data "huaweicloud_workspace_desktop_pools" "filter_by_enterprise_project_id" {
  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]

  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_enterprise_project_id.pools[*].enterprise_project_id : v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}

# filter by in maintenance mode 
data "huaweicloud_workspace_desktop_pools" "filter_by_maintenance_mode" {
  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]

  in_maintenance_mode = local.in_maintenance_mode
}

locals {
  in_maintenance_mode_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_maintenance_mode.pools[*].in_maintenance_mode : v == local.in_maintenance_mode
  ]
}

output "is_in_maintanence_mode_filter_useful" {
  value = length(local.in_maintenance_mode_filter_result) > 0 && alltrue(local.in_maintenance_mode_filter_result)
}
`, testAccDataSourceDesktopPools_base(rName), rName)
}
