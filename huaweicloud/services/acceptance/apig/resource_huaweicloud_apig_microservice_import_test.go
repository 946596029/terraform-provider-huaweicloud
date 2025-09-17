package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccApigMicroserviceImport_basic(t *testing.T) {
	var (
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_apig_microservice_import.test"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApigMicroserviceImport_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "vpc_channel_id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "imported_apis.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "imported_apis.0.name"),
				),
			},
		},
	})
}

func testAccMicroserviceImport_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_apig_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = data.huaweicloud_availability_zones.test.names
  enterprise_project_id = "0"

  name = "%[2]s"
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, common.TestVpc(name), name)
}

func testAccApigMicroserviceImport_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_microservice_import" "test" {
  instance_id     = huaweicloud_apig_instance.test.id
  service_type    = "CSE"
  protocol        = "HTTPS"
  backend_timeout = 5000
  auth_type       = "NONE"
  cors            = false
  
  group_info {
    group_id = huaweicloud_apig_group.test.id
  }
  
  apis {
    name       = "test_api"
    req_method = "GET"
    req_uri    = "/test"
    match_mode = "SWA"
  }
  
  cse_info {
    engine_id  = "test_engine_id"
    service_id = "test_service_id"
    version    = "1.0.0"
  }
}
`, testAccMicroserviceImport_base(name))
}
