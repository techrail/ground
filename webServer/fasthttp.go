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

// NewLocalServer creates a basic new local server and returns it.
// It can then be modified and started (or started as it is)
func NewLocalServer() *FastHttpServer {
	r := router.New()
	mws := map[string]MiddlewareSet{
		"Default": {
			middlewares.SetRequestId,
			middlewares.SetRandomVar,
			middlewares.CheckShutdownRequested,
			middlewares.CheckOpLogRequest,
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

// AddMiddlewareToSet will add a middleware to a single middleware set of a given name
// If a middleware set of that name does not exist, then one will be created.
func (s *FastHttpServer) AddMiddlewareToSet(m Middleware, name string) {
	if _, ok := s.middlewares[name]; !ok {
		s.middlewares[name] = MiddlewareSet{}
	}

	s.middlewares[name] = append(s.middlewares[name], m)
}

// AddMiddlewareSetToSet adds all the middlewares of one set to another set. Useful when you want to have some base
// middleware sets and you want to create others with those middlewares as common with the new ones
func (s *FastHttpServer) AddMiddlewareSetToSet(mwSet MiddlewareSet, name string) {
	for _, m := range mwSet {
		s.AddMiddlewareToSet(m, name)
	}
}

// EraseMiddlewares removes all the middlewares from the server so that you can create a whole new set from scratch
func (s *FastHttpServer) EraseMiddlewares() {
	s.middlewares = map[string]MiddlewareSet{}
}

// ListMiddlewareNames is supposed to return the list of middlewares in the server
func (s *FastHttpServer) ListMiddlewareNames() []string {
	mwArr := []string{}
	for setName, middlewareSet := range s.middlewares {
		for _, m := range middlewareSet {
			mwArr = append(mwArr, fmt.Sprintf("%v - %v", setName, utils.GetFunctionName(m, false)))
		}
	}
	return mwArr
}

// GetMiddlewareSet returns the middleware set by the given name. If there is no middleware set by that name, a blank
// one is returned instead
func (s *FastHttpServer) GetMiddlewareSet(name string) MiddlewareSet {
	if set, ok := s.middlewares[name]; ok {
		return set
	}
	return MiddlewareSet{}
}

// WithMiddlewareSet will return a request handler which would contain all the middlewares given by name applied
// to the handler supplied to this function.
func (s *FastHttpServer) WithMiddlewareSet(name string, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	mwSet := s.GetMiddlewareSet(name)
	if len(mwSet) == 0 {
		return handler
	} else {
		return chain(handler, mwSet...)
	}
}

// Start starts the web server according to given parameters
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

// Stop will stop the server. It does so by setting the current state. The manager will notice the change
// and stop the server gracefully
func (s *FastHttpServer) Stop() {
	s.currentState = StateShutdownRequested
}
