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

type CacheContextKey string

const CacheContextKeyValue = CacheContextKey("Cacher")

type cache struct {
	CachePath string `json:"cache_path"`
}

type NonCache struct{}

func (n *NonCache) Add(key string, data io.Reader) error {
	slog.Debug("adding to non cache, this is a noop")
	return nil
}

func (n *NonCache) Get(key string, data io.Writer) (bool, error) {
	slog.Debug("getting from non cache, this is a noop")
	return false, nil
}

func (n *NonCache) Clean() error {
	slog.Debug("cleaning the non cache, this is a noop")
	return nil
}

func (n *NonCache) GetSubCache(key string) (Cache, error) {
	slog.Debug("getting a subcache from the non cache, this is a noop")
	return n, nil
}

type Cache interface {
	Add(key string, data io.Reader) error
	Get(key string, dst io.Writer) (bool, error)
	Clean() error
	GetSubCache(key string) (Cache, error)
}

func (c *cache) Add(key string, data io.Reader) error {
	return telemetry.TimeWithError(fmt.Sprintf("add_to_cache_%s", key), func() error {
		cacheFile := path.Join(c.CachePath, key)
		slog.Debug("Writing to cache file", slog.String("cache_file", key))
		fi, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
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

func (c *cache) Get(key string, dst io.Writer) (bool, error) {
	success := false
	err := telemetry.TimeWithError(fmt.Sprintf("add_to_cache_%s", key), func() error {
		cacheFile := path.Join(c.CachePath, key)
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

func (c *cache) Clean() error {
	slog.Debug("removing cache directory", slog.String("cache_directory", c.CachePath))
	return os.RemoveAll(c.CachePath)
}

func (c *cache) GetSubCache(key string) (Cache, error) {
	cachePath := path.Join(c.CachePath, key)
	c2 := &cache{
		CachePath: cachePath,
	}
	err := os.MkdirAll(cachePath, 0744)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make cache dir")
	}
	return c2, nil
}

func New(base string) (Cache, error) {
	err := os.MkdirAll(base, 0744)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make cache dir")
	}
	c2 := &cache{
		CachePath: base,
	}
	return c2, nil

}
