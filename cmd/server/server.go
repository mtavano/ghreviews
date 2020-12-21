package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/mtavano/ghreviews/database"
	"github.com/mtavano/ghreviews/pkg/graph"
	"github.com/mtavano/ghreviews/pkg/service"
)

var (
	listenAddr string
	dbHostname string
	dbPort     int
	dbName     string
	dbUsername string
	dbPassword string
	dbSslMode  string
	dbUrl      string
	dbDriver   string
	dbLogs     bool
	env        string
)

func main() {
	parseFlags()

	logger := logrus.New()
	dataSourceName := getDatabaseUrl()

	store, err := database.NewStore(dbDriver, dataSourceName)
	if err != nil {
		panic(err)
	}

	isProduction := env == "production"
	if !isProduction {
		logger.SetLevel(logrus.DebugLevel)
	}

	reviewService := service.NewReviewService(logger, store)
	r := graph.NewResolver(logger, reviewService)
	graphqlServer := graph.NewServer(r, isProduction)

	router := http.NewServeMux()
	router.Handle("/query", graphqlServer)
	if env != "production" {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	logger.Printf("connect to http://%s/ for GraphQL playground", listenAddr)
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	w := logger.Writer()
	defer w.Close()
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      cors.AllowAll().Handler(router),
		ErrorLog:     log.New(w, "", 0),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Info("Server is ready to handle requests at ", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Info("Server stopped")
}

func getDatabaseUrl() string {
	if dbUrl == "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUsername, dbPassword, dbHostname, dbPort, dbName)
	}

	return dbUrl
}

func parseFlags() {
	flag.StringVar(&dbHostname, "db-hostname", "localhost", "database hostname")
	flag.IntVar(&dbPort, "db-port", 54320, "database port")
	flag.StringVar(&dbName, "db-name", "ghreviews", "database name")
	flag.StringVar(&dbUsername, "db-username", "apipath", "database username")
	flag.StringVar(&dbPassword, "db-password", "apipath", "database password")
	flag.BoolVar(&dbLogs, "db-verbose", false, "database print queries")
	flag.StringVar(&dbSslMode, "db-sslmode", "disabled", "database ssl mode")
	flag.StringVar(&dbUrl, "db-url", "", "database url - this flag preceds all the other db flags")
	flag.StringVar(&dbDriver, "db-driver", "postgres", "database driver")
	flag.StringVar(&listenAddr, "listen-addr", "localhost:8080", "server listen address")
	// TODO: improve value parsing and panic if invalid env is passed
	flag.StringVar(&env, "env", "development", "application environment (development, production)")
	flag.Parse()
}
