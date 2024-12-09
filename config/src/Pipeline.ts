import { Construct, IConstruct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";

export class Pipeline extends HarborConstruct {
	constructor(scope: Construct, id: string, actions: IConstruct[]) {
		super(scope, id, {
			kind: "harbor.dev/noop",
			options: {},
		});
		actions.reduce<IConstruct | null>((acc, construct) => {
			if (acc !== null) {
				construct.node.addDependency(acc);
			}
			return construct;
		}, null);
		this.node.addDependency(actions[0]);
	}
} 