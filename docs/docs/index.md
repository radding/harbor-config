# Welcome to Harbor

Harbor is the real programmable build chain built for simplicity and speed. Unlike other build chains, Harbor uses a well known programming langauge to drive its configurations (TypeScript), reducing cognitive overhead and context switching. Unlike other simpler tools, Harbor is also language agnostic, working on any kind of language you would like.

Harbor is designed from the groundup to be simple, avoid using magic, and just feel like normally programing a real program.

## Comparisons

### Bazel

Harbor does not build its own configuration language. If you can write it in Typescript, then you can write it in harbor. On top of that, making plugins is simpler than in Bazel and you can make plugins in any language you would like. Harbor is also designed outside of the traditional Monorepo mold, so it supports much more than Monorepos.

### Turbo

Like Turborepo, Harbor agressively caches and attempts to parralelize tasks as much as possible. This leads to a very fast and effecient run for all of your commands.

Unlike Turbo, Harbor does not rely on the Node ecosystem so you don't need to add a package.json files everywhere to support turbo. You also can script your build and tool chains so its smarter than just a JSON file.

### Lerna

Lerna's eco system is very agressively geared towards the node ecosystem, relying on package.json. Harbor is again outside of the node ecosystem, supporting a lot more than package.json files. Everything is scriptable.

### NX

While NX could be used outside of the Node eco system, its not often. NX relies on json files to hold its configuration, which makes it less flexible than Harbor.

### Pants

While Pants uses a similar approach to Bazel. In order to use Pants, you have to learn yet another configuration language. Or you could just use Typescript.

## Project Status

Harbor is still very early in its development. You can bootstrap a harbor project and have it work ([The repository that hosts these docs are a harbor project](https://github.com/radding/harbor-config)), but the ergonomics are still being worked out.

Write now you can basically declare depenedencies and call the `exec` command inside your application, which honestly will enable you to do almost everything, but with thorny experience.

**Some of these docs talk as if Harbor is perfect today, this is not the case!** I am writting these docs to mainly organize my thoughts around harbor. I will call out what things are and aren't read to be used.

You can checkout out road map here: [https://github.com/users/radding/projects/1](https://github.com/users/radding/projects/1)

## Getting help

If you have some problems with Harbor, feel free to [open an issue](https://github.com/radding/harbor-config/issues/new). Make sure that you label the issue properly!
