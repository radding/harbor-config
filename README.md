# Harbor: a workspace manager and runner [![Stability - Alpha](https://img.shields.io/badge/Stability-Alpha-CD4631)](#stability) [![Docs - Terrible](https://img.shields.io/badge/Docs-Terrible-CD4631)](#road-map)

Harbor is a workspace manager and task runner designed to make managing complex software a breeze. It works over monorepos, and will treat a poly repo like a monorepo. **This is currently a work in progress, see [roadmap](#roadmap)**

## Quick Start

### Install

**This package can not be installed from a package manager yet see [Building from native](#building-from-native)**.

### Your first `.harborrc.ts` file

At the root of your application, add a `.harborrc.ts` file. This file is how you configure your package. This is just a normal TS file, so you can use any native features of NodeJS here. Here is an example:

```typescript
import { ExecCommand, Package, PackageSetup } from "harbor-config";

const pkg = new Package("harbor-config", {
 repository: "https://github.com/radding/harbor",
 version: "1.0.0",
 path: "config/"
});

const install = new ExecCommand(pkg, "install", {
 executable: "yarn",
 args: [
  "install",
 ]
})

const buildCommand = new ExecCommand(pkg, "build", {
 executable: "yarn",
 args: [
  "build"
 ]
})

pkg.registerTask(buildCommand)

new PackageSetup(pkg, "setup", {
 actions: [
  install
 ]
})

export default pkg;
```

#### Breaking this down

##### `new Package`

This is the entrypoint to your application. This adds metadata for your application, collects all of the information for the application and is the root of what is called the construct tree. You can hook into this to add application and repo infromation to a service catalog, provide repo metadata to other release applications, etc.

##### `new ExecCommand`

This adds an executable target. This is what is known as a construct ([see Constructs Programming model](https://www.improving.com/thoughts/infrastructure-as-actual-code-an-introduction-to-the-construct/)). Other constructs can take this as a dependency, you can register a new command for harbor, etc. [See Built in Constructs](#built-in-constructs)

#### `pkg.registerTask(buildCommand)`

This creates a task for harbor to run. you can execute this command by running `harbor run build`. The command is going to be the `id` of the Construct.

#### `new PackageSetup`

This creates a set up task. A setup task is a way for an package to be setup periodically, such as package installing, setting versions, etc. This one setups an install to install all yarn depenedencies. Setups will be done before either the first time you run a command or if you run `harbor setup`.

#### `export default pkg`

This is required so that the harbor runner can have access to your application definition.

## Building from native

1. Clone this repo
2. go into `./config` and run `yarn i && yarn build`
3. go into `./harbor-runner` and run `go work vendor && go build -o harbor cmd/harbor/main.go`
4. You now have a bootstrapped harbor executable

## Built in Constructs

### Package

This is the base of the harbor package. This defines meta data around the package so harbor can make informed decisions around how to manage your project. Most of these don't do anything right now, but can be used to do things like add the package to a service catalog or publish the module to something like NPM.

```typescript
type PackageOptions = {
    // Meta for where the harbor cahce should operate
    meta: {
        harborPackageDirectory: string;
    };
    // Where the source code is hosted for this project
    repository: string;
    // The name of the package, if not defined, takes the name from first argument in the constructor
    name?: string | undefined;
    // If this is is a mono repo, the path where this project exsists. 
    path?: string | undefined;
    // Home page if applicable for the docs.
    homepage?: string | undefined;
    description?: string | undefined;
    // Where issues should be reported to.
    issues?: string | undefined;
    // What licenes this package is licensed under
    license?: string | undefined;
    // The current version of this package.
    version?: string | undefined;
    // How stable is this package (aka am I a mad man if I use this)
    stability?: "Beta" | "Generally Available" | "End of Life" | "Alpha" | "Pre-Alpha" | undefined;
    // Where to store artifacts of this package
    artifactsLocation?: string | undefined;
}
```

### HarborConstruct

This is the base of all constructs (except for `Package`). This provides all of the base functionality needed to build an custom construct. A `HarborConstruct` consists of a `kind` and `options`. The `kind` tells the harbor CLI what executor to use to execute this construct. All of the other builtin constructs will set the `kind` for you. `options` will be passed in as a JSON object into the executor so that it could be configured properly.

### ExecCommand (kind: `harbor.dev/ExecCommand`)

This is a simple executor that just calls `os.Exec` on your machine. This is not the same as a bash script or other such shell script, this literally calls `exec` on what ever the executable is.

```typescript
// Options for the exec command construct
type ExecCommandOpts = {
 // The executable to execute. 
 executable: string;
 // The arguments to pass to the executable
 args: string[];
 // Any environment variables you want to pass to the executable. 
 env?: Record<string, string> | typeof process.env;
 // Not used currently
 inputs?: string[]
}
```

### Dependencies

There are two dependencies: `RemoteDependency` and `LocalDependency`. Local dependencies are meant to be used in monorepos where the dependency is in the same repo. Remote dependencies are other packages that are hosted outside of this repository.

#### `Dependency`

This is generic interface for both Local and Remote Dependency. It provides a single method:

```typescript
export interface Dependency<T extends DependencyOptions = DependencyOptions> extends HarborConstruct{
  task(name: string): ITask;
}

```

`task` adds an entry for Harbor to go and execute a task inside that dependency. This allows you to wait on the results of another command inorder perform your command.

#### `LocalDependency`

```typescript
{
    name?: string; // The name of the dependecy if you want to give it one.
    path: string // where the dependency is located relative to this file.
    artifacts // not used yet.
}

// Most basic example:
new LocalDependency(pkg, "../config");
```

#### `RemoteDependecy`

**TODO**

### PackageSetup

These are commands that should be run on package setup. After `harbor setup` is run, these commands will be run. It is expected that everything needed to start working on the package will be done after the setup steps are done.

```typescript
type PackageSetupOpts = {
	actions: IConstruct[] // the actions to run at package setup
}
```

### Plugin

This is still in **TODO**. this is meant to define the plugins that this package needs to have to operate.

### RemoteResource

This is still in **TODO**. this is meant to be a resources hosted outside the workspace/repo that needs to be pulled in inorder for something in the package to work.

### Repository

This is still in **TODO**.

## Stability

### Alpha

Do not use this in real profesional projects. This is not stable at all. APIs and constructs may change with out warning. I am still shoring up some things, and things will break.

### Beta

You can use this. APIs and constructs are stable for the most part, there could be some weird bugs.

### Generally available

Use it. Everything is stable. Minimum bugs.

## Road Map

[See Road Map](https://github.com/users/radding/projects/1/views/1)

### Known bugs

[See Issues](https://github.com/radding/harbor-config/issues)
