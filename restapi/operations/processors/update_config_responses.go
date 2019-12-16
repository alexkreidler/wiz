// Code generated by go-swagger; DO NOT EDIT.

package processors

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// UpdateConfigOKCode is the HTTP code returned for type UpdateConfigOK
const UpdateConfigOKCode int = 200

/*UpdateConfigOK OK

swagger:response updateConfigOK
*/
type UpdateConfigOK struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewUpdateConfigOK creates UpdateConfigOK with default headers values
func NewUpdateConfigOK() *UpdateConfigOK {

	return &UpdateConfigOK{}
}

// WithPayload adds the payload to the update config o k response
func (o *UpdateConfigOK) WithPayload(payload interface{}) *UpdateConfigOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update config o k response
func (o *UpdateConfigOK) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// UpdateConfigBadRequestCode is the HTTP code returned for type UpdateConfigBadRequest
const UpdateConfigBadRequestCode int = 400

/*UpdateConfigBadRequest Bad Request

swagger:response updateConfigBadRequest
*/
type UpdateConfigBadRequest struct {

	/*
	  In: Body
	*/
	Payload *UpdateConfigBadRequestBody `json:"body,omitempty"`
}

// NewUpdateConfigBadRequest creates UpdateConfigBadRequest with default headers values
func NewUpdateConfigBadRequest() *UpdateConfigBadRequest {

	return &UpdateConfigBadRequest{}
}

// WithPayload adds the payload to the update config bad request response
func (o *UpdateConfigBadRequest) WithPayload(payload *UpdateConfigBadRequestBody) *UpdateConfigBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update config bad request response
func (o *UpdateConfigBadRequest) SetPayload(payload *UpdateConfigBadRequestBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateConfigBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}