import { ExecCommand, HarborConstruct, Package, PackageSetup } from "../config/dist";
import { Pipeline } from "../config/dist/Pipeline"
// @ts-expect-error
import { execSync }  from 'child_process';
const shell = (cmd) => execSync(cmd, { encoding: 'utf8' });

function executableIsAvailable(name){
  try{ shell(`which ${name}`); return true}
  catch(error){return false}
}

const pkg = new Package("harbor-docs", {
    repository: "https://github.com/radding/harbor-config.git",
    path: "/docs",
    issues: "https://github.com/radding/harbor-config/issues",
    homepage: "https://github.com/radding/harbor-config"
});

const actions: HarborConstruct[] = [];
// const pipe = new Pipeline(pkg, )

if (!executableIsAvailable("pyenv")) {
    const installPyenv = new ExecCommand(pkg, "pyenv", {
        executable: "brew",
        args: [
            "install",
            "pyenv",
        ]
    });
    actions.push(installPyenv);
}


const installCorrectVersion = new ExecCommand(pkg, "pyenv-version", {
    executable: "pyenv",
    args: [
        "install",
        "-s",
        "3.10.14",
    ]
}).then(new ExecCommand(pkg, "set-local", {
    executable: "pyenv",
    args: [
        "local",
        "3.10.14",
    ]
})).then(new ExecCommand(pkg, "create-env", {
    executable: "pyenv",
   args: [
    "exec",
    "python",
    "-m",
    "venv",
    ".venv",
   ]
})) as ExecCommand

actions.push(installCorrectVersion);

const pipe = new Pipeline(pkg, "setup-pipeline", actions);

const installMKdocs = new ExecCommand(pkg, "install-mk-docs", {
    executable: "./.venv/bin/pip",
    args: [
        "install", 
        "-r",
        "requirements.txt"
    ]
}).needs(pipe)

new PackageSetup(pkg, "setup", {
    actions: [
        installMKdocs
    ]
})

const serve = new ExecCommand(pkg, "serve", {
    executable: "./.venv/bin/mkdocs",
    args: [
        "serve",
    ]
})
const build = new ExecCommand(pkg, "build", {
    executable: "./.venv/bin/mkdocs",
    args: [
        "build"
    ]
})

pkg.registerTask(serve);
pkg.registerTask(build);

export default pkg;