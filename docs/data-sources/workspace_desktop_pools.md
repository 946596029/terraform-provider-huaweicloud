---
subcategory: "Workspace"
---

# huaweicloud_workspace_desktop_pools

Use this data source to get a list of Workspace desktop pools.

## Example Usage

```hcl
data "huaweicloud_workspace_desktop_pools" "test" {
  name = "pool_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the desktop pools.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the desktop pool.

* `type` - (Optional, String) Specifies the type of the desktop pool.
  The valid values are **DYNAMIC** and **STATIC**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `in_maintenance_mode` - (Optional, Bool) Specifies whether the desktop pool is in maintenance mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of desktop pools.

* `desktop_pools` - List of desktop pool details.
  The [desktop_pools](#workspace_desktop_pools_desktop_pools) structure is documented below.

<a name="workspace_desktop_pools_desktop_pools"></a>
The `desktop_pools` block supports:

* `id` - The ID of the desktop pool.

* `name` - The name of the desktop pool.

* `type` - The type of the desktop pool.

* `description` - The description of the desktop pool.

* `created_time` - The creation time of the desktop pool.

* `charging_mode` - The charging mode of the desktop pool.

* `desktop_count` - The total number of desktops in the pool.

* `desktop_used` - The number of used desktops in the pool.

* `availability_zone` - The availability zone of the desktop pool.

* `subnet_id` - The subnet ID of the desktop pool.

* `product` - The product information of the desktop pool.
  The [product](#workspace_desktop_pools_product) structure is documented below.

* `image_id` - The image ID used by the desktop pool.

* `image_name` - The image name used by the desktop pool.

* `image_os_type` - The OS type of the image.

* `image_os_version` - The OS version of the image.

* `image_os_platform` - The OS platform of the image.

* `image_product_code` - The product code of the image.

* `root_volume` - The root volume information of the desktop pool.
  The [root_volume](#workspace_desktop_pools_volume) structure is documented below.

* `data_volumes` - The data volumes information of the desktop pool.
  The [data_volumes](#workspace_desktop_pools_volume) structure is documented below.

* `security_groups` - The security groups of the desktop pool.
  The [security_groups](#workspace_desktop_pools_security_group) structure is documented below.

* `disconnected_retention_period` - The disconnected retention period in minutes.

* `enable_autoscale` - Whether auto scaling is enabled.

* `autoscale_policy` - The auto scaling policy of the desktop pool.
  The [autoscale_policy](#workspace_desktop_pools_autoscale_policy) structure is documented below.

* `status` - The status of the desktop pool.

* `enterprise_project_id` - The enterprise project ID.

* `in_maintenance_mode` - Whether the desktop pool is in maintenance mode.

* `desktop_name_policy_id` - The desktop name policy ID.

<a name="workspace_desktop_pools_product"></a>
The `product` block supports:

* `product_id` - The product ID.

* `flavor_id` - The flavor ID.

* `type` - The product type.

* `cpu` - The CPU specification.

* `memory` - The memory specification.

* `descriptions` - The product description.

* `charge_mode` - The charging mode.

<a name="workspace_desktop_pools_volume"></a>
The `root_volume` and `data_volumes` block supports:

* `id` - The volume ID.

* `type` - The volume type.

* `size` - The volume size in GB.

* `resource_spec_code` - The resource specification code.

<a name="workspace_desktop_pools_security_group"></a>
The `security_groups` block supports:

* `id` - The security group ID.

<a name="workspace_desktop_pools_autoscale_policy"></a>
The `autoscale_policy` block supports:

* `autoscale_type` - The auto scaling type.

* `max_auto_created` - The maximum number of auto-created desktops.

* `min_idle` - The minimum number of idle desktops.

* `once_auto_created` - The number of desktops to create in one auto scaling operation.
