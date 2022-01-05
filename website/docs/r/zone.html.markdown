---
layout: "ultradns"
page_title: "UltraDNS: ultradns_zone"
sidebar_current: "docs-ultradns-resource-zone"
description: |-
  Provides an UltraDNS Zone resource
---

# ultradns\_zone

Provides an UltraDNS Zone resource.

## Example Usage

```hcl
# Create a Zone
resource "ultradns_zone" "test" {
  name         = "tftest.example.com."
  type         = "PRIMARY"
  account      = "my_account"
}
```

## Argument References

See [related part of UltraDNS Docs](https://docs.ultradns.neustar/Content/REST%20API/Content/REST%20API/Zone%20API/Zone%20API%20DTOs.htm) for details about valid values.

The following arguments are supported:

* `name` - (Required) The name of the zone (as an FQDN, ending with a `.`)
* `type` - (Required) The type of zone to create. One of `PRIMARY`, `SECONDARY`, or `ALIAS`
* `account` - (Required) The name of the UltraDNS account (can be found at https://portal.ultradns.neustar/accounts)
* `create_type` - Designate whether zone is new (`NEW`) or copied from an existing zone (`COPY`)
* `original_zone` - If `create_type` is `COPY`, the name of the original zone to copy from
* `alias_target` - If `type` is `ALIAS`, the name of the target for the zone
