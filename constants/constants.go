package constants

const EmptyString = ""
const EmptyInt = int(0)

// TODO: move this value to config later

const OpLogRequestValue = "TECHRAIL_GROUND_OPLOG_REQUEST_VALUE"

const (
	RedisDefaultPassword                    = ""
	RedisDefaultAddr                        = "localhost:7006"
	RedisDefaultDbNumber                    = 0
	RedisDefaultProtocol                    = 3
	RedisDefaultOperationMode               = "auto"
	RedisDefaultUrl                         = "redis://127.0.0.1:7006"
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
