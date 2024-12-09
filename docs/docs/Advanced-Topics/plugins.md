# Plugins

Plugins in Harbor provide added functionality to Harbor. As of current, they can only provide a new executor for Harbor.

## How plugins work

Plugins are built on top of the gRPC protocol. The plugin architecture is based off of the [go-plugin](https://github.com/hashicorp/go-plugin) library. This allows you to make a plugin in any language, not just golang.

## Developing a plugin

Plugins must define a `manifest.json` file. This file contains two fields:

1. Supported Executors
    1. This is a list of strings that contains all of the executors this plugin supports. This should be prefixed with a unique identifier, for example `harbor.dev/`. This prevents collisions between plugins.
2. Executable Command
    1. This tells harbor how to start your plugin.
