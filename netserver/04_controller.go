package netserver

import (
	"net/http"

	"github.com/techrail/ground/constants/httpheaders"
)

type Controller struct {
	Server *NetHttpServer
}

func (c *Controller) Testz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(httpheaders.ContentType, "text/json; charset=utf-8")
	_, _ = w.Write([]byte(`{"test":"ok"}`))
}

func NewController(server *NetHttpServer) *Controller {
	if server == nil {
		panic("P#2RD6IF: Cannot create a new controller using a nil server")
	}

	return &Controller{
		Server: server,
	}
}

func NewControllerWithServer(port uint16, blockOnStart bool) (*NetHttpServer, *Controller) {
	// First create a new Server
	server := NewServer(port, blockOnStart)
	// Create the controller based on this server
	controller := NewController(server)

	return server, controller
}
