package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDesktopPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopPoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the desktop pools.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the desktop pool.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DYNAMIC", "STATIC"}, false),
				Description:  "The type of the desktop pool.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The enterprise project ID.",
			},
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the desktop pool is in maintenance mode.",
			},
			"pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolSchema(),
				Description: "The list of desktop pools.",
			},
		},
	}
}

func desktopPoolSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the desktop pool.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the desktop pool.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the desktop pool.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the desktop pool.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the desktop pool.",
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The charging mode of the desktop pool.",
			},
			"desktop_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of desktops in the pool.",
			},
			"desktop_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of used desktops in the pool.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability zone of the desktop pool.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subnet ID of the desktop pool.",
			},
			"product": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        productSchema(),
				Description: "The product information of the desktop pool.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image ID used by the desktop pool.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image name used by the desktop pool.",
			},
			"image_os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS type of the image.",
			},
			"image_os_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS version of the image.",
			},
			"image_os_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS platform of the image.",
			},
			"image_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product code of the image.",
			},
			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopVolumeSchema(),
				Description: "The root volume information of the desktop pool.",
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopVolumeSchema(),
				Description: "The data volumes information of the desktop pool.",
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the security group.",
						},
					},
				},
				Description: "The security groups of the desktop pool.",
			},
			"disconnected_retention_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The disconnected retention period in minutes.",
			},
			"enable_autoscale": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether auto scaling is enabled.",
			},
			"autoscale_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        autoscalePolicySchema(),
				Description: "The auto scaling policy of the desktop pool.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the desktop pool.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise project ID.",
			},
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the desktop pool is in maintenance mode.",
			},
			"desktop_name_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The desktop name policy ID.",
			},
		},
	}
}

func productSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the product.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor ID of the product.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the product.",
			},
			"cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CPU specification.",
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The memory specification.",
			},
			"descriptions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the product.",
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The charging mode of the product.",
			},
		},
	}
}

func desktopVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the volume.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the volume.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the volume in GB.",
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource specification code of the volume.",
			},
		},
	}
}

func autoscalePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"autoscale_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of auto scaling.",
			},
			"max_auto_created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of auto-created desktops.",
			},
			"min_idle": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum number of idle desktops.",
			},
			"once_auto_created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of desktops to create in one auto scaling operation.",
			},
		},
	}
}

func buildListDesktopPoolsParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("name"); ok {
		params = fmt.Sprintf("%s&name=%v", params, v)
	}
	if v, ok := d.GetOk("type"); ok {
		params = fmt.Sprintf("%s&type=%v", params, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		params = fmt.Sprintf("%s&enterprise_project_id=%v", params, v)
	}
	if v, ok := d.GetOk("in_maintenance_mode"); ok {
		params = fmt.Sprintf("%s&in_maintenance_mode=%v", params, v)
	}
	return params
}

func listDesktopPools(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktop-pools"
		offset  = 0
		limit   = 200
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s?limit=%d&offset=%d", listPath, limit, offset)
		listPathWithOffset += buildListDesktopPoolsParams(d)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		desktopPools := utils.PathSearch("desktop_pools", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, desktopPools...)
		if len(desktopPools) < limit {
			break
		}
		offset += len(desktopPools)
	}

	return result, nil
}

func dataSourceDesktopPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	desktopPools, err := listDesktopPools(client, d)
	if err != nil {
		return diag.Errorf("error querying desktop pools: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("pools", flattenDesktopPools(desktopPools)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProductInfo(productInfo interface{}) []interface{} {
	if productInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"product_id":   utils.PathSearch("product_id", productInfo, nil),
			"flavor_id":    utils.PathSearch("flavor_id", productInfo, nil),
			"type":         utils.PathSearch("type", productInfo, nil),
			"cpu":          utils.PathSearch("cpu", productInfo, nil),
			"memory":       utils.PathSearch("memory", productInfo, nil),
			"descriptions": utils.PathSearch("descriptions", productInfo, nil),
			"charge_mode":  utils.PathSearch("charge_mode", productInfo, nil),
		},
	}
}

