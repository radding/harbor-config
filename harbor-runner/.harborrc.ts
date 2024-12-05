import { LocalDependency, ExecCommand, Package, PackageSetup,  } from "../config/dist";

const GO_VERSION = "go1.22.1";

const pkg = new Package("harbor-code", {
    description: "The Package for the Harbor executable",
	repository: "https://github.com/radding/harbor",
	version: "1.0.0",
    path: "harbor-core/",
    stability: "Alpha", 
});

const config = new LocalDependency(pkg, "../config");


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

// const ensureGoVersion = new ExecCommand(pkg, "go-version", {
//     executable: "/bin/bash",
//     args: [
//         "-c",
//         `gvm_use ${GO_VERSION}`,
//     ],
//     env: {
//         // @ts-expect-error
//         "PATH": process.env.PATH
//     }
// })

let pipeline = installGoVersion
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

const goBuild = new ExecCommand(pkg, "go-build", {
    executable: "go",
    args: [
        "build",
        "-o",
        "harbor",
        "./cmd/harbor/main.go"
    ]
}).needs(
    config.task("build"),
    tests,
)

// const lint = new ExecCommand(pkg, "lint", {
//     executable: "golangci-lint",
//     args: [
//         "lint",
//     ]
// });
// pkg.registerTask(lint);

pkg.registerTask(new ExecCommand(pkg, "build", {
    executable: "codesign",
    args: [
        "--sign",
        "-",
        "--force",
        "--preserve-metadata=entitlements,requirements,flags,runtime",
        "./harbor"
    ]
}).needs(
   goBuild,
    // lint,
));

export default pkg;
