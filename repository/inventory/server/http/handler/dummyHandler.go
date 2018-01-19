package handler

import (
	"net/http"

	"github.com/ncrypthic/gocontainer"
)

//DummyHandler is a specific http handler with dummy response
type DummyHandler struct {
	h Handler
}

//Handle is a function for handling http request
func (dh *DummyHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	//dummy handler just writes dummy response
	w.Write([]byte("this is a dummy response"))
	return nil
}

//NewDummyHandler creates a new Dummy Handler and returns a pointer to it
func NewDummyHandler(sc gocontainer.ServiceContainer) *DummyHandler {
	newHandler := new(DummyHandler)
	//newHandler := &DummyHandler{Handle{}}
	newHandler.h.SetContainer(sc)
	newHandler.h.Handle = newHandler.Handle

	return newHandler
}

func (dh *DummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dh.h.ServeHTTP(w, r)
}

//StartUp allows the dummy handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (dh *DummyHandler) StartUp() {
	//dummy handler does not perform any startup/initialization/bootstrapping
}

//Shutdown allows the dummy handler to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (dh *DummyHandler) Shutdown() {
	//dummy handler does not perform any cleanup
}
