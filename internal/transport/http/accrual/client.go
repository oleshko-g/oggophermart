package accrual

import (
	"context"
	"fmt"
	"io"
	"net/http"
	stdHTTP "net/http"
	"net/url"
	"strconv"

	_ "github.com/oleshko-g/oggophermart/internal/gen/balance"
	_ "github.com/oleshko-g/oggophermart/internal/gen/http/balance/server"
	_ "github.com/oleshko-g/oggophermart/internal/gen/http/user/server"
	_ "github.com/oleshko-g/oggophermart/internal/gen/user"
	_ "github.com/oleshko-g/oggophermart/internal/service"
	_ "goa.design/clue/log"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"

	accrual "github.com/oleshko-g/oggophermart/internal/gen/accrual"
)

// Client lists the accrual service endpoint HTTP clients.
type Client struct {
	// GetOrderAccrual Doer is the HTTP client used to make requests to the
	// GetOrderAccrual endpoint.
	getOrderAccrualDoer goahttp.Doer
	scheme              string
	host                string
	decoder             func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the accrual service servers.
func NewClient(
	scheme string,
	host string,
) *Client {
	return &Client{
		getOrderAccrualDoer: &stdHTTP.Client{},
		scheme:              scheme,
		host:                host,
		decoder:             goahttp.ResponseDecoder,
	}
}

// GetOrderAccrual returns an endpoint that makes HTTP requests to the accrual
// service GetOrderAccrual server.
func (c *Client) GetOrderAccrual() goa.Endpoint {
	var (
		decodeResponse = DecodeGetOrderAccrualResponse(c.decoder)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildGetOrderAccrualRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.getOrderAccrualDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("accrual", "GetOrderAccrual", err)
		}
		return decodeResponse(resp)
	}
}

// BuildGetOrderAccrualRequest instantiates a HTTP request object with method
// and path set to call the "accrual" service "GetOrderAccrual" endpoint
func (c *Client) BuildGetOrderAccrualRequest(ctx context.Context, v any) (*http.Request, error) {
	var (
		number string
	)
	{
		p, ok := v.(*accrual.GetOrderAccrualPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("accrual", "GetOrderAccrual", "*accrual.GetOrderAccrualPayload", v)
		}
		number = string(p.Number)
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: getOrderAccrualAccrualPath(number)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("accrual", "GetOrderAccrual", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetOrderAccrualResponse returns a decoder for responses returned by
// the accrual GetOrderAccrual endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetOrderAccrualResponse may return the following errors:
//   - "The request rate limit has been exceeded" (type *AccrualError): http.StatusTooManyRequests
//   - "Internal service error" (type *AccrualError): http.StatusInternalServerError
//   - error: internal error
func DecodeGetOrderAccrualResponse(decoder func(*http.Response) goahttp.Decoder) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusNoContent:
			return nil, nil
		case http.StatusOK:
			var (
				body GetOrderAccrualOKResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("accrual", "GetOrderAccrual", err)
			}
			err = ValidateGetOrderAccrualOKResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("accrual", "GetOrderAccrual", err)
			}
			res := NewGetOrderAccrualResultOK(&body)
			return res, nil
		case http.StatusTooManyRequests:
			var (
				body GetOrderAccrualTheRequestRateLimitHasBeenExceededResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("accrual", "GetOrderAccrual", err)
			}
			var (
				retryAfter int
			)
			{
				retryAfterRaw := resp.Header.Get("Retry-After")
				if retryAfterRaw == "" {
					return nil, goahttp.ErrValidationError("accrual", "GetOrderAccrual", goa.MissingFieldError("retryAfter", "header"))
				}
				v, err2 := strconv.ParseInt(retryAfterRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("retryAfter", retryAfterRaw, "integer"))
				}
				retryAfter = int(v)
			}
			if retryAfter <= 0 {
				err = goa.MergeErrors(err, goa.InvalidRangeError("retryAfter", retryAfter, 0, true))
			}
			if err != nil {
				return nil, goahttp.ErrValidationError("accrual", "GetOrderAccrual", err)
			}
			return nil, NewGetOrderAccrualTheRequestRateLimitHasBeenExceeded(&body, retryAfter)
		case http.StatusInternalServerError:
			return nil, NewGetOrderAccrualInternalServiceError()
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("accrual", "GetOrderAccrual", resp.StatusCode, string(body))
		}
	}
}

