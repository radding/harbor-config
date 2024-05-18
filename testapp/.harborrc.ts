import { Dependency, Package, PackageSetup, Plugin } from "../src";
import { Command } from "../src/Command";

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

new Command(pkg, "build", {
	plugin: shell,
	options: {
		run: "ls -al",
		env: {},
	},
})

// console.log(JSON.stringify(pkg.synth(), null, 2));
const json: any = pkg.synth();

type GraphActions = {
	preOrder: (id: string, data: any) => void;
	postOrder: (id: string, data: any) => void;
}

const dfs = (id: string, actions: Partial<GraphActions>) => {
	const data = json.constructs[id];
	actions.preOrder && actions.preOrder(id, data);
	try {
		data.dependsOn.forEach(id => {
			dfs(id, actions);
		});
	} catch (e) {
		console.log("failed to dfs", id);
	}
	actions.postOrder && actions.postOrder(id, data);
};

const topoSetup: string[] = [];
const added = new Set();
json.setup.forEach(id => {
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
