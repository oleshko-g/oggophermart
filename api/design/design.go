package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
	})
})

var _ = Service("user", func() {
	Error("Invalid input parameter", ErrorType)
	Error("User is not authenticated", ErrorType)
	Error("Internal service error", ErrorType)
	Error("Not implemented", ErrorType)

	Method("register", func() {
		Payload(LoginPass)
		Result(userServiceResult)
		HTTP(func() {
			POST("/api/user/register")
			Response(StatusOK, func() {
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
			Response("Not implemented", StatusNotImplemented, func() {
				Body(Empty)
			})
		})

	})
	Method("login", func() {
		Payload(LoginPass)
		Result(userServiceResult)
		HTTP(func() {
			POST("/api/user/login")
			Response(StatusOK, func() {
				Body(Empty)
			})
			Response("Not implemented", StatusNotImplemented, func() {
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})

})

var _ = Service("balance", func() {
	Error("Invalid input parameter", ErrorType)
	Error("User is not authenticated", ErrorType)
	Error("Internal service error", ErrorType)
	Error("Not implemented", ErrorType)

	Method("post order", func() {
		Result(PostOrderResult)
		Error("The order belongs to another user", ErrorType)
		Error("Invalid order number", ErrorType)
		HTTP(func() {
			POST("/api/user/orders")
			Response(StatusOK, func() {
				Description("The order has been accepted for processing before.")
				Body(Empty)
			})
			Response(StatusAccepted, func() {
				Tag("accepted", "yes")
				Description("The order has been accepted for processing.")
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("The order belongs to another user", StatusConflict, func() {
				Body(Empty)
			})
			Response("Invalid order number", StatusUnprocessableEntity, func() {
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
			Response("Not implemented", StatusNotImplemented, func() {
				Body(Empty)
			})
		})
	})
})

var PostOrderResult = Type("PostOrderResult", func() {
	Attribute("accepted", func() {
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})
	Meta("openapi:example", "false")
})

var LoginPass = Type("LoginPass", func() {
	Attribute("login", String)
	Attribute("password", String)
	Required("login", "password")
})

var userServiceResult = Type("userServiceResult", func() {
	Attribute("statusCode", func() {
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})
})

var ErrorType = Type("GophermartError", func() {
	ErrorName("name", func() {
		Description("identifier to map an error to HTTP status codes")
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})
	Required("name")
	Meta("openapi:generate", "false")
	Meta("openapi:example", "false")
	Meta("struct:pkg:path", "service")
})
