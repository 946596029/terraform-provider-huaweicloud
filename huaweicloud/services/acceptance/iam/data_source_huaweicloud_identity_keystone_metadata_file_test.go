package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKeystoneMetadataFile_basic(t *testing.T) {
	dcName := "data.huaweicloud_identity_keystone_metadata_file.test"
	dc := acceptance.InitDataSourceCheck(dcName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeystoneMetadataFile_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "metadata_file"),
				),
			},
			{
				Config: testAccDataSourceKeystoneMetadataFile_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "metadata_file"),
				),
			},
		},
	})
}

const testAccDataSourceKeystoneMetadataFile_basic_step1 = `
data "huaweicloud_identity_keystone_metadata_file" "test" {}`

const testAccDataSourceKeystoneMetadataFile_basic_step2 = `
data "huaweicloud_identity_keystone_metadata_file" "test" {
  unsigned = true
}
`
