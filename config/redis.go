package config

import (
	"context"
	"fmt"
	"os"

	goredis "github.com/redis/go-redis/v9"
	"github.com/techrail/ground/constants"
	"github.com/techrail/ground/constants/exitCode"
)

type redis struct {
	Main       redisConfig
	Connection goredis.UniversalClient
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
	Address                            string
}

func init() {
	config.Redis = redis{
		Main: redisConfig{
			Enabled:                            constants.RedisEnabled,
			Url:                                constants.RedisDefaultUrl,
			Username:                           constants.RedisDefaultUsername,
			Password:                           constants.RedisDefaultPassword,
			OperationMode:                      constants.RedisDefaultOperationMode,
			MaxActiveConnections:               constants.RedisMaxActiveConnections,
			MaxIdleConnections:                 constants.RedisMaxIdleConnections,
			IdleTimeoutInSeconds:               constants.RedisIdleTimeoutInSeconds,
			CrashAppOnConnectionFailure:        constants.RedisCrashAppOnConnectionFailure,
			ConnectRetryIntervalInSeconds:      constants.RedisConnectRetryIntervalInSeconds,
			AutoExpireTopLevelKeysAfterSeconds: constants.RedisAutoExpireTopLevelKeysAfterSeconds,
			AppNamespace:                       constants.RedisAppNameSpace,
			Address:                            constants.RedisDefaultAddr,
		},
	}
}

