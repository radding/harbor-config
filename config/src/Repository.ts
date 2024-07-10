import { z } from "zod";
import { HarborConstruct } from "./HarborConstruct";
import { Construct } from "constructs";

const HarborRepositoryProps = z.object({
	protocol: z.string().optional().default("git"),
	url: z.string().url(),
});

export type HarborRepositoryProps = Partial<typeof HarborRepositoryProps["_input"]>

export class Repository extends HarborConstruct {
	constructor(scope: Construct, nameOrUrl: string, options: HarborRepositoryProps = {}) {
		options = HarborRepositoryProps.parse({
			url: nameOrUrl,
			...options,
		});
		super(scope, nameOrUrl, {
			kind: "harbor.dev/Repository",
			options,
		});
	}
}