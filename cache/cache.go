package cache

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
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

	if !config.Enabled {
		fmt.Println("P#1UF5RK - Enable redis to acquire a connection object.")
		return nil
	}

	if config.Enabled &&
		config.OperationMode != ModeCluster &&
		config.OperationMode != ModeStandalone &&
		config.OperationMode != ModeAuto {
		fmt.Printf("P#1MQUNR - Invalid redis operation mode. Cannot proceed.")
		return nil
	}

	if config.Url != "" {
		switch config.OperationMode {
		case ModeCluster:
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

		case ModeStandalone:
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

		case ModeAuto:
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
				config.OperationMode = ModeCluster
				fmt.Println("I#1PQR7T - Cluster mode detected! Redis is running in cluster mode.")

			} else {
				config.OperationMode = ModeStandalone
				fmt.Println("I#1PQR8I - Standalone mode detected! Redis is running in standalone mode.")
			}

		}

	} else {
		fmt.Println("P#1P87M6 - No connection url found for redis.")
		return nil
	}

	return &c
}

func hasClusterEnabled(c *Client) bool {
	err := c.Connection.ClusterNodes(ctx)
	return err.Err() == nil
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
