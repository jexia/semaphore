---
id: service_discovery.configuration
title: Service Discovery Client Configuration
sidebar_label: Service Discovery
slug: /service_discovery/configuration
---

## Intro

Semaphore can resolve service address not only through DNS, but using a specific service discovery server such as Consul by HashiCorp.

By default, Semaphore resolves service hosts  through the default DNS performing a regular query with no actions from the Semaphore side.
We can set up several service discovery configurations and use them as address resolvers in every service independently.

## Defining clients

To configure a service discovery client, we use `discovery` block:

```hcl
discovery "consul" {
  address = "http://localhost:8500"
}
```

By default, the block label (`"consul"` in the example above) defines not only the name of the configuration, but the provider type as well.
We still can define a configuration with a custom name, using `provider` field to tell Semaphore what adapter it should use:

```hcl
discovery "production-1" {
  address = "http://localhost:8500"
  provider = "consul"
}
``` 

Semaphore supports defining several discovery clients what might be useful in some rare cases:

```hcl
discovery "production-old" {
  address = "http://production-old:850"
  provider = "consul"
}

discovery "production-new" {
  address = "http://production-new:850"
  provider = "consul"
}

// and so on
```

## Using discovery clients

To use a discovery service, we should set `resolver` field in service configuration, referencing to the discovery client name:

```hcl
service "com.semaphore" "awesome-dogs" {
  transport = "http"
  codec     = "json"
  host      = "http://awesome-dogs"
  resolver  = "consul" // or "production-old"

  method "List" {
    response = "com.semaphore.Dogs"
    request = "com.semaphore.Void"

    options {
      endpoint = "/"
      method = "GET"
    }
  }
}
```

When we set `resolver` property, Semaphore uses `host` property as a service name and scheme. But the port should be received from the Service Discovery response.
In this case, `host` will be parsed, "awesome-dogs" will be used as the service name, and "http://" will be used as the scheme for the resolved address.

Let's assume, Consul returns "192.168.1.15" as the resolved address, and "8080" as the service port. 
The service URL will be `http://192.168.1.15:8080`. 

The another way to set a resolver for a service is using service selectors:

```hcl
services {
    select "com.semaphore.*" {
        resolver = "consul"
    }
}
```

## Example

You can find an example in [Semaphore](https://github.com/jexia/semaphore/tree/master/examples/consul) repositories, under `examples/consul` directory.