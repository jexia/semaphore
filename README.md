# Semaphore documentation

----

- Chat: [Discord](https://chat.jexia.com)
- Documentation: [Github pages](https://jexia.github.io/semaphore/) - [Branch](https://github.com/jexia/semaphore/tree/docs)
- Go package documentation: [GoDev](https://pkg.go.dev/github.com/jexia/semaphore)

> Take control of your data, connect with anything, and expose it anywhere through protocols such as HTTP, GraphQL, and gRPC!

This branch contains the documentation of the Semaphore project.
When wanting to suggest changes to the documentation please fork and clone this branch.
Changes to the documentation are automatically deployed and are hosted at Github pages.

## Developing locally

If you wish to work on the Semaphore documentations you'll first need [Node](https://nodejs.org/)
and [yarn](https://yarnpkg.com/)/[npm](npmjs.com) installed on your machine.
Clone the Semaphore repository and checkout the `docs` branch.

```sh
$ git clone https://github.com/jexia/semaphore.git
$ git checkout docs
```

You have to make sure that all dependencies are installed.

```sh
$ # install the project dependencies
$ npm install
```

The markdown source files are available inside `./src`.
Please do not perform manual changes inside the `./docs` directory.
The documentation is maintained and automatically generated through Github actions.

To run the development version of the Semaphore documentation, run `npm run dev` or `yarn dev`.
This will start a local development server that hot reloads on made changes.

```sh
$ npm run dev
```

## Contributing

Thank you for your interest in contributing to Semaphore! ❤
Check out the open projects and/or issues and feel free to join any ongoing discussion.

Everyone is welcome to contribute, whether it's in the form of code, documentation, bug reports, feature requests, or anything else. We encourage you to experiment with the project and make contributions to help evolve it to meet your needs!

See the [contributing guide](https://github.com/jexia/semaphore/blob/master/CONTRIBUTING.md) for more details.