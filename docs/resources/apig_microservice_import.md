---
subcategory: "API Gateway"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_microservice_import"
description: |-
  Use this resource to import microservice within HuaweiCloud.
---

# huaweicloud_apig_microservice_import

Use this resource to import microservice within HuaweiCloud.

-> This resource is only a one-time action resource for importing microservice. Deleting this resource will not clear
   the corresponding imported microservice, but will only remove the resource information from the tfstate file.

## Example Usage

### Import CSE microservice

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "cse_engine_id" {}
variable "cse_service_id" {}
variable "cse_version" {}

resource "huaweicloud_apig_microservice_import" "test" {
  instance_id     = var.instance_id
  service_type    = "CSE"
  protocol        = "HTTPS"
  backend_timeout = 5000
  auth_type       = "NONE"
  cors            = false

  group_info {
    group_id = var.group_id
  }
  
  apis {
    name       = "apiExampleName"
    req_method = "ANY"
    req_uri    = "/example/test"
    match_mode = "SWA"
  }
  
  cse_info {
    engine_id  = var.cse_engine_id
    service_id = var.cse_service_id
    version    = var.cse_version
  }
}
```

### Import CCE workload microservice

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "cluster_id" {}

resource "huaweicloud_apig_microservice_import" "test" {
  instance_id     = var.instance_id
  service_type    = "CCE"
  protocol        = "HTTPS"
  backend_timeout = 5000
  auth_type       = "NONE"
  cors            = false

  group_info {
    group_id = var.group_id
  }
  
  apis {
    name       = "example"
    req_method = "ANY"
    req_uri    = "/example/test"
    match_mode = "SWA"
  }
  
  cce_info {
    cluster_id    = var.cluster_id
    namespace     = "test"
    workload_type = "deployment"
    app_name      = "dp"
    port          = 80
  }
}
```

### Import Nacos microservice

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "service_name" {}
variable "microservice_ip_address" {}
variable "microservice_port" {}
variable "password" {}

