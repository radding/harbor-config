# Your First Harbor Package

Harbor refers to the most basic unit of a project as a "package". A package is a unit of code that has tasks associated with it (eg linting, running tests, releases, etc), potentionally produces artifacts (executables, docs, libraries, etc), and metadata.

## Defining package configuration

All configuration for harbor is held in a `.harborrc.ts` files. This is a normal typescript file, it uses your installed node version in order to render the configuration file. This means that you can use TS to check for the exsitance of files, make requests, etc.

This `.harborrc.ts` file should be held at the root of your project.

### Defining a package

In your `.harborrc.ts`, you must create a single Package construct. This serves as the root of your configuration. Everything else should be associated in one way or another on this package. Here is the most basic definition of a package:

```typescript
import { Package } from "harbor-config";

const pkg = new Package("harbor-docs");

export default pkg;
```

**NOTE:** you must export your package object as the default export from your file so that harbor can find it.

This creates a new package called `harbor-docs`.

#### Adding additional metadata

While the above example is simple, and gets you running, but what about metadata? How will people know what repository hosts this code? What about where they can get help?

You can pass a second argument to the Package constructor to pass in metadata. For example:

```typescript
const pkg = new Package("harbor-docs", {
    repository: "https://github.com/radding/harbor-config.git",
    path: "/docs",
    issues: "https://github.com/radding/harbor-config/issues",
    homepage: "https://github.com/radding/harbor-config"
});
```

See [Package Reference](/reference/package) for more information

## What is a package?

Harbor defines a package as a distinct collection of code that provides a unitary function. For example, a directory that produces a binary or makes a shared library are all examples of packages. In a polyrepo, a package will most likely correspond to a single repository, while in a monorepo, a package will most likely corespond to a single directory.

Packages consist of these major pieces exluding their code:

1. Artifacts
2. Tasks
3. Metadata
4. Dependencies

### Artifacts

Artifacts are the _output_ of the package. For a web app, this could be a docker image. For a CLI utility, this could be a binary. For a library, this could be the bundled assets. Artifacts are produced by Tasks.

### Tasks

Tasks are things that must be done in order to produce a package's artifacts. Examples of tasks are things like linting, testing, building, etc.

Tasks are also dependent on other tasks. For example, in order to build your package, you probably want to run testing. In order to run testing, you probably want to run linting.

### Metadata

Package metadata may not be important for an application to function, but its important for developers and humans to be able to find information about this package. Metadata includes pieces of information such as github repositories, package name, version, stability, where to get help, etc

### Dependencies

Dependencies are things that this package needs in order to produce some output. This could be something like a library, or a linting configurations.

## Next step

Read [Adding your first task](your-first-task.md)
