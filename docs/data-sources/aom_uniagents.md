---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_uniagents"
description: |-
  Use this data source to get the list of UniAgent hosts that have executed the UniAgent installation task.
---

# huaweicloud_aom_uniagents

Use this data source to get the list of UniAgent hosts that have executed the UniAgent installation task.

## Example Usage

```hcl
data "huaweicloud_aom_uniagents" "test" {
  page_size = 20
}
```

```hcl
data "huaweicloud_aom_uniagents" "filtered" {
  page_size     = 100
  ecs_id_list   = ["ecs-id-1", "ecs-id-2"]
  agent_id_list = ["agent-id-1"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `page` - (Optional, Int) Specifies the page number for pagination. Default value is 1.
  If specified, only the specified page will be queried. If not specified, all pages will be queried.

* `page_size` - (Optional, Int) Specifies the number of items per page. Default value is 20, maximum is 100.

* `ecs_id_list` - (Optional, List) Specifies the list of ECS IDs. Maximum 100 items are supported.

* `agent_id_list` - (Optional, List) Specifies the list of agent IDs. Maximum 100 items are supported.

* `coc_cmdb_id_list` - (Optional, List) Specifies the list of CMDB IDs. Maximum 100 items are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of UniAgents that matched the filter parameters.

* `uniagents` - The list of UniAgents that matched the filter parameters.  
  The [uniagents](#uniagents_struct) structure is documented below.

<a name="uniagents_struct"></a>
The `uniagents` block supports:

* `inner_ip` - The inner IP address of the host.

* `import_ip` - The import IP address of the host.

* `agent_id` - The ID of the agent.

* `host_name` - The host name.

* `os_type` - The operating system type.

* `agent_state` - The state of the UniAgent.

* `project_id` - The project ID to which the host belongs.

* `version` - The version of the UniAgent.

* `is_hw_cloud_host` - Whether the host is a Huawei Cloud host.

* `vpc_id` - The VPC ID.

* `cmdb_id` - The CMDB ID.

* `ecs_id` - The ECS ID, which is a unique value.

* `domain_id` - The domain ID to which the host belongs.

