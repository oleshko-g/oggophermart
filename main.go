package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	genAccrual "github.com/oleshko-g/oggophermart/internal/gen/accrual"
	genBalance "github.com/oleshko-g/oggophermart/internal/gen/balance"
	genAccrualHTTPClient "github.com/oleshko-g/oggophermart/internal/gen/http/accrual/client"
	genUser "github.com/oleshko-g/oggophermart/internal/gen/user"
	"github.com/oleshko-g/oggophermart/internal/service"
	balance "github.com/oleshko-g/oggophermart/internal/service/balance"
	user "github.com/oleshko-g/oggophermart/internal/service/user"
	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql"
	"github.com/oleshko-g/oggophermart/internal/transport/http"
	"goa.design/clue/debug"
	"goa.design/clue/log"
)

type gophermart struct {
	transport struct {
		http struct {
			http.Server
			http.Config
		}
	}
	service.Service
	accrual struct {
		genAccrual.Service
		genAccrualHTTPClient.Client
		// TODO: add Config
	}
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
//
// If successful it sets readyToRun flag
func (g *gophermart) setup() (err error) {
	if !g.configured {
		return errSetupGophermartNotConfigured
	}
	dbStorage, err := sql.New(&g.dbCfg)

	// 1. Sets the storage for each service
	// wrap concrete type [*sql.Storage] struct with interfaces
	g.Storage.User = dbStorage
	g.Storage.Balance = dbStorage

	//  2. Intanciates services with the set storage
	g.Service = service.Service{
		Balance: balance.New(g.Storage.Balance),
		User:    user.New(g.Storage.User),
	}

	//  3. Instanciates the HTTP server
	g.transport.http.Server = http.NewServer(g.transport.http.Config, g.Service)

	g.readyToRun = true
	return nil
}

// run launches gophermart
func (g *gophermart) run() (err error) {
	if !g.readyToRun {
		return errSetupGophermartNotReadyToRun
	}
	return g.transport.http.Server.ListenAndServe()
}

var errSetupGophermartNotConfigured = errors.New("can't setup. gophermart isn't configured")
var errSetupGophermartNotReadyToRun = errors.New("can't run. gophermart isn't set up")

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "8080", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logger. Replace logger with your own log package of choice.
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	if *dbgF {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}
	log.Print(ctx, log.KV{K: "http-port", V: *httpPortF})

	// Initialize the services.
	var (
		balanceSvc genBalance.Service
		userSvc    genUser.Service
	)
	{
		balanceSvc = balance.NewBalance()
		userSvc = user.NewUser()
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		balanceEndpoints *genBalance.Endpoints
		userEndpoints    *genUser.Endpoints
	)
	{
		balanceEndpoints = genBalance.NewEndpoints(balanceSvc)
		balanceEndpoints.Use(debug.LogPayloads())
		balanceEndpoints.Use(log.Endpoint)
		userEndpoints = genUser.NewEndpoints(userSvc)
		userEndpoints.Use(debug.LogPayloads())
		userEndpoints.Use(log.Endpoint)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://localhost:80"
			u, err := url.Parse(addr)
			if err != nil {
				log.Fatalf(ctx, err, "invalid URL %#v\n", addr)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					log.Fatalf(ctx, err, "invalid URL %#v\n", u.Host)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, balanceEndpoints, userEndpoints, &wg, errc, *dbgF)
		}

	default:
		log.Fatal(ctx, fmt.Errorf("invalid host argument: %q (valid hosts: localhost)", *hostF))
	}

	// Wait for signal.
	log.Printf(ctx, "exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	log.Printf(ctx, "exited")
}
