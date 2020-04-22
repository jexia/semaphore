# Maestro

Maestro is invoked from the command line. The CLI could be used to spin up a Maestro instance, validate flow definitions or manage your already running instances.
Execute `maestro` with the `--help` flag for more information.

```bash
$ maestro run --config config.yaml
```

## Installing

### Release (Preferred)

You are able to download a prebuild artifact from the [latest release](https://github.com/jexia/maestro/releases).
Feel free to open a new PR if you require a specific build for your CPU architecture.

### Docker images

Official docker images are available on Github. These images contain the Maestro CLI.

```
docker pull jexiacom/maestro-cli
```

### Using a Package Manager

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