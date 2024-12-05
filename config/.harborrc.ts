import { ExecCommand, Package, PackageSetup } from "./dist";

const pkg = new Package("harbor-config", {
	repository: "https://github.com/radding/harbor",
	version: "1.0.0",
    path: "config/"
});

const install = new ExecCommand(pkg, "install", {
	executable: "yarn",
	args: [
		"install",
	]
})

const buildCommand = new ExecCommand(pkg, "build", {
	executable: "yarn",
	args: [
		"build"
	]
})

pkg.registerTask(buildCommand)

new PackageSetup(pkg, "setup", {
	actions: [
		install
	]
})

export default pkg;
