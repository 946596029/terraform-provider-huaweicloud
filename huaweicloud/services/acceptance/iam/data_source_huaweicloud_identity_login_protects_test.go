package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLoginProtects_basic(t *testing.T) {
	var (
		userName = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		dcName = "data.huaweicloud_identity_login_protects.test"
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
				Config: testAccDataSourceLoginProtects_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "login_protects.#", "0"),
				),
			},
			{
				Config: testAccDataSourceLoginProtects_basic_step2(userName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "login_protects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "login_protects.0.user_id"),
					resource.TestCheckResourceAttr(dcName, "login_protects.0.enabled", "true"),
					resource.TestCheckResourceAttr(dcName, "login_protects.0.verification_method", "email"),
				),
			},
		},
	})
}

const testAccDataSourceLoginProtects_basic_step1 = `
data "huaweicloud_identity_login_protects" "test" {}
`

func testAccDataSourceLoginProtects_basic_step2(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = "%[2]s"
  enabled     = true
  email       = "%[1]s@abc.com"
  description = "tested by terraform"
  
  login_protect_verification_method = "email"
}

data "huaweicloud_identity_login_protects" "test" {
  user_id = huaweicloud_identity_user.test.id
}
`, name, password)
}