resource "huaweicloud_apig_microservice_import" "test" {
  instance_id  = var.instance_id
  service_type = "NACOS"
  protocol     = "HTTPS"
  
  group_info {
    group_id = var.group_id
  }
  
  apis {
    name       = "nacos_api"
    req_method = "GET"
    req_uri    = "/nacos/test"
    match_mode = "SWA"
  }
  
  nacos_info {
    service_name = var.service_name

    server_config {
      ip_address = var.microservice_ip_address
      port       = var.microservice_port
    }

    user_info {
      user_name = "terraform"
      password  = var.password
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dedicated instance is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance.

* `group_info` - (Required, List, NonUpdatable) Specifies the API group information for importing microservice.  
  The [group_info](#apig_microservice_group_info) structure is documented below.

* `service_type` - (Required, String, NonUpdatable) Specifies the type of microservice center.  
  The valid values are as follows:
  + **CSE**: CSE microservice registry center
  + **CCE**: CCE cloud container engine (workload)
  + **NACOS**: Nacos registry center

* `apis` - (Required, List, NonUpdatable) Specifies the list of APIs to be imported.  
  The [apis](#apig_microservice_apis) structure is documented below.

* `protocol` - (Optional, String, NonUpdatable) Specifies the protocol for the dedicated instance to access the
  microservice.  
  The valid values are **HTTP** and **HTTPS**.  
  Defaults to **HTTPS**.

* `backend_timeout` - (Optional, Int, NonUpdatable) Specifies the timeout for the dedicated instance to request backend
  services.  
  The upper limit that can be modified is `600000`, default `5000`.  
  Unit: milliseconds.  
  Defaults to **5000**.

* `auth_type` - (Optional, String, NonUpdatable) Specifies the authentication type of the API.  
  The valid values are as follow:
  + **NONE**: No authentication
  + **APP**: Authenticate using the credentials issued by APIG.
  + **IAM**: Authenticate using the credentials issued by IAM.
  Defaults to **NONE**.

* `cors` - (Optional, Bool, NonUpdatable) Specifies whether to support CORS.  
  Defaults to **false**.

* `cse_info` - (Optional, List, NonUpdatable) Specifies the CSE microservice information.  
  The [cse_info](#apig_microservice_cse_info) structure is documented below.  
  This parameter is required when `service_type` is **CSE**.

* `cce_info` - (Optional, List, NonUpdatable) Specifies the CCE microservice information.  
  The [cce_info](#apig_microservice_cce_info) structure is documented below.  
  This parameter is required when `service_type` is **CCE**.

* `nacos_info` - (Optional, List, NonUpdatable) Specifies the Nacos microservice information.  
  The [nacos_info](#apig_microservice_nacos_info) structure is documented below.  
  This parameter is required when `service_type` is **NACOS**.

<a name="apig_microservice_group_info"></a>
The `group_info` block supports:

* `group_id` - (Required, String) Specifies the ID of the API group.

* `group_name` - (Optional, String) Specifies the name of the API group.

* `app_id` - (Optional, String) Specifies the ID of the application.

<a name="apig_microservice_apis"></a>
The `apis` block supports:

* `req_uri` - (Required, String) Specifies the request URI of the API.

* `name` - (Optional, String) Specifies the name of the API.

* `req_method` - (Optional, String) Specifies the request method of the API.

* `match_mode` - (Optional, String) Specifies the matching mode of the API.

<a name="apig_microservice_cse_info"></a>
The `cse_info` block supports:

* `engine_id` - (Required, String) Specifies the ID of the CSE engine.

* `service_id` - (Required, String) Specifies the ID of the CSE service.

* `version` - (Required, String) Specifies the version of the CSE service.

<a name="apig_microservice_cce_info"></a>
The `cce_info` block supports:

* `cluster_id` - (Required, String) Specifies the cluster ID of the CCE application.

* `namespace` - (Required, String) Specifies the namespace of the CCE application.

* `workload_type` - (Required, String) Specifies the type of the CCE workload.

* `port` - (Required, Int) Specifies the port of the CCE application.

* `app_name` - (Optional, String) Specifies the name of the CCE application.

* `label_key` - (Optional, String) Specifies the label key of the CCE application.

* `label_value` - (Optional, String) Specifies the label value of the CCE application.

* `version` - (Optional, String) Specifies the version of the CCE application.

* `labels` - (Optional, List) Specifies the labels of the CCE workload.  
  The [labels](#apig_microservice_cce_info_labels) structure is documented below.

<a name="apig_microservice_cce_info_labels"></a>
The `labels` block supports:

* `name` - (Required, String) Specifies the name of the label.

* `value` - (Required, String) Specifies the value of the label.

<a name="apig_microservice_nacos_info"></a>
The `nacos_info` block supports:

* `service_name` - (Required, String) Specifies the name of Nacos microservice.

* `server_config` - (Required, List) Specifies the server config of Nacos microservice.
  The [server_config](#apig_microservice_server_config) structure is documented below.

* `user_info` - (Required, List) Specifies the user information of Nacos microservice.
  The [user_info](#apig_microservice_user_info) structure is documented below.

* `namespace` - (Optional, String) Specifies the namespace of Nacos microservice.

* `group_name` - (Optional, String) Specifies the group name of Nacos microservice.

* `cluster_name` - (Optional, String) Specifies the cluster name of Nacos microservice.

<a name="apig_microservice_server_config"></a>
The `server_config` block supports:

* `ip_address` - (Required, String) Specifies the IP of the Nacos server.

* `port` - (Required, String) Specifies the port of the Nacos server.

* `grpc_port` - (Optional, String) Specifies the grpc port of the Nacos server.

<a name="apig_microservice_user_info"></a>
The `user_info` block supports:

* `user_name` - (Required, String) Specifies the name of the Nacos user.

* `password` - (Required, String) Specifies the password of the Nacos user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `vpc_channel_id` - The ID of the VPC channel created for the microservice.

* `api_group_id` - The ID of the API group.

* `imported_apis` - The list of imported APIs with generated IDs.  
  The [imported_apis](#apig_microservice_imported_apis) structure is documented below.

<a name="apig_microservice_imported_apis"></a>
The `imported_apis` block contains:

* `id` - The ID of the imported API.

* `name` - The name of the imported API.

* `req_uri` - The request URI of the imported API.

* `req_method` - The request method of the imported API.

* `match_mode` - The matching mode of the imported API.
