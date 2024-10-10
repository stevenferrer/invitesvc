package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/openapi"
	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/token"
)

//go:embed static
var embededFiles embed.FS

// List of defalt server configs
const (
	defaultDSN  = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	defaultHost = "localhost"
	defaultPort = 8000
)

func main() {
	var (
		host = flag.String("host", defaultHost, "server host")
		port = flag.Int("port", defaultPort, "server port")
		dsn  = envStr("DSN", defaultDSN)
	)

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// connect to database
	db, err := sql.Open("postgres", envStr("DSN", dsn))
	if err != nil {
		logger.Fatal().Err(err).Msg("open database")
	}
	defer db.Close()

	err = retry(logger, 10, time.Second, db.Ping)
	if err != nil {
		logger.Fatal().Err(err).Msg("ping database")
	}

	// migrate the database
	err = postgres.Migrate(db)
	if err != nil {
		logger.Fatal().Err(err).Msg("migrate database")
	}

	// initialize services
	var tokenSvc token.Service
	{
		tokenRepo := postgres.NewTokenRepository(db)
		tokenSvc = token.NewService(tokenRepo)
	}

	var authSvc authn.Service
	{
		authRepo := postgres.NewAuthRepository(db)
		authSvc = authn.NewAuthService(authRepo)
	}

	ctx := context.Background()

	// generate initial auth key
	// TODO: don't generate new auth keys when there are more than 1 keys already
	authKey, err := authSvc.GenerateAuthKey(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("generate new auth key")
	}

	logger.Info().Str("authKey", string(authKey)).Msg("initial auth key")

	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// static files
	staticFiles, err := getStaticFiles()
	if err != nil {
		logger.Fatal().Err(err).Msg("get static files")
	}
	assetHandler := http.FileServer(staticFiles)
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	// openapi3 spec
	openapi.InitOpenAPI3Routes(e)

	// admin and public routes
	token.InitAdminRoutes(e, tokenSvc, authSvc)
	token.InitPublicRoutes(e, tokenSvc)

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", *host, *port),
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start server
	go func() {
		logger.Info().Msgf("listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal().Err(err).Msg("listen and serve")
		}
	}()

	// setup signal capturing
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// wait for SIGINT
	<-c

	// shutdown server
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server shutdown")
	}
}

func envStr(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func retry(
	logger zerolog.Logger,
	attempts int,
	sleep time.Duration,
	f func() error,
) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			logger.Error().Err(err).Msg("retrying after error")
			time.Sleep(sleep)
			sleep *= 2
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func getStaticFiles() (http.FileSystem, error) {
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		return nil, err
	}

	return http.FS(fsys), nil
}
