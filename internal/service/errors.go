// Package service is the package shared by all oggophermart services
package service

import "errors"

var (
	ErrInvalidFormat = errors.New("invalid request format")
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
)
