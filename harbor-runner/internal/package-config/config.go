package packageconfig

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/radding/harbor-runner/internal/telemetry"
)

type Construct struct {
	Kind      string          `json:"kind"`
	Options   json.RawMessage `json:"options"`
	DependsOn []string        `json:"dependsOn"`
}

type PackageInfo struct {
	Repository string `json:"repository"`
	Version    string `json:"version"`
}

type Config struct {
	hash        string
	Constructs  map[string]Construct `json:"constructs"`
	Tasks       map[string]string    `json:"tasks"`
	Setup       []string             `json:"setup"`
	PackageInfo PackageInfo          `json:"packageInfo"`
}

var configs map[string]Config = map[string]Config{}

func LoadConfig(fileName string) (Config, error) {
	telemetry.Trace(fmt.Sprintf("loading %s config", fileName))
	if conf, ok := configs[fileName]; ok {
		slog.Debug(fmt.Sprintf("%s was found in our cache returning it", fileName))
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
	slog.Debug(fmt.Sprintf("loading config from %s", configPath))

	configBytes, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		slog.Debug(fmt.Sprintf("%s does not exsist, executing now", configPath))
	} else if err != nil {
		slog.Warn(fmt.Sprintf("error reading %s", configPath), slog.String("error", err.Error()))
	}

	conf := Config{}
	conf.hash = hashedFile
	err = json.Unmarshal(configBytes, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to unserialize configuration: %w", err)
	}
	configs[fileName] = conf
	return conf, nil
}
