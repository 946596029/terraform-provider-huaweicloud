package vpn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var p2cGatewayNonUpdatableParams = []string{
	"vpc_id",
	"connect_subnet",
	"eip.*.id",
	"eip.*.type",
	"eip.*.charge_mode",
	"eip.*.bandwidth_id",
	"eip.*.bandwidth_size",
	"eip.*.bandwidth_name",
	"flavor",
	"availability_zone_ids",
	"max_connection_number",
	"enterprise_project_id",
	"tags",
}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}
// @API VPN PUT /v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}
// @API VPN DELETE /v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}
func ResourceP2CGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceP2CGatewayCreate,
		ReadContext:   resourceP2CGatewayRead,
		UpdateContext: resourceP2CGatewayUpdate,
		DeleteContext: resourceP2CGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(p2cGatewayNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the P2C VPN gateway is located.`,
			},

			// Required parameters.
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC used by the P2C VPN gateway.`,
			},
			"connect_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC subnet used by the P2C VPN gateway.`,
			},
			"eip": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        p2cGatewayEipSchema(),
				Description: `The EIP used by the P2C VPN gateway.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the P2C VPN gateway.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The flavor of the P2C VPN gateway.`,
			},
			"availability_zone_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of availability zone IDs.`,
			},
			"max_connection_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of simultaneously online client connections.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    20,
				Elem:        p2cGatewayTagSchema(),
				Description: `The tags of the P2C VPN gateway.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the P2C VPN gateway.`,
			},
			"current_connection_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The current number of client connections.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order ID.`,
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the P2C VPN gateway is frozen.`,
			},
			"frozen_effect": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The frozen effect.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the P2C VPN gateway.`,
			},
			"upgrade_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The upgrade information of the P2C VPN gateway.`,
			},
			"applied_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The effective time, in RFC3339 format.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					"Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.",
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func p2cGatewayEipSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the existing EIP.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the EIP.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The charge mode of the bandwidth.`,
			},
			"bandwidth_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the shared bandwidth.`,
			},
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The bandwidth size in Mbit/s.`,
			},
			"bandwidth_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The bandwidth name.`,
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The EIP version.`,
			},
			"ip_billing_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order information of the EIP.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IPv4 address of the EIP.`,
			},
			"bandwidth_billing_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order information of the bandwidth.`,
			},
			"share_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bandwidth share type.`,
			},
		},
	}
}

