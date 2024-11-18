package cache

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/telemetry"
)

type Cache struct {
	cachePath string
}

func (c *Cache) Add(key string, data io.Reader) error {
	return telemetry.TimeWithError(fmt.Sprintf("add_to_cache_%s", key), func() error {
		cacheFile := path.Join(c.cachePath, key)
		slog.Debug("Writing to cache file", slog.String("cache_file", key))
		fi, err := os.OpenFile(cacheFile, os.O_CREATE, 0666)
		if err != nil {
			return errors.Wrap(err, "failed to open cahce file")
		}
		defer fi.Close()
		num, err := io.Copy(fi, data)
		if err != nil {
			return errors.Wrap(err, "failed to write to cache file")
		}
		slog.Debug(fmt.Sprintf("Wrote %d bytes to cache file", num), slog.String("cache_file", key))
		return nil
	})
}

func (c *Cache) Get(key string, dst io.Writer) (bool, error) {
	success := false
	err := telemetry.TimeWithError(fmt.Sprintf("add_to_cache_%s", key), func() error {
		cacheFile := path.Join(c.cachePath, key)
		slog.Debug("trying to get cache file", slog.String("cache_file", key))
		fi, err := os.OpenFile(cacheFile, os.O_RDONLY, 0666)
		if err != nil && os.IsNotExist(err) {
			slog.Debug("Cache file not found", slog.String("cache_file", key))
			return nil
		} else if err != nil {
			return errors.Wrap(err, "failed to open cache file")
		}
		defer fi.Close()

		num, err := io.Copy(dst, fi)
		if err != nil {
			return errors.Wrap(err, "failed to copy cached item")
		}
		slog.Debug(fmt.Sprintf("copied %d bytes from cached item", num), slog.String("cache_file", key))
		success = true
		return nil
	})
	return success, err
}

func (c *Cache) GetSubCache(key string) (*Cache, error) {
	cachePath := path.Join(c.cachePath, key)
	c2 := &Cache{
		cachePath: cachePath,
	}
	err := os.MkdirAll(cachePath, 0744)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make cache dir")
	}
	return c2, nil
}

func New(base string) (*Cache, error) {
	err := os.MkdirAll(base, 0744)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make cache dir")
	}
	c2 := &Cache{
		cachePath: base,
	}
	return c2, nil

}
