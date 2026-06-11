package dli

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v2/batches"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSparkJobResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI v2 client: %s", err)
	}
	return batches.Get(c, state.Primary.ID)
}

// waitForElasticResourcePoolDeletionCooldown waits for 15 minutes after elastic resource pool creation
// before deletion can be performed (DLI service requirement)
func waitForElasticResourcePoolDeletionCooldown() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// After elastic resource pool is created, it cannot be deleted within 15 minutes.
		// lintignore:R018
		time.Sleep(15 * time.Minute)
		return nil
	}
}

func TestAccDliSparkJobV2_basic(t *testing.T) {
	var job batches.CreateResp

	rName := acceptance.RandomAccResourceName()
	dashName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dli_spark_job.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&job,
		getSparkJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSparkJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliSparkJob_basic(rName, dashName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "queue_name"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
					resource.TestMatchResourceAttr(resourceName, "created_at", 
						regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})`)),
				),
			},
			{
				// Wait for 15 minutes before deletion to satisfy DLI elastic resource pool deletion cooldown requirement
				Config: testAccDliSparkJob_basic(rName, dashName),
				Check: resource.ComposeTestCheckFunc(
					waitForElasticResourcePoolDeletionCooldown(),
				),
			},
		},
	})
}

func testAccCheckDliSparkJobDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.DliV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Dli v2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dli_spark_job" {
			continue
		}

		resp, err := batches.GetState(client, rs.Primary.ID)
		// If the status of the spark job is "dead" or "success", it means that the life cycle of the job has ended.
		if err == nil && resp != nil && (resp.State != batches.StateDead && resp.State != batches.StateSuccess) {
			return fmt.Errorf("spark job (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccDliSparkJob_base(name, dashName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name   = "%[2]s_pool"
  min_cu = 64
  max_cu = 128
}

resource "huaweicloud_dli_queue" "test" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1
  name                       = "%[2]s_queue"
  cu_count                   = 16
  queue_type                 = "general"
}

resource "huaweicloud_dli_package" "test" {
  depends_on  = [huaweicloud_obs_bucket_object.test]
  group_name  = "%[3]s"
  type        = "pyFile"
  object_path = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/dli/packages/simple_pyspark_test_DLF_refresh.py"
}
`, testAccDliPackage_base(dashName), name, dashName)
}

func testAccDliSparkJob_basic(name, dashName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_spark_job" "test" {
  queue_name = huaweicloud_dli_queue.test.name
  name       = "%[2]s"
  app_name   = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"

  depends_on = [
    huaweicloud_obs_bucket.test,
    huaweicloud_obs_bucket_object.test,
    huaweicloud_dli_queue.test,
  ]
}
`, testAccDliSparkJob_base(name, dashName), name)
}