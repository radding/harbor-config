# Your First Task

Tasks are the underpinning of how Harbor understands how to run commands inside your packages. With out tasks, Harbor would not understand how to build your executable, or lint you project or run your tests. Harbor makes this simple to do.

## Defining a Task

Tasks are simple constructs in Harbor. They tell Harbor how to do something. There are many different types of tasks, the most basic of which are `ExecCommand`s. An exec command is a simple task that simply calls the OS's `exec` utility. This enables you to call an executable you have installed on your machine, For example:

```typescript
const build = new ExecCommand(pkg, "build", {
    executable: "./.venv/bin/mkdocs",
    args: [
        "build"
    ]
})
```

When this task is executed, Harbor will execute `./.venv/bin/mkdocs` and pass in `build` to the executable. This is how this site is built.

### Other tasks

There are a few built-ins for Harbor tasks, see [Config Reference](/reference/config)

### Running tasks

To run this task, all you have to do is run:

```sh
harbor run build
```

and look at that, you get an error.

The reason you get an error is because you have to register the task with the package. Registering a command is as easy as:

``` typescript
pkg.registerTask(build);
```

Now when you run `harbor run build`, harbor will actually execute your build step.

#### Why not register when you define the task?

Ideally, your package only needs a handful of entrypoints, but may need some complex pipelines to execute those entry points. A good example of this are setup tasks (more on this later). You don't want to overload your team mates with to many commands, so you really only want the commands that are actually useful to be registered.

Typically "helpful" commands are things like:

1. build
2. develop
3. deploy:staging
4. deploy:prod

Why polute this more commands that you probably won't use? **K**eep **I**t **S**imple **S**illy.

## Setup tasks

When you first clone a package, there are often configurations you need to set, dependencies you need to have, setting specific versions of toolchains, etc.

Harbor solves this problem by introducing the concept of setup tasks.

### Defining a setup task

A setup task is the same as a normal task, only instead of eventually calling `pkg.registerTask(...)`, you add it to another special construct: `PackageSetup`. This Construct informs Harbor that this command must be run before the package is "ready". For example:

```typescript
const installMKdocs = new ExecCommand(pkg, "install-mk-docs", {
    executable: "./.venv/bin/pip",
    args: [
        "install", 
        "mkdocs-material"
    ]
})

new PackageSetup(pkg, "setup", {
    actions: [
        installMKdocs
    ]
})
```

### Running Setup tasks

Setup tasks are run once before the first task is run. You don't have to run a special command to run them. For example, the first time you run `harbor run build`, Harbor will first check if needs to run setup, run it if it needs to, then it will run the `build` task.

However, sometimes you do have run setup periodically. For that, Harbor provides a command:

```sh
harbor setup
```

That will run setup if you need it. If you want to force it to run irregardless of its needed, run `harbor setup --force`

## Task Dependencies

Often, we need to have dependencies between tasks. For example, you often want to run linting before you test, and you want to test before building your app, and you want to build your app before you deploy it.

Harbor gives you the ability to model these dependencies in a couple of ways. Tasks have a `needs` and a `then` method on them.

`taskA.needs(taskB)` will tell Harbor that you need to run `taskB` before you can run `taskA`, while `taskA.then(taskB)` will do the opposite (run `taskA` before `taskB`).

For example, in the setup task from above, we don't know if `pip` is even installed at `./.venv/bin`. We need a task that sets this up before we run `installMkdocs`! Here is how this is modeled:

```typescript
const installCorrectVersion = new ExecCommand(pkg, "pyenv-version", {
    executable: "pyenv",
    args: [
        "install",
        "-s",
        "3.10.14",
    ]
}).then(new ExecCommand(pkg, "set-local", {
    executable: "pyenv",
    args: [
        "local",
        "3.10.14",
    ]
})).then(new ExecCommand(pkg, "create-env", {
    executable: "pyenv",
   args: [
    "exec",
    "python",
    "-m",
    "venv",
    ".venv",
   ]
})) as ExecCommand

actions.push(installCorrectVersion);

const pipe = new Pipeline(pkg, "setup-pipeline", actions);

const installMKdocs = new ExecCommand(pkg, "install-mk-docs", {
    executable: "./.venv/bin/pip",
    args: [
        "install", 
        "mkdocs-material"
    ]
}).needs(pipe)
```

**Note:** `Pipeline` is a different way to model the dependencies. This is analogous to iterating over the `actions` array and calling `.then` in order.

This tells harbor that to setup this package what you need to to do is:

1. Ensure that python version 3.10.14 is installed
2. Then make sure that `pyenv` sets the current directory's python version to 3.10.14
3. After that, we want to create a venv at `.venv`.
4. Finally after all of that is complete, we want install the actuall dependency. Now the package is ready!

## Conditional tasks

Sometimes, you want to change the task you are running, either due to the OS you are on, or the availibility of tools on the computer, or just want to add some logic. Since the `.harborrc.ts` file is literally just TypeScript, this is trivial. Just add an if statement:

```typescript
if (!executableIsAvailable("pyenv")) {
    const installPyenv = new ExecCommand(pkg, "pyenv", {
        executable: "brew",
        args: [
            "install",
            "pyenv",
        ]
    });
    actions.push(installPyenv);
}
```

Basically, if `pyenv` is installed on this computer, this task will be skipped.

## All together now

At this point, your `.harborrc.ts` file should look like this:

```typescript
import { ExecCommand, HarborConstruct, Package, PackageSetup } from "../config/dist";
import { Pipeline } from "../config/dist/Pipeline"

const pkg = new Package("harbor-docs", {
    repository: "https://github.com/radding/harbor-config.git",
    path: "/docs",
    issues: "https://github.com/radding/harbor-config/issues",
    homepage: "https://github.com/radding/harbor-config"
});

const actions: HarborConstruct[] = [];
if (!executableIsAvailable("pyenv")) {
    const installPyenv = new ExecCommand(pkg, "pyenv", {
        executable: "brew",
        args: [
            "install",
            "pyenv",
        ]
    });
    actions.push(installPyenv);
}

const installCorrectVersion = new ExecCommand(pkg, "pyenv-version", {
    executable: "pyenv",
    args: [
        "install",
        "-s",
        "3.10.14",
    ]
}).then(new ExecCommand(pkg, "set-local", {
    executable: "pyenv",
    args: [
        "local",
        "3.10.14",
    ]
})).then(new ExecCommand(pkg, "create-env", {
    executable: "pyenv",
   args: [
    "exec",
    "python",
    "-m",
    "venv",
    ".venv",
   ]
})) as ExecCommand

actions.push(installCorrectVersion);

const pipe = new Pipeline(pkg, "setup-pipeline", actions);

const installMKdocs = new ExecCommand(pkg, "install-mk-docs", {
    executable: "./.venv/bin/pip",
    args: [
        "install", 
        "mkdocs-material"
    ]
}).needs(pipe)

new PackageSetup(pkg, "setup", {
    actions: [
        installMKdocs
    ]
})

const build = new ExecCommand(pkg, "build", {
    executable: "./.venv/bin/mkdocs",
    args: [
        "build"
    ]
})

pkg.registerTask(build);

export default pkg;
```

## Next steps

Read [Your first Package Dependency](your-first-pkg-dep.md)
