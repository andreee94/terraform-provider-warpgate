---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "warpgate_http_target Resource - terraform-provider-warpgate"
subcategory: ""
description: |-
  
---

# warpgate_http_target (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `options` (Attributes) (see [below for nested schema](#nestedatt--options))

### Read-Only

- `id` (String) Id of the http target in warpgate

<a id="nestedatt--options"></a>
### Nested Schema for `options`

Required:

- `tls` (Attributes) (see [below for nested schema](#nestedatt--options--tls))
- `url` (String)

Optional:

- `external_host` (String)
- `headers` (Map of String)

<a id="nestedatt--options--tls"></a>
### Nested Schema for `options.tls`

Required:

- `mode` (String)
- `verify` (Boolean)

