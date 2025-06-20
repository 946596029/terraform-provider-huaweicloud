---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_pools"
description: |-
  Use this data source to query the desktop pools under a specified region within HuaweiCloud.
---

# huaweicloud_workspace_desktop_pools

Use this data source to query the desktop pools under a specified region within HuaweiCloud.

## Example Usage

### Querying all desktop pools

```hcl
data "huaweicloud_workspace_desktop_pools" "test" {}
```

### Querying desktop pools by name

```hcl
variable "pool_name" {}

data "huaweicloud_workspace_desktop_pools" "test" {
  name = var.pool_name
}
```

### Querying desktop pools by type

```hcl
data "huaweicloud_workspace_desktop_pools" "test" {
  type = "STATIC"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the desktop pools.

* `name` - (Optional, String) The name of the desktop pool.

* `type` - (Optional, String) The type of the desktop pool. Valid values are: **DYNAMIC**, **STATIC**.

* `enterprise_project_id` - (Optional, String) The enterprise project ID.

* `in_maintenance_mode` - (Optional, Bool) Whether the desktop pool is in maintenance mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pools` - The list of desktop pools.  
  The [pool](#workspace_desktop_pool) structure is documented below.

<a name="workspace_desktop_pool"></a>
The `pool` block supports:

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
  The [product](#workspace_desktop_pool_product) structure is documented below.

* `image_id` - The image ID used by the desktop pool.

* `image_name` - The image name used by the desktop pool.

* `image_os_type` - The OS type of the image.

* `image_os_version` - The OS version of the image.

* `image_os_platform` - The OS platform of the image.

* `image_product_code` - The product code of the image.

* `root_volume` - The root volume information of the desktop pool.  
  The [volume](#workspace_desktop_pool_volume) structure is documented below.

* `data_volumes` - The data volumes information of the desktop pool.  
  The [volume](#workspace_desktop_pool_volume) structure is documented below.

* `security_groups` - The security groups of the desktop pool.  
  The [security_group](#workspace_desktop_pool_security_group) structure is documented below.

* `disconnected_retention_period` - The disconnected retention period in minutes.

* `enable_autoscale` - Whether auto scaling is enabled.

* `autoscale_policy` - The auto scaling policy of the desktop pool.  
  The [autoscale_policy](#workspace_desktop_pool_autoscale_policy) structure is documented below.

* `status` - The status of the desktop pool.

* `enterprise_project_id` - The enterprise project ID.

* `in_maintenance_mode` - Whether the desktop pool is in maintenance mode.

* `desktop_name_policy_id` - The desktop name policy ID.

<a name="workspace_desktop_pool_product"></a>
The `product` block supports:

* `product_id` - The ID of the product.

* `flavor_id` - The flavor ID of the product.

* `type` - The type of the product.

* `cpu` - The CPU specification.

* `memory` - The memory specification.

* `descriptions` - The description of the product.

* `charge_mode` - The charging mode of the product.

<a name="workspace_desktop_pool_volume"></a>
The `volume` block supports:

* `id` - The ID of the volume.

* `type` - The type of the volume.

* `size` - The size of the volume in GB.

* `resource_spec_code` - The resource specification code of the volume.

<a name="workspace_desktop_pool_security_group"></a>
The `security_group` block supports:

* `id` - The ID of the security group.

<a name="workspace_desktop_pool_autoscale_policy"></a>
The `autoscale_policy` block supports:

* `autoscale_type` - The type of auto scaling.

* `max_auto_created` - The maximum number of auto-created desktops.

* `min_idle` - The minimum number of idle desktops.

* `once_auto_created` - The number of desktops to create in one auto scaling operation.
