---
id: flows.proxy
title: Proxy Forwarding
sidebar_label: Proxy Forwarding
slug: /flows/proxy
---

Proxy forwarding allows the entire request to be forwarded to other services. A proxy forward is unable to switch protocol and forwards the entire request body to the targeted service. Requests could be made before forwarding a request. The input body could not be used in any of the configured resources.

[Check out the hubs example inside the git repo.](https://github.com/jexia/semaphore/tree/master/examples/multiple-gateways)

```hcl
proxy "forward" {
	resource "authenticate" {
		request "auth" "User" {
		}
	}

	forward "hub" {
	}
}
```

:::important
Not all protocols support proxy forwarding, please check the protocol documentation for more information
:::