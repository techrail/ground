package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/techrail/ground/constants/exitCode"
)

type RedisConfig struct {
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

type Client struct {
	Connection goredis.UniversalClient
}

var ctx = context.Background()

func CreateNewRedisClient(config RedisConfig) *Client {
	c := Client{}
	// Try to connect to redis, check if operation mode is cluster, if yes then try all possible ways to connect.
	if config.OperationMode == ModeCluster {
		if config.Url != "" {
			opts, errConnect := goredis.ParseURL(config.Url)
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
		} else if config.Username != "" {
			c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
				Username: config.Username,
				Password: config.Password,
			})
			err := c.Connection.ClusterNodes(ctx).Err()
			if err != nil {
				fmt.Println("E#1P87HK - Could not connect to redis cluster. Check the username provided.")
			} else {
				fmt.Println("I#1P87I1 - Connection to redis in cluster mode established successfully.")
			}
		}
		// Try to connect to redis, check if operation mode is standalone, if yes then try all possible ways to connect.
	} else if config.OperationMode == ModeStandalone {
		if config.Url != "" {
			opts, err := goredis.ParseURL(config.Url)
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
		} else if config.Username != "" {
			c.Connection = goredis.NewClient(&goredis.Options{
				Username: config.Username,
				Password: config.Password,
			})
			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1P87KX - Could not connect to redis server. Check the username provided.")
			} else {
				fmt.Println("I#1P87L8 - Connection to redis in standalone mode established successfully.")
			}
		}
		// redis operation mode is auto, we will detect the redis mode and try to return the connection object accordingly
	} else if config.OperationMode == ModeAuto {
		if config.Url != "" {
			opts, errConnect := goredis.ParseURL(config.Url)
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
				config.OperationMode = "cluster"
				fmt.Println("I#1PQR7T - Cluster mode detected! Redis is running in cluster mode.")

			} else {
				config.OperationMode = "standalone"
				fmt.Println("I#1PQR8I - Standalone mode detected! Redis is running in standalone mode.")
			}
		} else if config.Username != "" {
			c.Connection = goredis.NewUniversalClient(&goredis.UniversalOptions{
				Username: config.Username,
				Password: config.Password,
			})

			err := c.Connection.Ping(ctx).Err()
			if err != nil {
				fmt.Println("E#1PQRG2 - Could not connect to redis. Check the credentials provided.")
			} else {
				fmt.Println("I#1PQRG9 - Connection established to redis using credentials. Detecting redis mode...")
			}

			if hasClusterEnabled(&c) {
				c.Connection = goredis.NewClusterClient(&goredis.ClusterOptions{
					Addrs: []string{config.Url},
				})
				config.OperationMode = "cluster"
				fmt.Println("I#1PQRGF - Cluster mode detected! Redis is running in cluster mode.")

			} else {
				config.OperationMode = "standalone"
				fmt.Println("I#1PQRGK - Standalone mode detected! Redis is running in standalone mode.")
			}
		}
	} else {
		fmt.Println("P#1P87M6 - No operation mode for redis found, exiting.")
		os.Exit(exitCode.RedisConnectionFailed)
	}

	return &c
}

func hasClusterEnabled(c *Client) bool {
	err := c.Connection.ClusterNodes(ctx)
	if err.Err() != nil {
		return false
	}
	return true
}

func (c *Client) Info() *goredis.StringCmd {
	return c.Connection.Info(ctx)
}

func (c *Client) ClusterNodes() *goredis.StringCmd {
	return c.Connection.ClusterNodes(ctx)
}

func (c *Client) Get(key string) *goredis.StringCmd {
	return c.Connection.Get(ctx, key)
}

func (c *Client) Set(key, value string) *goredis.StatusCmd {
	return c.Connection.Set(ctx, key, value, 0)
}

// Using with versions lower than 6.2.0 to delete string after specified time
func (c *Client) SetStringWithExpiry(key string, value string, expiration time.Duration) *goredis.StatusCmd {
	return c.Connection.SetEx(ctx, key, value, expiration)
}

func (c *Client) SetListContents(key string, values interface{}) *goredis.IntCmd {
	return c.Connection.LPush(ctx, key, values)
}

func (c *Client) GetListRange(key string, start, end int64) *goredis.StringSliceCmd {
	return c.Connection.LRange(ctx, key, start, end)
}

func (c *Client) SetHash(key string, values interface{}) *goredis.IntCmd {
	return c.Connection.HSet(ctx, key, values)
}

func (c *Client) GetHashVals(key string) *goredis.StringSliceCmd {
	return c.Connection.HVals(ctx, key)
}

func (c *Client) SetAdd(key string, members interface{}) *goredis.IntCmd {
	return c.Connection.SAdd(ctx, key, members)
}

func (c *Client) GetSetMembers(key string) *goredis.StringSliceCmd {
	return c.Connection.SMembers(ctx, key)
}

func (c *Client) DeleteListElements(key string) *goredis.StringCmd {
	return c.Connection.RPop(ctx, key)
}

// Works only with Redis >= 6.2.0, use `SetStringWithExpiry` with versions lower than 6.2
func (c *Client) DeleteString(key string) *goredis.StringCmd {
	return c.Connection.GetDel(ctx, key)
}

func (c *Client) DeleteHash(key string, fields ...string) *goredis.IntCmd {
	return c.Connection.HDel(ctx, key, fields...)
}

func (c *Client) DeleteSet(key string) *goredis.StringCmd {
	return c.Connection.SPop(ctx, key)
}
