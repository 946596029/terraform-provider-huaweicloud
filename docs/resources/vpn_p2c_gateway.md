---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_p2c_gateway"
description: |-
  Use this resource to manage a P2C VPN gateway within HuaweiCloud.
---

# huaweicloud_vpn_p2c_gateway

Use this resource to manage a P2C VPN gateway within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "eip_id" {}
variable "availability_zone_ids" {
  type = list(string)
}

resource "huaweicloud_vpn_p2c_gateway" "test" {
  name                  = var.name
  vpc_id                = var.vpc_id
  connect_subnet        = var.subnet_id
  availability_zone_ids = var.availability_zone_ids

  eip {
    id = var.eip_id
  }
}
```

### Creating a P2C VPN gateway with creating new EIP

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "availability_zone_ids" {
  type = list(string)
}
variable "eip_bandwidth_size" {}
variable "eip_bandwidth_name" {}

resource "huaweicloud_vpn_p2c_gateway" "test" {
  name                  = var.name
  vpc_id                = var.vpc_id
  connect_subnet        = var.subnet_id
  availability_zone_ids = var.availability_zone_ids

  eip {
    type           = "5_bgp"
    bandwidth_size = var.eip_bandwidth_size
    charge_mode    = "traffic"
    bandwidth_name = var.eip_bandwidth_name
  }
}
```

### Creating a P2C VPN gateway with full parameters

```hcl
variable "name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "max_connection_number" {}
variable "availability_zone_ids" {
  type = list(string)
}
variable "eip_bandwidth_size" {}
variable "eip_bandwidth_name" {}

resource "huaweicloud_vpn_p2c_gateway" "test" {
  name                  = var.name
  vpc_id                = var.vpc_id
  connect_subnet        = var.subnet_id
  flavor                = "Professional1"
  max_connection_number = var.max_connection_number
  availability_zone_ids = var.availability_zone_ids

  eip {
    type           = "5_bgp"
    bandwidth_size = var.eip_bandwidth_size
    charge_mode    = "traffic"
    bandwidth_name = var.eip_bandwidth_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the P2C VPN gateway is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC used by the P2C VPN gateway.

* `connect_subnet` - (Required, String, NonUpdatable) Specifies the ID of the VPC subnet used by the P2C VPN gateway.

* `eip` - (Required, List) Specifies the EIP used by the P2C VPN gateway.  
  The [object](#vpn_p2c_gateway_eip) structure is documented below.

* `name` - (Optional, String) Specifies the name of the P2C VPN gateway.  
  The valid length is limited from `1` to `64`.

* `flavor` - (Optional, String, NonUpdatable) Specifies the flavor of the P2C VPN gateway.  
  Defaults to **Professional1**.

* `availability_zone_ids` - (Optional, List, NonUpdatable) Specifies the list of availability zone IDs.  

* `max_connection_number` - (Optional, Int, NonUpdatable) Specifies the maximum number of simultaneously online client
  connections.  
  The valid value ranges from `1` to `500`.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

* `tags` - (Optional, List, NonUpdatable) Specifies the tags of the P2C VPN gateway.  
  The [object](#vpn_p2c_gateway_tags) structure is documented below.

<a name="vpn_p2c_gateway_eip"></a>
The `eip` block supports:

* `id` - (Optional, String, NonUpdatable) Specifies the ID of the existing EIP.

* `type` - (Optional, String, NonUpdatable) Specifies the type of the EIP.

* `charge_mode` - (Optional, String, NonUpdatable) Specifies the charge mode of the bandwidth.  
  The valid value are as follows:
  + **bandwidth**
  + **traffic**

* `bandwidth_id` - (Optional, String, NonUpdatable) Specifies the ID of the shared bandwidth.

* `bandwidth_size` - (Optional, Int, NonUpdatable) Specifies the bandwidth size in Mbit/s.  
  The valid value ranges from `1` to `300`.

* `bandwidth_name` - (Optional, String, NonUpdatable) Specifies the bandwidth name.

  ~> You can use `id` to specify an existing EIP or use `type`, `bandwidth_name`, `bandwidth_size` and
    `charge_mode` to create a new EIP.

<a name="vpn_p2c_gateway_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `value` - (Required, String) Specifies the value of the tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the P2C VPN gateway.

* `current_connection_number` - The current number of client connections.

* `order_id` - The order ID.

* `admin_state_up` - Whether the P2C VPN gateway is frozen.

* `frozen_effect` - The frozen effect.

* `version` - The version of the P2C VPN gateway.

* `upgrade_info` - The upgrade information of the P2C VPN gateway.

* `created_at` - The creation time, in RFC3339 format.

* `applied_at` - The effective time, in RFC3339 format.

* `updated_at` - The last update time, in RFC3339 format.

* `eip` - The EIP used by the P2C VPN gateway.  
  The [object](#vpn_p2c_gateway_eip_attr) structure is documented below.

<a name="vpn_p2c_gateway_eip_attr"></a>
The `eip` block supports:

* `ip_version` - The EIP version.

* `ip_billing_info` - The order information of the EIP.

* `ip_address` - The public IPv4 address of the EIP.

* `bandwidth_billing_info` - The order information of the bandwidth.

* `share_type` - The bandwidth share type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is `10` minutes.
* `update` - Default is `10` minutes.
* `delete` - Default is `10` minutes.

## Import

The P2C VPN gateway can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpn_p2c_gateway.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `vpc_id`, `connect_subnet`, `eip`,
`flavor`, `availability_zone_ids`, `max_connection_number`, `enterprise_project_id`, `tags`.
It is generally recommended running `terraform plan` after importing a VPN P2C gateway.
You can then decide if changes should be applied to the VPN P2C gateway, or the resource definition
should be updated to align with the VPN P2C gateway. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpn_p2c_gateway" "test" {
  ...

  lifecycle {
    ignore_changes = [
      vpc_id,
      connect_subnet,
      eip,
      flavor,
      availability_zone_ids,
      max_connection_number,
      enterprise_project_id,
      tags,
    ]
  }
}
