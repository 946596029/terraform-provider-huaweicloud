// Generated by PMS #687
package dcs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceDcsInstanceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsInstanceNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the node information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logical_node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the logical node ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node status.`,
						},
						"az_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the AZ code.`,
						},
						"node_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node role.`,
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates  the node master/standby role.`,
						},
						"node_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node IP.`,
						},
						"node_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node port.`,
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node ID.`,
						},
						"priority_weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates  the replica promotion priority.`,
						},
						"is_access": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the IP address of the node can be directly accessed.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance shard ID.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance shard name.`,
						},
						"is_remove_ip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the IP address is removed from the read-only domain name.`,
						},
						"replication_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance replica ID.`,
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the monitoring metric dimension of the replica.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the monitoring dimension name.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the dimension value.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type InstanceNodesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newInstanceNodesDSWrapper(d *schema.ResourceData, meta interface{}) *InstanceNodesDSWrapper {
	return &InstanceNodesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDcsInstanceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newInstanceNodesDSWrapper(d, meta)
	shoNodInfRst, err := wrapper.ShowNodesInformation()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.showNodesInformationToSchema(shoNodInfRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API DCS GET /v2/{project_id}/instances/{instance_id}/logical-nodes
func (w *InstanceNodesDSWrapper) ShowNodesInformation() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dcs")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/instances/{instance_id}/logical-nodes"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *InstanceNodesDSWrapper) showNodesInformationToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("nodes", schemas.SliceToList(body.Get("nodes"),
			func(nodes gjson.Result) any {
				return map[string]any{
					"logical_node_id": nodes.Get("logical_node_id").Value(),
					"name":            nodes.Get("name").Value(),
					"status":          nodes.Get("status").Value(),
					"az_code":         nodes.Get("az_code").Value(),
					"node_role":       nodes.Get("node_role").Value(),
					"node_type":       nodes.Get("node_type").Value(),
					"node_ip":         nodes.Get("node_ip").Value(),
					"node_port":       nodes.Get("node_port").Value(),
					"node_id":         nodes.Get("node_id").Value(),
					"priority_weight": nodes.Get("priority_weight").Value(),
					"is_access":       nodes.Get("is_access").Value(),
					"group_id":        nodes.Get("group_id").Value(),
					"group_name":      nodes.Get("group_name").Value(),
					"is_remove_ip":    nodes.Get("is_remove_ip").Value(),
					"replication_id":  nodes.Get("replication_id").Value(),
					"dimensions": schemas.SliceToList(nodes.Get("dimensions"),
						func(dimensions gjson.Result) any {
							return map[string]any{
								"name":  dimensions.Get("name").Value(),
								"value": dimensions.Get("value").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
