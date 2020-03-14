# CLI

Use the standard CLI or develop your own implementation.

## Installing

### Docker images

Official docker images are available on Github. These images contain the Maestro CLI and could be used to set up a server.

```
docker pull docker.pkg.github.com/jexia/maestro/cli:latest
```

### Using a Package Manager (Preferred)

> 🚧 Currently it is only possible to download + build the CLI from source

### Building the Development Version from Source

If you have a Go environment
configured, you can install the development version of `maestro` from
the command line.

```
go install github.com/jexia/maestro/cmd/maestro
```

While the development version is a good way to take a peek at
`maestro`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.