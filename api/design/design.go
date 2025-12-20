package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
		Response("unAuthorized", StatusUnauthorized)
	})

	Error("badRequest", ErrorResult, "invalid request format.")
	Error("unAuthorized", ErrorResult, "invalid request format.")
})

var _ = Service("oggophermart", func() {
	HTTP(func() {
		Response("badRequest", StatusBadRequest, func() {
			Description("the user is not authorized")
		})
	})
	Method("post order", func() {
		Result(PostOrderResult)
		HTTP(func() {
			POST("/api/user/orders")
			Response(StatusOK, func() {
				Tag("successCode", "OK")
				Description("The order number has already been uploaded by this user.")
				Body(Empty)
			})
			Response(StatusAccepted, func() {
				Tag("successCode", "Accepted")
				Description("The new order number has been processed.")
				Body(Empty)
			})
			Response(StatusInternalServerError, func() {
				Body(Empty)
			})
			// 409 — the order number has already been uploaded by another user;
			// 422 - invalid order number format;
			// 500 — internal server error.
		})
	})
})

var PostOrderResult = Type("PostOrderResult", func() {
	Attribute("successCode", String, func() {
		Enum("OK")
		Enum("Accepted")
	})
})

var APIResponse = Type("APIErrResponse", func() {
	Attribute("resErrCode", String, func() {
		Enum("badRequest")
		Enum("unAuthorized")
	})
})
