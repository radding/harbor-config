import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";
import z from "zod";

export const CredentialsConfig = z.discriminatedUnion("type", [
	z.object({
		type: z.literal("process"),
		name: z.string(),
	}),
	z.object({
		type: z.literal("env"),
		userNameVar: z.string(),
		password: z.string(),
	}),
	z.object({
		type: z.literal("key"),
		path: z.string(),
	}),
]);
export type CredentialsConfig = z.infer<typeof CredentialsConfig>;

const RemoteResourceOpts = z.object({
	url: z.string().url().optional(),
	destination: z.string().optional(),
	name: z.string().optional(),
	credentials: CredentialsConfig.optional(),
	addToGlobalCache: z.boolean().optional().default(true),
	postDownload: z.instanceof(Construct).optional(),
	preDownload: z.instanceof(Construct).optional(),
})

export type RemoteResourceOpts = typeof RemoteResourceOpts["_input"]

export class RemoteResource extends HarborConstruct {
	constructor(scope: Construct, idOrLocation: string, opts?: RemoteResourceOpts) {
		opts = RemoteResourceOpts.parse(opts);
		super(scope, idOrLocation, {
			kind: "harbor.dev/RemoteResource",
			options: {
				credentials: opts?.credentials,
				destination: opts?.destination,
				url: opts?.url ?? idOrLocation,
				addToGlobalCache: opts?.addToGlobalCache ?? true,
				name: opts?.name ?? idOrLocation,
			}
		});
		this.json.options.destination = opts?.destination || this.package.location;

		if (opts?.preDownload) {
			this.node.addDependency(opts.preDownload);
		}

		if (opts?.postDownload) {
			opts.postDownload.node.addDependency(this);
		}

	}
} 