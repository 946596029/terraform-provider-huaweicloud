// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	deleteNotExistMsg            = "logical cluster is not existed"
	deleteFirstLogicalClusterMsg = "the first logical cluster can't be deleted"
	createDuplicateNameMsg       = "logical cluster already existed"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
}

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/logical-clusters
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters
// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters/volumes
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}
func ResourceLogicalCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogicalClusterCreate,
		ReadContext:   resourceLogicalClusterRead,
		UpdateContext: resourceLogicalClusterUpdate,
		DeleteContext: resourceLogicalClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLogicalClusterImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			// fix
			Create: schema.DefaultTimeout(3 * time.Minute),

			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Required
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the logical cluster name.`,
			},
			"cluster_rings": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        logicalClusterRingsSchema(),
				Description: `Specifies the DWS cluster ring list information.`,
			},

			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DWS logical cluster status.`,
			},
			"first_logical_cluster": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is the first logical cluster. The first logical cluster cannot be deleted.`,
			},
			"edit_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether editing is allowed.`,
			},
			"restart_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to allow restart.`,
			},
			"delete_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether deletion is allowed.`,
			},
			"volume_usage": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        logicalClusterVolumeUsageSchema(),
				Description: `The volume usage information of the logical cluster.`,
			},

			// Internal parameters
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func logicalClusterRingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ring_hosts": {
				Type:        schema.TypeSet,
				Elem:        logicalRingHostsSchema(),
				Required:    true,
				Description: `Indicates the cluster host ring information.`,
			},
		},
	}
	return &sc
}

func logicalRingHostsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the host name.`,
			},
			"back_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the backend IP address.`,
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of CPU cores.`,
			},
			"memory": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: `Specifies the host memory.`,
			},
			"disk_size": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: `Specifies the host disk size.`,
			},
		},
	}
	return &sc
}

func logicalClusterVolumeUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk usage of the logical cluster.`,
			},
			"total": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The total disk capacity of the logical cluster.`,
			},
			"percent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk usage percentage of the logical cluster.`,
			},
		},
	}
}

func buildLogicalRingHostsRequestBody(rawParams interface{}) []map[string]interface{} {
	if rawSet, ok := rawParams.(*schema.Set); ok && rawSet.Len() > 0 {
		rst := make([]map[string]interface{}, 0, rawSet.Len())
		for _, v := range rawSet.List() {
			raw, isMap := v.(map[string]interface{})
			if !isMap {
				continue
			}

			rst = append(rst, map[string]interface{}{
				"host_name": raw["host_name"],
				"back_ip":   raw["back_ip"],
				"cpu_cores": raw["cpu_cores"],
				"memory":    raw["memory"],
				"disk_size": raw["disk_size"],
			})
		}
		return rst
	}
	return nil
}

func buildLogicalClusterRingsRequestBody(rawParams interface{}) []map[string]interface{} {
	if rawSet, ok := rawParams.(*schema.Set); ok && rawSet.Len() > 0 {
		rst := make([]map[string]interface{}, 0, rawSet.Len())
		for _, v := range rawSet.List() {
			raw, isMap := v.(map[string]interface{})
			if !isMap {
				continue
			}

			rst = append(rst, map[string]interface{}{
				"ring_hosts": buildLogicalRingHostsRequestBody(raw["ring_hosts"]),
			})
		}
		return rst
	}
	return nil
}

func buildCreateLogicalClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"logical_cluster": map[string]interface{}{
			"logical_cluster_name": d.Get("logical_cluster_name"),
			"cluster_rings":        buildLogicalClusterRingsRequestBody(d.Get("cluster_rings")),
		},
	}
}

// When an error occurs when calling the API, the creation is considered failed and there is no need to retry.
// When the `error_code` is equal to `DWS.0000`, it means the creation is successful.
// When the "error_code" is not equal to "DWS.0000", it means that the creation failed and needs to be retried.
func buildCreateRetryFunc(client *golangsdk.ServiceClient, createPath string, createOpt *golangsdk.RequestOpts) common.RetryFunc {
	retryFunc := func() (interface{}, bool, error) {

		log.Printf("[MY-DEBUG] create logical cluster retry | buildCreateRetryFunc: createPath: %s, createOpt: %#v", createPath, createOpt.JSONBody)

		createResp, err := client.Request("POST", createPath, createOpt)
		if err != nil {
			return nil, false, fmt.Errorf("error creating DWS logical cluster: %s", err)
		}

		createRespBody, err := utils.FlattenResponse(createResp)
		if err != nil {
			return nil, false, err
		}

		errCode := utils.PathSearch("error_code", createRespBody, "").(string)
		log.Printf("[MY-DEBUG] create logical cluster retry | buildCreateRetryFunc: errCode: %s", errCode)
		if errCode == "DWS.0000" {
			return nil, false, nil
		}

		errMsg := utils.PathSearch("error_msg", createRespBody, "").(string)
		// Stop retrying create operations when names are duplicated
		if errMsg == createDuplicateNameMsg {
			return nil, false, fmt.Errorf("error creating DWS logical cluster: %s", errMsg)
		}

		return nil, true, fmt.Errorf("error creating DWS logical cluster: error code: %s, error message: %s", errCode, errMsg)
	}
	return retryFunc
}

func waitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) (interface{}, error) {
	clusterName := d.Get("logical_cluster_name").(string)
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_name=='%s']|[0]", clusterName)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {

			log.Printf("[MY-DEBUG] waiting for state completed | waitingForStateCompleted: clusterName: %s", clusterName)
			log.Printf("[MY-DEBUG] waiting for state completed | waitingForStateCompleted: expression: %s", expression)

			clusterRespBody, err := readLogicalClusters(client, d)
			if err != nil {
				return nil, "ERROR", err
			}


			cluster := utils.PathSearch(expression, clusterRespBody, nil)
			log.Printf("[MY-DEBUG] waiting for state completed | waitingForStateCompleted: cluster: %#v", cluster!=nil)
			if cluster == nil {
				return nil, "PENDING", nil
			}

			completed := utils.PathSearch("action_info.completed", cluster, false).(bool)
			result := utils.PathSearch("action_info.result", cluster, "").(string)
			if completed && result == "success" {
				return cluster, "COMPLETED", nil
			}

			if completed && result == "failed" {
				return cluster, "ERROR", fmt.Errorf("the DWS logical cluster (%s) is failed", clusterName)
			}
			return cluster, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func resourceLogicalClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/clusters/{cluster_id}/logical-clusters"
		product     = "dws"
		clusterName = d.Get("logical_cluster_name").(string)
		clusterId   = d.Get("cluster_id").(string)
	)

	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	log.Printf("[MY-DEBUG] create logical cluster start | createLogicalClusterCreate: clusterId: %s, clusterName: %s", clusterId, clusterName)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         buildCreateLogicalClusterBodyParams(d),
		OkCodes:          []int{200, 417},
	}

	// Multiple logical clusters cannot be created in parallel and need to wait for retry.
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    buildCreateRetryFunc(client, createPath, &createOpt),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 30 * time.Second,
		PollInterval: 30 * time.Second,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for the created logical cluster to be stable and obtain the stable target logical cluster.
	clusterRespBody, err := waitingForStateCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for DWS logical cluster (%s) creation to complete: %s", clusterName, err)
	}

	logicalClusterId := utils.PathSearch("logical_cluster_id", clusterRespBody, "").(string)
	if logicalClusterId == "" {
		return diag.Errorf("unable to find the DWS logical cluster ID from the API response")
	}
	d.SetId(logicalClusterId)

	return resourceLogicalClusterRead(ctx, d, meta)
}

func buildUpdateLogicalClusterBodyParams(clusterRings interface{}) map[string]interface{} {
	return map[string]interface{}{
		"cluster_rings": buildLogicalClusterRingsRequestBody(clusterRings),
		"mode": "insert",
		"waiting_for_killing": 0,
	}
}

func updateLogicalClusterRings(client *golangsdk.ServiceClient, clusterId, logicalClusterId string, clusterRings interface{}) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{logical_cluster_id}", logicalClusterId)

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateLogicalClusterBodyParams(clusterRings)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func queryLogicalClusterById(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	clusterRespBody, err := readLogicalClusters(client, d)
	if err != nil {
		return nil, err
	}

	logicalClusterId := d.Id()
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", logicalClusterId)
	cluster := utils.PathSearch(expression, clusterRespBody, nil)
	if cluster == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return cluster, nil
}

func flattenResponseBodyClusterRings(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("cluster_rings", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"ring_hosts": flattenRingHosts(v),
		}
	}
	return rst
}

func hostInfoToUniqueId(host map[string]interface{}) string {
	return fmt.Sprintf("%v:%v:%v:%.2f:%.2f",
		utils.PathSearch("host_name", host, nil),
		utils.PathSearch("back_ip", host, nil),
		utils.PathSearch("cpu_cores", host, nil),
		utils.PathSearch("memory", host, float64(0)).(float64),
		utils.PathSearch("disk_size", host, float64(0)).(float64),
	)
}

