package ground

import (
	"github.com/techrail/ground/bgRoutine"
	"github.com/techrail/ground/cache"
	"github.com/techrail/ground/dbcodegen"
	"github.com/techrail/ground/logger"
	"github.com/techrail/ground/typs/appError"
	"github.com/techrail/ground/typs/jsonObject"
	"github.com/techrail/ground/webServer"
)

func GiveMeACodeGenerator(config dbcodegen.CodegenConfig) (dbcodegen.Generator, appError.Typ) {
	return dbcodegen.NewCodeGenerator(config)
}

func GiveMeAWebServer() *webServer.FastHttpServer {
	return webServer.NewLocalServer()
}

func GiveMeARoutineManager() bgRoutine.Manager {
	return bgRoutine.NewManager()
}

func GiveMeABarkSLogger() logger.Logger {
	return logger.NewSloggerClient()
}

func GiveMeADirectToDbBarkClient(dbUrl, defaultLogLvl, svcName, sessName string, enableSlog bool) logger.Logger {
	return logger.NewDirectToDbBarkClient(dbUrl, defaultLogLvl, svcName, sessName, enableSlog)
}

func GiveMeANewDirectToDbBarkClientCustomSchemaTable(dbUrl, schemaName, tableName, defaultLogLvl, svcName, sessName string, enableSlog bool) logger.Logger {
	return logger.NewDirectToDbBarkClientCustomSchemaTable(dbUrl, schemaName, tableName, defaultLogLvl, svcName, sessName, enableSlog)
}

func GiveMeABarkRemoteServerClient(remoteServerUrl, defaultLogLvl, svcName, sessName string, enableSlog bool,
	enableBulkSend bool,
) logger.Logger {
	return logger.NewBarkClient(remoteServerUrl, defaultLogLvl, svcName, sessName, enableSlog, enableBulkSend)
}

func GiveMeABlankJsonObject() jsonObject.Typ {
	return jsonObject.EmptyNotNullJsonObject()
}

func GiveMeACacheManager(config cache.RedisConfig) *cache.Client {
	return cache.CreateNewRedisClient(config)
}
