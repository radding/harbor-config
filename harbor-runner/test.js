"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var dist_1 = require("../config/dist");
var GO_VERSION = "go1.22.1";
var pkg = new dist_1.Package("harbor-code", {
  description: "The Package for the Harbor executable",
  repository: "https://github.com/radding/harbor",
  version: "1.0.0",
  path: "harbor-core/",
  stability: "Alpha",
});
var nodego = new dist_1.Dependency(pkg, "https://github.com/radding/harbor", {
  path: "nodego/",
});
var vendorModules = new dist_1.ExecCommand(pkg, "vendor-modules", {
  executable: "go",
  args: ["work", "vendor"],
});
var tidyModules = new dist_1.ExecCommand(pkg, "tidy-modules", {
  executable: "go",
  args: ["mod", "tidy"],
});
var installGoVersion = new dist_1.ExecCommand(pkg, "go-install", {
  executable: "gvm",
  args: ["install", GO_VERSION],
});
var ensureGoVersion = new dist_1.ExecCommand(pkg, "go-version", {
  executable: "gvm",
  args: ["use", GO_VERSION],
});
var pipeline = installGoVersion
  .then(ensureGoVersion)
  .then(tidyModules)
  .then(vendorModules);
new dist_1.PackageSetup(pkg, "setup-harbor-core", {
  actions: [pipeline],
});
var tests = new dist_1.ExecCommand(pkg, "test", {
  executable: "go",
  args: ["test", "-v", "./..."],
});
pkg.registerTask(tests);
var lint = new dist_1.ExecCommand(pkg, "lint", {
  executable: "golangci-lint",
  args: ["lint"],
});
pkg.registerTask(lint);
pkg.registerTask(
  new dist_1.ExecCommand(pkg, "build", {
    executable: "go",
    args: ["build", "-o", "harbor", "./cmd/harbor/main.go"],
  }).needs(nodego.task("test"), tests, lint)
);
exports.default = pkg;
const fs = require("fs");
fs.writeFileSync(
  "/Users/raddi/AppData/LocalTemp/fi-3034fabb8b9524545f5855ae011abe0c66a656808f683153087947fec615913a257621142",
  JSON.stringify(exports.default.createTree())
);
