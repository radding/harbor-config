# Harbor's Architecture

Harbor has two components: the configuration and the executor. The configuration is written in TS and utilizes the Construct based programming model, while the executor takes the output of the configurations and builds out task trees and setup tasks in order to actually execute the project.

## Configuration

Harbor Configuration is written in Typescript, and uses [the construct programming model](https://www.improving.com/thoughts/infrastructure-as-actual-code-an-introduction-to-the-construct/) that AWS CDK pioneered. This model is also used in adjacent tools such as [projen](https://projen.io/docs/introduction/), [CDK8s](https://cdk8s.io/), [CDKTF](https://developer.hashicorp.com/terraform/cdktf), and others.

Each declared Construct inside of your configuration file (.harborrc.ts), is added onto the construct tree. The construct Tree starts at the root (the [`Package`](../Reference/Package.md) construct), and goes down to the leaf nodes. Nodes in the construct tree could be executors with dependencies, a reference to an executor, or a no-op.

For example, a node could be a custom construct built to always add some plugins, configuration, or other logic that is re-usable. Here is how I would make a simple nodejs consturct node:

```typescript
import { Package, PackageSetup, ExecCommand, PackageOptions } from "harbor-config";

export type NodeJSPackageOpts = {
    manager: "npm" | "yarn" | "pnpmp"
} & Partial<PackageOptions>

export class NodeJSPackage extends Package {

    constructor(name: string, opts: NodeJSPackageOpts) {
        super(name, opts);

        const install = new ExecCommand(this, "install-packages", {
            executable: opts.manager,
            args: [
                "i"
            ]
        });

        new PackageSetup(this, "setup-node", {
            actions: [install]
        });

        const packageJSON = require("package.json")

        Object.entries(packageJSON.scripts).forEach(([name, command]) => {
            const commandParts = command.split(" ");
            const task = new ExecCommand(this, name, {
                executable: commandParts[0],
                args: commandParts.slice(1),
            })
            this.registerTask(task);
        });
    }
}
```

**NOTE**: the above code is hypothetical and not tested, don't blindly copy and paste it.

This code will setup a node JS package for you. It will add the install step on setup so you don't have to. This will also read the package.json file and add all of its scripts into Harbor so harbor could call them for you.

This powerful arrangement allows you to provide really powerful abstracts and setup without needing to create new plugins. You can create Configuration primatives that combine together to provide the functionality you would like.

### Transforming into useable configuration

When harbor is loaded, harbor will take the configuration file and execute it. This execution will then result in a Construct tree. This Construct tree is then converted into a JSON file. The JSON file will have four root keys:

1. Constructs
2. Tasks
3. Setup
4. packageInfo

#### Constructs

This is probably the most important part of the config file. This is the construct tree. This is a key value object where the Key is node id of the construct and the value is the config of the Construct.

The Config will always have three keys:

1. Kind
2. Options
3. dependsOn

Kind tells Harbor which Executor to use to execute this Construct. This should be a URI, and look something like: `harbor.dev/ExecCommand`

Options is an any type that is spececif to each executor. This is basically Executor specific configuration. For example, the `ExecCommand` construct will have options that look like this:

```json
{
    "executable": "pyenv",
    "args": [
        "exec",
        "python",
        "-m",
        "venv",
        ".venv"
    ]
}
```

DependOn is a list of node ids that this construct depends on. Basically these nodes must be executed first before this node can be executed. This is how Harbor can construct it's trees. _**NOTE:** These are not in any particular order, and harbor will execute all of these nodes in parallel if it can._

#### Tasks

Tasks are a simple key value store where the key is the command that gets executed via `harbor run <command>`, and the value is the node id of the first construct in the commands run tree.

#### Setup

Set up is a simple list, where each entry in the list is a node id pointing to the first construct in the tree that must be run before the package is ready. _**NOTE:** These are not in any particular order, and harbor will execute all of these nodes in parallel if it can._

#### PackageInfo

This is just human readable package info. Harbor doesn't really use this information yet. We could use this information to open up PRs via the CLI if we would like.

## The Executor

Once the configuration is executed and the JSON constructed, Harbor will then take that JSON and figure out how to execute certain tasks.

### Constructing the Task Trees

It will first create a setup task tree. The root node is not transparent to users. Its essentiall a noop node whos only job is to collect all of the setup tasks defined in `setup` in the config.

After the setup task tree is built, harbor will then create a task tree for each of the tasks defined in `tasks`. The entry point of these tasks are the constructs defined as the entry point for the task. Harbor then iterates over each object in `dependsOn` and creates a node for that tree.

### Executing the Task Trees

Once the Task tree is complete, Harbor starts are the root, then performs a post-order Depth First Search to execute each of its children before executing its self. The children are executed in parallel.

## Caching

In order to provide fast execution, Harbor caches logs, artifacts and others inside of the `.harbor` directory (you should `.gitignore` this file). The cache Key presently the sha 256 hash of the `.harborrc.ts` file. The Cacher will create a directory inside of `.harbor` were the name is the hash. Inside of this directory, the cached execution of `.harborrc.ts` is stored, and each executed task will get its own directory inside of it to store artifacts.

### My comentary on caching

Caching is currently not great. I would like to revisit this at somepoint. Two problems with my current approach:

1. The cache key is based off of the result of configuration. This breaks `if` statements inside the config file.
2. This ignores changes to source files that should trigger re-runs
