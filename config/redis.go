package config

import (
	"context"
	"fmt"
	"os"

	goredis "github.com/redis/go-redis/v9"
	"github.com/techrail/ground/constants/exitCode"
)

type redis struct {
	Main        redisConfig
	RedisClient connection
}

type Cache struct {
	config redis
	client *connection
}

type connection struct {
	SingleClient  *goredis.Client
	ClusterClient *goredis.ClusterClient
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
			Enabled:                            true,
			Url:                                "redis://127.0.0.1:6379",
			Username:                           "",
			Password:                           "",
			OperationMode:                      "cluster",
			MaxActiveConnections:               10,
			MaxIdleConnections:                 10,
			IdleTimeoutInSeconds:               60,
			CrashAppOnConnectionFailure:        false,
			ConnectRetryIntervalInSeconds:      10,
			AutoExpireTopLevelKeysAfterSeconds: 0,
			AppNamespace:                       "ground-based-serviceNs",
			Address:                            ":7000",
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
			// Check using single/multiple addresses (host:port combinations)
			if config.Redis.Main.Address != "" {
				config.Redis.RedisClient.ClusterClient = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{config.Redis.Main.Address},
				})
				err := config.Redis.RedisClient.ClusterClient.ClusterNodes(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87GS - Could not connect to redis cluster. Check the address provided.")
				} else {
					fmt.Println("I#1P87H7 - Connection established to redis cluster using node(s) address.")
				}
			} else if config.Redis.Main.Username != "" {
				config.Redis.RedisClient.ClusterClient = goredis.NewClusterClient(&goredis.ClusterOptions{
					Username: config.Redis.Main.Username,
				})
				err := config.Redis.RedisClient.ClusterClient.ClusterNodes(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87HK - Could not connect to redis cluster. Check the username provided.")
				} else {
					fmt.Println("I#1P87I1 - Connection established to redis cluster using node username.")
				}
			} else if config.Redis.Main.Username != "" && config.Redis.Main.Password != "" {
				config.Redis.RedisClient.ClusterClient = goredis.NewClusterClient(&goredis.ClusterOptions{
					Username: config.Redis.Main.Username,
					Password: config.Redis.Main.Password,
				})
				err := config.Redis.RedisClient.ClusterClient.ClusterNodes(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87JF - Could not connect to redis cluster. Check the username and password provided.")
				} else {
					fmt.Println("I#1P87JK - Connection established to redis cluster using node username and password.")
				}
			}
			// Try to connect to redis, check if operation mode is standalone, if yes then try all possible ways to connect.
			// If not we will connect using auto mode, which is, standalone mode.
		} else if config.Redis.Main.OperationMode == "standalone" || config.Redis.Main.OperationMode == "auto" {
			if config.Redis.Main.Address != "" {
				config.Redis.RedisClient.SingleClient = goredis.NewClient(&goredis.Options{
					Addr: config.Redis.Main.Address,
				})
				err := config.Redis.RedisClient.SingleClient.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87JU - Could not connect to redis server. Check the address provided.")
				} else {
					fmt.Println("I#1P87K1 - Connection established to redis server using server address.")
				}
			} else if config.Redis.Main.Url != "" {
				opts, err := goredis.ParseURL(config.Redis.Main.Url)
				if err != nil {
					fmt.Println("E#1OEMOC - Could not parse the connection URL.")
					panic(err)
				}
				config.Redis.RedisClient.SingleClient = goredis.NewClient(opts)
				errPing := config.Redis.RedisClient.SingleClient.Ping(ctx).Err()
				if errPing != nil {
					fmt.Println("E#1P87KH - Could not connect to redis server. Check the connection url provided.")
				} else {
					fmt.Println("I#1P87KQ - Connection established to redis server using connection url.")
				}
			} else if config.Redis.Main.Username != "" {
				config.Redis.RedisClient.SingleClient = goredis.NewClient(&goredis.Options{
					Username: config.Redis.Main.Username,
				})
				err := config.Redis.RedisClient.SingleClient.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87KX - Could not connect to redis server. Check the username provided.")
				} else {
					fmt.Println("I#1P87L8 - Connection established to redis server using instance username.")
				}
			} else if config.Redis.Main.Username != "" && config.Redis.Main.Password != "" {
				config.Redis.RedisClient.SingleClient = goredis.NewClient(&goredis.Options{
					Username: config.Redis.Main.Username,
					Password: config.Redis.Main.Password,
				})
				err := config.Redis.RedisClient.SingleClient.Ping(ctx).Err()
				if err != nil {
					fmt.Println("E#1P87LK - Could not connect to redis server. Check the username and password provided.")
				} else {
					fmt.Println("I#1P87LS - Connection established to redis cluster using instance username and password.")
				}
			}
		} else {
			fmt.Println("P#1P87M6 - No operation mode for redis found, exiting.")
			os.Exit(exitCode.RedisConnectionFailed)
		}
	}
}

// Cluster management methods

func (c *Cache) Info() *goredis.StringCmd {
	if c.client.ClusterClient != nil {
		return config.Redis.RedisClient.ClusterClient.ClusterInfo(ctx)
	} else {
		return config.Redis.RedisClient.SingleClient.Info(ctx)
	}
}

// func (c *Cache) ClusterNodes() *goredis.StringCmd {
// 	if c.client.ClusterClient != nil {
// 		return config.Redis.RedisClient.ClusterClient.ClusterNodes(ctx)
// 	} else {
// 		cmd := NewStringCmd(ctx, "cluster", "cluster not enabled")
// 		return &goredis.StringCmd{cmd}
// 	}
// }

func (c *Cache) PoolStats() *goredis.PoolStats {
	if c.client.ClusterClient != nil {
		return config.Redis.RedisClient.ClusterClient.PoolStats()
	} else {
		return config.Redis.RedisClient.SingleClient.PoolStats()
	}

}

// func (c *Cache) ForEachShard(client *goredis.Client,fn func(contxt context.Context, goredis.Client)) goredis.Error {
// 	if c.client.ClusterClient != nil {
// 		return config.Redis.RedisClient.ClusterClient.ForEachShard(ctx, fn, )
// 	} else {
//
// 	}
// }
