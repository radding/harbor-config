import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";
import { z } from "zod";
import { ITask, RemoteTask, Task } from "./Task";
import fs from "fs";
import path from "path"

const LocalDependencyOpts =  z.object({
	name: z.string().optional(),
	path: z.string(),
	artifacts: z.string().url().optional(),
});
type LocalDependencyOpts = z.infer<typeof LocalDependencyOpts>

const RemoteDependencyOpts = z.object({
	url: z.string().url(),
	name: z.string().optional(),
	localPath: z.string().optional(),
	path: z.string().optional(),
	artifacts: z.string().url().optional(),
});

type RemoteDependencyOpts = z.infer<typeof RemoteDependencyOpts>


type DependencyOptions = RemoteDependencyOpts | LocalDependencyOpts;

export interface Dependency<T extends DependencyOptions = DependencyOptions> extends HarborConstruct{
	readonly options: T;
	task(name: string): ITask;
}

export class RemoteDependency extends HarborConstruct implements Dependency<RemoteDependencyOpts>{
	private depName: string;
	public readonly options: RemoteDependencyOpts;

	constructor(scope: Construct, idOrRepo: string, opts: Partial<RemoteDependencyOpts> = {} as any) {
		const options = RemoteDependencyOpts.parse({
			url: idOrRepo,
			...opts,
		})
		super(scope, idOrRepo, {
			kind: "harbor.dev/Dependency",
			options,
		});
		this.options = options;
		this.depName = options.name ?? options.url;
	}

	task(name: string): ITask {
		return new RemoteTask(this, `dep-${this.depName}-${name}`, {
			dependency: this,
			taskName: name,
			options: {},
			isDepenedencyLocal: false
		});
	}
}

export class LocalDependency extends HarborConstruct implements Dependency<LocalDependencyOpts> {
	private depName: string;
	public readonly options: LocalDependencyOpts;

	constructor(scope: Construct, idOrPath: string, opts: Partial<LocalDependencyOpts> = {} as any) {

		const options = LocalDependencyOpts.parse({
			path: idOrPath,
			...opts,
		})
		super(scope, idOrPath, {
			kind: "harbor.dev/LocalDependency",
			options,
		});
		console.log(this.package.root);
		const pathStr = path.join(this.package.root, options.path)
		try {
			const info = fs.statSync(pathStr);
			console.log(pathStr)
			if (!info.isDirectory()) {
				throw new Error(`local dependency does not exsist at path ${pathStr}!`)
			}
		} catch (e) {
			throw new Error(`error stating at path: ${pathStr}!, error: ${e}`)
		}
		this.options = options;
		this.depName = options.name ?? (options as any).url ?? options.path;
	}

	task(name: string): ITask {
		return new RemoteTask(this, `dep-${this.depName}-${name}`, {
			dependency: this,
			taskName: name,
			options: {},
			isDepenedencyLocal: true
		});
	}
}