func hostsToUniqueIds(hosts []interface{}) []string {
	hostList := make([]string, 0, len(hosts))
	for _, host := range hosts {
		hostList = append(hostList, hostInfoToUniqueId(host.(map[string]interface{})))
	}
	return hostList
}

func logicalClusterRingsUpdateRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := queryLogicalClusterById(client, d)
		if err != nil {
			return nil, "ERROR", err
		}

		completed := utils.PathSearch("action_info.completed", cluster, false).(bool)
		result := utils.PathSearch("action_info.result", cluster, "").(string)

		if completed && result == "failed" {
			return cluster, "ERROR", errors.New("the logical cluster update failed")
		}

		expectedUniqueIds := hostsToUniqueIds(d.Get("cluster_rings").(*schema.Set).List())
		currentUniqueIds := hostsToUniqueIds(flattenResponseBodyClusterRings(cluster))
		isAllContain := true
		for _, expectedId := range expectedUniqueIds {
			if !utils.StrSliceContains(currentUniqueIds, expectedId) {
				isAllContain = false
				break
			}
		}

		if !completed || !isAllContain {
			return cluster, "PENDING", nil
		}

		return cluster, "COMPLETED", nil
	}
}

func waitingForClusterRingsUpdated(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      logicalClusterRingsUpdateRefreshFunc(client, d),
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceLogicalClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		clusterId        = d.Get("cluster_id").(string)
		logicalClusterId = d.Id()
		clusterRings     = d.Get("cluster_rings")
	)

	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChange("cluster_rings") {
		err = updateLogicalClusterRings(client, clusterId, logicalClusterId, clusterRings)
		if err != nil {
			return diag.Errorf("error updating logical cluster: %s", err)
		}

		err = waitingForClusterRingsUpdated(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for logical cluster (%s) update to complete: %s", logicalClusterId, err)
		}
	}

	return resourceLogicalClusterRead(ctx, d, meta)
}

func listLogicalClusterVolumes(client *golangsdk.ServiceClient, clusterId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/logical-clusters/volumes?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{cluster_id}", clusterId)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))

	listOpts := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		volumes := utils.PathSearch("volumes", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, volumes...)
		if len(volumes) < limit {
			break
		}
		offset += len(volumes)
	}

	return result, nil
}

func flattenVolumeUsage(volume interface{}) []map[string]interface{} {
	if volume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"usage":   utils.PathSearch("usage", volume, nil),
			"total":   utils.PathSearch("total", volume, nil),
			"percent": utils.PathSearch("percent", volume, nil),
		},
	}
}

func resourceLogicalClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr      *multierror.Error
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	clusterRespBody, err := readLogicalClusters(client, d)
	// The list API response status code is `404` when the cluster does not exist (standard UUID format).
	// "DWS.0001": The cluster ID is a non-standard UUID, the status code is 400.
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ClusterIdIllegalErrCode),
			"error retrieving DWS logical cluster")
	}

	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", d.Id())
	cluster := utils.PathSearch(expression, clusterRespBody, nil)
	if cluster == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	volumes, err := listLogicalClusterVolumes(client, clusterId)
	if err != nil {
		log.Printf("[WARN] error retrieving volume usage: %s", err)
	}
	volumeUsage := utils.PathSearch(fmt.Sprintf("[?logical_cluster_name=='%s']|[0]", d.Get("logical_cluster_name").(string)), volumes, nil)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("logical_cluster_name", utils.PathSearch("logical_cluster_name", cluster, nil)),
		d.Set("cluster_rings", flattenResponseBodyClusterRings(cluster)),
		d.Set("status", utils.PathSearch("status", cluster, nil)),
		d.Set("first_logical_cluster", utils.PathSearch("first_logical_cluster", cluster, nil)),
		d.Set("edit_enable", utils.PathSearch("edit_enable", cluster, nil)),
		d.Set("restart_enable", utils.PathSearch("restart_enable", cluster, nil)),
		d.Set("delete_enable", utils.PathSearch("delete_enable", cluster, nil)),
		d.Set("volume_usage", flattenVolumeUsage(volumeUsage)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// waitingForDeleteStateEnable This method is used to wait for operable status before deleting.
// Deleting operations can only be performed when `delete_enable` is true.
func waitingForDeleteStateEnable(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) (interface{}, error) {
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FIRST_LOGICAL_CLUSTER"},
		Refresh: func() (interface{}, string, error) {
			clusterRespBody, err := readLogicalClusters(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			cluster := utils.PathSearch(expression, clusterRespBody, nil)
			if cluster == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			// The last logical cluster cannot be deleted for versions after 820.
			// The only two types of logical cluster names are "elastic_group" and custom names.
			clusters := utils.PathSearch("logical_clusters[?logical_cluster_name!='elastic_group']", clusterRespBody, make([]interface{}, 0))
			if len(clusters.([]interface{})) <= 1 {
				return "last_logical_cluster", "COMPLETED", golangsdk.ErrDefault404{}
			}

			enable := utils.PathSearch("delete_enable", cluster, false).(bool)
			if enable {
				return enable, "COMPLETED", nil
			}

			// When `first_logical_cluster` is true, field `delete_enable` will always be false.
			// The `first_logical_cluster` always false for versions after 820.
			isFirstCluster := utils.PathSearch("first_logical_cluster", cluster, false).(bool)
			if isFirstCluster {
				return enable, "FIRST_LOGICAL_CLUSTER", nil
			}

			return enable, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

// When an error occurs when calling the API, the deletion is deemed to have failed and there is no need to retry.
// When the "error_code" is equal to "DWS.0000", it means the deletion is successful.
// When the "error_code" is not equal to "DWS.0000", we need to use "error_msg" to determine the next operation.
func buildDeleteRetryFunc(client *golangsdk.ServiceClient, deletePath string, deleteOpt *golangsdk.RequestOpts) common.RetryFunc {
	retryFunc := func() (interface{}, bool, error) {
		deleteResp, err := client.Request("DELETE", deletePath, deleteOpt)
		if err != nil {
			return nil, false, fmt.Errorf("error deleting DWS logical cluster: %s", err)
		}

		deleteRespBody, err := utils.FlattenResponse(deleteResp)
		if err != nil {
			return nil, false, err
		}

		errCode := utils.PathSearch("error_code", deleteRespBody, "").(string)
		if errCode == "DWS.0000" {
			return nil, false, nil
		}

		errMsg := utils.PathSearch("error_msg", deleteRespBody, "").(string)
		// Stop retrying deletion when the resource does not exist or the current resource is the first logical cluster.
		if errMsg == deleteNotExistMsg || errMsg == deleteFirstLogicalClusterMsg {
			return errMsg, false, nil
		}
		return nil, true, fmt.Errorf("error deleting DWS logical cluster: error code: %s, error message: %s", errCode, errMsg)
	}
	return retryFunc
}

// waitingForStateDeleted This method is used to wait for delete to complete.
func waitingForStateDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			clusterRespBody, err := readLogicalClusters(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			cluster := utils.PathSearch(expression, clusterRespBody, nil)
			if cluster == nil {
				obj := map[string]string{"code": "COMPLETED"}
				return obj, "COMPLETED", nil
			}
			return cluster, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceLogicalClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}"
		product   = "dws"
		clusterId = d.Get("cluster_id").(string)
	)
	// Cannot be deleted when there are other tasks being executed.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	rst, err := waitingForDeleteStateEnable(ctx, client, d, d.Timeout(schema.TimeoutDelete))
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		if rst == "last_logical_cluster" {
			errorMsg := `The last logical cluster can't be deleted. Deleting this resource will only remove
		    the resource information from the tfstate file, but it remains in the cloud.`
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  errorMsg,
				},
			}
		}
		return common.CheckDeletedDiag(d, err, "error deleting DWS logical cluster")
	}

	if err != nil {
		return diag.Errorf("error waiting for DWS logical cluster (%s) to become delete enable: %s", d.Id(), err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)
	deletePath = strings.ReplaceAll(deletePath, "{logical_cluster_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		OkCodes:          []int{200, 202, 204, 417},
	}

	// When the cluster is operated concurrently, the deletion operation may also fail and needs to be retried.
	errMsg, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    buildDeleteRetryFunc(client, deletePath, &deleteOpt),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 30 * time.Second,
		PollInterval: 30 * time.Second,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	if errMsg == deleteNotExistMsg {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	if errMsg == deleteFirstLogicalClusterMsg {
		errMessage := "The first logical cluster can't be deleted. The project is only removed from the state," +
			" but it remains in the cloud."
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errMessage,
			},
		}
	}

	err = waitingForStateDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for DWS logical cluster (%s) deletion to complete: %s", d.Id(), err)
	}

	return nil
}

func resourceLogicalClusterImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	partLength := len(parts)

	if partLength == 2 {
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("cluster_id", parts[0])
	}
	return nil, fmt.Errorf("invalid format specified for import ID, must be <cluster_id>/<id>")
}
