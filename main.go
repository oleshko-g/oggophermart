package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	userCfg    user.Config
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

	secretKey := g.userCfg.SecretAuthKey()
	if v, ok := os.LookupEnv("JWT_SECRET"); ok {
		err = secretKey.Set(v)
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
	log.Infof(g.loggingCtx, "Connected the storage")

	// 1. Sets the storage for each service
	// wrap concrete type [*sql.Storage] struct with interfaces
	g.Storage.User = dbStorage
	log.Infof(g.loggingCtx, "set User service storage")
	g.Storage.Balance = dbStorage
	log.Infof(g.loggingCtx, "set Balance service storage")

	// 2. Intanciates services with the set storage
	userSvc := user.New(&g.userCfg, g.Storage.User)
	g.Service = service.Service{
		User:    userSvc,
		Balance: balance.New(g.Storage.Balance, userSvc),
	}

	// 3. Instanciates the HTTP server
	g.transport.http.Server = http.NewServer(g.loggingCtx, g.transport.http.Config, g.Service)

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
	if !g.readyToRun {
		return errSetupGophermartNotReadyToRun
	}

	var wg sync.WaitGroup

	runCtx, runCancel := context.WithCancelCause(g.loggingCtx)
	defer runCancel(nil)

	osSignals := make(chan os.Signal, 1)
	wg.Go(func() {
		log.Infof(g.loggingCtx, "in os.Signal goroutine")
		signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
		err = fmt.Errorf("recieved OS signal: %s", <-osSignals)
		runCancel(err)
	})

	wg.Go(func() {
		log.Infof(g.loggingCtx, "in HTTP server")
		errHTTP := g.transport.http.Server.ListenAndServe()
		if errHTTP != nil {
			runCancel(errHTTP)
		}
		log.Printf(g.loggingCtx, "gophermart HTTP server is listening on %s", g.transport.http.Address().String())

		<-runCtx.Done()

		ctx := context.Background()
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel()
		g.transport.http.Server.Shutdown(shutdownCtx)
		err = runCtx.Err()
	})

	wg.Wait()

	return err
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
