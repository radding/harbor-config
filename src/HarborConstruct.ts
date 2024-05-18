import { Construct, IConstruct } from "constructs";
import { Package } from "./package";

export interface HarborConstructOptions {
	kind: string;
	options: any;
}

const harborConstructSym = Symbol.for("HarborConstruct");

export class HarborConstruct extends Construct {
	public readonly __kind = harborConstructSym;
	protected readonly package: Package;

	constructor(parent: Construct, id: string, protected readonly json: HarborConstructOptions) {
		super(parent, id);

		this.package = Package.of(this);
	}

	public static of(construct: any): construct is HarborConstruct {
		return Object.prototype.hasOwnProperty.call(construct, "__kind") && construct.__kind === harborConstructSym;
	}

	synth(): any {
		const dependencies = [...this.node.dependencies].map(dep => dep.node.path)

		const baseJson = {
			[this.node.path]: {
				kind: this.json.kind,
				options: this.json.options,
				dependsOn: dependencies,
			}
		};

		return this.node.children.filter(HarborConstruct.of).reduce((acc, child) => {
			return {
				...acc,
				...child.synth(),
			}
		}, baseJson)
	}
}