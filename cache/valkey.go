// valkey_cache.go
//
// Idiomatic Go wrapper around github.com/valkey-io/valkey-go for standalone Valkey servers.
// Provides utility methods for string and list operations, with context support, error handling,
// and extensibility for future data types.
//
// Usage example:
//
//   cache, err := NewValkeyCache("localhost", "6379")
//   if err != nil {
//       log.Fatalf("Failed to connect: %v", err)
//   }
//   defer cache.Close()
//
//   ctx := context.Background()
//   err = cache.Set(ctx, "foo", "bar")
//   val, err := cache.Get(ctx, "foo")
//   n, err := cache.LPush(ctx, "mylist", "a", "b")
//   vals, err := cache.LRange(ctx, "mylist", 0, -1)
//

package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

// ErrNotFound is returned when a key does not exist in Valkey.
var ErrNotFound = errors.New("valkey: key not found")

// ValkeyCacheAPI defines the interface for the cache wrapper.
// Useful for dependency injection and testing.
type ValkeyCacheAPI interface {
	Set(ctx context.Context, key, value string, opts ...SetOption) error
	Get(ctx context.Context, key string) (string, error)
	LPush(ctx context.Context, key string, values ...string) (int64, error)
	RPush(ctx context.Context, key string, values ...string) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Close()
	// Underlying returns the underlying valkey.Client for advanced use.
	Underlying() valkey.Client
}

// ValkeyCache is an idiomatic Go wrapper around valkey-go for standalone Valkey servers.
type ValkeyCache struct {
	client valkey.Client
}

// NewValkeyCache creates a new ValkeyCache instance connected to the given host and port.
// Supports optional configuration via functional options (e.g., WithAuth, WithTLS).
func NewValkeyCache(host, port string, opts ...Option) (*ValkeyCache, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	cfg := &config{
		addresses: []string{address},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	clientOpt := valkey.ClientOption{
		InitAddress: cfg.addresses,
		Username:    cfg.username,
		Password:    cfg.password,
		TLSConfig:   cfg.tlsConfig,
		// Add more options as needed (e.g., ClientName, DisableCache, etc.)
	}
	client, err := valkey.NewClient(clientOpt)
	if err != nil {
		return nil, fmt.Errorf("valkey: failed to connect: %w", err)
	}
	return &ValkeyCache{client: client}, nil
}

// Option configures ValkeyCache.
type Option func(*config)

type config struct {
	addresses []string
	username  string
	password  string
	tlsConfig *tls.Config
}

// WithAuth sets the username and password for Valkey AUTH.
func WithAuth(username, password string) Option {
	return func(cfg *config) {
		cfg.username = username
		cfg.password = password
	}
}

// WithTLS enables TLS with the given configuration.
func WithTLS(tlsConfig *tls.Config) Option {
	return func(cfg *config) {
		cfg.tlsConfig = tlsConfig
	}
}

// SetOption configures the Set operation (e.g., expiration).
type SetOption func(*setOptions)

type setOptions struct {
	expiration time.Duration
}

// WithExpiration sets the expiration for the key.
func WithExpiration(d time.Duration) SetOption {
	return func(opts *setOptions) {
		opts.expiration = d
	}
}

// Set sets the string value for a key, optionally with expiration.
func (c *ValkeyCache) Set(ctx context.Context, key, value string, opts ...SetOption) error {
	var so setOptions
	for _, opt := range opts {
		opt(&so)
	}
	builder := c.client.B().Set().Key(key).Value(value)
	if so.expiration > 0 {
		b := builder.Ex(time.Duration(so.expiration.Seconds()))
		cmd := b.Build()
		return c.client.Do(ctx, cmd).Error()
	}
	cmd := builder.Build()
	return c.client.Do(ctx, cmd).Error()
}

// Get retrieves the string value for a key.
// Returns ErrNotFound if the key does not exist.
func (c *ValkeyCache) Get(ctx context.Context, key string) (string, error) {
	cmd := c.client.B().Get().Key(key).Build()
	val, err := c.client.Do(ctx, cmd).ToString()
	if valkey.IsValkeyNil(err) {
		return "", ErrNotFound
	}
	return val, err
}

// LPush pushes values to the head of the list at key.
// Returns the new length of the list.
func (c *ValkeyCache) LPush(ctx context.Context, key string, values ...string) (int64, error) {
	if len(values) == 0 {
		return 0, errors.New("valkey: LPush requires at least one value")
	}
	cmd := c.client.B().Lpush().Key(key).Element(values...).Build()
	return c.client.Do(ctx, cmd).AsInt64()
}

// RPush pushes values to the tail of the list at key.
// Returns the new length of the list.
func (c *ValkeyCache) RPush(ctx context.Context, key string, values ...string) (int64, error) {
	if len(values) == 0 {
		return 0, errors.New("valkey: RPush requires at least one value")
	}
	cmd := c.client.B().Rpush().Key(key).Element(values...).Build()
	return c.client.Do(ctx, cmd).AsInt64()
}

// LPop pops a value from the head of the list at key.
// Returns ErrNotFound if the list is empty or the key does not exist.
func (c *ValkeyCache) LPop(ctx context.Context, key string) (string, error) {
	cmd := c.client.B().Lpop().Key(key).Build()
	val, err := c.client.Do(ctx, cmd).ToString()
	if valkey.IsValkeyNil(err) {
		return "", ErrNotFound
	}
	return val, err
}

// RPop pops a value from the tail of the list at key.
// Returns ErrNotFound if the list is empty or the key does not exist.
func (c *ValkeyCache) RPop(ctx context.Context, key string) (string, error) {
	cmd := c.client.B().Rpop().Key(key).Build()
	val, err := c.client.Do(ctx, cmd).ToString()
	if valkey.IsValkeyNil(err) {
		return "", ErrNotFound
	}
	return val, err
}

// LRange returns the elements of the list at key between start and stop (inclusive).
// Returns an empty slice if the key does not exist.
func (c *ValkeyCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	cmd := c.client.B().Lrange().Key(key).Start(start).Stop(stop).Build()
	vals, err := c.client.Do(ctx, cmd).AsStrSlice()
	if valkey.IsValkeyNil(err) {
		return []string{}, nil
	}
	return vals, err
}

// Close closes the underlying Valkey client and releases resources.
func (c *ValkeyCache) Close() {
	c.client.Close()
}

// Underlying returns the underlying valkey.Client for advanced use.
func (c *ValkeyCache) Underlying() valkey.Client {
	return c.client
}

// --- Usage Example ---

/*
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "yourmodule/valkeycache"
)

func main() {
    cache, err := valkeycache.NewValkeyCache("localhost", "6379")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer cache.Close()

    ctx := context.Background()

    // Set and Get
    err = cache.Set(ctx, "foo", "bar", valkeycache.WithExpiration(10*time.Second))
    if err != nil {
        log.Fatalf("Set failed: %v", err)
    }
    val, err := cache.Get(ctx, "foo")
    if err != nil {
        log.Fatalf("Get failed: %v", err)
    }
    fmt.Println("foo =", val)

    // List operations
    _, _ = cache.LPush(ctx, "mylist", "a", "b", "c")
    vals, _ := cache.LRange(ctx, "mylist", 0, -1)
    fmt.Println("mylist =", vals)
}
*/
