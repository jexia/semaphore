# Go API

Semaphore exposes most of the internally used methods. This allows developers to create their own custom implementations. Please check out the [go docs](https://pkg.go.dev/github.com/jexia/maestro).

<CodeSwitcher :languages="{go:'Go lang',bash:'Bash'}">
<template v-slot:bash>

```bash
go get github.com/jexia/maestro
```

</template>
<template v-slot:go>

```go
package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol/graphql"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
)

func main() {
    protobuffers, err := protoc.Collect([]string{"./"}, "./*")
    if err != nil {
        // handle err
    }

    client, err := maestro.New(
        maestro.WithListener(graphql.NewListener(":9090", specs.Options{})),
        maestro.WithListener(http.NewListener(":8080", specs.Options{})),
        maestro.WithDefinitions(hcl.DefinitionResolver("./*")),
        maestro.WithSchema(protobuffers),
        maestro.WithCodec(json.NewConstructor()),
        maestro.WithCodec(proto.NewConstructor()),
        maestro.WithCaller(http.NewCaller()),
    )
    
    if err != nil {
        // handle err
    }
    
    err = client.Serve()
    if err != nil {
        // handle err
    }
}
```
</template>
</CodeSwitcher>

## CLI
You could create your own Semaphore implementation or use the predefined CLI. The official CLI has most implementations available.

```bash
semaphore daemon --http=:8080 --proto=./* --flow=./*
```

::: tip
Enterprise features are only available inside the enterprise CLI
:::

Multiple flow paths could be given as arguments. Definitions could be looked up recursively with wildcards.

<CodeSwitcher :languages="{bash:'Bash'}">
<template v-slot:bash>

``` bash
semaphore daemon --http=:8080 --proto=./* --flow=/var/flows/** --flow=/var/endpoints/**
```

</template>
</CodeSwitcher>

## Installation
You are able to download a prebuild artifact from the [latest release](https://github.com/jexia/semaphore/releases). Feel free to open a new PR if you require a specific build for your CPU architecture.

## Linux / MacOS
A installer script is available. By default the latest binaries for your operating system will be pulled and stored in `/usr/local/bin`. Arguments could be given to pull a specific version and/or store the binary inside a specific directory.

```bash
# pull latest version
curl https://raw.githubusercontent.com/jexia/maestro/master/install.sh | sh
# pull version v2.0.0 and store it in ./bin
curl https://raw.githubusercontent.com/jexia/maestro/master/install.sh | sh -s -- -b ./bin v2.0.0
```

## Homebrew
A Maestro Homebrew installer is available inside the Jexia cask. Simply tap into the cask and install Maestro.

```bash
brew tap jexia/cask
brew install maestro
```

## Docker image
Official docker images are available on Github and Docker hub. These images contain the Maestro CLI.

```bash
docker pull jexiacom/maestro-cli
```

## Building the Development Version from Source
If you have a Go environment configured, you can install the development version of `semaphore` from the command line.

```bash
go install github.com/jexia/maestro/cmd/semaphore
```

While the development version is a good way to take a peek at `maestro's` latest features before they get released, be aware that it may have bugs. Officially released versions will generally be more stable.

## Validate flows
Flow definitions could easily be checked with the `validate` command. This validates the flow property types with the configured schema(s).
```bash
semaphore validate --proto=./* --flow=/var/flows/** --flow=/var/endpoints/**
```

## Configuration files
CLI configurations could be defined inside a hcl file. Configuration files could include other hcl definitions and access environment variables. Service selectors could be defined which override the configuration of services. It is often adviced to store these service selectors inside a seperate file to include them inside your specific environment `include = ["services.$ENV.hcl"].`

```hcl
config.hcl
log_level = "$LOG_LEVEL"
protobuffers = ["../../annotations", "./proto/*.proto"]
​
include = ["flow.hcl"]
​
graphql {
    address = "$GRAPHQL_ADDRESS"
}
​
http {
    address = "$HTTP_ADDRESS"
}
​
prometheus {
    address = "$PROMETHEUS_ADDRESS"
}
​
services {
    select "proto.users.*" {
        host = "api.jexia.com"
    }
​
    select "proto.projects.*" {
        host = "api.jexia.com"
    }
}
```

`hcl` files could be passed to Maestro with the `--file or (-f)` flag.

```bash
semaphore daemon --file config.hcl
```

## Advanced pattern matching
Maestro supports advanced pattern matching for most paths. Patterns could be used to target specific files inside a given directory.

```bash
./**/*.e2e.hcl
```

## Docker
General docker images are available containing the Semaphore CLI configured as entrypoint. Please check out the CLI documentation for more information.

```bash
docker pull jxapp/semaphore
```

You could build your own docker image containing your flows and schema definitions. Simply pull the `jxapp/semaphore` image and copy your definitions.

```bash
FROM jxapp/semaphore

# include local flow definitions
COPY . .
```

## Examples
Various examples are available inside the git repository. Most examples could be ran using the Semaphore CLI. Some examples require Go to be installed on your local machine. Check out the official Go website for more information on how to install Go on your machine.