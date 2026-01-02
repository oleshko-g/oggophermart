// Package http implements http server for the gophermart API
package http //revive:disable-line:var-naming

import (
	"errors"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	"github.com/oleshko-g/oggophermart/internal/service"
)


type ogGophermartHTTPErrorStatus string

func (or *oggophermartHTTPstatus) Status(err error) int {
	switch {
		case errors.Is(err, service.ErrInvalidFormat):
		return http.StatusBadRequest
	}

}
var ()

//		Tag("statusCode", "400")
//		Body(Empty)
//	})
//
//	Response(StatusConflict, func() {
//		Tag("statusCode", "409")
//		Description("The order number has already been uploaded by another user")
//		Body(Empty)
//	})
//
//	Response(StatusUnprocessableEntity, func() {
//		Tag("statusCode", "422")
//		Description("Invalid order number format.")
//		Body(Empty)
//	})
//
//	Response(StatusUnauthorized, func() {
//		Tag("statusCode", "401")
//		Description("User is not authenticated")
//		Body(Empty)
//	})
// )
