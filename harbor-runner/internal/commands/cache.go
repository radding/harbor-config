package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/spf13/cobra"
)

func boolPtr(b bool) *bool {
	return &b
}

var cleanAllCache = boolPtr(false)
var cleanOnlyOld = boolPtr(false)

func init() {
	rootCmd.AddCommand(Cache)
	Cache.AddCommand(CacheClean)
	CacheClean.Flags().BoolVarP(cleanAllCache, "all", "a", false, "Clean all cached elements, both old and new")
	CacheClean.Flags().BoolVarP(cleanOnlyOld, "old", "o", false, "Only clean the old items out")
	Cache.AddCommand(CacheInfo)
}

var Cache = &cobra.Command{
	Use:   "cache",
	Short: "manipulate the harbor cache",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

var CacheClean = &cobra.Command{
	Use:   "clean",
	Short: "clean the cache",
	Long:  "By default this, command just removes the current cache elements",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("cleaning the cache")
		cfg := packageconfig.GetConfig()
		if *cleanAllCache {
			return os.RemoveAll("./.harbor")
		}
		if *cleanOnlyOld {
			fileInfos, err := os.ReadDir("./.harbor")
			slog.Debug(fmt.Sprintf("found %d files", len(fileInfos)))
			if err != nil && !os.IsNotExist(err) {
				return errors.Wrap(err, "failed to get .harbor dir")
			}
			for _, info := range fileInfos {
				if info.Name() != cfg.GetHash() {
					slog.Debug(fmt.Sprintf("deleting cache element %q", info.Name()))
					err := os.RemoveAll(info.Name())
					if err != nil {
						return err
					}
				} else {
					slog.Info("skipping current directory")
				}

			}
		}
		return cfg.GetCache().Clean()
	},
}

var CacheInfo = &cobra.Command{
	Use:   "info",
	Short: "get info on the cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("Getting Cache information")
		cfg := packageconfig.GetConfig()
		slog.Info(fmt.Sprintf("Base cache hash: %s", cfg.GetHash()))
		return nil
	},
}
