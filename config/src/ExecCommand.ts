import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";
import { ITask } from "./Task";

type ExecCommandOpts = {
	executable: string;
	args: string[];
	env?: Record<string, string> | typeof process.env;
	inputs?: string[]
}

export class ExecCommand extends HarborConstruct implements ITask {
	constructor(scope: Construct, public readonly id: string, opts: ExecCommandOpts) {
		super(scope, id, {
			kind: "harbor.dev/ExecCommand",
			options: {
				...opts,
			},
		});
	}

	needs(...deps: Construct[]): ITask {
		deps.forEach(dep => this.node.addDependency(dep));
		return this;
	}

	public then(construct: ITask): ITask {
		construct.node.addDependency(this);
		return construct;
	}
}