package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
	})
})

var _ = Service("user", func() {
	Method("register", func() {
		Payload(LoginPass)
		Result(userServiceResult)
		HTTP(func() {
			POST("/api/user/register")
			Response(StatusOK, func() {
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
		})
	})

})

var _ = Service("balance", func() {
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
			Response(StatusBadRequest, func() {
				Tag("statusCode", "400")
				Description("invalid request format")
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
