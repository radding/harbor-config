import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";

type ExecCommandOpts = {
	executable: string;
	args: string[];
	env?: Record<string, string> | typeof process.env;
}

export class ExecCommand extends HarborConstruct {
	constructor(scope: Construct, id: string, opts: ExecCommandOpts) {
		super(scope, id, {
			kind: "harbor.dev/ExecCommand",
			options: {
				...opts,
			},
		});
	}
}