func flattenVolumeInfo(volumeInfo interface{}) map[string]interface{} {
	if volumeInfo == nil {
		return nil
	}

	return map[string]interface{}{
		"id":                 utils.PathSearch("id", volumeInfo, nil),
		"type":               utils.PathSearch("type", volumeInfo, nil),
		"size":               utils.PathSearch("size", volumeInfo, nil),
		"resource_spec_code": utils.PathSearch("resource_spec_code", volumeInfo, nil),
	}
}

func flattenDesktopPoolsRootVolume(rootVolume interface{}) []map[string]interface{} {
	if rootVolume == nil {
		return nil
	}

	return []map[string]interface{}{flattenVolumeInfo(rootVolume)}
}

func flattenDesktopPoolsDataVolumes(dataVolumes interface{}) []map[string]interface{} {
	dataVolumeArray := dataVolumes.([]interface{})
	if len(dataVolumeArray) == 0 {
		return nil
	}

	volumes := make([]map[string]interface{}, 0)
	for _, volume := range dataVolumeArray {
		volumes = append(volumes, flattenVolumeInfo(volume))
	}

	return volumes
}

func flattenSecurityGroups(securityGroups interface{}) []interface{} {
	if securityGroups == nil {
		return nil
	}

	groups := make([]interface{}, 0)
	for _, group := range securityGroups.([]interface{}) {
		groups = append(groups, map[string]interface{}{
			"id": utils.PathSearch("id", group, nil),
		})
	}

	return groups
}

func flattenAutoscalePolicy(autoscalePolicy interface{}) []interface{} {
	if autoscalePolicy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"autoscale_type":    utils.PathSearch("autoscale_type", autoscalePolicy, nil),
			"max_auto_created":  utils.PathSearch("max_auto_created", autoscalePolicy, nil),
			"min_idle":          utils.PathSearch("min_idle", autoscalePolicy, nil),
			"once_auto_created": utils.PathSearch("once_auto_created", autoscalePolicy, nil),
		},
	}
}

func flattenDesktopPools(pools []interface{}) []interface{} {
	if len(pools) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(pools))
	for _, pool := range pools {
		result = append(result, map[string]interface{}{
			"id":                            utils.PathSearch("id", pool, nil),
			"name":                          utils.PathSearch("name", pool, nil),
			"type":                          utils.PathSearch("type", pool, nil),
			"description":                   utils.PathSearch("description", pool, nil),
			"created_time":                  utils.PathSearch("created_time", pool, nil),
			"charging_mode":                 utils.PathSearch("charging_mode", pool, nil),
			"desktop_count":                 utils.PathSearch("desktop_count", pool, nil),
			"desktop_used":                  utils.PathSearch("desktop_used", pool, nil),
			"availability_zone":             utils.PathSearch("availability_zone", pool, nil),
			"subnet_id":                     utils.PathSearch("subnet_id", pool, nil),
			"product":                       flattenProductInfo(utils.PathSearch("product", pool, nil)),
			"image_id":                      utils.PathSearch("image_id", pool, nil),
			"image_name":                    utils.PathSearch("image_name", pool, nil),
			"image_os_type":                 utils.PathSearch("image_os_type", pool, nil),
			"image_os_version":              utils.PathSearch("image_os_version", pool, nil),
			"image_os_platform":             utils.PathSearch("image_os_platform", pool, nil),
			"image_product_code":            utils.PathSearch("image_product_code", pool, nil),
			"root_volume":                   flattenDesktopPoolsRootVolume(utils.PathSearch("root_volume", pool, nil)),
			"data_volumes":                  flattenDesktopPoolsDataVolumes(utils.PathSearch("data_volumes", pool, nil)),
			"security_groups":               flattenSecurityGroups(utils.PathSearch("security_groups", pool, nil)),
			"disconnected_retention_period": utils.PathSearch("disconnected_retention_period", pool, nil),
			"enable_autoscale":              utils.PathSearch("enable_autoscale", pool, nil),
			"autoscale_policy":              flattenAutoscalePolicy(utils.PathSearch("autoscale_policy", pool, nil)),
			"status":                        utils.PathSearch("status", pool, nil),
			"enterprise_project_id":         utils.PathSearch("enterprise_project_id", pool, nil),
			"in_maintenance_mode":           utils.PathSearch("in_maintenance_mode", pool, nil),
			"desktop_name_policy_id":        utils.PathSearch("desktop_name_policy_id", pool, nil),
		})
	}
	return result
}
