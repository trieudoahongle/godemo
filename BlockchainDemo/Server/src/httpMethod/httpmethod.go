package httpMethod

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type GenericHandler func(http.ResponseWriter, *http.Request)

/*
type HttpMethodHandle interface {
	checkMethod(*http.Request) bool
	Handler(http.ResponseWriter, *http.Request)
}
*/

type HttpHandler struct {
	//	 Methods [] GenericHandler
	HandlerMap map[string]GenericHandler
}

func (h *HttpHandler) checkMethod(req *http.Request) bool {
	fmt.Println("checkMethod:" + req.Method)
	return h.HandlerMap[req.Method] != nil
}

func NewHttpHandler() *HttpHandler {
	methods := make(map[string]GenericHandler)
	return &HttpHandler{methods}
}
func (h *HttpHandler) Handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handler")
	if !h.checkMethod(req) {
		fmt.Println("Method is not register.")
		return
	}
	h.HandlerMap[req.Method](w, req)
}

func (h *HttpHandler) SetHandlerMethod(httpMethodName string, handler GenericHandler) error {
	fmt.Println("SetHandlerMethod")
	if h.HandlerMap[httpMethodName] != nil {
		return errors.New("Hanler has been registered !")
	}
	h.HandlerMap[httpMethodName] = handler
	return nil
}
func (h *HttpHandler) GetHandlerMethod(httpMethodName string) (GenericHandler, error) {
	if h.HandlerMap[httpMethodName] == nil {
		return nil, errors.New("Hanler is not registered !")
	}
	return h.HandlerMap[httpMethodName], nil
}

func defaultGetHanler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("default Get Hanler")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Write([]byte("This is an example server.\n"))
}

func defaultPostHanler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("default Post Hanler")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Write([]byte("This is an example server.\n"))
}

func GetDefaultMuxServer() *http.ServeMux {
	fmt.Println("GetDefaultMuxServer")
	mux := http.NewServeMux()
	h := NewHttpHandler()
	h.SetHandlerMethod(http.MethodGet, defaultGetHanler)
	h.SetHandlerMethod(http.MethodPost, defaultPostHanler)
	mux.HandleFunc("/", h.Handler)
	return mux
}
func StartDefaultServer() *http.Server {
	fmt.Println("StartDefaultServer")
	mux := GetDefaultMuxServer()
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
	return srv
}
