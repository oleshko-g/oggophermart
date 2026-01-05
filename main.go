package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/joho/godotenv"
	genAccrualHTTPClient "github.com/oleshko-g/oggophermart/internal/gen/http/accrual/client"
	"github.com/oleshko-g/oggophermart/internal/service"
	balance "github.com/oleshko-g/oggophermart/internal/service/balance"
	user "github.com/oleshko-g/oggophermart/internal/service/user"
	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql"
	"github.com/oleshko-g/oggophermart/internal/transport/http"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
	"golang.org/x/sync/errgroup"
)

type gophermart struct {
	transport struct {
		http struct {
			http.Server
			http.Config
			client struct {
				accrual *genAccrualHTTPClient.Client
			}
		}
	}
	service.Service
	storage.Storage
	dbCfg      db.Config
	loggingCtx context.Context
	configured bool
	readyToRun bool
}

var g = gophermart{
	loggingCtx: newLoggingCtx(),
}

// newLoggingCtx returns the context with goa logger
func newLoggingCtx() context.Context {
	ctx := context.Background()

	opts := []log.LogOption{
		log.WithFormat(log.FormatTerminal),
		log.WithDebug(),
		log.WithFileLocation(),
	}

	// puts logger in ctx
	logCtx := log.Context(ctx, opts...)
	return logCtx
}

// configure sets the gophermart config parameters in the following priority:
//  1. command line flags
//  2. env vars
//  3. default values
//
// If successful cofigure set the [slog.Logger] field and [configured] flag
func (g *gophermart) cofigure() (err error) {

	err = loadEnvVarsFromFile()
	if err != nil {
		return err
	}

	// HTTP host address
	aF := g.transport.http.Address()
	err = aF.Set("localhost:8080") // default
	if err != nil {
		return err
	}

	flag.Var(aF, "a", "The host address of the gophermart")

	v, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		err = aF.Set(v) // override the default
		if err != nil {
			return err
		}
	}

	// DB
	dF := g.dbCfg.DSN()
	flag.Var(dF, "d", "Database connection address")

	v, ok = os.LookupEnv("DATABASE_URI")
	if ok {
		err = dF.Set(v) // override the default
		if err != nil {
			return err
		}
	}

	// The accrual system host address
	rF := g.transport.http.AccrualAddress()
	err = rF.Set("localhost:8081") // default
	if err != nil {
		return err
	}

	flag.Var(rF, "r", "Address of the accrual system")

	if v, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		err = rF.Set(v) // override the default
		if err != nil {
			return err
		}
	}

	flag.Parse() // if any of the flags are set they override the defaults or env vars
	log.Printf(g.loggingCtx, "gophermart host address is set to %s from the %s", aF.String(), aF.Source)
	log.Printf(g.loggingCtx, "gophermart database connection is set to %s from the %s", dF.DriverName.String(), dF.Source)
	log.Printf(g.loggingCtx, "gophermart address of the accrual system is set to %s from the %s", rF.String(), rF.Source)
	g.configured = true
	return nil
}

// loadEnvVarsFromFile opens or if not exists creates the env file and loads env vars from it.
func loadEnvVarsFromFile() (err error) {
	const (
		envFileName                    = ".env"
		envFilePermissions os.FileMode = 0o644
	)

	_, err = os.OpenFile(envFileName, os.O_RDONLY|os.O_CREATE, envFilePermissions)
	if err != nil {
		return err
	}

	err = godotenv.Load(envFileName)
	if err != nil {
		return err
	}
	return nil
}

// setup readies the gopheramart to run.
// It does the following:
//  1. Sets the storage for each service
//  2. Intanciates services with the set storage
//  3. Instanciates the HTTP server
//  4. Instanicates the Accrual system HTTP client
//
// If successful it sets readyToRun flag
func (g *gophermart) setup() (err error) {
	if !g.configured {
		return errSetupGophermartNotConfigured
	}
	dbStorage, err := sql.New(&g.dbCfg)
	if err != nil {
		return err
	}

	// 1. Sets the storage for each service
	// wrap concrete type [*sql.Storage] struct with interfaces
	g.Storage.User = dbStorage
	g.Storage.Balance = dbStorage

	// 2. Intanciates services with the set storage
	g.Service = service.Service{
		Balance: balance.New(g.Storage.Balance),
		User:    user.New(g.Storage.User),
	}

	// 3. Instanciates the HTTP server
	g.transport.http.Server = http.NewServer(g.transport.http.Config, g.Service)

	// 4. Instanicates the Accrual system HTTP client
	// err = genAccrual.NewGetOrderEndpoint(a)
	g.transport.http.client.accrual = genAccrualHTTPClient.NewClient(
		g.transport.http.AccrualAddress().Host,
		g.transport.http.AccrualAddress().Port,
		&http.Client{},
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		true,
	)

	g.readyToRun = true
	return nil
}

// run launches gophermart
func (g *gophermart) run() (err error) {
	errGroup, runCtx := errgroup.WithContext(g.loggingCtx)

	errGroup.Go(func() error {
		log.Infof(g.loggingCtx, "in HTTP server")
		if !g.readyToRun {
			return errSetupGophermartNotReadyToRun
		}
		return g.transport.http.Server.ListenAndServe()
	})

	errGroup.Go(func() error {
		log.Infof(g.loggingCtx, "in accrual worker")
		if !g.readyToRun {
			return errSetupGophermartNotReadyToRun
		}
		// TODO: wrap into worker func
		getOrder := g.transport.http.client.accrual.GetOrder()
		for {
			// TODO: storage.GetOrderToProcess
			ordersToProcess := []string{"1"}

			errGroup.Go(func() error {
				for _, order := range ordersToProcess {
					req, err := genAccrualHTTPClient.BuildGetOrderPayload(order)
					if err != nil {
						return err
					}

					// TODO: store res
					// TODO: handle err
					res, err := getOrder(runCtx, req)
					_, _ = res, err
				}
				return nil
			})

		}
	})

	// TODO: add os.Signal handling
	// TODO: add graceful shutdown

	log.Printf(g.loggingCtx, "gophermart HTTP server is listening on %s", g.transport.http.Address().String())
	return errGroup.Wait()
}

var errSetupGophermartNotConfigured = errors.New("can't setup. gophermart isn't configured")
var errSetupGophermartNotReadyToRun = errors.New("can't run. gophermart isn't set up")

func main() {
	var err error
	if err = g.cofigure(); err != nil {
		log.Fatal(g.loggingCtx, err)
	}
	if err = g.setup(); err != nil {
		log.Fatal(g.loggingCtx, err)
	}
	if err = g.run(); err != nil {
		log.Fatal(g.loggingCtx, err)
	}
}
