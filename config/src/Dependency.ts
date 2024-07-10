import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";
import { z } from "zod";

const DependencyOptions = z.object({
	url: z.string().url(),
	name: z.string().optional(),
	artifacts: z.string().url().optional(),
})

type DependencyOptions = z.infer<typeof DependencyOptions>;

export class Dependency extends HarborConstruct {
	constructor(scope: Construct, idOrRepo: string, options: DependencyOptions = {} as any) {
		options = DependencyOptions.parse({
			// @ts-expect-error
			url: idOrRepo,
			...options,
		})
		super(scope, idOrRepo, {
			kind: "harbor.dev/Dependency",
			options,
		});
	}
}