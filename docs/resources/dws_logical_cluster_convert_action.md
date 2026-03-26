---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_cluster_convert_action"
description: |-
  Use this resource to convert a physical DWS cluster to logical cluster mode within HuaweiCloud.
---

# huaweicloud_dws_logical_cluster_convert_action

Use this resource to convert a physical DWS cluster to logical cluster mode within HuaweiCloud.

-> This resource is only a one-time action resource for the conversion. Deleting this resource will not undo the
   conversion or clear the corresponding request record, but will only remove the resource information from the tfstate file.

-> Use only on a disposable physical cluster. The operation is irreversible.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "logical_cluster_name" {}

resource "huaweicloud_dws_logical_cluster_convert_action" "test" {
  cluster_id             = var.dws_cluster_id
  logical_cluster_name   = var.logical_cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DWS cluster is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the DWS cluster to convert to logical cluster mode.

* `logical_cluster_name` - (Required, String, NonUpdatable) Specifies the name to set for the logical cluster created from
  the conversion.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the convert action resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is `30` minutes.
