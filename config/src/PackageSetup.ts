import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";

type PackageSetupOpts = {
	actions: HarborConstruct[]
}

export class PackageSetup extends HarborConstruct {
	constructor(scope: Construct, id: string, opts: PackageSetupOpts) {
		super(scope, id, {
			kind: "harbor.dev/PackageSetup",
			options: {},
		});

		this.node.addDependency(...opts.actions);
		this.package.addSetup(this.node);
	}
}