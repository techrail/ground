package ground

import (
	"github.com/techrail/ground/bgRoutine"
	"github.com/techrail/ground/dbCodegen"
	"github.com/techrail/ground/typs/appError"
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
