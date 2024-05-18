import { Construct, IConstruct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";

type CommandOpts = {
	plugin: IConstruct;
	name?: string;
	options: any;
	dependencies?: IConstruct[],
}

export class Command extends HarborConstruct {
	constructor(scope: Construct, id: string, opts: CommandOpts) {
		super(scope, id, {
			kind: "Harbor.dev/Command",
			options: {
				...opts.options,
			},
		});

		const commandName = opts.name ?? id;

		this.node.addDependency(opts.plugin);
		this.node.addDependency(...opts.dependencies ?? []);

		this.package.addCommand(commandName, this.node);
	}
}