func p2cGatewayTagSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The key of the tag.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value of the tag.`,
			},
		},
	}
}

func buildP2CGatewayEipParams(eip []interface{}) map[string]interface{} {
	if len(eip) == 0 {
		return nil
	}

	return map[string]interface{}{
		"id":             utils.ValueIgnoreEmpty(utils.PathSearch("id", eip[0], nil)),
		"type":           utils.ValueIgnoreEmpty(utils.PathSearch("type", eip[0], nil)),
		"charge_mode":    utils.ValueIgnoreEmpty(utils.PathSearch("charge_mode", eip[0], nil)),
		"bandwidth_id":   utils.ValueIgnoreEmpty(utils.PathSearch("bandwidth_id", eip[0], nil)),
		"bandwidth_size": utils.ValueIgnoreEmpty(utils.PathSearch("bandwidth_size", eip[0], nil)),
		"bandwidth_name": utils.ValueIgnoreEmpty(utils.PathSearch("bandwidth_name", eip[0], nil)),
	}
}

func buildP2CGatewayTagsParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, ""),
			"value": utils.PathSearch("value", tag, ""),
		})
	}
	return result
}

func buildP2CGatewayCreateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"p2c_vpn_gateway": map[string]interface{}{
			"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
			"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
			"connect_subnet":        utils.ValueIgnoreEmpty(d.Get("connect_subnet")),
			"flavor":                utils.ValueIgnoreEmpty(d.Get("flavor")),
			"availability_zone_ids": utils.ValueIgnoreEmpty(d.Get("availability_zone_ids").([]interface{})),
			"eip":                   buildP2CGatewayEipParams(d.Get("eip").([]interface{})),
			"max_connection_number": utils.ValueIgnoreEmpty(d.Get("max_connection_number")),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"tags":                  buildP2CGatewayTagsParams(d.Get("tags").([]interface{})),
		},
	}
}

func GetP2CGatewayById(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/p2c-vpn-gateways/{id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func waitForP2CGatewayState(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration, targetState string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			var (
				cfg    = meta.(*config.Config)
				region = cfg.GetRegion(d)
			)

			client, err := cfg.NewServiceClient("vpn", region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			respBody, err := GetP2CGatewayById(client, d.Id())
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					if targetState == "DELETED" {
						return "Resource Not Found", "COMPLETED", nil
					}
					return nil, "ERROR", fmt.Errorf("the P2C VPN gateway (%s) has been deleted unexpectedly", d.Id())
				}
				return nil, "ERROR", err
			}

			status := utils.PathSearch("p2c_vpn_gateway.status", respBody, "").(string)
			if utils.StrSliceContains([]string{targetState}, status) {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains([]string{
				"PENDING_CREATE", "PENDING_UPDATE", "PENDING_DELETE", "UPGRADING",
				"ROLLING_BACK", "PENDING_UPGRADE_CONFIRM",
			}, status) {
				return respBody, "PENDING", nil
			}

			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func createP2CGateway(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	createPath := client.Endpoint + "v5/{project_id}/p2c-vpn-gateways"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildP2CGatewayCreateBodyParams(d, cfg)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceP2CGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	respBody, err := createP2CGateway(client, d, cfg)
	if err != nil {
		return diag.Errorf("error creating P2C VPN gateway: %s", err)
	}

	id := utils.PathSearch("p2c_vpn_gateway.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating P2C VPN gateway: ID is not found in API response")
	}
	d.SetId(id)

	err = waitForP2CGatewayState(ctx, d, meta, d.Timeout(schema.TimeoutCreate), "ACTIVE")
	if err != nil {
		return diag.Errorf("error waiting for the status of P2C VPN gateway (%s) to become ACTIVE: %s", d.Id(), err)
	}
	return resourceP2CGatewayRead(ctx, d, meta)
}

func flattenP2CGatewayEip(eip map[string]interface{}) []map[string]interface{} {
	if len(eip) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":                     utils.PathSearch("id", eip, nil),
			"type":                   utils.PathSearch("type", eip, nil),
			"charge_mode":            utils.PathSearch("charge_mode", eip, nil),
			"bandwidth_id":           utils.PathSearch("bandwidth_id", eip, nil),
			"bandwidth_size":         utils.PathSearch("bandwidth_size", eip, nil),
			"bandwidth_name":         utils.PathSearch("bandwidth_name", eip, nil),
			"ip_version":             utils.PathSearch("ip_version", eip, nil),
			"ip_billing_info":        utils.PathSearch("ip_billing_info", eip, nil),
			"ip_address":             utils.PathSearch("ip_address", eip, nil),
			"bandwidth_billing_info": utils.PathSearch("bandwidth_billing_info", eip, nil),
			"share_type":             utils.PathSearch("share_type", eip, nil),
		},
	}
}

func flattenP2CGatewayTags(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}
	return result
}

func resourceP2CGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	respBody, err := GetP2CGatewayById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving P2C VPN gateway")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", utils.PathSearch("p2c_vpn_gateway.vpc_id", respBody, nil)),
		d.Set("connect_subnet", utils.PathSearch("p2c_vpn_gateway.connect_subnet", respBody, nil)),
		d.Set("eip", flattenP2CGatewayEip(
			utils.PathSearch("p2c_vpn_gateway.eip", respBody, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("name", utils.PathSearch("p2c_vpn_gateway.name", respBody, nil)),
		d.Set("flavor", utils.PathSearch("p2c_vpn_gateway.flavor", respBody, nil)),
		d.Set("availability_zone_ids", utils.PathSearch("p2c_vpn_gateway.availability_zone_ids", respBody, nil)),
		d.Set("max_connection_number", utils.PathSearch("p2c_vpn_gateway.max_connection_number", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("p2c_vpn_gateway.enterprise_project_id", respBody, nil)),
		d.Set("tags", flattenP2CGatewayTags(utils.PathSearch("p2c_vpn_gateway.tags", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("p2c_vpn_gateway.status", respBody, nil)),
		d.Set("current_connection_number", utils.PathSearch("p2c_vpn_gateway.current_connection_number", respBody, nil)),
		d.Set("order_id", utils.PathSearch("p2c_vpn_gateway.order_id", respBody, nil)),
		d.Set("admin_state_up", utils.PathSearch("p2c_vpn_gateway.admin_state_up", respBody, nil)),
		d.Set("frozen_effect", utils.PathSearch("p2c_vpn_gateway.frozen_effect", respBody, nil)),
		d.Set("version", utils.PathSearch("p2c_vpn_gateway.version", respBody, nil)),
		d.Set("upgrade_info", utils.PathSearch("p2c_vpn_gateway.upgrade_info", respBody, nil)),
		d.Set("created_at", utils.PathSearch("p2c_vpn_gateway.created_at", respBody, nil)),
		d.Set("applied_at", utils.PathSearch("p2c_vpn_gateway.applied_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("p2c_vpn_gateway.updated_at", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildP2CGatewayUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	eips := d.Get("eip").([]interface{})
	return map[string]interface{}{
		"p2c_vpn_gateway": map[string]interface{}{
			"name":   utils.ValueIgnoreEmpty(d.Get("name")),
			"eip_id": utils.ValueIgnoreEmpty(utils.PathSearch("id", eips[0], nil)),
		},
	}
}

func updateP2CGateway(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updatePath := client.Endpoint + "v5/{project_id}/p2c-vpn-gateways/{id}"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildP2CGatewayUpdateBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}
	return nil
}

func resourceP2CGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	if d.HasChanges("name", "eip") {
		err := updateP2CGateway(client, d)
		if err != nil {
			return diag.Errorf("error updating P2C VPN gateway: %s", err)
		}

		err = waitForP2CGatewayState(ctx, d, meta, d.Timeout(schema.TimeoutUpdate), "ACTIVE")
		if err != nil {
			return diag.Errorf("error waiting for the status of P2C VPN gateway (%s) to become ACTIVE: %s", d.Id(), err)
		}
	}

	return resourceP2CGatewayRead(ctx, d, meta)
}

func deleteP2CGateway(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	deletePath := client.Endpoint + "v5/{project_id}/p2c-vpn-gateways/{id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return err
	}
	return nil
}

func resourceP2CGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	err = deleteP2CGateway(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting P2C VPN gateway")
	}

	err = waitForP2CGatewayState(ctx, d, meta, d.Timeout(schema.TimeoutDelete), "DELETED")
	if err != nil {
		return diag.Errorf("error waiting for P2C VPN gateway (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}
