package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/techrail/ground/config"
	"github.com/techrail/ground/constants/exitCode"
)

type CacheClient struct {
	Connection goredis.UniversalClient
}

var ctx = context.Background()

func CreateNewRedisClient() *CacheClient {
	c := CacheClient{}
	// Try to connect to redis, check if operation mode is cluster, if yes then try all possible ways to connect.
	if config.Store().Redis.Main.OperationMode == "cluster" {
		if config.Store().Redis.Main.Url != "" {
			opts, errConnect := goredis.ParseURL(config.Store().Redis.Main.Url)
			if errConnect != nil {
				fmt.Println("E#1OEMOC - Could not parse the connection URL.")
				panic(errConnect)
			}
			c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
				Addrs: []string{opts.Addr},
			})
			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1PQSM2 - Could not connect to redis server. Check the url provided.")
			} else {
				fmt.Println("I#1PQSLR - Connection to redis in cluster mode established successfully.")
			}
		} else if config.Store().Redis.Main.Address != "" {
			c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
				Addrs: []string{config.Store().Redis.Main.Address},
			})
			err := c.Connection.ClusterNodes(ctx).Err()
			if err != nil {
				fmt.Println("E#1P87GS - Could not connect to redis cluster. Check the address provided.")
			} else {
				fmt.Println("I#1P87H7 - Connection to redis in cluster mode established successfully.")
			}
		} else if config.Store().Redis.Main.Username != "" {
			c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
				Username: config.Store().Redis.Main.Username,
				Password: config.Store().Redis.Main.Password,
			})
			err := c.Connection.ClusterNodes(ctx).Err()
			if err != nil {
				fmt.Println("E#1P87HK - Could not connect to redis cluster. Check the username provided.")
			} else {
				fmt.Println("I#1P87I1 - Connection to redis in cluster mode established successfully.")
			}
		}
		// Try to connect to redis, check if operation mode is standalone, if yes then try all possible ways to connect.
	} else if config.Store().Redis.Main.OperationMode == "standalone" {
		if config.Store().Redis.Main.Address != "" {
			c.Connection = goredis.NewClient(&goredis.Options{
				Addr: config.Store().Redis.Main.Address,
			})
			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1P87JU - Could not connect to redis server. Check the address provided.")
			} else {
				fmt.Println("I#1P87K1 - Connection to redis in standalone mode established successfully.")
			}
		} else if config.Store().Redis.Main.Url != "" {
			opts, err := goredis.ParseURL(config.Store().Redis.Main.Url)
			if err != nil {
				fmt.Println("E#1OEMOC - Could not parse the connection URL.")
				panic(err)
			}
			c.Connection = goredis.NewClient(opts)
			errPing := c.Connection.Ping(ctx).Err()
			if errPing != nil {
				fmt.Println("E#1P87KH - Could not connect to redis server. Check the connection url provided.")
			} else {
				fmt.Println("I#1P87KQ - Connection to redis in standalone mode established successfully.")
			}
		} else if config.Store().Redis.Main.Username != "" {
			c.Connection = goredis.NewClient(&goredis.Options{
				Username: config.Store().Redis.Main.Username,
				Password: config.Store().Redis.Main.Password,
			})
			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1P87KX - Could not connect to redis server. Check the username provided.")
			} else {
				fmt.Println("I#1P87L8 - Connection to redis in standalone mode established successfully.")
			}
		}
		// redis operation mode is auto, we will detect the redis mode and try to return the connection object accordingly
	} else if config.Store().Redis.Main.OperationMode == "auto" {

		if config.Store().Redis.Main.Url != "" {
			opts, errConnect := goredis.ParseURL(config.Store().Redis.Main.Url)
			if errConnect != nil {
				fmt.Println("E#1OEMOC - Could not parse the connection URL.")
				panic(errConnect)
			}

			c.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
				Addrs: []string{opts.Addr},
			})
			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1PQR6V - Could not connect to redis. Check the url provided.")
			} else {
				fmt.Println("I#1PQR74 - Connection established to redis using url. Detecting redis mode...")
			}

			if hasClusterEnabled(&c) {
				c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{opts.Addr},
				})
				config.Store().Redis.Main.OperationMode = "cluster"
				fmt.Println("I#1PQR7T - Cluster mode detected! Redis is running in cluster mode.")

			} else {
				config.Store().Redis.Main.OperationMode = "standalone"
				fmt.Println("I#1PQR8I - Standalone mode detected! Redis is running in standalone mode.")
			}

		} else if config.Store().Redis.Main.Address != "" {
			c.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
				Addrs: []string{config.Store().Redis.Main.Address},
			})

			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1PQRCO - Could not connect to redis. Check the address provided.")
			} else {
				fmt.Println("I#1PQRCT - Connection established to redis using address. Detecting redis mode...")
			}

			if hasClusterEnabled(&c) {
				c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{config.Store().Redis.Main.Address},
				})
				config.Store().Redis.Main.OperationMode = "cluster"
				fmt.Println("I#1PQRCZ - Cluster mode detected! Redis is running in cluster mode.")

			} else {
				config.Store().Redis.Main.OperationMode = "standalone"
				fmt.Println("I#1PQRD5 - Standalone mode detected! Redis is running in standalone mode.")
			}

		} else if config.Store().Redis.Main.Username != "" {
			c.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
				Username: config.Store().Redis.Main.Username,
				Password: config.Store().Redis.Main.Password,
			})

			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1PQRG2 - Could not connect to redis. Check the credentials provided.")
			} else {
				fmt.Println("I#1PQRG9 - Connection established to redis using credentials. Detecting redis mode...")
			}

			if hasClusterEnabled(&c) {
				c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{config.Store().Redis.Main.Address},
				})
				config.Store().Redis.Main.OperationMode = "cluster"
				fmt.Println("I#1PQRGF - Cluster mode detected! Redis is running in cluster mode.")

			} else {
				config.Store().Redis.Main.OperationMode = "standalone"
				fmt.Println("I#1PQRGK - Standalone mode detected! Redis is running in standalone mode.")
			}
		}

	} else {
		fmt.Println("P#1P87M6 - No operation mode for redis found, exiting.")
		os.Exit(exitCode.RedisConnectionFailed)
	}

	return &c
}

