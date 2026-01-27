package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGroup_basic(t *testing.T) {
	var (
		rName    = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		dcName = "data.huaweicloud_identity_group.test"
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
				Config: testAccDataSourceGroup_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "name", rName),
					resource.TestCheckResourceAttr(dcName, "users.#", "0"),
				),
			},
			{
				Config: testAccDataSourceGroup_basic_step2(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "name", rName),
					resource.TestCheckResourceAttr(dcName, "users.#", "2"),
					resource.TestCheckResourceAttrSet(dcName, "users.0.id"),
				),
			},
		},
	})
}

func testAccDataSourceGroup_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%[1]s"
  description = "An ACC test group"
}
`, rName)
}

func testAccDataSourceGroup_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_group" "test" {
  name = huaweicloud_identity_group.test.name
}
`, testAccDataSourceGroup_base(rName))
}

func testAccDataSourceGroup_basic_step2(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "test" {
  count    = 2
  name     = "%[2]s-${count.index}"
  password = "%[3]s"
  enabled  = true
}

resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = [
    for user in huaweicloud_identity_user.test : user.id
  ]
}

data "huaweicloud_identity_group" "test" {
  name = huaweicloud_identity_group.test.name

  depends_on = [
    huaweicloud_identity_group_membership.test
  ]
}
`, testAccDataSourceGroup_base(rName), rName, password)
}
