package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkspaceDesktopPoolsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_workspace_desktop_pools.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceDesktopPoolsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.desktop_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.desktop_used"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.product.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.product.0.product_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.product.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.product.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.product.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.product.0.memory"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.root_volume.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.root_volume.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.root_volume.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.security_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktop_pools.0.status"),
				),
			},
		},
	})
}

func testAccWorkspaceDesktopPoolsDataSource_basic (rName string) string {
	return fmt.Sprintf(`
%[1]s 

data "huaweicloud_workspace_desktop_pools" "test" {}

data "huaweicloud_workspace_desktop_pools" "test" {
  name = "%s[2]s"
}
`, rName, rName)
}