func hasClusterEnabled(c *CacheClient) bool {
	err := c.Connection.ClusterNodes(ctx)
	if err.Err() != nil {
		return false
	}
	return true
}

func (c *CacheClient) Info() *goredis.StringCmd {
	return c.Connection.Info(ctx)
}

func (c *CacheClient) ClusterNodes() *goredis.StringCmd {
	return c.Connection.ClusterNodes(ctx)
}

func (c *CacheClient) GetString(key string) *goredis.StringCmd {
	return c.Connection.Get(ctx, key)
}

func (c *CacheClient) SetString(key string, value interface{}) *goredis.StatusCmd {
	return c.Connection.Set(ctx, key, value, 0)
}

// Using with versions lower than 6.2.0 to delete string after specified time
func (c *CacheClient) SetStringWithExpiry(key string, value string, expiration time.Duration) *goredis.StatusCmd {
	return c.Connection.SetEx(ctx, key, value, expiration)
}

func (c *CacheClient) SetListContents(key string, values interface{}) *goredis.IntCmd {
	return c.Connection.LPush(ctx, key, values)
}

func (c *CacheClient) GetListRange(key string, start, end int64) *goredis.StringSliceCmd {
	return c.Connection.LRange(ctx, key, start, end)
}

func (c *CacheClient) SetHash(key string, values interface{}) *goredis.IntCmd {
	return c.Connection.HSet(ctx, key, values)
}

func (c *CacheClient) GetHashVals(key string) *goredis.StringSliceCmd {
	return c.Connection.HVals(ctx, key)
}

func (c *CacheClient) SetAdd(key string, members interface{}) *goredis.IntCmd {
	return c.Connection.SAdd(ctx, key, members)
}

func (c *CacheClient) GetSetMembers(key string) *goredis.StringSliceCmd {
	return c.Connection.SMembers(ctx, key)
}

func (c *CacheClient) DeleteListElements(key string) *goredis.StringCmd {
	return c.Connection.RPop(ctx, key)
}

// Works only with Redis >= 6.2.0, use `SetStringWithExpiry` with versions lower than 6.2
func (c *CacheClient) DeleteString(key string) *goredis.StringCmd {
	return c.Connection.GetDel(ctx, key)
}

func (c *CacheClient) DeleteHash(key string, fields ...string) *goredis.IntCmd {
	return c.Connection.HDel(ctx, key, fields...)
}

func (c *CacheClient) DeleteSet(key string) *goredis.StringCmd {
	return c.Connection.SPop(ctx, key)
}
