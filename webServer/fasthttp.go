package webServer

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/techrail/ground/typs/appError"
	"github.com/techrail/ground/utils"
	"github.com/techrail/ground/webServer/middlewares"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
)

const (
	StateNotStarted        = "NotStarted"
	StateStarted           = "Started"
	StateShutdownRequested = "ShutdownRequested"
	StateShutdownCompleted = "ShutdownCompleted"
)

type FastHttpServer struct {
	Router       *router.Router
	Server       fasthttp.Server
	BindPort     int
	EnableIpv6   bool
	BlockOnStart bool   // Should we block on start or not
	currentState string // What is the current state of the server
	middlewares  map[string]MiddlewareSet
}

func NewLocalServer() *FastHttpServer {
	r := router.New()
	mws := map[string]MiddlewareSet{
		"Default": {
			middlewares.SetRequestId,
			middlewares.SetRandomVar,
			middlewares.CheckShutdownRequested,
		},
	}
	return &FastHttpServer{
		Router:       r,
		Server:       fasthttp.Server{Handler: r.Handler},
		BindPort:     8080,
		EnableIpv6:   false,
		BlockOnStart: false,
		currentState: StateNotStarted,
		middlewares:  mws,
	}
}

func (s *FastHttpServer) AddMiddlewareToSet(m Middleware, name string) {
	if _, ok := s.middlewares[name]; !ok {
		s.middlewares[name] = MiddlewareSet{}
	}

	s.middlewares[name] = append(s.middlewares[name], m)
}

func (s *FastHttpServer) AddMiddlewareSetToSet(mwSet MiddlewareSet, name string) {
	for _, m := range mwSet {
		s.AddMiddlewareToSet(m, name)
	}
}

func (s *FastHttpServer) EraseMiddlewares() {
	s.middlewares = map[string]MiddlewareSet{}
}

func (s *FastHttpServer) ListMiddlewareNames() []string {
	mwArr := []string{}
	for _, m := range s.middlewares {
		mwArr = append(mwArr, utils.GetFunctionName(m, false))
	}
	return mwArr
}

func (s *FastHttpServer) GetMiddleware(name string) MiddlewareSet {
	if set, ok := s.middlewares[name]; ok {
		return set
	}
	return MiddlewareSet{}
}

func (s *FastHttpServer) WithMiddlewareSet(name string, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	mwSet := s.GetMiddleware(name)
	if len(mwSet) == 0 {
		return handler
	} else {
		return chain(handler, mwSet...)
	}
}

func (s *FastHttpServer) Start() appError.Typ {
	s.Server = fasthttp.Server{
		Handler: s.Router.Handler,
	}

	var listener net.Listener
	var err error

	if s.EnableIpv6 {
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(s.BindPort))
	} else {
		listener, err = net.Listen("tcp4", ":"+strconv.Itoa(s.BindPort))
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Printf("E#1MOCDG - Could not close listener. Error: %v", err)
		}
	}(listener)

	if err != nil {
		return appError.NewError(
			appError.Error,
			"1MHI99",
			"Can't create the listener. Error: "+err.Error())
	}

	fn := func() appError.Typ {
		e := s.Server.Serve(listener)
		if e == nil {
			return appError.BlankError
		} else {
			return appError.NewError(
				appError.Error,
				"1MHOPJ",
				"Something went wrong when trying to start the server. Error: "+e.Error())
		}
	}

	s.currentState = StateStarted

	if s.BlockOnStart {
		return fn()
	} else {
		go func() {
			e := fn()
			if e.IsNotBlank() {
				fmt.Printf("E#1MOV7B - Server failed to either start or finish properly. Error: %v", e)
			}
		}()
	}

	return appError.BlankError
}

func (s *FastHttpServer) Stop() {
	s.currentState = StateShutdownRequested
}
