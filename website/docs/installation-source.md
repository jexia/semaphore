---
id: installation.source
title: Build from source
sidebar_label: Build from source
slug: /installation/source
---

If you have a Go environment configured, you can install the development version of `semaphore` from
the command line. This will build a binary for the machines CPU architecture and environment.
While the development version is a good way to take a peek at
`semaphore`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.

```sh
$ git clone https://github.com/jexia/semaphore.git
$ go build -o semaphore ./cmd/semaphore
```