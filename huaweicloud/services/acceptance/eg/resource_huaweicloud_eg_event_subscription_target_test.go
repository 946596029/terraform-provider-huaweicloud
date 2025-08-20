package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eg"
)

func getEventSubscriptionTargetFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("eg", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG client: %s", err)
	}

	return eg.GetEventSubscriptionTargetById(client, state.Primary.Attributes["subscription_id"], state.Primary.ID)
}

func TestAccEventSubscriptionTarget_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		withTest   = "huaweicloud_eg_event_subscription_target.test"
		target     interface{}
		rcWithTest = acceptance.InitResourceCheck(withTest, &target, getEventSubscriptionTargetFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPreCheckEgSubscriptionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithTest.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testEventSubscriptionTarget_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithTest.CheckResourceExists(),
					resource.TestCheckResourceAttr(withTest, "name", name),
					resource.TestCheckResourceAttr(withTest, "provider_type", "CUSTOM"),
					resource.TestCheckResourceAttr(withTest, "retry_times", "3"),
					resource.TestCheckResourceAttrSet(withTest, "created_at"),
					resource.TestCheckResourceAttrSet(withTest, "updated_at"),
				),
			},
			{
				Config: testEventSubscriptionTarget_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithTest.CheckResourceExists(),
					resource.TestCheckResourceAttr(withTest, "name", name+"_update"),
					resource.TestCheckResourceAttr(withTest, "retry_times", "5"),
					resource.TestCheckResourceAttrSet(withTest, "updated_at"),
				),
			},
			{
				ResourceName:      withTest,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testEventSubscriptionTarget_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_event_subscription_target" "test" {
  subscription_id  = "%[2]s"
  name             = "HTTPS"
  provider_type    = "OFFICIAL"
  retry_times      = 16

  key_transform {
    type = "ORIGINAL"
  }

  detail = jsonencode({
    url = "http://example.com/webhook"
  })
}
`, name, acceptance.HW_EG_SUBCRIPTION_ID)
}

func testEventSubscriptionTarget_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_event_subscription_target" "test" {
  subscription_id  = "%[2]s"
  name             = "HTTPS"
  provider_type    = "OFFICIAL"
  retry_times      = 20

  key_transform {
    type = "ORIGINAL"
  }

  detail = jsonencode({
    url = "http://example.com/webhook"
  })
}
`, name, acceptance.HW_EG_SUBCRIPTION_ID)
}
