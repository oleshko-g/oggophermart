package design

import . "goa.design/goa/v3/dsl"

var _ = Service("oggophermart", func() {
	Method("post order", func() {
		HTTP(func() {
			POST("/api/user/orders")
			Response(StatusOK)
		})
		Result(Empty)
	})
})
