package webServer

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/techrail/goround/typs/appError"
	"github.com/techrail/goround/typs/errLevel"
	"github.com/valyala/fasthttp"
	"net"
	"strconv"
)

type FastHttpServer struct {
	Router       *router.Router
	Server       fasthttp.Server
	BindPort     int
	EnableIpv6   bool
	BlockOnStart bool // Should we block on start or not
}

func NewLocalServer() *FastHttpServer {
	r := router.New()
	return &FastHttpServer{
		Router:       r,
		Server:       fasthttp.Server{Handler: r.Handler},
		BindPort:     8080,
		EnableIpv6:   false,
		BlockOnStart: false,
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
			errLevel.Error,
			"1MHI99",
			"Can't create the listener. Error: "+err.Error())
	}

	fn := func() appError.Typ {
		e := s.Server.Serve(listener)
		if e == nil {
			return appError.BlankError
		} else {
			return appError.NewError(
				errLevel.Error,
				"1MHOPJ",
				"Something went wrong when trying to start the server. Error: "+e.Error())
		}
	}

	if s.BlockOnStart {
		return fn()
	} else {
		go func() {
			e := fn()
			fmt.Printf("E#1MOV7B - Server failed to either start or finish properly. Error: %v", e)
		}()
	}

	return appError.BlankError
}
