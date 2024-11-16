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

export interface ITask extends IConstruct {
	needs(...dep: Construct[]): ITask;
	then(next: ITask): ITask;
	readonly id: string;
}

export class Task extends HarborConstruct implements ITask{
	constructor(scope: Construct, public readonly id: string, opts: TaskOpts) {
		let inputs = opts.inputs ?? [];
		if (typeof inputs === "function") {
			inputs = inputs();
		}
		super(scope, id, {
			kind: "Harbor.dev/Task",
			options: {
				plugin: opts.plugin.node.id,
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
	then(next: ITask): ITask {
		next.node.addDependency(this);
		return next;
	}
	
	public needs(...dependency: Construct[]): ITask {
		this.node.addDependency(...dependency);
		return this;
	}
}