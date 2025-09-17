package apig

import (
	"context"
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

var apigMicroserviceImportNonUpdatableParams = []string{"instance_id", "group_info", "service_type", "apis",
	"protocol", "backend_timeout", "auth_type", "cors", "cse_info", "cce_info", "nacos_info"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/microservice/import
func ResourceMicroserviceImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceImportCreate,
		ReadContext:   resourceMicroserviceImportRead,
		UpdateContext: resourceMicroserviceImportUpdate,
		DeleteContext: resourceMicroserviceImportDelete,

		CustomizeDiff: config.FlexibleForceNew(apigMicroserviceImportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dedicated instance is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance.`,
			},
			"group_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        microserviceGroupInfoSchema(),
				Description: `The API group information for importing microservice.`,
			},
			"service_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of microservice center.`,
			},
			"apis": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    50,
				Elem:        microserviceApiSchema(),
				Description: `The list of APIs to be imported.`,
			},

			// Optional parameters.
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "HTTPS",
				Description: `The protocol for the dedicated instance to access the microservice.`,
			},
			"backend_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5000,
				Description: `The timeout for the dedicated instance to request backend services.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
				Description: `The authentication type of the API.`,
			},
			"cors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to support CORS.`,
			},
			"cse_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        microserviceCseInfoSchema(),
				Description: `The CSE microservice information.`,
			},
			"cce_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        microserviceCceInfoSchema(),
				Description: `The CCE microservice information.`,
			},
			"nacos_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        microserviceNacosInfoSchema(),
				Description: `The Nacos microservice information.`,
			},

			// Attributes.
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC channel created for the microservice.`,
			},
			"api_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the API group.`,
			},
			"imported_apis": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        microserviceImportedApisSchema(),
				Description: `The list of imported APIs with generated IDs.`,
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

func microserviceGroupInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API group.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the API group.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the application.`,
			},
		},
	}
}

func microserviceApiSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"req_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request URI of the API.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the API.`,
			},
			"req_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The request method of the API.`,
			},
			"match_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The matching mode of the API.`,
			},
		},
	}
}

func microserviceCseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the CSE engine.`,
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the CSE service.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the CSE service.`,
			},
		},
	}
}

func microserviceCceLabelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the label.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value of the label.`,
			},
		},
	}
}

func microserviceCceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cluster ID of the CCE application.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The namespace of the CCE application.`,
			},
			"workload_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the CCE workload.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The port of the CCE application.`,
			},

			// Optional
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the CCE application.`,
			},
			"label_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The label key of the CCE application.`,
			},
			"label_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The label value of the CCE application.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version of the CCE application.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        microserviceCceLabelSchema(),
				Description: `The labels of the CCE workload.`,
			},
		},
	}
}

func microserviceNacosServerConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The IP of the Nacos server.`,
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The port of the Nacos server.`,
			},
			"grpc_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The grpc port of the Nacos server.`,
			},
		},
	}
}

func microserviceNacosUserInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Nacos user.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The password of the Nacos user.`,
			},
		},
	}
}

func microserviceNacosInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of Nacos microservice.`,
			},
			"server_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    9,
				Elem:        microserviceNacosServerConfigSchema(),
				Description: `The server config of Nacos microservice.`,
			},
			"user_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        microserviceNacosUserInfoSchema(),
				Description: `The user information of Nacos microservice.`,
			},

			// Optional
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The namespace of Nacos microservice.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The group name of Nacos microservice.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The cluster name of Nacos microservice.`,
			},
		},
	}
}

func microserviceImportedApisSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the imported API.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the imported API.`,
			},
			"req_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request URI of the imported API.`,
			},
			"req_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request method of the imported API.`,
			},
			"match_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The matching mode of the imported API.`,
			},
		},
	}
}

func buildMicroserviceGroupInfo(groupInfos []interface{}) map[string]interface{} {
	if len(groupInfos) < 1 {
		return nil
	}

	return map[string]interface{}{
		"group_id":   utils.PathSearch("group_id", groupInfos[0], nil),
		"group_name": utils.ValueIgnoreEmpty(utils.PathSearch("group_name", groupInfos[0], nil)),
		"app_id":     utils.ValueIgnoreEmpty(utils.PathSearch("app_id", groupInfos[0], nil)),
	}
}

func buildMicroserviceApis(apis []interface{}) []interface{} {
	if len(apis) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(apis))
	for _, api := range apis {
		result = append(result, map[string]interface{}{
			"req_uri":    utils.PathSearch("req_uri", api, nil),
			"name":       utils.ValueIgnoreEmpty(utils.PathSearch("name", api, nil)),
			"req_method": utils.ValueIgnoreEmpty(utils.PathSearch("req_method", api, nil)),
			"match_mode": utils.ValueIgnoreEmpty(utils.PathSearch("match_mode", api, nil)),
		})
	}

	return result
}

func buildMicroserviceCseInfo(cseInfos []interface{}) map[string]interface{} {
	if len(cseInfos) < 1 {
		return nil
	}

	return map[string]interface{}{
		"engine_id":  utils.PathSearch("engine_id", cseInfos[0], nil),
		"service_id": utils.PathSearch("service_id", cseInfos[0], nil),
		"version":    utils.PathSearch("version", cseInfos[0], nil),
	}
}

