package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
)

func getP2CGatewayResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("vpn", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN client: %s", err)
	}

	return vpn.GetP2CGatewayById(client, state.Primary.ID)
}

func TestAccVpnP2CGateway_basic(t *testing.T) {
	var (
		rName = "huaweicloud_vpn_p2c_gateway.test"
		name  = acceptance.RandomAccResourceName()

		obj interface{}
		rc  = acceptance.InitResourceCheck(
			rName,
			&obj,
			getP2CGatewayResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpnP2CGateway_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "flavor", "Professional1"),
					resource.TestCheckResourceAttr(rName, "max_connection_number", "100"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "eip.0.type", "5_bgp"),
					resource.TestCheckResourceAttr(rName, "eip.0.bandwidth_size", "10"),
					resource.TestCheckResourceAttr(rName, "eip.0.charge_mode", "traffic"),
					resource.TestCheckResourceAttr(rName, "tags.0.key", "foo"),
					resource.TestCheckResourceAttr(rName, "tags.0.value", "bar"),
				),
			},
			{
				Config: testAccVpnP2CGateway_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"vpc_id",
					"connect_subnet",
					"eip",
					"flavor",
					"availability_zone_ids",
					"max_connection_number",
					"enterprise_project_id",
					"tags",
				},
			},
		},
	})
}

func testAccVpnP2CGateway_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpn_p2c_gateway_availability_zones" "test" {
  flavor = "Professional1"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}
`, name)
}

func testAccVpnP2CGateway_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_p2c_gateway" "test" {
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  connect_subnet        = huaweicloud_vpc_subnet.test.id
  flavor                = "Professional1"
  max_connection_number = 100
  availability_zone_ids = [
    data.huaweicloud_vpn_p2c_gateway_availability_zones.test.availability_zones[0],
    data.huaweicloud_vpn_p2c_gateway_availability_zones.test.availability_zones[1]
  ]

  eip {
    type           = "5_bgp"
    bandwidth_size = 10
    charge_mode    = "traffic"
    bandwidth_name = "%[2]s-bw"
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}
`, testAccVpnP2CGateway_base(name), name)
}

func testAccVpnP2CGateway_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_p2c_gateway" "test" {
  name                  = "%[2]s-update"
  vpc_id                = huaweicloud_vpc.test.id
  connect_subnet        = huaweicloud_vpc_subnet.test.id
  flavor                = "Professional1"
  max_connection_number = 100
  availability_zone_ids = [
    data.huaweicloud_vpn_p2c_gateway_availability_zones.test.availability_zones[0],
    data.huaweicloud_vpn_p2c_gateway_availability_zones.test.availability_zones[1]
  ]

  eip {
    type           = "5_bgp"
    bandwidth_size = 10
    charge_mode    = "traffic"
    bandwidth_name = "%[2]s-bw"
  }

  tags {
    key   = "foo"
    value = "bar"
  }
}
`, testAccVpnP2CGateway_base(name), name)
}
