# Installing Harbor

## Quick Installation (**This is not done yet**)

To install harbor, simply install from homebrew:

```bash
brew install harbor
```

## Installing from source (**This is Done**)

### Pre-reqs

At minimum, Harbor requires two dependecies:

1. Golang version 1.22.1
2. Node version 20

### Building `config/`

`config/` is the typescript definition for all of the configuration elements for Harbor. In order to build Harbor, first you need to build these.

1. First install yarn elements: `yarn i`
2. Then run `yarn build`

This will create a `dist/` folder in here that has the built JS that contains the constructs needed to configure harbor.

### Building the executable

1. cd into `harbor-runner/`
2. run `go work vendor`
3. run `go build -o harbor cmd/harbor/main.go` to build the executable
4. for windows, run this command: `go build -o harbor.exe cmd/harbor/main.go`

### Move the executable into your path

1. Just add the executable somewhere in your path

## Next steps

Read [Your First Harbor Package](./your-first-harbor-package.md)
