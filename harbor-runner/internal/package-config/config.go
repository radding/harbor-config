package packageconfig

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/cache"
	"github.com/radding/harbor-runner/internal/telemetry"
)

type Construct struct {
	Kind      string          `json:"kind"`
	Options   json.RawMessage `json:"options"`
	DependsOn []string        `json:"dependsOn"`
}

type PackageInfo struct {
	Meta struct {
		HarborPackageDirectory string `json:"harborPackageDirectory"`
	} `json:"meta"`
	Repository        string `json:"repository"`
	Version           string `json:"version"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	Homepage          string `json:"homepage"`
	Description       string `json:"description"`
	Issues            string `json:"issues"`
	License           string `json:"license"`
	Stability         string `json:"stability"`
	ArtifactsLocation string `json:"artifactsLocation"`
}

type Config struct {
	hash           string
	cachedLocation string
	workingDir     string
	Constructs     map[string]Construct `json:"constructs"`
	Tasks          map[string]string    `json:"tasks"`
	Setup          []string             `json:"setup"`
	PackageInfo    PackageInfo          `json:"packageInfo"`
	WasSetupRun    bool                 `json:"was_setup_run"`
	cacher         *cache.Cache
}

func (c *Config) Save() error {
	bts, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshal configuration")
	}
	err = os.WriteFile(c.cachedLocation, bts, 0666)
	if err != nil {
		return errors.Wrap(err, "failed to save configuration")
	}
	return nil
}

func (c *Config) GetCache() *cache.Cache {
	return c.cacher
}

func (c *Config) ConfigureContext(ct context.Context) context.Context {
	ctx := context.WithValue(ct, "CacheLocation", path.Dir(c.cachedLocation))
	ctx = context.WithValue(ctx, "WorkingDir", c.workingDir)
	ctx = context.WithValue(ctx, "Cache", c.cacher)
	return ctx
}

var configs map[string]Config = map[string]Config{}

func LoadConfig(fileName string) (Config, error) {
	telemetry.Trace(fmt.Sprintf("loading %s config", fileName))
	if conf, ok := configs[fileName]; ok {
		slog.Debug("config already loaded into memory, returning it now", slog.String("file", fileName))
		return conf, nil
	}
	hasher := sha256.New()
	s, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, fmt.Errorf("error opening file: %w", err)
	}
	hasher.Write(s)
	hashedFile := hex.EncodeToString(hasher.Sum(nil))
	configPath := path.Join("./.harbor", hashedFile, "config.json")
	var config = Config{
		WasSetupRun:    false,
		cachedLocation: configPath,
		workingDir:     path.Dir(fileName),
	}
	config.cacher, err = cache.New(path.Dir(configPath))
	if err != nil {
		return config, errors.Wrap(err, "failed to create cache")
	}
	slog.Debug(fmt.Sprintf("loading config from %s", configPath))
	buff := new(bytes.Buffer)
	success, err := config.cacher.Get("config.json", buff)
	if err != nil {
		return config, errors.Wrap(err, "failed to read from cache")
	}
	if !success {
		slog.Debug("config isn't cached, creating it now", slog.String("CachedPath", configPath))
		err := telemetry.TimeWithError("compile config", func() error {
			buffer := new(bytes.Buffer)
			err := CompileAndExecute(fileName, buffer)
			if err != nil {
				return errors.Wrap(err, "failed to execute config file")
			}
			bts := buffer.Bytes()
			configResults := string(bts)
			telemetry.Trace("Got config results", slog.String("results", configResults))
			if err = json.Unmarshal(bts, &config); err != nil {
				return errors.Wrap(err, "failed to unmarshal resulting config")
			}
			telemetry.Trace("writing the config file to cache", slog.String("cachedfile", configPath))
			err = config.cacher.Add("config.json", bytes.NewBuffer(bts))
			if err != nil {
				return errors.Wrap(err, "failed to add to cache")
			}
			return nil
		})
		if err != nil {
			slog.Error("Faild to execute configuration file", slog.String("error", err.Error()))
			return config, nil
		}

	} else {
		if err = json.Unmarshal(buff.Bytes(), &config); err != nil {
			return config, errors.Wrap(err, "failed to unmarshal cached config")
		}
	}
	config.hash = hashedFile
	configs[fileName] = config

	return config, nil
}

func tryFindConfigBase(pathName, fileName string, maxRecursion int64) (string, error) {
	if maxRecursion < 0 || pathName == "/" {
		return "", fmt.Errorf("max recursion to find package config")
	}
	maybeFile := path.Join(pathName, fileName)
	_, err := os.Stat(maybeFile)
	if err != nil && os.IsNotExist(err) {
		return tryFindConfigBase(filepath.Dir(pathName), fileName, maxRecursion-1)
	} else if err != nil {
		return "", errors.Wrap(err, "couldn't stat potential config file")
	}
	return pathName, nil
}

var pkg *Config

func GetConfig() *Config {
	return pkg
}

type Lifecycle struct{}

type NotFoundError struct {
	err error
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("no configuration file found: %s", n.err.Error())
}

func IsNotFoundError(e error) bool {
	_, ok := e.(*NotFoundError)
	return ok
}

func (l *Lifecycle) Initialize() error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not get working directory")
	}
	base, err := tryFindConfigBase(wd, ".harborrc.ts", 100)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to find the base of the project: %s", err))
		return nil
	}
	p, err := LoadConfig(path.Join(base, ".harborrc.ts"))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load configuration of the project: %s", err))
		return err
	}
	pkg = &p
	return nil
}
