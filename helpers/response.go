/*
response provides an easy way to produce standardized & generic API responses

This file should be converted to easyjson, and cleaned up a bit. This was transferred from a
previous project.
I encourage having built in error codes within each response, rather than using HTTP errors, as
these errors should be reserved for reporting issues assoicated with the API service.
*/

package helpers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

/*
Response is a standardized response for this API which
is encoded in the body of a HTTP response
*/
type Response struct {
	wtr       http.ResponseWriter
	req       *http.Request
	Ok        bool         `json:"ok"`
	ErrorCode int          `json:"errorCode"`
	Message   string       `json:"message"`
	Body      ResponseBody `json:"body"`
}

// ResponseBody is the body of a response belonging to Response
type ResponseBody struct {
	ValidationErrors []ValidationError `json:"validationErrors"`
	Data             interface{}       `json:"data"`
}

// ValidationError represents a form or field error
type ValidationError struct {
	Err  string `json:"error"`
	Path string `json:"path"`
}

const (
	ErrorNoError         = iota
	ErrorValidationError = iota
	ErrorInternalError   = iota
	ErrorNotAuthorized   = iota
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NewResponse creates a new reply
func NewResponse(w http.ResponseWriter, r *http.Request) Response {
	resp := Response{}
	resp.wtr = w
	resp.req = r
	resp.ErrorCode = ErrorNoError
	return resp
}

// Error sends an error reply to the client, and logs it
func (resp Response) Error(message string, errorCode int) {
	resp.Ok = false
	resp.Message = message
	resp.ErrorCode = errorCode

	log.WithFields(log.Fields{
		"remote":    resp.req.RemoteAddr,
		"errorcode": errorCode,
	}).Errorf("%s", message)

	resp.Commit(nil)
}

func (resp Response) Success(message string, data interface{}) {
	resp.Ok = true
	resp.Message = message
	resp.ErrorCode = ErrorNoError
	resp.Commit(data)
}

// Commit sends the Response as the HTTP body
func (resp Response) Commit(data interface{}) {
	resp.wtr.Header().Set("Content-Type", "application/json")
	resp.Body.Data = data

	err := json.NewEncoder(resp.wtr).Encode(resp)
	if err != nil {
		resp.Error("There was a problem preparing the data", ErrorInternalError)
	}
}

// AddValidationError adds a validation error to a given Response
func (resp *Response) AddValidationError(path string, err string) {
	verr := ValidationError{Err: err, Path: path}
	resp.Body.ValidationErrors = append(resp.Body.ValidationErrors, verr)
}

// HasValidationErrors checks if a Response has any validation errors
func (resp *Response) HasValidationErrors() bool {
	if len(resp.Body.ValidationErrors) < 1 {
		return false
	}
	return true
}
