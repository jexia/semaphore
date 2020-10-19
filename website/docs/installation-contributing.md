---
id: installation.contributing
title: Contributing
sidebar_label: Contributing
slug: /installation/contributing
---

If you wish to work on Semaphore itself or any of its built-in systems, you'll
first need [Go](https://www.golang.org) installed on your machine. Go version
1.13.7+ is *required*.

For local dev first make sure Go is properly installed, including setting up a
[GOPATH](https://golang.org/doc/code.html#GOPATH). Ensure that `$GOPATH/bin` is in
your path as some distributions bundle old version of build tools. Next, clone this
repository. Semaphore uses [Go Modules](https://github.com/golang/go/wiki/Modules),
so it is recommended that you clone the repository ***outside*** of the GOPATH.
You can then download any required build tools by bootstrapping your environment:

```sh
$ make bootstrap
...
```

To compile a development version of Semaphore, run `make` or `make dev`. This will
put the Semaphore binary in the `bin` folders:

```sh
$ make dev
...
$ bin/semaphore
...
```

To run tests, type `make test`. If
this exits with exit status 0, then everything is working!

```sh
$ make test
...
```