// GetOrderAccrualOKResponseBody is the type of the "accrual" service
// "GetOrderAccrual" endpoint HTTP response body.
type GetOrderAccrualOKResponseBody struct {
	Order   *string  `form:"order,omitempty" json:"order,omitempty" xml:"order,omitempty"`
	Status  *string  `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	Accrual *float64 `form:"accrual,omitempty" json:"accrual,omitempty" xml:"accrual,omitempty"`
}

// GetOrderAccrualTheRequestRateLimitHasBeenExceededResponseBody is the type of
// the "accrual" service "GetOrderAccrual" endpoint HTTP response body for the
// "The request rate limit has been exceeded" error.
type GetOrderAccrualTheRequestRateLimitHasBeenExceededResponseBody struct {
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
}

// NewGetOrderAccrualResultOK builds a "accrual" service "GetOrderAccrual"
// endpoint result from a HTTP "OK" response.
func NewGetOrderAccrualResultOK(body *GetOrderAccrualOKResponseBody) *accrual.GetOrderAccrualResult {
	v := &accrual.GetOrderAccrualResult{
		Order:   accrual.OrderNumber(*body.Order),
		Status:  *body.Status,
		Accrual: body.Accrual,
	}

	return v
}

// NewGetOrderAccrualTheRequestRateLimitHasBeenExceeded builds a accrual
// service GetOrderAccrual endpoint The request rate limit has been exceeded
// error.
func NewGetOrderAccrualTheRequestRateLimitHasBeenExceeded(body *GetOrderAccrualTheRequestRateLimitHasBeenExceededResponseBody, retryAfter int) *AccrualError {
	v := &AccrualError{
		Message: *body.Message,
	}
	v.RetryAfter = retryAfter

	return v
}

// NewGetOrderAccrualInternalServiceError builds a accrual service
// GetOrderAccrual endpoint Internal service error error.
func NewGetOrderAccrualInternalServiceError() *AccrualError {
	v := &AccrualError{}

	return v
}

// ValidateGetOrderAccrualOKResponseBody runs the validations defined on
// GetOrderAccrualOKResponseBody
func ValidateGetOrderAccrualOKResponseBody(body *GetOrderAccrualOKResponseBody) (err error) {
	if body.Order == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("order", "body"))
	}
	if body.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "body"))
	}
	if body.Order != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.order", *body.Order, "[1-9][0-9]*"))
	}
	if body.Status != nil {
		if !(*body.Status == "REGISTERED" || *body.Status == "INVALID" || *body.Status == "PROCESSING" || *body.Status == "PROCESSED") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.status", *body.Status, []any{"REGISTERED", "INVALID", "PROCESSING", "PROCESSED"}))
		}
	}
	if body.Accrual != nil {
		if *body.Accrual <= 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.accrual", *body.Accrual, 0, true))
		}
	}
	return
}

// getOrderAccrualAccrualPath returns the URL path to the accrual service GetOrderAccrual HTTP endpoint.
func getOrderAccrualAccrualPath(number string) string {
	return fmt.Sprintf("/api/orders/%v", number)
}

// BuildGetOrderAccrualPayload builds the payload for the accrual
// GetOrderAccrual endpoint from CLI flags.
func BuildGetOrderAccrualPayload(accrualGetOrderAccrualNumber string) (*accrual.GetOrderAccrualPayload, error) {
	var err error
	var number string
	{
		number = accrualGetOrderAccrualNumber
		err = goa.MergeErrors(err, goa.ValidatePattern("number", number, "[1-9][0-9]*"))
		if err != nil {
			return nil, err
		}
	}
	v := &accrual.GetOrderAccrualPayload{}
	v.Number = accrual.OrderNumber(number)

	return v, nil
}
