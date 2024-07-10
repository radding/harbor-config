import { Dependency, Package, PackageSetup, Plugin, Task } from "../src";

const pkg = new Package("Harbor-Test", {
	repository: "https://github.com/radding/harbor",
	version: "1.0.0",
});

const shell = new Plugin(pkg, "shell");
new Plugin(pkg, "nodejs");
new Dependency(pkg, "https://github.com/radding/something")

new PackageSetup(pkg, "ensure_node_env", {
	actions: []
});

new Task(pkg, "build", {
	plugin: shell,
	options: {
		run: "ls -al",
		env: {},
	},
	artifacts: [],
})

// console.log(JSON.stringify(pkg.synth(), null, 2));
pkg.synth(true);
const json: any = pkg.createTree();

// END ACTUAL CODE


type GraphActions = {
	preOrder: (id: string, data: any) => void;
	postOrder: (id: string, data: any) => void;
}

const dfs = (id: string, actions: Partial<GraphActions>) => {
	const data = json.constructs[id];
	actions.preOrder && actions.preOrder(id, data);
	try {
		data.dependsOn.forEach((id: string) => {
			dfs(id, actions);
		});
	} catch (e) {
		console.log("failed to dfs", id);
	}
	actions.postOrder && actions.postOrder(id, data);
};

const topoSetup: string[] = [];
const added = new Set();
json.setup.forEach((id: string) => {
	dfs(id, {
		postOrder: (id, data) => {
			if (added.has(id)) {
				return;
			}
			added.add(id)
			topoSetup.push(id)
		}
	})
});

console.log(topoSetup);
