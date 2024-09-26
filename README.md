### [Discussions](https://github.com/rocketblend/rocketblend-desktop/discussions) │ [Latest Release](https://github.com/rocketblend/rocketblend-desktop/releases/latest)

# RocketBlend Desktop

[![Github tag](https://badgen.net/github/tag/rocketblend/rocketblend-desktop)](https://github.com/rocketblend/rocketblend-desktop/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/rocketblend/rocketblend-desktop)](https://goreportcard.com/report/github.com/rocketblend/rocketblend-desktop)
[![GitHub](https://img.shields.io/github/license/rocketblend/rocketblend-desktop)](https://github.com/rocketblend/rocketblend-desktop/blob/master/LICENSE)

RocketBlend Desktop is an open-source desktop application for [RocketBlend](https://github.com/rocketblend/rocketblend), available on Windows, Mac, and Linux. Improve your [Blender](https://www.blender.org/) workflow by adding dependency management and enahanced project discoverability.

> [!NOTE]  
> **Important:** RocketBlend Desktop is currently under active development and still evolving. As such, expect significant changes, potential bugs, and incomplete features.

![Image of RocketBlend desktop application](docs/assets/rocketblend-desktop-dev.png)

## Features

- **Project Exploration and Search**: Easily explore and search through all your Blender projects, making management and organization simple.
- **Build and Add-on Management**: Discover, download, and manage Blender builds and add-ons (note: add-ons are currently a work in progress).
- **Custom Blender Environment**: Assign specific Blender builds and add-ons to your projects to ensure consistency and prevent dependency issues.
- **Easy Project Access**: Set RocketBlend Desktop as the default application for `.blend` files to open projects with the correct builds and add-ons, preventing project breaks.
- **Custom Package Definition**: Create and manage custom builds and add-ons specific to your project's needs.
- **Network Storage and Team Collaboration**: Work seamlessly across shared network storage, enabling easy sharing across machines and team members. Configurations are stored in `.json` files alongside your `.blend` files, suitable for shared drives, cloud storage, or Git repositories.
- **Lightweight and Flexible**: Use RocketBlend's [CLI](https://github.com/rocketblend/rocketblend) directly without needing the desktop UI. Ideal for professionals setting up render pipelines or server tasks, while the desktop application offers an intuitive interface for regular users.
- **Visual Project Gallery**: Navigate and manage your work using the visual project gallery, supporting both images and videos.

## Development Requirements

- [Golang 1.21.x](https://go.dev/dl/)
- [NodeJS 20.x.x](https://nodejs.org/en/) Recommended to use [nvm](https://github.com/nvm-sh/nvm#installing-and-updating) or [windows-nvm](https://github.com/coreybutler/nvm-windows#installation--upgrades) to manage NodeJS versions.
- [Wails 2.8.0](https://wails.io/docs/gettingstarted/installation#platform-specific-dependencies)
  - Then run `wails doctor` to ensure you have all the correct system-level dependencies installed.
- [Git LFS](https://git-lfs.com/) to manage binary files (images/fonts).
- [Mage](https://magefile.org/)

### Developing Locally

To run in live development mode:

- `wails dev` in the project directory
  - This will run a Vite development server that will provide very fast hot reload of your frontend changes.

If you want to develop in a browser and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect to this in your browser, and you can call your Go code from devtools.

### Building

To build a redistributable, production mode package, use `wails build`.

### Code Signing (Github Action)



#### Windows (Azure Trusted)

See Azure.

#### Windows (Classic)

- `WIN_CERTIFICATE` - `.p12` formatted certificate encoded in base64.
- `WIN_CERTIFICATE_PASSWORD` - certificate password.

#### Apple

- `AC_CERTIFICATE` - `.p12` Apple certificate encoded in base64.
- `AC_CERTIFICATE_PASSWORD ` - Apple certificate password.
- `AC_USERNAME` - The Apple ID username, typically an email address. 
- `AC_PASSWORD` - The password for the associated Apple ID.
- `AC_PROVIDER` - The App Store Connect provider when using multiple teams within App Store Connect
- `AC_APPLICATION_IDENTITY` - The name or ID of the "Developer ID Application" certificate to use to sign applications

## See Also

- [RocketBlend](https://github.com/rocketblend/rocketblend) - CLI tool that powers the build and addon management for blender.
- [Official Library](https://github.com/rocketblend/official-library) - Collection of packages for rocketblend.

## Special Thanks

- [Wails](https://wails.io/) - A Go library for building Web applications.
- [Sveltekit](https://kit.svelte.dev/) - The framework that powers the frontend of RocketBlend Desktop.
- [Skeleton](https://www.skeleton.dev/) - UI toolkit for Svelte and Tailwind.
- [Bleve](https://github.com/blevesearch/bleve) - A modern text indexing library for Go.
- [Abjrcode](https://github.com/abjrcode) - For providing a good example pipeline and cross platform build system for wails applications.

## License

RocketBlend Desktop is licensed under the [AGPL License](LICENSE).