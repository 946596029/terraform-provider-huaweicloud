package dws

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var logicalClusterConvertActionNonUpdatableParams = []string{"cluster_id", "logical_cluster_name"}

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/convert-to-logical-cluster/{name}
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceLogicalClusterConvertAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogicalClusterConvertActionCreate,
		ReadContext:   resourceLogicalClusterConvertActionRead,
		UpdateContext: resourceLogicalClusterConvertActionUpdate,
		DeleteContext: resourceLogicalClusterConvertActionDelete,

		CustomizeDiff: config.FlexibleForceNew(logicalClusterConvertActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the DWS cluster is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DWS cluster to convert to logical cluster mode.`,
			},
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name to set for the logical cluster created from the conversion.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceLogicalClusterConvertActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		clusterId   = d.Get("cluster_id").(string)
		logicalName = d.Get("logical_cluster_name").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	httpUrl := "v2/{project_id}/clusters/{cluster_id}/convert-to-logical-cluster/{name}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)
	path = strings.ReplaceAll(path, "{name}", url.PathEscape(logicalName))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return diag.Errorf("error converting DWS cluster (%s) to logical cluster: %s", clusterId, err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := waitClusterTaskStateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), clusterId); err != nil {
		return diag.Errorf("error waiting for DWS cluster (%s) convert task to complete: %s", clusterId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceLogicalClusterConvertActionRead(ctx, d, meta)
}

func resourceLogicalClusterConvertActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLogicalClusterConvertActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLogicalClusterConvertActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for converting a physical DWS cluster to 
a logical cluster. Deleting this resource will not clear the corresponding request record, but will only remove 
the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
