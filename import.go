package ground

import (
	"github.com/techrail/ground/bgRoutine"
	"github.com/techrail/ground/dbCodegen"
	"github.com/techrail/ground/logger"
	"github.com/techrail/ground/typs/appError"
	"github.com/techrail/ground/typs/jsonObject"
	"github.com/techrail/ground/webServer"
)

func GiveMeACodeGenerator(config dbCodegen.CodegenConfig) (dbCodegen.Generator, appError.Typ) {
	return dbCodegen.NewCodeGenerator(config)
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

func GiveMeABarkEmbeddedServerClient(dbUrl, defaultLogLvl, svcName, sessName string, enableSlog bool) logger.Logger {
	return logger.NewEmbeddedServerBarkClient(dbUrl, defaultLogLvl, svcName, sessName, enableSlog)
}

func GiveMeABarkRemoteServerClient(remoteServerUrl, defaultLogLvl, svcName, sessName string, enableSlog bool,
	enableBulkSend bool) logger.Logger {
	return logger.NewBarkClient(remoteServerUrl, defaultLogLvl, svcName, sessName, enableSlog, enableBulkSend)
}

func GiveMeABlankJsonObject() jsonObject.Typ {
	return jsonObject.EmptyNotNullJsonObject()
}
