import { Construct, IConstruct, Node } from "constructs";
import _ from "lodash";
import * as path from "path"
import { HarborConstruct } from "./HarborConstruct";
import { z } from "zod";

// export interface PackageOptions {
// 	/**
// 	 * The directory where harbor caches resources related to this package
// 	 * 
// 	 * @default $PWD/.harbor
// 	 */
// 	harborPackageDirectory: string;
// }

const PackageOptions = z.object({
	meta: z.object({
		harborPackageDirectory: z.string(),
	}),

	repository: z.string(),
	homepage: z.string().url().optional(),
	description: z.string().optional(),
	issues: z.string().url().optional(),
	license: z.string().url().optional(),
	version: z.string().optional(),
	stability: z.enum(["Beta", "Generally Available", "End of Life", "Alpha", "Pre-Alpha"]).optional(),
	artifactsLocation: z.string().url().optional(),
});

export type PackageOptions = z.infer<typeof PackageOptions>

const defaultOptions: Pick<PackageOptions, "meta"> = {
	meta: {
		harborPackageDirectory: path.join(path.dirname(require.main?.filename!), ".harbor/"),
	},
}

const pkgSymbol = Symbol.for("Package");

export class Package extends Construct {
	public readonly __kind = pkgSymbol;

	public readonly location: string;
	private readonly commands: Record<string, string> = {};
	private readonly setup: string[] = [];
	public readonly packageInfo: Omit<PackageOptions, "meta">

	constructor(name: string, opts: Partial<PackageOptions> = {}) {
		super(null as any, name);
		const options = PackageOptions.parse(_.merge(defaultOptions, opts));
		const { meta, ...rest } = options;
		this.location = meta.harborPackageDirectory;
		this.packageInfo = rest;
	}

	/**
	 * static isconstruct: IConstruct: construct is Package
	 */
	public static is(construct: IConstruct): construct is Package {
		return Object.prototype.hasOwnProperty.call(construct, "__kind") && (construct as any).__kind === pkgSymbol;
	}

	/**
	 * static of(construct: IConstruct) finds the root package
	 */
	public static of(construct: IConstruct | undefined): Package {
		if (construct === undefined) {
			throw new Error("Could not find package");

		}
		if (Package.is(construct)) {
			return construct;
		}
		return Package.of(construct.node.scope);
	}

	/**
	 * addCommand
	 */
	public addCommand(commandName: string, commandNode: Node) {
		this.commands[commandName] = commandNode.path;
	}

	public addSetup(node: Node) {
		this.setup.push(node.path);
	}

	synth(): object {
		const constructs = this.node.children.filter(HarborConstruct.of).reduce((acc, child) => {
			return {
				...acc,
				...child.synth(),
			}
		}, {});
		return {
			constructs,
			commands: this.commands,
			setup: this.setup,
			packageInfo: this.packageInfo,
		}
	}
}