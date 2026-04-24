package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceArchitectureProcess_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_architecture_process.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_architecture_process.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byOwner   = "data.huaweicloud_dataarts_architecture_process.filter_by_owner"
		dcByOwner = acceptance.InitDataSourceCheck(byOwner)

		byCreator   = "data.huaweicloud_dataarts_architecture_process.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArchitectureProcess_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "processes.#"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByOwner.CheckResourceExists(),
					resource.TestCheckOutput("is_owner_filter_useful", "true"),
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceArchitectureProcess_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_process" "test" {}

data "huaweicloud_dataarts_architecture_process" "filter_by_name" {
  name = "test"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_name.processes[*].name: v != null
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0
}

data "huaweicloud_dataarts_architecture_process" "filter_by_owner" {
  owner = "test_owner"
}

locals {
  owner_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_owner.processes[*].owner: v == "test_owner"
  ]
}

output "is_owner_filter_useful" {
  value = length(local.owner_filter_result) > 0 && alltrue(local.owner_filter_result)
}

data "huaweicloud_dataarts_architecture_process" "filter_by_creator" {
  create_by = "test_creator"
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_creator.processes[*].create_by: v == "test_creator"
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}
`)
}
