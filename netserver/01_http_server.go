package netserver

import "fmt"

type NetHttpServer struct {
	Name         string // Name of the server  (used to identify it against another one, in case it is needed)
	Router       *Router
	Port         uint16
	BlockOnStart bool   // Should we block on start or not
	currentState string // What is the current state of the server
}

func NewServer(port uint16) *NetHttpServer {
	router := NewRouter()
	return &NetHttpServer{
		Port:   port,
		Router: router,
	}
}

func (h *NetHttpServer) PortString() string {
	return fmt.Sprintf("%d", h.Port)
}

// File ends here
