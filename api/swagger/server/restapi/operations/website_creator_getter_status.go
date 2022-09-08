// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// WebsiteCreatorGetterStatusHandlerFunc turns a function with the right signature into a website creator getter status handler
type WebsiteCreatorGetterStatusHandlerFunc func(WebsiteCreatorGetterStatusParams) middleware.Responder

// Handle executing the request and returning a response
func (fn WebsiteCreatorGetterStatusHandlerFunc) Handle(params WebsiteCreatorGetterStatusParams) middleware.Responder {
	return fn(params)
}

// WebsiteCreatorGetterStatusHandler interface for that can handle valid website creator getter status params
type WebsiteCreatorGetterStatusHandler interface {
	Handle(WebsiteCreatorGetterStatusParams) middleware.Responder
}

// NewWebsiteCreatorGetterStatus creates a new http.Handler for the website creator getter status operation
func NewWebsiteCreatorGetterStatus(ctx *middleware.Context, handler WebsiteCreatorGetterStatusHandler) *WebsiteCreatorGetterStatus {
	return &WebsiteCreatorGetterStatus{Context: ctx, Handler: handler}
}

/* WebsiteCreatorGetterStatus swagger:route GET /websiteCreator/state/{contractAddress} websiteCreatorGetterStatus

WebsiteCreatorGetterStatus website creator getter status API

*/
type WebsiteCreatorGetterStatus struct {
	Context *middleware.Context
	Handler WebsiteCreatorGetterStatusHandler
}

func (o *WebsiteCreatorGetterStatus) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewWebsiteCreatorGetterStatusParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
