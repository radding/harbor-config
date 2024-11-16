package taskgraph

import (
	"encoding/json"
	"testing"

	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/stretchr/testify/assert"
)

func matchTree(t *Task, cfg *packageconfig.Config, assert *assert.Assertions) {
	construct := cfg.Constructs[t.ID]
	for ndx, task := range t.Dependencies {
		assert.Equal(task.ID, construct.DependsOn[ndx])
		matchTree(task, cfg, assert)
	}
}

func TestCreatingATree(t *testing.T) {
	assert := assert.New(t)
	cfg := &packageconfig.Config{}
	json.Unmarshal([]byte(testConfig), cfg)
	tree, err := CreateTreeFromConfig(cfg, &MockExecutor{})
	assert.NoError(err)
	assert.NotNil(tree)
	assert.Equal(tree.setupTask.ID, "setup-mock-task")
	assert.Equal(tree.setupTask.Kind, "harbor.dev/noop")
	assert.Len(tree.setupTask.Dependencies, len(cfg.Setup))
	for ndx, task := range tree.setupTask.Dependencies {
		assert.Equal(cfg.Setup[ndx], task.ID)
		matchTree(task, cfg, assert)
	}

	assert.Len(tree.tasks, len(cfg.Tasks))
	for key, task := range tree.tasks {
		assert.Equal(task.ID, cfg.Tasks[key])
		matchTree(task, cfg, assert)
	}
}

var testConfig = `{
    "constructs": {
        "harbor-code/build": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "go",
                "args": [
                    "build",
                    "-o",
                    "harbor",
                    "./cmd/harbor/main.go"
                ]
            },
            "dependsOn": [
                "harbor-code/https:----github.com--radding--harbor/dep-https:----github.com--radding--harbor-test",
                "harbor-code/test",
                "harbor-code/lint"
            ]
        },
        "harbor-code/go-install": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "gvm",
                "args": [
                    "install",
                    "go1.22.1"
                ]
            },
            "dependsOn": []
        },
        "harbor-code/go-version": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "gvm",
                "args": [
                    "use",
                    "go1.22.1"
                ]
            },
            "dependsOn": [
                "harbor-code/go-install"
            ]
        },
        "harbor-code/https:----github.com--radding--harbor": {
            "kind": "harbor.dev/Dependency",
            "options": {
                "url": "https://github.com/radding/harbor",
                "path": "nodego/"
            },
            "dependsOn": []
        },
        "harbor-code/https:----github.com--radding--harbor/dep-https:----github.com--radding--harbor-test": {
            "kind": "Harbor.dev/Task",
            "options": {
                "plugin": "remote-executor",
                "run": "test",
                "inputs": []
            },
            "dependsOn": [
                "harbor-code/remote-executor",
                "harbor-code/https:----github.com--radding--harbor"
            ]
        },
        "harbor-code/lint": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "golangci-lint",
                "args": [
                    "lint"
                ]
            },
            "dependsOn": []
        },
        "harbor-code/remote-executor": {
            "kind": "harbor.dev/RemoteExecutor",
            "options": {},
            "dependsOn": []
        },
        "harbor-code/setup-harbor-core": {
            "kind": "harbor.dev/PackageSetup",
            "options": {},
            "dependsOn": [
                "harbor-code/vendor-modules"
            ]
        },
        "harbor-code/test": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "go",
                "args": [
                    "test",
                    "-v",
                    "./..."
                ]
            },
            "dependsOn": []
        },
        "harbor-code/tidy-modules": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "go",
                "args": [
                    "mod",
                    "tidy"
                ]
            },
            "dependsOn": [
                "harbor-code/go-version"
            ]
        },
        "harbor-code/vendor-modules": {
            "kind": "harbor.dev/ExecCommand",
            "options": {
                "executable": "go",
                "args": [
                    "work",
                    "vendor"
                ]
            },
            "dependsOn": [
                "harbor-code/tidy-modules"
            ]
        }
    },
    "tasks": {
        "build": "harbor-code/build",
        "dep-https://github.com/radding/harbor-test": "harbor-code/https:----github.com--radding--harbor/dep-https:----github.com--radding--harbor-test",
        "lint": "harbor-code/lint",
        "test": "harbor-code/test"
    },
    "setup": [
        "harbor-code/setup-harbor-core"
    ],
    "packageInfo": {
        "meta": {
            "harborPackageDirectory": ""
        },
        "repository": "https://github.com/radding/harbor",
        "version": "1.0.0",
        "name": "harbor-code",
        "path": "harbor-core/",
        "homepage": "",
        "description": "The Package for the Harbor executable",
        "issues": "",
        "license": "",
        "stability": "Alpha",
        "artifactsLocation": ""
    },
    "was_setup_run": true
}`
