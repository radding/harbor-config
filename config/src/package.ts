import { Construct, IConstruct, Node } from "constructs";
import _ from "lodash";
import * as path from "path"
import { HarborConstruct } from "./HarborConstruct";
import { z } from "zod";
import * as fs from "fs";
import { PackageSetup } from "./PackageSetup";
import { ITask } from "./Task";
import crypto from "crypto"

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
	name: z.string().optional(),
	repository: z.string(),
	path: z.string().optional(),
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
		harborPackageDirectory: path.join(path.dirname(require.main?.filename ?? process.env.HARBORJS_HARBOR_LOC ?? ""), ".harbor/"),
	},
}

class RemoteExecutorPlugin extends HarborConstruct {
	constructor(scope: Construct, id: string) {
		super(scope, id, {
			kind: "harbor.dev/RemoteExecutor",
			options: {},
		});
	}
}

const pkgSymbol = Symbol.for("Package");

export class Package extends Construct {
	public readonly __kind = pkgSymbol;

	public readonly location: string;
	private readonly tasks: Record<string, string> = {};
	private readonly setup: string[] = [];
	public readonly packageInfo: Omit<PackageOptions, "meta">
	public readonly remoteExcecutor: IConstruct;

	constructor(name: string, opts: Partial<PackageOptions> = {}) {
		super(null as any, name);
		const options = PackageOptions.parse(_.merge(defaultOptions, opts));
		const { meta, ...rest } = options;
		this.location = meta.harborPackageDirectory;
		this.packageInfo = rest;
		this.packageInfo.name = this.packageInfo.name ?? name;
		this.remoteExcecutor = new RemoteExecutorPlugin(this, "remote-executor");
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
	 *This returns the root of the package, aka where the .harborrc.ts file is found
	 *
	 * @readonly
	 * @type {string}
	 * @memberof Package
	 */
	public get root(): string {
		return path.dirname(this.location);
	}

	/**
	 * addCommand
	 */
	public addTask(taskName: string, taskNode: Node) {
		this.tasks[taskName] = taskNode.path;
	}

	public registerTask(task: ITask) {
		this.addTask(task.id, task.node);
	}

	public addSetup(node: Node) {
		this.setup.push(node.path);
	}

	createTree(): object {
		const constructs = this.node.children.filter(HarborConstruct.of).reduce((acc, child) => {
			return {
				...acc,
				...child.synth(),
			}
		}, {});
		return {
			constructs,
			tasks: this.tasks,
			setup: this.setup,
			packageInfo: this.packageInfo
		}

	}

	synth(prettyFormat: boolean = false): void {
		const tree = this.createTree();
		const content = fs.readFileSync(require.main?.filename ?? process.env.HARBORJS_HARBOR_LOC ?? "")
		const hash = crypto.createHash("sha256").update(content).digest("hex").toString();
		fs.mkdirSync(path.join(this.location, hash), { recursive: true})
		fs.writeFileSync(path.join(this.location, hash,`config.json`), JSON.stringify(tree, undefined, prettyFormat ? 2 : undefined));
	}
}