# Your first Package Dependency

A package dependency is another, seperate package that your current package needs in order to work properly. This could be a library, a build tool, what ever.

Harbor has two constructs to model dependencies: a `LocalDependency` and a `RemoteDependency`

## Local Dependecies

A local dependency is probably the easiest dependecy to understand. This is something that is assumed to live along side the current package. In a monorepo, almost all of the dependencies will be Local Dependencies. These dependencies must exsist next to the application _before_ harbor runs. Harbor will not create or clone code with local dependencies.

### Example

```Typescript
const config = new LocalDependency(pkg, "../config");
```

## Remote Depenencies (STILL A WORK IN PROGRESS)

In a poly repo, most dependencies will be modeled as an Remote Dependency. Remote Dependencies are dependencies that live outside of the scope of the current harbor workspace/package. If the Dependency is cloned to the current package, Harbor will treat it as a local dependency.

A Remote Dependency has two flavors: an optional and required dependency. An Optional Dependency is a dependency where you only need the artifact, and not the code in order to operate, while a required dependency's repository will always been cloned down before Harbor runs any commands.

### Optional Dependencies

Optional Dependencies are useful in cases where you need the artifact but don't really care about the code of the dependency. An example would be a Node Module, or a pip package. However, sometimes having a deeper connection to the dependency's code is important. For example, if you are solving a bug in a dependency that is impacting your code, or you need to fork a dependency.

When these cases arise, you can call `harbor dock <dependency name>` and harbor will clone the dependency and host it along side your package. Then, you can just define your tasks as being dependent on the dependeny's task. After this, harbor will treat the remote dependency as a local dependency.

### Required Dependencies

Required Dependencies are dependencies that must be cloned durring package setup. After the dependency is cloned, harbor will then treat the dependency as a local depedency.

## Taking Dependencies on Tasks from Package Dependencies

After you declare a dependency, you'll often want to say "ensure that the dependency is built before building my package". You can do this pretty easily by simply doing:

```typescript
const configTask = config.task("build");

const goBuild = new ExecCommand(pkg, "go-build", {
    executable: "go",
    args: [
        "build",
        "-o",
        "harbor",
        "./cmd/harbor/main.go"
    ]
}).needs(
    configTask,
    tests,
)
```

This will ensure that the build process for `config` is completed before build of this project is done.

## Next Steps

Read our [Reference](../Reference/introduction.md)
