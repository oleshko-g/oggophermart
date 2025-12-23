// Package http implements http server for the gophermart API
package http //revive:disable-line:var-naming

import (
	"context"
	"net/http"
	"net/url"

	balance "github.com/oleshko-g/oggophermart/internal/gen/balance"
	balanceHTTP "github.com/oleshko-g/oggophermart/internal/gen/http/balance/server"
	userHTTP "github.com/oleshko-g/oggophermart/internal/gen/http/user/server"
	"github.com/oleshko-g/oggophermart/internal/gen/user"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
)

type Server struct {
	goahttp.Server
	balanceSvc balance.Service
	balanceSrv balanceHTTP.Server
}

func NewMux(ctx context.Context,
	u *url.URL, balanceEndpoints *balance.Endpoints, userEndpoints *user.Endpoints, errc chan error) goahttp.Muxer {
	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and mount debug and profiler
	// endpoints in debug mode.
	mux := goahttp.NewMuxer()

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.

	// TODO: read code error handler does
	// eh := errorHandler(ctx)
	balanceServer := balanceHTTP.New(balanceEndpoints, mux, dec, enc /*eh */, nil, nil)
	userServer := userHTTP.New(userEndpoints, mux, dec, enc /*eh */, nil, nil)

	// Configure the mux.
	balanceHTTP.Mount(mux, balanceServer)
	userHTTP.Mount(mux, userServer)

	var handler http.Handler = mux
	handler = log.HTTP(ctx)(handler)
	return mux

}
