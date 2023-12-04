package config

import (
	"context"
	"fmt"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/techrail/ground/constants"
	"github.com/techrail/ground/constants/exitCode"
)

type redis struct {
	Main redisConfig
}

type redisConfig struct {
	Enabled                            bool
	Url                                string
	Username                           string
	Password                           string
	OperationMode                      string
	MaxActiveConnections               int
	MaxIdleConnections                 int
	IdleTimeoutInSeconds               int
	CrashAppOnConnectionFailure        bool
	ConnectRetryIntervalInSeconds      int
	AutoExpireTopLevelKeysAfterSeconds int
	AppNamespace                       string
}

func init() {
	config.Redis = redis{
		Main: redisConfig{
			Enabled:                            false,
			Url:                                "redis://127.0.0.1:6379",
			Username:                           "",
			Password:                           "",
			OperationMode:                      "auto",
			MaxActiveConnections:               10,
			MaxIdleConnections:                 10,
			IdleTimeoutInSeconds:               60,
			CrashAppOnConnectionFailure:        false,
			ConnectRetryIntervalInSeconds:      10,
			AutoExpireTopLevelKeysAfterSeconds: 0,
			AppNamespace:                       "ground-based-serviceNs",
		},
	}
}

func initializeRedisConfig() {
	config.Redis.Main.Enabled = envOrViperOrDefaultBool("redis.main.enabled", config.Redis.Main.Enabled)
	config.Redis.Main.Url = envOrViperOrDefaultString("redis.main.url", config.Redis.Main.Url)
	config.Redis.Main.Username = envOrViperOrDefaultString("redis.main.username", config.Redis.Main.Username)
	config.Redis.Main.Password = envOrViperOrDefaultString("redis.main.password", config.Redis.Main.Password)
	config.Redis.Main.OperationMode = envOrViperOrDefaultString("redis.main.operationmode", config.Redis.Main.OperationMode)
	config.Redis.Main.MaxActiveConnections = int(envOrViperOrDefaultInt64("redis.main.maxActiveConnections", int64(config.Redis.Main.MaxActiveConnections)))
	config.Redis.Main.MaxIdleConnections = int(envOrViperOrDefaultInt64("redis.main.maxIdleConnections", int64(config.Redis.Main.MaxIdleConnections)))
	config.Redis.Main.IdleTimeoutInSeconds = int(envOrViperOrDefaultInt64("redis.main.idleTimeoutInSeconds", int64(config.Redis.Main.IdleTimeoutInSeconds)))
	config.Redis.Main.CrashAppOnConnectionFailure = envOrViperOrDefaultBool("redis.main.crashAppOnConnectionFailure", config.Redis.Main.CrashAppOnConnectionFailure)
	config.Redis.Main.ConnectRetryIntervalInSeconds = int(envOrViperOrDefaultInt64("redis.main."+
		"connectRetryIntervalInSeconds", int64(config.Redis.Main.ConnectRetryIntervalInSeconds)))
	config.Redis.Main.AutoExpireTopLevelKeysAfterSeconds = int(envOrViperOrDefaultInt64("redis.main.autoExpireTopLevelKeysAfterSeconds", int64(config.Redis.Main.AutoExpireTopLevelKeysAfterSeconds)))
	config.Redis.Main.AppNamespace = envOrViperOrDefaultString("redis.main.appNamespace", config.Redis.Main.AppNamespace)

	// TODO: work on it and implement it
	if config.Redis.Main.Enabled &&
		config.Redis.Main.OperationMode != "auto" &&
		config.Redis.Main.OperationMode != "cluster" &&
		config.Redis.Main.OperationMode != "standalone" {
		fmt.Printf("P#1MQUNR - Invalid redis operation mode. Cannot proceed.")
		os.Exit(exitCode.RedisConnectionFailed)
	}
}

func NewClientWithUrl(url string) *goredis.Client {
	opts, err := goredis.ParseURL(url)
	if err != nil {
		fmt.Println("E#1OELAY - Could not parse the Redis URL provided.")
		panic(err)
	}
	rdb := goredis.NewClient(opts)
	return rdb
}

var rdb *goredis.Client

var ctx = context.Background()

