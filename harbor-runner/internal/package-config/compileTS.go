package packageconfig

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/clarkmcc/go-typescript"
	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/telemetry"
	// v8harbor "github.com/radding/harbor-runner/internal/v8"
)

var compileOpts map[string]interface{} = map[string]interface{}{}

func init() {
	compileOpts["module"] = "commonjs"
}

func CompileAndExecute(fiName string, resultWriter io.Writer) error {
	h := sha256.New()
	h.Write([]byte(fiName))
	d := fmt.Sprintf("fi-%x", h.Sum(nil))
	tempFi, err := os.CreateTemp("", d)
	if err != nil {
		return errors.Wrap(err, "failed to create a temp file")
	}
	tempFi.Close()
	fi, err := os.ReadFile(fiName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read typescript file %s", fiName))
	}
	var res string
	err = telemetry.TimeWithError("transpiling", func() error {
		var err error
		res, err = typescript.TranspileString(string(fi), typescript.WithCompileOptions(compileOpts))
		return err
	})
	if err != nil {
		return errors.Wrap(err, "failed to transpile typescript")
	}
	res = fmt.Sprintf("%s\n const fs = require(\"fs\");fs.writeFileSync(\"%s\", JSON.stringify(exports.default.createTree()));", res, filepath.ToSlash(tempFi.Name()))
	slog.Debug(fmt.Sprintf("COMPILED SCRIPT: \n%s\nEND COMPILED SCRIPT", res))
	cmd := exec.Command("node", "-e", res)
	cmd.Dir = path.Dir(fiName)
	cmd.Env = append(os.Environ(), "HARBORJS_IS_IN_RUNNER=true", fmt.Sprintf("HARBORJS_HARBOR_LOC=%s", fiName))
	cmd.Stderr = NewPipedLogger(slog.Error, slog.String("file", fiName))
	cmd.Stdout = NewPipedLogger(slog.Info, slog.String("file", fiName))
	err = telemetry.TimeWithError("execute node", cmd.Run)
	if err != nil {
		return errors.Wrap(err, "failed to run Node")
	}
	tempFi, err = os.Open(tempFi.Name())
	if err != nil {
		return errors.Wrap(err, "failed to reop tempfi")
	}
	defer tempFi.Close()
	io.Copy(resultWriter, tempFi)

	return nil
}

type PipedLogger struct {
	logger func(msg string, args ...any)
	attrs  []any
}

func (p *PipedLogger) Write(b []byte) (int, error) {
	p.logger(string(b), p.attrs...)
	return len(b), nil
}

func NewPipedLogger(logger func(msg string, args ...any), attrs ...any) *PipedLogger {
	return &PipedLogger{
		logger: logger,
		attrs:  attrs,
	}
}
