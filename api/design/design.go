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
			Response("Login is taken already", StatusConflict, func() {
				Description("Login is taken already")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
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
				Description("User is not authenticated")
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
	Error("missing_field")
	HTTP(func() {
		Header("Authorization", func() {
			Example(func() {
				Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
			})
		})

	})

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
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Attribute("OrderNumber", String, func() {
				Description("Unique user order number")
				Pattern("[1-9][0-9]*")
			})
			Required("Authorization", "OrderNumber")
		})
		Error("The order belongs to another user", ErrorType)
		Error("Invalid order number", ErrorType)
		HTTP(func() {
			POST("/user/orders")
			Body("OrderNumber", func() {
				Example(func() {
					Value("12345678903")
				})
			})
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
				Description("The order belongs to another user")
				Body(Empty)
			})
			Response("Invalid order number", StatusUnprocessableEntity, func() {
				Description("Invalid order number")
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
	Method("ListUserOrder", func() {
		Description("List user orders")
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Required("Authorization")
		})
		Result(func() {
			Attribute("orders", ArrayOf(Order), func() {
				Example(func() {
					Value([]Val{
						{
							"number":      "9278923470",
							"status":      "PROCESSED",
							"accrual":     500,
							"uploaded_at": "2020-12-10T15:15:45+03:00",
						},
						{
							"number":      "12345678903",
							"status":      "PROCESSING",
							"uploaded_at": "2020-12-10T15:12:01+03:00",
						},
						{
							"number":      "346436439",
							"status":      "INVALID",
							"uploaded_at": "2020-12-09T16:09:53+03:00",
						}})
				})
			})
			Attribute("no orders", func() {
				Meta("struct:tag:json", "-")
				Meta("openapi:generate", "false")
				Meta("openapi:example", "false")
			})
		})
		HTTP(func() {
			GET("/user/orders")
			Response(StatusOK, func() {
				Body("orders")
			})
			Response(StatusNoContent, func() {
				Tag("no orders", "yes")
				Description("No orders available")
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
	Method("GetUserBalance", func() {
		Description("Get user balance")
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Required("Authorization")
		})
		Result(func() {
			Attribute("current", UInt)
			Attribute("withdrawn", UInt)
			Example(func() {
				Value(Val{
					"current":   500.5,
					"withdrawn": 42,
				})
			})
		})
		HTTP(func() {
			GET("/user/balance")
			Response(StatusOK)
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
				Body(Empty)
			})
			Response("missing_field", StatusUnauthorized, func() {
				Description("Missing or empty Authorization header")
				Body(Empty)
			})
			Response("Internal service error", StatusInternalServerError, func() {
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
	Example(func() {
		Value(Val{
			"login":    "<login>",
			"password": "<password>",
		})
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

var JWTAuth = JWTSecurity("jwt", func() {
	Description("Secures an endpoint by requiring a valid JWT token.")
})

var JWTToken = Type("JWTToken", func() {
	Token("authToken", String, func() {
		Description("A JWT token used to authenticate a request")
		Example(func() {
			Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
		})
	})
	Required("authToken")
	Meta("struct:pkg:path", "service")
})

var Order = Type("Order", func() {
	Attribute("number", String, func() {
		Pattern("[1-9][0-9]*")
	})
	Attribute("status", String, func() {
		Enum("NEW", "PROCESSING", "INVALID", "PROCESSED")
	})
	Attribute("accrual", UInt)
	Attribute("uploaded_at", String, func() {
		Format(FormatDateTime)
	})
	Required("number", "status", "uploaded_at")
	// Example("PROCESSED", func() {
	// 	Value(Val{
	// 		"number":      "9278923470",
	// 		"status":      "PROCESSED",
	// 		"accrual":     500,
	// 		"uploaded_at": "2020-12-10T15:15:45+03:00",
	// 	})
	// })
	// Example("PROCESSING", func() {
	// 	Value(Val{
	// 		"number":      "12345678903",
	// 		"status":      "PROCESSING",
	// 		"uploaded_at": "2020-12-10T15:12:01+03:00",
	// 	})
	// })
	// Example("INVALID", func() {
	// 	Value(Val{
	// 		"number":      "346436439",
	// 		"status":      "INVALID",
	// 		"uploaded_at": "2020-12-09T16:09:53+03:00",
	// 	})
	// })
})