func buildMicroserviceCceInfoLabels(cceInfoLabels []interface{}) []interface{} {
	if len(cceInfoLabels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(cceInfoLabels))
	for _, label := range cceInfoLabels {
		result = append(result, map[string]interface{}{
			"label_name":  utils.PathSearch("name", label, nil),
			"label_value": utils.PathSearch("value", label, nil),
		})
	}

	return result
}

func buildMicroserviceCceInfo(cceInfos []interface{}) map[string]interface{} {
	if len(cceInfos) < 1 {
		return nil
	}

	return map[string]interface{}{
		"cluster_id":    utils.PathSearch("cluster_id", cceInfos[0], nil),
		"namespace":     utils.PathSearch("namespace", cceInfos[0], nil),
		"workload_type": utils.PathSearch("workload_type", cceInfos[0], nil),
		"port":          utils.PathSearch("port", cceInfos[0], nil),
		"app_name":      utils.ValueIgnoreEmpty(utils.PathSearch("app_name", cceInfos[0], nil)),
		"label_key":     utils.ValueIgnoreEmpty(utils.PathSearch("label_key", cceInfos[0], nil)),
		"label_value":   utils.ValueIgnoreEmpty(utils.PathSearch("label_value", cceInfos[0], nil)),
		"version":       utils.ValueIgnoreEmpty(utils.PathSearch("version", cceInfos[0], nil)),
		"labels": utils.ValueIgnoreEmpty(buildMicroserviceCceInfoLabels(
			utils.PathSearch("labels", cceInfos[0], make([]interface{}, 0)).([]interface{}))),
	}
}

func buildMicroserviceNacosServerConfig(serverConfigs []interface{}) map[string]interface{} {
	if len(serverConfigs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"ip_address": utils.PathSearch("ip_address", serverConfigs[0], nil),
		"port":       utils.PathSearch("port", serverConfigs[0], nil),
		"grpc_port":  utils.ValueIgnoreEmpty(utils.PathSearch("grpc_port", serverConfigs[0], nil)),
	}
}

func buildMicroserviceNacosUserInfo(userInfos []interface{}) map[string]interface{} {
	if len(userInfos) < 1 {
		return nil
	}

	return map[string]interface{}{
		"user_name": utils.PathSearch("user_name", userInfos[0], nil),
		"password":  utils.PathSearch("password", userInfos[0], nil),
	}
}

func buildMicroserviceNacosInfo(nacosInfos []interface{}) map[string]interface{} {
	if len(nacosInfos) < 1 {
		return nil
	}

	return map[string]interface{}{
		"service_name": utils.PathSearch("service_name", nacosInfos[0], nil),
		"server_config": buildMicroserviceNacosServerConfig(
			utils.PathSearch("server_config", nacosInfos[0], make([]interface{}, 0)).([]interface{})),
		"user_info": buildMicroserviceNacosUserInfo(
			utils.PathSearch("user_info", nacosInfos[0], make([]interface{}, 0)).([]interface{})),
		"namespace":    utils.ValueIgnoreEmpty(utils.PathSearch("namespace", nacosInfos[0], nil)),
		"group_name":   utils.ValueIgnoreEmpty(utils.PathSearch("group_name", nacosInfos[0], nil)),
		"cluster_name": utils.ValueIgnoreEmpty(utils.PathSearch("cluster_name", nacosInfos[0], nil)),
	}
}

func flattenMicroserviceImportedApis(apis []interface{}) []map[string]interface{} {
	if len(apis) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(apis))
	for _, api := range apis {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", api, nil),
			"name":       utils.PathSearch("name", api, nil),
			"req_uri":    utils.PathSearch("req_uri", api, nil),
			"req_method": utils.PathSearch("req_method", api, nil),
			"match_mode": utils.PathSearch("match_mode", api, nil),
		})
	}

	return result
}

func buildMicroserviceImportBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"group_info":      buildMicroserviceGroupInfo(d.Get("group_info").([]interface{})),
		"service_type":    d.Get("service_type"),
		"apis":            buildMicroserviceApis(d.Get("apis").([]interface{})),
		"protocol":        utils.ValueIgnoreEmpty(d.Get("protocol")),
		"backend_timeout": utils.ValueIgnoreEmpty(d.Get("backend_timeout")),
		"auth_type":       utils.ValueIgnoreEmpty(d.Get("auth_type")),
		"cors":            utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "cors"),
		"cse_info":        utils.ValueIgnoreEmpty(buildMicroserviceCseInfo(d.Get("cse_info").([]interface{}))),
		"cce_info":        utils.ValueIgnoreEmpty(buildMicroserviceCceInfo(d.Get("cce_info").([]interface{}))),
		"nacos_info":      utils.ValueIgnoreEmpty(buildMicroserviceNacosInfo(d.Get("nacos_info").([]interface{}))),
	}
}

func resourceMicroserviceImportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/microservice/import"
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildMicroserviceImportBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error importing microservice: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("vpc_channel_id", utils.PathSearch("vpc_channel_id", respBody, nil)),
		d.Set("api_group_id", utils.PathSearch("api_group_id", respBody, nil)),
		d.Set("imported_apis", flattenMicroserviceImportedApis(
			utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{}))),
	)
	err = mErr.ErrorOrNil()
	if err != nil {
		return diag.Errorf("error setting microservice import attributes: %s", err)
	}

	return resourceMicroserviceImportRead(ctx, d, meta)
}

func resourceMicroserviceImportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMicroserviceImportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMicroserviceImportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for importing microservice. Deleting this resource 
will not clear the corresponding imported microservice, but will only remove the resource information from the tfstate 
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
