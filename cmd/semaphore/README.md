# Semaphore

Semaphore is invoked from the command line. The CLI could be used to spin up a Semaphore instance, validate flow definitions or manage your already running instances.
Execute `semaphore` with the `--help` flag for more information.

```bash
$ semaphore daemon -f config.hcl
```

## Installing

You are able to download a prebuild artifact from the [latest release](https://github.com/jexia/semaphore/releases).
Feel free to open a new PR if you require a specific build for your CPU architecture.

### Homebrew

A Semaphore Homebrew installer is available inside the Jexia cask.
Simply tap into the cask and install Semaphore.

```sh
$ brew tap jexia/cask
$ brew install semaphore
```

### Unix

A installer script is available.
By default the latest binaries for your operating system will be pulled and stored in `/usr/local/bin`.
Arguments could be given to pull a specific version and/or store the binary inside a specific directory.

```sh
$ # pull latest version
$ curl https://raw.githubusercontent.com/jexia/semaphore/master/install.sh | sh
$ # pull version v2.0.0 and store it in ./bin
$ curl https://raw.githubusercontent.com/jexia/semaphore/master/install.sh | sh -s -- -b ./bin v2.0.0
```

### Docker images

Official docker images are available on Github and Docker hub. These images contain the Semaphore CLI.

```sh
$ docker pull jxapp/semaphore-cli
```

### Building the Development Version from Source

If you have a Go environment
configured, you can install the development version of `semaphore` from
the command line.

```sh
$ git clone https://github.com/jexia/semaphore.git
$ go build -o semaphore ./cmd/semaphore
```

This will build a binary for the machines CPU architecture and environment.
While the development version is a good way to take a peek at
`semaphore`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.
