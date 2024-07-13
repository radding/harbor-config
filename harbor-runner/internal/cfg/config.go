package cfg

import (
	"os"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/telemetry"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type ConfigLifeCycle struct {
}

func (c *ConfigLifeCycle) Initialize() error {
	viper.SetConfigName("harbor_cfg")
	viper.SetConfigType("json")
	viper.SetEnvPrefix("harbor")
	viper.AddConfigPath("$HOME/.harbor")
	slog.Debug("Reading configuration")
	err := viper.ReadInConfig()

	viper.SetDefault("log_level", telemetry.InfoLevel)
	viper.SetDefault("log_format_json", false)
	viper.SetDefault("plugin_cache", "$HOME/.harbor/plugins/cache")

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		} else {
			return errors.Wrap(err, "failed to read configuration file")
		}
	}
	return nil
}

func (c *ConfigLifeCycle) Clean() error {
	slog.Debug("Saving configuration")
	err := os.MkdirAll(os.ExpandEnv("$HOME/.harbor"), 0700)
	if err != nil {
		return errors.Wrap(err, "failed to create harbor home")
	}
	err = viper.SafeWriteConfig()
	if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok && err != nil {
		return err
	}
	return nil
}
