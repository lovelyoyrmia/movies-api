package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lovelyoyrmia/movies-api/domain/pb"
	"github.com/lovelyoyrmia/movies-api/domain/services"
	"github.com/lovelyoyrmia/movies-api/internal/db"
	"github.com/lovelyoyrmia/movies-api/pkg/config"
	"github.com/lovelyoyrmia/movies-api/pkg/logger"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// Load configuration
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("error occurred : %v", err)
	}

	// Add interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	// Initialize Database
	database, err := db.NewDatabase(ctx, c)
	if err != nil {
		log.Error().Msg("cannot connect to database")
	}

	store := db.NewStore(database.DB)
	// Register wait group
	waitGroup, ctx := errgroup.WithContext(ctx)

	// Run Server
	runHTTPGatewayServer(waitGroup, ctx, c, store)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runHTTPGatewayServer(
	waitGroup *errgroup.Group,
	ctx context.Context,
	config config.Config,
	store db.Store,
) {

	// Declare json option for grpc request and response
	jsonOpt := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:     true,
			EmitUnpopulated:   true,
			EmitDefaultValues: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Create grpc serve mux
	grpcMux := runtime.NewServeMux(jsonOpt)

	// Register movies services
	movieService := services.NewMovieService(store)
	err := pb.RegisterMovieServicesHandlerServer(ctx, grpcMux, movieService)
	if err != nil {
		log.Fatal().Msgf("error occured : %v", err)
	}

	// Create http serve mux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Register middleware Redoc for API Spec
	ops := middleware.RedocOpts{SpecURL: fmt.Sprintf("%s%s", "./docs/swagger", "/movies.swagger.json")}
	sh := middleware.Redoc(ops, nil)
	mux.Handle("/docs", sh)

	// Serve API Docs
	fs := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle(fmt.Sprintf("/docs/swagger%s", "/movies.swagger.json"), http.StripPrefix("/docs/swagger/", fs))

	// Register http server
	s := &http.Server{
		Addr:         config.GatewayAddress,
		Handler:      logger.HTTPLogger(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	log.Info().Msgf("Starting HTTP server on port: %s", s.Addr)
	// Create go routines to serve http
	waitGroup.Go(func() error {
		err := s.ListenAndServe()
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Fatal().Msgf("cannot start GRPC server: %v", err)
			return err
		}
		return nil
	})

	// Waiting server to gracefully shutdown
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := s.Shutdown(context.Background())
		if err != nil {
			log.Error().Msg("failed to shutdown http gateway server")
			return err
		}
		log.Info().Msg("HTTP Gateway server is stopped")
		return nil
	})
}
