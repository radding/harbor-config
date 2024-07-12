import { ExecCommand, Package, PackageSetup,  } from "../config/src";

const GO_VERSION = "go1.22.1";

const pkg = new Package("harbor-code", {
    description: "The Package for the Harbor executable",
	repository: "https://github.com/radding/harbor",
	version: "1.0.0",
    path: "harbor-core/",
    stability: "Alpha", 
});


const vendorModules = new ExecCommand(pkg, "vendor-modules", {
    executable: "go",
    args: [
        "work",
        "vendor",
    ],
})

const tidyModules = new ExecCommand(pkg, "tidy-modules", {
    executable: "go",
    args: [
        "mod",
        "tidy",
    ],
});

const installGoVersion = new ExecCommand(pkg, "go-install", {
    executable: "gvm",
    args: [
        "install",
        GO_VERSION,
    ]
})

const ensureGoVersion = new ExecCommand(pkg, "go-version", {
    executable: "gvm",
    args: [
        "use",
        GO_VERSION,
    ]
})

let pipeline = installGoVersion
    .then(ensureGoVersion)
    .then(tidyModules)
    .then(vendorModules);

new PackageSetup(pkg, "setup-harbor-core", {
    actions: [pipeline],
});

const tests = new ExecCommand(pkg, "test", {
    executable: "go",
    args: [
        "test",
        "-v",
        "./...",
    ]
});
pkg.registerTask(tests);

const lint = new ExecCommand(pkg, "lint", {
    executable: "golangci-lint",
    args: [
        "lint",
    ]
});
pkg.registerTask(lint);

pkg.registerTask(new ExecCommand(pkg, "build", {
    executable: "go",
    args: [
        "build",
        "-o",
        "harbor",
        "./cmd/harbor/main.go"
    ]
}).needs(
    tests,
    lint,
));

pkg.synth();
