package aom

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v1/uniagent-console/agent-list/all
func DataSourceUniAgents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUniAgentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the UniAgents are located.`,
			},

			// Optional parameters.
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The page number for pagination. Default value is 1.`,
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of items per page. Default value is 20, maximum is 100.`,
			},
			"ecs_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of ECS IDs. Maximum 100 items are supported.`,
			},
			"agent_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of agent IDs. Maximum 100 items are supported.`,
			},
			"coc_cmdb_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of CMDB IDs. Maximum 100 items are supported.`,
			},

			// Attributes.
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of UniAgents that matched the filter parameters.`,
			},
			"uniagents": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        uniagentSchema(),
				Description: `The list of UniAgents that matched the filter parameters.`,
			},
		},
	}
}

func uniagentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"inner_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The inner IP address of the host.`,
			},
			"import_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The import IP address of the host.`,
			},
			"agent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the agent.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host name.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operating system type.`,
			},
			"agent_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The state of the UniAgent.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID to which the host belongs.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the UniAgent.`,
			},
			"is_hw_cloud_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether the host is a Huawei Cloud host.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID.`,
			},
			"cmdb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CMDB ID.`,
			},
			"ecs_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ECS ID, which is a unique value.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID to which the host belongs.`,
			},
		},
	}
}

func buildUniAgentsRequestBody(d *schema.ResourceData, page, pageSize int) map[string]interface{} {
	body := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	if v, ok := d.GetOk("ecs_id_list"); ok {
		body["ecs_id_list"] = utils.ValueIgnoreEmpty(v)
	}
	if v, ok := d.GetOk("agent_id_list"); ok {
		body["agent_id_list"] = utils.ValueIgnoreEmpty(v)
	}
	if v, ok := d.GetOk("coc_cmdb_id_list"); ok {
		body["coc_cmdb_id_list"] = utils.ValueIgnoreEmpty(v)
	}

	return body
}

func listUniAgents(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, int, error) {
	var (
		httpUrl     = "v1/uniagent-console/agent-list/all"
		pageSize    = 20
		currentPage = 1
		result      = make([]interface{}, 0)
		totalCount  = 0
	)

	// Get page_size from schema if provided
	if v, ok := d.GetOk("page_size"); ok {
		pageSize = v.(int)
		if pageSize > 100 {
			pageSize = 100
		}
	}

	// Get initial page from schema if provided
	if v, ok := d.GetOk("page"); ok {
		currentPage = v.(int)
	}

	listPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// If page is specified, only query that page
	if _, ok := d.GetOk("page"); ok {
		body := buildUniAgentsRequestBody(d, currentPage, pageSize)
		opt.JSONBody = body

		requestResp, err := client.Request("POST", listPath, &opt)
		if err != nil {
			return nil, 0, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, 0, err
		}

		agents := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		result = append(result, agents...)

		return result, totalCount, nil
	}

	// Otherwise, query all pages
	for {
		body := buildUniAgentsRequestBody(d, currentPage, pageSize)
		opt.JSONBody = body

		requestResp, err := client.Request("POST", listPath, &opt)
		if err != nil {
			return nil, 0, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, 0, err
		}

		agents := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if totalCount == 0 {
			totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		}

		result = append(result, agents...)

		// If current page has fewer items than page_size, we've reached the last page
		if len(agents) < pageSize {
			break
		}

		// Calculate total pages
		totalPages := (totalCount + pageSize - 1) / pageSize
		if currentPage >= totalPages {
			break
		}

		currentPage++
	}

	return result, totalCount, nil
}

func flattenUniAgents(agents []interface{}) []map[string]interface{} {
	if len(agents) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(agents))
	for _, agent := range agents {
		result = append(result, map[string]interface{}{
			"inner_ip":         utils.PathSearch("inner_ip", agent, nil),
			"import_ip":        utils.PathSearch("import_ip", agent, nil),
			"agent_id":         utils.PathSearch("agent_id", agent, nil),
			"host_name":        utils.PathSearch("host_name", agent, nil),
			"os_type":          utils.PathSearch("os_type", agent, nil),
			"agent_state":      utils.PathSearch("agent_state", agent, nil),
			"project_id":       utils.PathSearch("project_id", agent, nil),
			"version":          utils.PathSearch("version", agent, nil),
			"is_hw_cloud_host": utils.PathSearch("is_hw_cloud_host", agent, nil),
			"vpc_id":           utils.PathSearch("vpc_id", agent, nil),
			"cmdb_id":          utils.PathSearch("cmdb_id", agent, nil),
			"ecs_id":           utils.PathSearch("ecs_id", agent, nil),
			"domain_id":        utils.PathSearch("domain_id", agent, nil),
		})
	}

	return result
}

func dataSourceUniAgentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	agents, totalCount, err := listUniAgents(client, d)
	if err != nil {
		return diag.Errorf("error querying AOM UniAgents: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("uniagents", flattenUniAgents(agents)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