var ctx = context.Background()

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
	config.Redis.Main.Address = envOrViperOrDefaultString("redis.main.address", config.Redis.Main.Address)

	// TODO: work on it and implement it
	if config.Redis.Main.Enabled &&
		config.Redis.Main.OperationMode != "auto" &&
		config.Redis.Main.OperationMode != "cluster" &&
		config.Redis.Main.OperationMode != "standalone" {
		fmt.Printf("P#1MQUNR - Invalid redis operation mode. Cannot proceed.")
		os.Exit(exitCode.RedisConnectionFailed)
	} else {
		// Try to connect to redis, check if operation mode is cluster, if yes then try all possible ways to connect.
		if config.Redis.Main.OperationMode == "cluster" {
			if config.Redis.Main.Url != "" {
				opts, errConnect := goredis.ParseURL(config.Redis.Main.Url)
				if errConnect != nil {
					fmt.Println("E#1OEMOC - Could not parse the connection URL.")
					panic(errConnect)
				}
				config.Redis.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{opts.Addr},
				})
				err := config.Redis.Connection.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1PQSM2 - Could not connect to redis server. Check the url provided.")
				} else {
					fmt.Println("I#1PQSLR - Connection to redis in cluster mode established successfully.")
				}
			} else if config.Redis.Main.Address != "" {
				config.Redis.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{config.Redis.Main.Address},
				})
				err := config.Redis.Connection.ClusterNodes(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87GS - Could not connect to redis cluster. Check the address provided.")
				} else {
					fmt.Println("I#1P87H7 - Connection to redis in cluster mode established successfully.")
				}
			} else if config.Redis.Main.Username != "" {
				config.Redis.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Username: config.Redis.Main.Username,
					Password: config.Redis.Main.Password,
				})
				err := config.Redis.Connection.ClusterNodes(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87HK - Could not connect to redis cluster. Check the username provided.")
				} else {
					fmt.Println("I#1P87I1 - Connection to redis in cluster mode established successfully.")
				}
			}
			// Try to connect to redis, check if operation mode is standalone, if yes then try all possible ways to connect.
		} else if config.Redis.Main.OperationMode == "standalone" {
			if config.Redis.Main.Address != "" {
				config.Redis.Connection = goredis.NewClient(&goredis.Options{
					Addr: config.Redis.Main.Address,
				})
				err := config.Redis.Connection.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87JU - Could not connect to redis server. Check the address provided.")
				} else {
					fmt.Println("I#1P87K1 - Connection to redis in standalone mode established successfully.")
				}
			} else if config.Redis.Main.Url != "" {
				opts, err := goredis.ParseURL(config.Redis.Main.Url)
				if err != nil {
					fmt.Println("E#1OEMOC - Could not parse the connection URL.")
					panic(err)
				}
				config.Redis.Connection = goredis.NewClient(opts)
				errPing := config.Redis.Connection.Ping(ctx).Err()
				if errPing != nil {
					fmt.Println("E#1P87KH - Could not connect to redis server. Check the connection url provided.")
				} else {
					fmt.Println("I#1P87KQ - Connection to redis in standalone mode established successfully.")
				}
			} else if config.Redis.Main.Username != "" {
				config.Redis.Connection = goredis.NewClient(&goredis.Options{
					Username: config.Redis.Main.Username,
					Password: config.Redis.Main.Password,
				})
				err := config.Redis.Connection.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87KX - Could not connect to redis server. Check the username provided.")
				} else {
					fmt.Println("I#1P87L8 - Connection to redis in standalone mode established successfully.")
				}
			}
			// redis operation mode is auto, we will detect the redis mode and try to return the connection object accordingly
		} else if config.Redis.Main.OperationMode == "auto" {

			if config.Redis.Main.Url != "" {
				opts, errConnect := goredis.ParseURL(config.Redis.Main.Url)
				if errConnect != nil {
					fmt.Println("E#1OEMOC - Could not parse the connection URL.")
					panic(errConnect)
				}

				config.Redis.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
					Addrs: []string{opts.Addr},
				})
				err := config.Redis.Connection.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1PQR6V - Could not connect to redis. Check the url provided.")
				} else {
					fmt.Println("I#1PQR74 - Connection established to redis using url. Detecting redis mode...")
				}

				if hasClusterEnabled() {
					config.Redis.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
						Addrs: []string{opts.Addr},
					})
					config.Redis.Main.OperationMode = "cluster"
					fmt.Println("I#1PQR7T - Cluster mode detected! Redis is running in cluster mode.")

				} else {
					config.Redis.Main.OperationMode = "standalone"
					fmt.Println("I#1PQR8I - Standalone mode detected! Redis is running in standalone mode.")
				}

			} else if config.Redis.Main.Address != "" {
				config.Redis.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
					Addrs: []string{config.Redis.Main.Address},
				})

				err := config.Redis.Connection.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1PQRCO - Could not connect to redis. Check the address provided.")
				} else {
					fmt.Println("I#1PQRCT - Connection established to redis using address. Detecting redis mode...")
				}

				if hasClusterEnabled() {
					config.Redis.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
						Addrs: []string{config.Redis.Main.Address},
					})
					config.Redis.Main.OperationMode = "cluster"
					fmt.Println("I#1PQRCZ - Cluster mode detected! Redis is running in cluster mode.")

				} else {
					config.Redis.Main.OperationMode = "standalone"
					fmt.Println("I#1PQRD5 - Standalone mode detected! Redis is running in standalone mode.")
				}

			} else if config.Redis.Main.Username != "" {
				config.Redis.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
					Username: config.Redis.Main.Username,
					Password: config.Redis.Main.Password,
				})

				err := config.Redis.Connection.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1PQRG2 - Could not connect to redis. Check the credentials provided.")
				} else {
					fmt.Println("I#1PQRG9 - Connection established to redis using credentials. Detecting redis mode...")
				}

				if hasClusterEnabled() {
					config.Redis.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
						Addrs: []string{config.Redis.Main.Address},
					})
					config.Redis.Main.OperationMode = "cluster"
					fmt.Println("I#1PQRGF - Cluster mode detected! Redis is running in cluster mode.")

				} else {
					config.Redis.Main.OperationMode = "standalone"
					fmt.Println("I#1PQRGK - Standalone mode detected! Redis is running in standalone mode.")
				}
			}

		} else {
			fmt.Println("P#1P87M6 - No operation mode for redis found, exiting.")
			os.Exit(exitCode.RedisConnectionFailed)
		}
	}
}

func hasClusterEnabled() bool {
	err := config.Redis.Connection.ClusterNodes(ctx)
	if err.Err() != nil {
		return false
	}
	return true
}

// Cluster management methods

func (r *redis) Info() *goredis.StringCmd {
	return config.Redis.Connection.Info(ctx)
}

func (r *redis) ClusterNodes() *goredis.StringCmd {
	return config.Redis.Connection.ClusterNodes(ctx)

}

func (r *redis) Get(key string) *goredis.StringCmd {
	return config.Redis.Connection.Get(ctx, key)
}

func (r *redis) Set(key, value string) *goredis.StatusCmd {
	return config.Redis.Connection.Set(ctx, key, value, 0)
}
