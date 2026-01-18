package accrual

import (
	"context"
	"fmt"
	"io"
	"net/http"
	stdHTTP "net/http"
	"net/url"
	"strconv"

	"github.com/oleshko-g/oggophermart/internal/transport"
	_ "goa.design/clue/log"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the accrual service endpoint HTTP clients.
type Client struct {
	// FetchOrderAccrual Doer is the HTTP client used to make requests to the
	// FetchOrderAccrual endpoint.
	reqDoer goahttp.Doer
	scheme  string
	host    string
}

// NewClient instantiates HTTP clients for all the accrual service servers.
func NewClient(
	scheme string,
	host string,
) *Client {
	return &Client{
		reqDoer: &stdHTTP.Client{},
		scheme:  scheme,
		host:    host,
	}
}

// FetchOrderAccrual returns an endpoint that makes HTTP requests to the accrual
// service FetchOrderAccrual server.
func (c *Client) FetchOrderAccrual(ctx context.Context, payload transport.FetchOrderAccrualPayload) (*transport.FetchOrderAccrualResult, error) {
	{
		req, err := c.buildFetchOrderAccrualRequest(ctx, payload.Number)
		if err != nil {
			return nil, err
		}
		resp, err := c.reqDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("accrual", "FetchOrderAccrual", err)
		}
		return decodeFetchOrderAccrualResponse(resp)
	}
}

// BuildFetchOrderAccrualRequest instantiates a HTTP request object with method
// and path set to call the "accrual" service "FetchOrderAccrual" endpoint
func (c *Client) buildFetchOrderAccrualRequest(ctx context.Context, orderNumber string) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: fetchOrderAccrualAccrualPath(orderNumber)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("accrual", "FetchOrderAccrual", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeFetchOrderAccrualResponse returns a decoder for responses returned by
// the accrual FetchOrderAccrual endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeFetchOrderAccrualResponse may return the following errors:
//   - "The request rate limit has been exceeded" (type *AccrualError): http.StatusTooManyRequests
//   - "Internal service error" (type *AccrualError): http.StatusInternalServerError
//   - error: internal error
func decodeFetchOrderAccrualResponse(resp *http.Response) (*transport.FetchOrderAccrualResult, error) {
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil, nil
	case http.StatusOK:
		var (
			body FetchOrderAccrualOKResponseBody
			err  error
		)
		err = goahttp.ResponseDecoder(resp).Decode(&body)
		if err != nil {
			return nil, goahttp.ErrDecodingError("accrual", "FetchOrderAccrual", err)
		}
		err = validateFetchOrderAccrualOKResponseBody(&body)
		if err != nil {
			return nil, goahttp.ErrValidationError("accrual", "FetchOrderAccrual", err)
		}
		res := newFetchOrderAccrualResultOK(&body)
		return res, nil
	case http.StatusTooManyRequests:
		var (
			body string
			err  error
		)
		err = goahttp.ResponseDecoder(resp).Decode(&body)
		if err != nil {
			return nil, goahttp.ErrDecodingError("accrual", "FetchOrderAccrual", err)
		}
		var (
			retryAfter int
		)
		{
			retryAfterRaw := resp.Header.Get("Retry-After")
			if retryAfterRaw == "" {
				return nil, goahttp.ErrValidationError("accrual", "FetchOrderAccrual", goa.MissingFieldError("retryAfter", "header"))
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
			return nil, goahttp.ErrValidationError("accrual", "FetchOrderAccrual", err)
		}
		return nil, newFetchOrderAccrualTheRequestRateLimitHasBeenExceeded(body, retryAfter)
	case http.StatusInternalServerError:
		return nil, newFetchOrderAccrualInternalServiceError()
	default:
		body, _ := io.ReadAll(resp.Body)
		return nil, goahttp.ErrInvalidResponse("accrual", "FetchOrderAccrual", resp.StatusCode, string(body))
	}
}

// FetchOrderAccrualOKResponseBody is the type of the "accrual" service
// "FetchOrderAccrual" endpoint HTTP response body.
type FetchOrderAccrualOKResponseBody struct {
	Order   *string  `form:"order,omitempty" json:"order,omitempty" xml:"order,omitempty"`
	Status  *string  `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	Accrual *float64 `form:"accrual,omitempty" json:"accrual,omitempty" xml:"accrual,omitempty"`
}

// NewFetchOrderAccrualResultOK builds a "accrual" service "FetchOrderAccrual"
// endpoint result from a HTTP "OK" response.
func newFetchOrderAccrualResultOK(body *FetchOrderAccrualOKResponseBody) *transport.FetchOrderAccrualResult {
	v := &transport.FetchOrderAccrualResult{
		Order:   *body.Order,
		Status:  *body.Status,
		Accrual: body.Accrual,
	}

	return v
}

// NewFetchOrderAccrualTheRequestRateLimitHasBeenExceeded builds a accrual
// service FetchOrderAccrual endpoint The request rate limit has been exceeded
// error.
func newFetchOrderAccrualTheRequestRateLimitHasBeenExceeded(body string, retryAfter int) *transport.AccrualError {
	return &transport.AccrualError{
		Message:    body,
		RetryAfter: retryAfter,
	}
}

// NewFetchOrderAccrualInternalServiceError builds a accrual service
// FetchOrderAccrual endpoint Internal service error error.
func newFetchOrderAccrualInternalServiceError() *transport.AccrualError {
	v := &transport.AccrualError{
		Message: "Internal service",
	}

	return v
}

// ValidateFetchOrderAccrualOKResponseBody runs the validations defined on
// FetchOrderAccrualOKResponseBody
func validateFetchOrderAccrualOKResponseBody(body *FetchOrderAccrualOKResponseBody) (err error) {
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

// fetchOrderAccrualAccrualPath returns the URL path to the accrual service FetchOrderAccrual HTTP endpoint.
func fetchOrderAccrualAccrualPath(number string) string {
	return fmt.Sprintf("/api/orders/%v", number)
}
