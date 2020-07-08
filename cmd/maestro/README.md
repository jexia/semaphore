# Maestro

Maestro is invoked from the command line. The CLI could be used to spin up a Maestro instance, validate flow definitions or manage your already running instances.
Execute `maestro` with the `--help` flag for more information.

```bash
$ maestro daemon -f config.hcl
```

## Installing

You are able to download a prebuild artifact from the [latest release](https://github.com/jexia/maestro/releases).
Feel free to open a new PR if you require a specific build for your CPU architecture.

### Homebrew

A Maestro Homebrew installer is available inside the Jexia cask.
Simply tap into the cask and install Maestro.

```sh
$ brew tap jexia/cask
$ brew install maestro
```

### Unix

A installer script is available.
By default the latest binaries for your operating system will be pulled and stored in `/usr/local/bin`.
Arguments could be given to pull a specific version and/or store the binary inside a specific directory.

```sh
$ # pull latest version
$ curl https://raw.githubusercontent.com/jexia/maestro/master/install.sh | sh
$ # pull version v2.0.0 and store it in ./bin
$ curl https://raw.githubusercontent.com/jexia/maestro/master/install.sh | sh -s -- -b ./bin v2.0.0
```

### Docker images

Official docker images are available on Github and Docker hub. These images contain the Maestro CLI.

```sh
$ docker pull jexiacom/maestro-cli
```

### Building the Development Version from Source

If you have a Go environment
configured, you can install the development version of `maestro` from
the command line.

```sh
$ git clone https://github.com/jexia/maestro.git
$ go build -o maestro ./cmd/maestro
```

This will build a binary for the machines CPU architecture and environment.
While the development version is a good way to take a peek at
`maestro`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.
