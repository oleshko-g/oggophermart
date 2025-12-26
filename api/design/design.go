package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
	})
})

var _ = Service("user", func() {
	Error("invalidInputParameter", OggophermartErrorType)
	Method("register", func() {
		Payload(LoginPass)
		Result(userServiceResult)
		HTTP(func() {
			POST("/api/user/register")
			Response(StatusOK, func() {
				Body(Empty)
			})
			Response("invalidInputParameter", func() {
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
			Response("invalidInputParameter", func() {
				Body(Empty)
			})
		})
	})

})

var _ = Service("balance", func() {
	Error("invalid input parameter", OggophermartErrorType)
	Method("post order", func() {
		Result(PostOrderResult)
		HTTP(func() {
			POST("/api/user/orders")
			Response(StatusOK, func() {
				Tag("statusCode", "200")
				Description("The order number has already been uploaded by this user.")
				Body(Empty)
			})
			Response(StatusAccepted, func() {
				Tag("statusCode", "202")
				Description("The new order number has been processed.")
				Body(Empty)
			})
			Response("invalid input parameter", StatusBadGateway, func() {
				Description("aaaa")
				Body(Empty)
			})
			Response(StatusConflict, func() {
				Tag("statusCode", "409")
				Description("The order number has already been uploaded by another user")
				Body(Empty)
			})
			Response(StatusUnprocessableEntity, func() {
				Tag("statusCode", "422")
				Description("Invalid order number format.")
				Body(Empty)
			})
			Response(StatusUnauthorized, func() {
				Tag("statusCode", "401")
				Description("User is not authenticated")
				Body(Empty)
			})
			Response(StatusInternalServerError, func() {
				// not tag for default case
				Description("Internal server error")
				Body(Empty)
			})
		})
	})
})

var PostOrderResult = Type("PostOrderResult", func() {
	Attribute("statusCode", func() {
		Meta("struct:tag:json", "-") // hide from response
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false") // hide from swagger
	})
	Meta("openapi:example", "false") // hide from swagger
})

var LoginPass = Type("LoginPass", func() {
	Attribute("login", String)
	Attribute("password", String)
	Required("login", "password")
})

var userServiceResult = Type("userServiceResult", func() {
	Attribute("statusCode", func() {
		Meta("struct:tag:json", "-") // hide from response
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false") // hide from swagger
	})
})

var OggophermartErrorType = Type("OggophermartError", func(){
	ErrorName("name", func(){
		Description("identifier to map an error to HTTP status codes")
		Meta("struct:tag:json", "-") // hide from response
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false") // hide from swagger
	})
	Meta("openapi:generate", "false")
	Meta("openapi:example", "false") // hide from swagger
	Required("name")
})

// var TestType = Type("TestType", func(){
// 	Attribute("att", String)
// })


	   // var CustomErrorType = Type("CustomError", func() {
	   //     // The "name" attribute is used to select the error response.
	   //     // name should be set to either "internal_error" or "bad_request" by
	   //     // the service method returning the error.
	   //     ErrorName("name", String, "Name of error.")
	   //     Attribute("message", String, "Message of error.")
	   //     Attribute("occurred_at", String, "Time error occurred.", func() {
	   //         Format(FormatDateTime)
	   //     })
	   //     Required("name", "message", "occurred_at")
	   // })
