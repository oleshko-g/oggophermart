package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
	})
})

var _ = Service("user", func() {
	Error("invalid input parameter", OggophermartErrorType)
	Error("unauthorized", OggophermartErrorType)
	Error("internal service error", OggophermartErrorType)

	Method("register", func() {
		Payload(LoginPass)
		Result(userServiceResult)
		HTTP(func() {
			POST("/api/user/register")
			Response(StatusOK, func() {
				Body(Empty)
			})
			Response("invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("internal service error", StatusInternalServerError, func() {
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
			Response("invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})

})

var _ = Service("balance", func() {
	Error("invalid input parameter", OggophermartErrorType)
	Error("unauthorized", OggophermartErrorType)
	Error("internal service error", OggophermartErrorType)

	Method("post order", func() {
		Result(PostOrderResult)
		Error("already uploaded", OggophermartErrorType)
		Error("invalid order number", OggophermartErrorType)
		HTTP(func() {
			POST("/api/user/orders")
			Response(StatusOK, func() {
				Description("The order number has already been uploaded by this user.")
				Body(Empty)
			})
			Response(StatusAccepted, func() {
				Tag("uploadedBefore", "yes")
				Description("The new order number has been processed.")
				Body(Empty)
			})
			Response("invalid input parameter", StatusBadRequest, func() {
				Body(Empty)
			})
			Response("already uploaded", StatusConflict, func() {
				Description("The order number has already been uploaded by another user")
				Body(Empty)
			})
			Response("invalid order number", StatusUnprocessableEntity, func() {
				Description("invalid format")
				Body(Empty)
			})
			Response("unauthorized", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
})

var PostOrderResult = Type("PostOrderResult", func() {
	Attribute("uploadedBefore", func() {
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

var OggophermartErrorType = Type("OggophermartError", func(){
	ErrorName("name", func(){
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
