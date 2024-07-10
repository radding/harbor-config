import { Construct, IConstruct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";

export type TaskOpts = {
	plugin: IConstruct;
	name?: string;
	options: any;
	dependencies?: IConstruct[],
	artifacts?: string[],
	inputs?: string | string[] | (() => (string | string[]))
}

export class Task extends HarborConstruct {
	constructor(scope: Construct, id: string, opts: TaskOpts) {
		let inputs = opts.inputs ?? [];
		if (typeof inputs === "function") {
			inputs = inputs();
		}
		super(scope, id, {
			kind: "Harbor.dev/Task",
			options: {
				...opts.options,
				artifacts: opts.artifacts,
				inputs,
			},
		});

		const taskName = opts.name ?? id;

		this.node.addDependency(opts.plugin);
		this.node.addDependency(...opts.dependencies ?? []);

		this.package.addTask(taskName, this.node);
	}
	
	public needs(dependency: Task) {
		this.node.addDependency(dependency);
	}
}