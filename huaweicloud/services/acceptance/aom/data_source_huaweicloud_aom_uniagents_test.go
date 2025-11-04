package aom

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceUniAgents_basic(t *testing.T) {
	var (
		all            = "data.huaweicloud_aom_uniagents.all"
		dcForAll       = acceptance.InitDataSourceCheck(all)

		withPageSize      = "data.huaweicloud_aom_uniagents.with_page_size"
		dcForWithPageSize = acceptance.InitDataSourceCheck(withPageSize)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUniAgents_basic,
				Check: resource.ComposeTestCheckFunc(
					dcForAll.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "total_count"),
					resource.TestMatchResourceAttr(all, "uniagents.#", regexp.MustCompile(`^[0-9]*$`)),
					// Check the attributes.
					resource.TestCheckResourceAttrSet(all, "uniagents.0.inner_ip"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.import_ip"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.agent_id"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.host_name"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.os_type"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.agent_state"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.project_id"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.version"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.is_hw_cloud_host"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.vpc_id"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.cmdb_id"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.ecs_id"),
					resource.TestCheckResourceAttrSet(all, "uniagents.0.domain_id"),
					// With page_size.
					dcForWithPageSize.CheckResourceExists(),
					resource.TestMatchResourceAttr(withPageSize, "uniagents.#", regexp.MustCompile(`^[0-9]*$`)),
					resource.TestCheckResourceAttrSet(withPageSize, "uniagents.0.agent_id"),
				),
			},
		},
	})
}

const testAccDataSourceUniAgents_basic = `
data "huaweicloud_aom_uniagents" "all" {}

data "huaweicloud_aom_uniagents" "with_page_size" {
  page_size = 20
}
`

