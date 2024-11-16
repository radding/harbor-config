package v8harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	v8 "rogchap.com/v8go"
)

func loadModule(fileName string) *v8.Value {
	contents, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading %s", fileName), slog.String("error", err.Error()))
		return nil
	}
	// fmt.Printf("%s", contents) // when the JS function is called this Go callback will execute
	_, _, exported, err := RunScript(string(contents), fileName)
	if err != nil {
		slog.Error(fmt.Sprintf("Error evaluating %s", fileName), slog.String("error", err.Error()))
	}
	if exported != nil {
		fmt.Println("object: ", exported.DetailString())
	}
	return exported // you can return a value back to the JS caller if required
}

func loadAsFile(module string) *v8.Value {
	_, err := os.Stat(module)
	if err != nil && os.IsNotExist(err) {
		module = fmt.Sprintf("%s.js", module)
		_, err = os.Stat(module)
		if err != nil && os.IsNotExist(err) {
			slog.Debug("Failed to find module as file", slog.String("module", module))
			return nil
		} else if err != nil {
			slog.Warn(fmt.Sprintf("failed to stat module %s", module), slog.String("error", err.Error()))
			return nil
		}

	} else if err != nil {
		slog.Error(fmt.Sprintf("Could not stat module %s", module), slog.String("error", err.Error()))
		return nil
	}
	return loadModule(module)
}

func loadIndex(dir string) *v8.Value {
	loadPath := filepath.Join(dir, "index.js")
	return loadModule(loadPath)
}

type packageInfo struct {
	Main *string `json:"main"`
}

func loadAsDirectory(dir string) *v8.Value {
	pkgFile := filepath.Join(dir, "package.json")
	info, err := os.Stat(pkgFile)
	if os.IsNotExist(err) || (err == nil && !info.IsDir()) {
		return loadIndex(dir)
	}
	bs, err := os.ReadFile(pkgFile)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read dir %s", dir), slog.String("error", err.Error()))
		return nil
	}

	pkg := &packageInfo{}
	err = json.Unmarshal(bs, pkg)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read pacakge.json %s", pkgFile), slog.String("error", err.Error()))
		return nil
	}
	if pkg.Main == nil {
		return loadIndex(dir)
	}
	loadLocation := filepath.Join(dir, *pkg.Main)
	data := loadAsFile(loadLocation)
	if data == nil {
		data = loadIndex(loadLocation)
	}
	return data
}
func require(ctx context.Context, info *v8.FunctionCallbackInfo) *v8.Value {
	scope := ctx.Value(MODULE_SCOPE_REF).(moduleScope)
	loc := scope.__dirname
	toImport := info.Args()[0].String()
	if strings.HasPrefix(toImport, "/") {
		loc = toImport
	}
	if strings.HasPrefix(toImport, "./") || strings.HasPrefix(toImport, "../") {
		p := filepath.Join(loc, toImport)
		if strings.HasPrefix(p, "/") || strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") {
			d := loadAsFile(p)
			if d == nil {
				d = loadAsDirectory(p)
			}
			if d == nil {
				slog.Warn(fmt.Sprintf("failed to load module %s", p))
			}
			return d
		}

	}
	if strings.HasPrefix(toImport, "#") {
		slog.Warn("This method is not supported")
	}
	// statInfo, err := os.Stat(p)
	// if err != nil && !os.IsNotExist(err) {
	// 	slog.Error(fmt.Sprintf("Error importing %s", toImport), slog.String("error", err.Error()))
	// 	return nil
	// }
	// if statInfo != nil && statInfo.IsDir() {
	// 	p = filepath.Join(p, "index.js")
	// } else {
	// 	p += ".js"
	// }
	// contents, err := os.ReadFile(p)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Error importing %s", toImport), slog.String("error", err.Error()))
	// }
	// fmt.Printf("%s", contents) // when the JS function is called this Go callback will execute
	// // d, err := info.Context().RunScript(string(contents), p)
	// d, err := RunScript(string(contents), p)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Error evaluating %s", toImport), slog.String("error", err.Error()))
	// }
	return nil // you can return a value back to the JS caller if required
}
