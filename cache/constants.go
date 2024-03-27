package cache

const (
	ModeStandalone = "standalone"
	ModeCluster    = "cluster"
	ModeAuto       = "auto"
)

const (
	DefaultRedisUrl = "redis://localhost:6379"
)

const (
	RedisDefaultPassword                    = ""
	RedisDefaultAddr                        = "localhost:6379"
	RedisDefaultDbNumber                    = 0
	RedisDefaultProtocol                    = 3
	RedisDefaultOperationMode               = "auto"
	RedisDefaultUrl                         = "redis://localhost:6379"
	RedisEnabled                            = true //set to false once config manager starts working
	RedisDefaultUsername                    = ""
	RedisMaxActiveConnections               = 10
	RedisMaxIdleConnections                 = 10
	RedisIdleTimeoutInSeconds               = 60
	RedisCrashAppOnConnectionFailure        = false
	RedisConnectRetryIntervalInSeconds      = 10
	RedisAutoExpireTopLevelKeysAfterSeconds = 0
	RedisAppNameSpace                       = "ground-based-serviceNs"
)