func (r *redis) NewClientFromConfig() *goredis.Client {

	// If caching is not enabled, connecting attempt will not go through.
	if !r.Main.Enabled {
		fmt.Println("E#1OELCB - Enable Redis in the configuration file to create a client.")
		panic("Redis is not enabled. Enable it to start using Redis.")
	}

	// If username and password are provided, use the provided credentials to connect.
	if r.Main.Username != "" && r.Main.Password != "" {
		rdb = goredis.NewClient(&goredis.Options{
			Username: r.Main.Username,
			Addr:     constants.RedisDefaultAddr,
			Password: r.Main.Password,
			DB:       constants.RedisDefaultDbNumber,
			Protocol: constants.RedisDefaultProtocol,
		})

		fmt.Println("I#1OEMDA - Connected to Redis server using provided credentials.")
		return rdb

	}

	// If only password is provided, use the default username to connect.
	if r.Main.Password != "" {
		rdb = goredis.NewClient(&goredis.Options{
			Addr:     constants.RedisDefaultAddr,
			Password: r.Main.Password,
			DB:       constants.RedisDefaultDbNumber,
			Protocol: constants.RedisDefaultProtocol,
		})

		fmt.Println("I#1OEO9R - Connected to Redis server using default host and port.")
		return rdb

	}

	// Check if connection URL is available, if yes, then connect using that.
	if r.Main.Url != "" {
		opts, err := goredis.ParseURL(r.Main.Url)
		if err != nil {
			fmt.Println("E#1OEMOC - Could not parse the Redis URL provided in configuration file.")
			panic(err)
		}
		rdb = goredis.NewClient(opts)
		fmt.Println("I#1OEO90 - Connected to Redis server using connection url.")
		return rdb
	}
	// TODO - More options to be added

	//Nothing is provided, try with default config
	rdb = goredis.NewClient(&goredis.Options{
		Addr:     constants.RedisDefaultAddr,
		Password: constants.RedisDefaultPassword,
		DB:       constants.RedisDefaultDbNumber,
		Protocol: constants.RedisDefaultProtocol,
	})

	return rdb

}

// TODO - func for clusters and TLS connections

// We want this function to take any value (string, hash, set, list) and get it from cache,
// if the key does not exist, we want to set it in cache
func GetOrSet(key string, value []any) {

}

// This function would get the value for provided key from cache - works with all datatypes
func GetFromCache() {

}

// This function would set the provided value in cache - works with all datatypes
func SetInCache() {

}

// This function can take in a key, look up in cache, if not present in cache,
// the function would call the provided database function, to get the value from DB
// This should also work with all major types of data structures
func GetFromCacheOrDb() {

}

// Gets the string value associated with the given key.
func GetString(key string) string {
	value, err := rdb.Get(ctx, key).Result()
	if err == goredis.Nil {
		fmt.Println("E#1OIBZH - Key does not exist. " + err.Error())
	}
	return value
}

// Sets a string value with the provided expiry duration.
func SetStringWithExpiry(key, value string, ttl time.Duration) {
	err := rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		fmt.Println("E#1OICUW - Error occurred while setting string " + err.Error())
	}
}

// Sets a string value with no expiry.
func SetStringNoExpiry(key, value string) {
	err := rdb.SetNX(ctx, key, value, -1).Err()
	if err != nil {
		fmt.Println("E#1OICT1 - Error occurred while setting string " + err.Error())
	}

}

// This function appends the provided string to the existing string value mapped with the key.
func AppendToString(key, value string) {
	err := rdb.Append(ctx, key, value).Err()
	if err != nil {
		fmt.Println("E#1OJ0IC - Error occurred while appending to string " + err.Error())
	}
}

// This function increments the value associated with key by 1
func IncrementString(key string) {
	err := rdb.Incr(ctx, key).Err()
	if err != nil {
		fmt.Println("E#1OJ0SV - Error occurred while incrementing value. " + err.Error())
	}

}

func IncrementStringBy(key string, value int64) {
	err := rdb.IncrBy(ctx, key, value).Err()
	if err != nil {
		fmt.Println("E#1OJ16B - Error occurred while incrementing value by the given integer. " + err.Error())
	}

}

func IncrementStringByFloat(key string, value float64) {
	err := rdb.IncrByFloat(ctx, key, value).Err()
	if err != nil {
		fmt.Println("E#1OJ22L - Error occurred while incrementing value by given float. " + err.Error())
	}
}
