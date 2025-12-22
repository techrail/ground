package netserver

import (
	"fmt"
	"net/http"

	"github.com/techrail/ground/constants"
	"github.com/techrail/ground/typs/appError"
)

type NetHttpServer struct {
	Name         string    // Name of the server  (used to identify it against another one, in case it is needed)
	Router       *Router   // Associated Router
	Render       *Renderer // Renderer attached to the server. It is here only for ease-of-use
	BlockOnStart bool      // Should we block on start or not
	port         uint16    // The port on which the Server will start listening
	currentState string    // What is the current state of the server
}

func NewServer(port uint16, blockOnStart bool) *NetHttpServer {
	router := NewRouter()
	return &NetHttpServer{
		port:         port,
		Router:       router,
		BlockOnStart: blockOnStart,
	}
}

func (s *NetHttpServer) PortString() string {
	return fmt.Sprintf("%d", s.port)
}

func (s *NetHttpServer) Start() appError.Typ {
	fmt.Println("I#2R2NLB: About to start ", s.Name)

	// Ensure that the Router was configured
	if s.Router == nil {
		return appError.NewError(appError.Error, "2R2Q6W", "Router was nil. Cannot proceed.")
	}

	// Make sure that the Router has routes according to how it was built
	if !s.Router.HasRoutes() {
		appError.NewError(
			appError.Error, "2R2QY3",
			fmt.Sprintf("Router is %v and has no routes in that category", s.Router.Type()))
	}

	if s.BlockOnStart {
		err := http.ListenAndServe(":"+s.PortString(), s.Router)
		if err != nil {
			fmt.Printf("E#%v: %v\n", constants.ErrWebServerStartFailed, err)
			return appError.NewFromExisting(err, "2R2QGK")
		}
	} else {
		go func() {
			err := http.ListenAndServe(":"+s.PortString(), s.Router)
			if err != nil {
				fmt.Printf("E#%v: %v\n", constants.ErrWebServerStartFailed, err)
				return
			}
		}()
	}

	return appError.BlankError
}

// File ends here
