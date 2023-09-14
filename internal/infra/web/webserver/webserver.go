package webserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var ErrHandlerAlreadyExists = errors.New("a handler is already registered with that http method")

type WebServerHandler struct {
	Handler http.HandlerFunc
	Method  string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]*WebServerHandler
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]*WebServerHandler),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, webServerHandler *WebServerHandler) error {
	webHandler, found := s.Handlers[path]
	if found {
		if webServerHandler.Method == webHandler.Method {
			return ErrHandlerAlreadyExists
		}
	}
	s.Handlers[path] = webServerHandler
	return nil
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, webServerHandler := range s.Handlers {
		requestMethod := strings.TrimSpace(strings.ToUpper(webServerHandler.Method))
		switch requestMethod {
		case http.MethodPost:
			s.Router.Post(path, webServerHandler.Handler)
		case http.MethodGet:
			s.Router.Get(path, webServerHandler.Handler)
		case http.MethodPatch:
			s.Router.Patch(path, webServerHandler.Handler)
		case http.MethodPut:
			s.Router.Put(path, webServerHandler.Handler)
		case http.MethodDelete:
			s.Router.Delete(path, webServerHandler.Handler)
		case http.MethodOptions:
			s.Router.Options(path, webServerHandler.Handler)
		case http.MethodHead:
			s.Router.Head(path, webServerHandler.Handler)
		case http.MethodTrace:
			s.Router.Trace(path, webServerHandler.Handler)
		case http.MethodConnect:
			s.Router.Connect(path, webServerHandler.Handler)
		default:
			s.Router.HandleFunc(path, webServerHandler.Handler)
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
