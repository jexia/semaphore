# Maestro

Maestro is invoked from the command line. The CLI could be used to spin up a Maestro instance, validate flow definitions or manage your already running instances.
Execute `maestro` with the `--help` flag for more information.

```bash
$ maestro daemon -f config.hcl
```

## Installing

### Brew

You are able to download a prebuild artifact from the [latest release](https://github.com/jexia/maestro/releases).
Feel free to open a new PR if you require a specific build for your CPU architecture.

```sh
$ brew tap jexia/cask
$ brew install maestro
```

### Docker images

Official docker images are available on Github. These images contain the Maestro CLI.

```
docker pull jexiacom/maestro-cli
```

### Building the Development Version from Source

If you have a Go environment
configured, you can install the development version of `maestro` from
the command line.

```
go build -o maestro ./*.go
```

This will build a binary for the machines CPU architecture and environment.
While the development version is a good way to take a peek at
`maestro`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.
