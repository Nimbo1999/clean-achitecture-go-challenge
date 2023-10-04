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

const PATH_SEPARATOR = "__"

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]*WebServerHandler),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) getHandlerKey(path string, webServerHandler *WebServerHandler) string {
	return path + PATH_SEPARATOR + webServerHandler.Method
}

func (s *WebServer) AddHandler(path string, webServerHandler *WebServerHandler) error {
	webHandler, found := s.Handlers[s.getHandlerKey(path, webServerHandler)]
	if found {
		if webServerHandler.Method == webHandler.Method {
			return ErrHandlerAlreadyExists
		}
	}
	s.Handlers[s.getHandlerKey(path, webServerHandler)] = webServerHandler
	return nil
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, webServerHandler := range s.Handlers {
		splitted := strings.Split(path, PATH_SEPARATOR)
		if len(splitted) < 1 {
			panic(errors.New("could not retrieve the handler path"))
		}
		requestMethod := strings.TrimSpace(strings.ToUpper(webServerHandler.Method))
		switch requestMethod {
		case http.MethodPost:
			s.Router.Post(splitted[0], webServerHandler.Handler)
		case http.MethodGet:
			s.Router.Get(splitted[0], webServerHandler.Handler)
		case http.MethodPatch:
			s.Router.Patch(splitted[0], webServerHandler.Handler)
		case http.MethodPut:
			s.Router.Put(splitted[0], webServerHandler.Handler)
		case http.MethodDelete:
			s.Router.Delete(splitted[0], webServerHandler.Handler)
		case http.MethodOptions:
			s.Router.Options(splitted[0], webServerHandler.Handler)
		case http.MethodHead:
			s.Router.Head(splitted[0], webServerHandler.Handler)
		case http.MethodTrace:
			s.Router.Trace(splitted[0], webServerHandler.Handler)
		case http.MethodConnect:
			s.Router.Connect(splitted[0], webServerHandler.Handler)
		default:
			s.Router.HandleFunc(splitted[0], webServerHandler.Handler)
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
