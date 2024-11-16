import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";
import { z } from "zod";
import { ITask, Task } from "./Task";
import { Pipeline } from "./Pipeline";

const DependencyOptions = z.object({
	url: z.string().url(),
	name: z.string().optional(),
	localPath: z.string().optional(),
	path: z.string().optional(),
	artifacts: z.string().url().optional(),
})

type DependencyOptions = z.infer<typeof DependencyOptions>;


export class Dependency extends HarborConstruct {
	private depName: string;
	constructor(scope: Construct, idOrRepo: string, opts: Partial<DependencyOptions> = {} as any) {
		const options = DependencyOptions.parse({
			url: idOrRepo,
			...opts,
		})
		super(scope, idOrRepo, {
			kind: "harbor.dev/Dependency",
			options,
		});

		this.depName = options.name ?? options.url;
	}

	task(name: string): ITask {
		return new Task(this, `dep-${this.depName}-${name}`, {
			plugin: this.package.remoteExcecutor,
			dependencies: [this],
			options: {
				run: name,
			}
		});
	}
}