# DevOps
Running Semaphore is pretty simple. The daemon by default spawns a production ready broker. The Semaphore daemon is stateless and multiple instances could be ran at the same time to provide scalability.

## Service hosts
Services could have different hosts when running Maestro in multiple environments. Service configurations could be overridden through service selectors. It is adviced to store your service selectors inside a separate file and use a environment variable to include a specific service configuration.

config.hcl
```hcl
include = ["services.$ENVIRONMENT.hcl"]
```

service.production.hcl
```hcl
services {
    select "org.users.*" {
        host = "prod.org-users.com"
    }

    select "org.projects.*" {
        host = "prod.org-projects.com"
    }
}
```

## Service certificates
Root certificates could be included to provide secure connections. Certificates could be passed as options or be overridden through service selectors.

```flow
services {
    select "proto.users.*" {
			host = "api.jexia.com"
			insecure = "false"
			ca_file = "/etc/ca.crt"
    }

    select "proto.projects.*" {
      host = "api.jexia.com"
			insecure = "true"
    }
}
```

## Prometheus
A Prometheus metrics endpoint could be set-up. This endpoint exposes metrics such as flow executions, executed rollbacks and flow latency. The Prometheus agent starts its own HTTP server and requires a separate port.

```hcl
config.hcl
prometheus {
    address = ":5050"
}
```