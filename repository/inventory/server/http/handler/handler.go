//Package handler defines the http handlers
//Taken and adapted from: https://elithrar.github.io/article/http-handler-error-handling-revisited/
package handler

import (
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/ncrypthic/gocontainer"
	log "github.com/sirupsen/logrus"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code          int           //represents the HTTP status code to return to http client
	Err           *errors.Error //the error itself
	ReturnMessage string        //the message to return to the http client
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

//Status is a function for returning HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

//Handler is a struct representing http handler
type Handler struct {
	Sc     gocontainer.ServiceContainer                       //the service container object (servicemanager)
	Handle func(w http.ResponseWriter, r *http.Request) error //http handler function, returns an implementation of the built in error interface
}

//SetContainer allows our handler to satisfy gocontainer.ContainerAware interface
func (h *Handler) SetContainer(sc gocontainer.ServiceContainer) {
	h.Sc = sc
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Handle == nil {
		log.WithFields(log.Fields{
			"code":            http.StatusInternalServerError,
			"returnedMessage": "nil handler error",
			"errorMessage":    "nil handler",
			"trace":           "nil handler",
		}).Error("nil handler")
		http.Error(w, "nil Handler", http.StatusInternalServerError)
		return
	}
	err := h.Handle(w, r)
	if err != nil {
		switch e := err.(type) {
		case *StatusError:
			//retrieve and log the specific HTTP status code
			log.WithFields(log.Fields{
				"code":            e.Code,
				"returnedMessage": e.ReturnMessage,
				"errorMessage":    e.Err.Error(),
				"trace":           e.Err.ErrorStack(),
			}).Error(e.Err.Error())
			http.Error(w, e.ReturnMessage, e.Code)
			fmt.Printf("Return message: %+v\n", e.ReturnMessage)
			fmt.Printf("Stack trace: %+v\n", e.Err.ErrorStack())
		default:
			log.WithFields(log.Fields{
				"code":            http.StatusInternalServerError,
				"returnedMessage": err.Error(),
			}).Error(err.Error())
			//For any other types, set response HTTP status to 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Printf("Return message: %+v\n", err.Error())
		}
	}
	//NOTE: it is also possible to log every requests handled and the returned responses here
}
