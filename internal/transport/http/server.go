// Package http implements http server for the gophermart API
package http //revive:disable-line:var-naming
import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	balance "github.com/oleshko-g/oggophermart/internal/gen/balance"
	genBalanceHTTPSrv "github.com/oleshko-g/oggophermart/internal/gen/http/balance/server"
	genUserHTTPSvr "github.com/oleshko-g/oggophermart/internal/gen/http/user/server"
	user "github.com/oleshko-g/oggophermart/internal/gen/user"
	"github.com/oleshko-g/oggophermart/internal/service"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
)

type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type server struct {
	goa struct {
		goahttp.Server
	}
	Server
}

func newHandlers(loggingCtx context.Context, balanceEndpoints *balance.Endpoints, userEndpoints *user.Endpoints) http.Handler {
	var (
		reqDecoder func(r *http.Request) goahttp.Decoder
		resEncoder func(ctx context.Context, res http.ResponseWriter) goahttp.Encoder
		mux        goahttp.Muxer
		errHandler func(context.Context, http.ResponseWriter, error)
	)
	{
		reqDecoder = goahttp.RequestDecoder
		resEncoder = goahttp.ResponseEncoder
		errHandler = errorHandler
		mux = goahttp.NewMuxer()
	}

	// create HTTP servers
	balanceServer := genBalanceHTTPSrv.New(balanceEndpoints, mux, reqDecoder, resEncoder, errHandler, nil)
	userServer := genUserHTTPSvr.New(userEndpoints, mux, reqDecoder, resEncoder, errHandler, nil)

	// mount HTTP endpoint onto mux
	balanceServer.Mount(mux)
	userServer.Mount(mux)

	loggingMiddleware := log.HTTP(loggingCtx)
	var handlers = loggingMiddleware(mux)

	return handlers

}

func NewServer(loggingCtx context.Context, cfg Config, svc service.Service) Server {
	var (
		balanceEndpoints *balance.Endpoints
		userEndpoints    *user.Endpoints
		handlers         http.Handler
	)
	{
		balanceEndpoints = balance.NewEndpoints(svc.Balance)
		userEndpoints = user.NewEndpoints(svc.User)
		handlers = newHandlers(loggingCtx, balanceEndpoints, userEndpoints)
	}

	return &server{
		Server: &http.Server{
			Addr:              cfg.Address().String(),
			Handler:           handlers,
			ReadHeaderTimeout: time.Second * 60},
	}
}

// errorHandler is the handler which is called when ther was an HTTP response encoding error
func errorHandler(ctx context.Context, res http.ResponseWriter, err error) {
	if res == nil {
		err = fmt.Errorf("%w: %s", errResponseWithError, errors.New("nil responseWriter"))
		slog.Error(err.Error())
		return
	}
	res.Header().Set("Content-Type", "text/plain")

	var statusCode int
	if err == nil {
		// replace the err and HTTP status code and still write the response
		err = fmt.Errorf("%w: %s", errResponseWithError, errors.New("nil err"))
		statusCode := http.StatusInternalServerError
		slog.Error(err.Error())
		http.Error(res, err.Error(), statusCode)
		return
	}

	http.Error(res, err.Error(), statusCode)
}

var errResponseWithError = errors.New("failed to response with error")

type Client = http.Client
