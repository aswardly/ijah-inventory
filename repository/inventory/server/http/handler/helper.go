package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-errors/errors"
)

var (
	//ErrCodeSuccessful is the error code for successful operation
	ErrCodeSuccessful = "S"
	//ErrCodeFailed is the error code for failed operation
	ErrCodeFailed = "F"
)

//SimpleResponseStruct is representation of simple response returned to the http client
type SimpleResponseStruct struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//composeJSONResponse is a helper function for composing JSON string for http response (will be displayed to user's browser)
//input is expected to be a struct to be marshalled to json
func composeJSONResponse(input interface{}) (string, *StatusError) {
	responseJSON, err := json.Marshal(input)
	if err != nil {
		return "", &StatusError{
			Code:          http.StatusInternalServerError,
			ReturnMessage: string("JSON marshalling error: " + err.Error()),
			Err:           errors.Wrap(err, 0),
		}
	}
	return string(responseJSON), nil
}

//composeJSONError is a helper function for composing JSON string for http response and StatusError object with the JSON string as its return message (will be displayed to user's browser)
//input is expected to be a struct to be marshalled to json
func composeJSONError(input interface{}) *StatusError {
	//compose return json format
	var returnMessage string
	responseJSON, err := json.Marshal(input)
	if err != nil {
		returnMessage = string("JSON marshalling error: " + err.Error())
	} else {
		returnMessage = string(responseJSON)
	}
	return &StatusError{
		Code:          http.StatusInternalServerError,
		ReturnMessage: returnMessage,
		Err:           errors.Wrap(err, 0),
	}
}

//composeError returns a StatusError from a given error object
func composeError(err error) *StatusError {
	response := SimpleResponseStruct{}
	response.Code = ErrCodeFailed
	response.Message = "Error: " + err.Error()
	statusErr := composeJSONError(response)
	return statusErr
}
