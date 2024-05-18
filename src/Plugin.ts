import { Construct } from "constructs";
import { HarborConstruct } from "./HarborConstruct";
import { CredentialsConfig, RemoteResource, RemoteResourceOpts } from "./RemoteResource";
import path from "path";
import { ExecCommand } from "./ExecCommand";
import z from "zod";
import { PackageSetup } from "./PackageSetup";
import { Pipeline } from "./Pipeline";

const PluginOpts = z.discriminatedUnion("defaultRepository", [
	z.object({
		repository: z.string().url(),
		credentials: CredentialsConfig.optional(),
		defaultRepository: z.literal(false),
		name: z.string(),
	}),
	z.object({
		defaultRepository: z.literal(true),
		credentials: CredentialsConfig.optional(),
		name: z.string(),
	})
]).default({
	defaultRepository: true,
	name: "",
})

export type PluginOpts = z.infer<typeof PluginOpts>;

export class Plugin extends HarborConstruct {
	constructor(scope: Construct, pluginName: string, opts?: PluginOpts) {
		opts = PluginOpts.parse({
			//@ts-expect-error
			name: pluginName,
			//@ts-expect-error
			defaultRepository: true,
			...(opts || {}),
		})

		const url = opts.defaultRepository === true ? `https://artifacts.harbor.dev/plugins/${opts.name}.tar.gz` : `${opts.repository}/${opts.name}.tar.gz`;
		super(scope, pluginName, {
			kind: "harbor.dev/Plugin",
			options: {

			},
		});

		const downloadDestination = path.join(this.package.location, "packed/plugins", opts.name!);
		const unpackedDestination = path.join(this.package.location, "plugins", opts.name!);

		const unpack = new ExecCommand(this, "unpack", {
			executable: "tar",
			args: ["-xvzf", downloadDestination, "-C", unpackedDestination],
		});

		const verifyChecksum = new ExecCommand(this, "checksum", {
			executable: "harbor",
			args: ["plugin:verify", "--binary-location", unpackedDestination, "---source", url]
		});

		const executeInstallScript = new ExecCommand(this, "install_script", {
			executable: "harbor",
			args: ["plugin:install", "--wd", this.package.location, "--plugin", unpackedDestination]
		})

		const pipeline = new Pipeline(this, "plugin_install", [unpack, verifyChecksum, executeInstallScript]);

		new RemoteResource(this, "binary", {
			url,
			credentials: opts.credentials,
			destination: downloadDestination,
			postDownload: pipeline,
		});

		this.node.addDependency(pipeline);

		new PackageSetup(this, "setup", {
			actions: [this],
		})
	}
}