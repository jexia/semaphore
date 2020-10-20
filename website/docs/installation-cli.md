---
id: installation.cli
title: Semaphore CLI
sidebar_label: CLI
slug: /installation/cli
---

There are variouse sources available to download and install the [⚡ Semaphore CLI](https://github.com/jexia/semaphore/tree/master/cmd/semaphore). Semaphore is invoked from the command line. The CLI could be used to spin up a Semaphore instance, validate flow definitions or manage your already running instances.
Execute `semaphore` with the `--help` flag for more information.

```bash
$ semaphore daemon -f config.hcl
```

:::tip
Enterprise features are only available inside the enterprise CLI
:::

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
