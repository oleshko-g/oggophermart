package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	Version("0.2")
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
		Payload(LoginPassword)
		Result(JWTToken)
		Error("Login is taken already", ErrorType)
		HTTP(func() {
			POST("/user/register")
			Response(StatusOK, func() {
				Header("authToken:Authorization")
				Body(Empty)
			})
			Response("Invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
			Response("Login is taken already", StatusConflict, func() {
				Body(Empty)
			})
		})

	})
	Method("login", func() {
		Payload(LoginPassword)
		Result(JWTToken)
		HTTP(func() {
			POST("/user/login")
			Response(StatusOK, func() {
				Header("authToken:Authorization")
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
	Security(JWTAuth)
	Error("Invalid input parameter", ErrorType)
	Error("User is not authenticated", ErrorType)
	Error("Internal service error", ErrorType)
	Error("Not implemented", ErrorType)

	Method("UploadUserOrder", func() {
		Description("Upload user order")
		Result(func() {
			Attribute("accepted", func() {
				Meta("struct:tag:json", "-")
				Meta("openapi:generate", "false")
				Meta("openapi:example", "false")
			})
			Meta("openapi:example", "false")
		})
		Payload(func() {
			Token("JWTToken", String, "A JWT token used to authenticate a request")
			Attribute("OrderNumber", String, func() {
				Description("Unique user order number")
				Pattern("[1-9][0-9]*")
			})
			Required("JWTToken", "OrderNumber")
		})
		Error("The order belongs to another user", ErrorType)
		Error("Invalid order number", ErrorType)
		HTTP(func() {
			POST("/user/orders")
			Body("OrderNumber")
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

var _ = Service("accrual", func() {
	Method("GetOrder", func() {
		Error("Internal service error", ErrorType)
		Result(GetOrderResult)
		Payload(func() {
			Attribute("number", String, func() {
			})
			Required("number")
		})
		HTTP(func() {
			GET("GET /orders/{number}")
			Param("number", String)
			Response(StatusOK)
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})

	})
})

var UploadUserOrderResult = Type("PostOrderResult", func() {
	Attribute("accepted", func() {
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})
	Meta("openapi:example", "false")
})
var GetOrderResult = Type("GetOrderResult", func() {
	Attribute("order", String)
	Attribute("status", String)
	Attribute("accrual", UInt)
})

var LoginPassword = Type("LoginPassword", func() {
	Attribute("login", String)
	Attribute("password", String)
	Required("login", "password")
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

var JWTAuth = JWTSecurity("jwt", func() {
	Description("Secures an endpoint by requiring a valid JWT token.")
})

var JWTToken = Type("JWTToken", func() {
	Token("authToken", String, func() {
		Description("A JWT token used to authenticate a request")
	})
	Required("authToken")
	Meta("struct:pkg:path", "service")
})
