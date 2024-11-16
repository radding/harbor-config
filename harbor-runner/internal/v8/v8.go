package v8harbor

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	v8 "rogchap.com/v8go"
)

var vm *V8VM = &V8VM{}

type V8VM struct {
	iso    *v8.Isolate
	global *v8.ObjectTemplate
}

func GetVm() *V8VM {
	return vm
}

func (v *V8VM) Initialize() error {
	v.iso = v8.NewIsolate()
	v.global = v8.NewObjectTemplate(v.iso)
	return nil
}

func (v *V8VM) Clean() error {
	// v.iso.TerminateExecution()
	// v.iso.Dispose()
	return nil
}

type moduleScope struct {
	__dirname  string
	__filename string
}

var MODULE_SCOPE_REF moduleScope = moduleScope{}

func RunScript(script string, filename string) (*v8.Context, *v8.Value, *v8.Value, error) {
	ctxBg := context.Background()
	loc := filepath.Dir(filename)
	m := moduleScope{
		__dirname:  loc,
		__filename: filename,
	}
	ctxBg = context.WithValue(ctxBg, MODULE_SCOPE_REF, m)
	iso := vm.iso // create a new VM
	exports := v8.NewObjectTemplate(iso)
	requireFn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		return require(ctxBg, info)
	})
	global := v8.NewObjectTemplate(iso) // a template that represents a JS Object
	global.Set("exports", exports)
	global.Set("require", requireFn)
	ctx := v8.NewContext(iso, global)
	res, err := ctx.RunScript(script, filename) // will execute the Go callback with a single argunent 'foo'
	if err != nil {
		return ctx, nil, nil, errors.Wrap(err, "failed to execute script")
	}
	exported, err := ctx.Global().Get("exports")
	return ctx, res, exported, err
}
