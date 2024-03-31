### [Discussions](https://github.com/rocketblend/rocketblend-desktop/discussions) │ [Latest Release](https://github.com/rocketblend/rocketblend-desktop/releases/latest)

# Rocketblend Desktop

[![Github tag](https://badgen.net/github/tag/rocketblend/rocketblend-desktop)](https://github.com/rocketblend/rocketblend-desktop/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocketblend/rocketblend-desktop)](https://goreportcard.com/report/github.com/rocketblend/rocketblend-desktop)
[![GitHub](https://img.shields.io/github/license/rocketblend/rocketblend-desktop)](https://github.com/rocketblend/rocketblend-desktop/blob/master/LICENSE)

> [!NOTE]  
> Please note that RocketBlend Desktop is still in development and may undergo significant changes as it continues to evolve and is not yet ready for general use.

![Image of RocketBlend desktop application](docs/assets/rocketblend-desktop-dev.png)

## Development Requirements

- [Golang 1.21.x](https://go.dev/dl/)
- [NodeJS 20.x.x](https://nodejs.org/en/) Recommended to use [nvm](https://github.com/nvm-sh/nvm#installing-and-updating) or [windows-nvm](https://github.com/coreybutler/nvm-windows#installation--upgrades) to manage NodeJS versions.
- [Wails 2.7.1](https://wails.io/docs/gettingstarted/installation#platform-specific-dependencies)
  - Then run `wails doctor` to ensure you have all the correct system-level dependencies installed.
- [Mage](https://magefile.org/)

### Developing Locally

To run in live development mode:

- `wails dev` in the project directory
  - This will run a Vite development server that will provide very fast hot reload of your frontend changes.

If you want to develop in a browser and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect to this in your browser, and you can call your Go code from devtools.

### Building

To build a redistributable, production mode package, use `wails build`.

## See Also

- [RocketBlend](https://github.com/rocketblend/rocketblend) - CLI tool that powers the build and addon management for blender.
- [RocketBlend Companion](https://github.com/rocketblend/rocketblend-companion) - Blender addon to aid with working with RocketBlend. **NOTE: WIP**
- [Official Library](https://github.com/rocketblend/official-library) - Collection of packages for rocketblend.

## Special Thanks

- [Wails](https://wails.io/) - A Go library for building Web applications.
- [Sveltekit](https://kit.svelte.dev/) - The framework that powers the frontend of RocketBlend Desktop.
- [Skeleton](https://www.skeleton.dev/) - UI toolkit for Svelte and Tailwind.
- [Bleve](https://github.com/blevesearch/bleve) - A modern text indexing library for Go.
- [Abjrcode](https://github.com/abjrcode) - For the pipeline and build system help with wails.

## License

RocketBlend Desktop is licensed under the [AGPL License](LICENSE).