import { Dependency, Package, PackageSetup, Plugin, Task } from "./src";

const pkg = new Package("Harbor-Config", {
	repository: "https://github.com/radding/harbor",
	version: "1.0.0",
    path: "config/"
});

