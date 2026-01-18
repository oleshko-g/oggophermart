package design

import . "goa.design/goa/v3/dsl"

var _ = API("gophermart", func() {
	Version("0.2")
	HTTP(func() {
		Path("/api")
		Consumes("text/plain", "application/json")
	})
})

// INFO: User
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

// INFO: Balance
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
	Method("ListUserOrders", func() {
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
			Attribute("current", Float64, func() {
				Minimum(0)
			})
			Attribute("withdrawn", Float64, func() {
				Minimum(0)
			})
			Example(func() {
				Value(Val{
					"current":   500.5,
					"withdrawn": 42,
				})
			})
			Required("current", "withdrawn")
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
	Method("WithdrawUserBalance", func() {
		Payload(func() {
			Token("Authorization", String, "A JWT token used to authenticate a request", func() {
				Example(func() {
					Value("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
				})
			})
			Attribute("order", String, func() {
				Pattern("[1-9][0-9]*")
			})
			Attribute("sum", Float64, func() {
				ExclusiveMinimum(0)
			})
			Required("Authorization", "order", "sum")
			Example(func() {
				Value(
					Val{
						"order": "2377225624",
						"sum":   751,
					})
			})
		})
		Error("Insufficient funds", ErrorType)
		Error("Invalid order number", ErrorType)
		HTTP(func() {
			POST("/user/balance/withdraw")
			Response(StatusOK, func() {
				Body(Empty)
			})
			Response("User is not authenticated", StatusUnauthorized, func() {
				Description("User is not authenticated")
			})
			Response("missing_field", StatusUnauthorized, func() {
				Body(Empty)
				Description("Missing or empty Authorization header")
			})
			Response("Insufficient funds", StatusPaymentRequired, func() {
				Body(Empty)
				Description("Insufficient funds")
			})
			Response("Invalid order number", StatusUnprocessableEntity, func() {
				Body(Empty)
				Description("Invalid order number")
			})
			Response("Internal service error", StatusInternalServerError, func() {
				Body(Empty)
			})
		})
	})
})

// INFO: Accrual
var _ = Service("accrual", func() {
	Method("GetOrderAccrual", func() {

		Payload(func() {
			Attribute("number", OrderNumber)
			Required("number")
		})

		Result(func() {
			Attribute("order", OrderNumber)
			Attribute("status", String, func() {
				Enum("REGISTERED", "INVALID", "PROCESSING", "PROCESSED")
			})
			Attribute("accrual", Float64, func() {
				ExclusiveMinimum(0)
			})

			Required("order", "status")
			Example(func() {
				Value(Val{
					"order":   "42",
					"status":  "PROCESSED",
					"accrual": 500.1,
				})
				Value(Val{
					"order":  "32",
					"status": "INVALID",
				})
				Value(Val{
					"order":  "2",
					"status": "PROCESSING",
				})
				Value(Val{
					"order":  "4",
					"status": "REGISTERED",
				})
			})
		})

		Error("Internal service error", AccrualErrorType)
		Error("The request rate limit has been exceeded", AccrualErrorType, func() {
			Required("name", "retryAfter", "message")
		})

		HTTP(func() {
			GET("/orders/{number}")
			Param("number", String)
			Response(StatusOK)
			Response(StatusNoContent, func() {
				Body(Empty)
			})
			Response("The request rate limit has been exceeded", StatusTooManyRequests, func() {
				Description("The request rate limit has been exceeded")
				Header("retryAfter:Retry-After")
				ContentType("text/plain")
			})
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
})

var OrderNumber = Type("OrderNumber", String, func() {
	Description("Unique user order number")
	Pattern("[1-9][0-9]*")
})

var AccrualErrorType = Type("AccrualError", func() {
	ErrorName("name", func() {
		Description("identifier to map an error to HTTP status codes")
		Meta("struct:tag:json", "-")
		Meta("openapi:generate", "false")
		Meta("openapi:example", "false")
	})

	Attribute("retryAfter", Int, func() {
		ExclusiveMinimum(0)
	})

	Attribute("message", String)

	Meta("openapi:generate", "false")
	Meta("openapi:example", "false")
	Meta("struct:pkg:path", "service")

	Required("name")
})
