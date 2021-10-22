# OpenFaaS Waypoint Plugin

This is a prototype OpenFaaS Waypoint plugin. It currently has _very_ basic support for deploying a function to an
OpenFaaS instance. **It should be considered experimental.**

See the [Waypoint documentation](https://www.waypointproject.io/docs/plugins#using-an-external-plugin)
on installing external plugins for more detail on installing external plugins.

## Configuration

A Waypoint configuration of this might look like:

```hcl
project = "my-functions"

app "wordcount" {
  build {
    use "docker" {}
    registry {
      use "docker" {
        image = "wordcount"
        tag   = "latest"
      }
    }
  }

  deploy {
    use "openfaas" {
      gateway  = "https://<your gateway endpoint>/"
      username = "<your gateway username>"
      password = "<your gateway password>"
    }
  }
}
```

## Development

### Building

To build the plugin, run:

```shell
make
```

This will regenerate the protos and build binaries for multiple platforms.

### Installation

To install the binary to `${HOME}/.config/waypoint/plugins/` run:

```shell
make install
```

### Building with Docker

To build plugins for release you can use the `build-docker` Makefile target, this will build your plugin for all
architectures and create zipped artifacts which can be uploaded to an artifact manager such as GitHub releases.

The built artifacts will be output in the `./releases` folder.

### Building and releasing with GitHub Actions

When cloning the template a default GitHub Action is created at the path `.github/workflows/build-plugin.yaml`. You can
use this action to automatically build and release your plugin.

The action has two main phases:

1. **Build** - This phase builds the plugin binaries for all the supported architectures. It is triggered when pushing
   to a branch or on pull requests.
1. **Release** - This phase creates a new GitHub release containing the built plugin. It is triggered when pushing tags
   which starting with `v`, for example `v0.1.0`.