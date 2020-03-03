# CLI

Use the standard CLI or develop your own implementation.

## Installing

### Using a Package Manager (Preferred)

> 🚧 Currently it is only possible to download + build the CLI from source

### Downloading a Release from GitHub

> TLDR; one-liner where `<version>` is the full semantic version, e.g., `1.17.0`.

```
VERSION=0.0.0; curl -sL "https://github.com/jexia/maestro/releases/download/v$VERSION/maestro-$VERSION-linux-amd64.tar.gz" | tar -xzv; chmod +x ./maestro
```

---

Visit the [Releases
page](https://github.com/jexia/maestro/releases) for the
[`maestro` GitHub project](https://github.com/jexia/maestro), and find the
appropriate archive for your operating system and architecture.
Download the archive from from your browser or copy its URL and
retrieve it to your home directory with `wget` or `curl`.

For example, with `wget`:

```
cd ~
wget https://github.com/jexia/maestro/releases/download/v<version>/maestro-<version>-linux-amd64.tar.gz
```

Or with `curl`:

```
cd ~
curl -OL https://github.com/jexia/maestro/releases/download/v<version>/maestro-<version>-linux-amd64.tar.gz
```

Extract the binary:

```
tar xf ~/maestro-<version>-linux-amd64.tar.gz
```

where `<version>` is the full semantic version, e.g., `1.17.0`.

On Windows systems, you should be able to double-click the zip archive to extract the `doctl` executable.

Move the `doctl` binary to somewhere in your path. For example, on GNU/Linux and OS X systems:

```
sudo mv ~/doctl /usr/local/bin
```

Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `doctl` to their `PATH`.

### Building the Development Version from Source

If you have a Go environment
configured, you can install the development version of `maestro` from
the command line.

```
go get github.com/jexia/maestro/cli
```

While the development version is a good way to take a peek at
`maestro`'s latest features before they get released, be aware that it
may have bugs. Officially released versions will generally be more
stable.