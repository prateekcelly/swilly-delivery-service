package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"swilly-delivery-service/config"
	"swilly-delivery-service/internal/app"
	"swilly-delivery-service/internal/pkg/log"
	"swilly-delivery-service/internal/pkg/middleware"
	"time"

	"github.com/gocraft/work"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func StartServer() {
	idleConsClosed := make(chan struct{})
	server, err := newServer()
	if err != nil {
		log.Fatal("unable to create the server", zap.Error(err))
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.stop(ctx); err != nil {
			log.Error("server shutdown error", zap.Error(err))
		}
		idleConsClosed <- struct{}{}
		close(idleConsClosed)
	}()
	ctx := context.Background()
	_ = middleware.ProcessWithRecovery(ctx, func(ctx context.Context) error {
		server.start(ctx)
		return nil
	})

	<-idleConsClosed
}

func (s *Server) start(ctx context.Context) {
	fp, err := NewFileProcessor(config.AppConfig.DirectoryPath, work.NewEnqueuer("delivery", app.AppDependency.Redis))
	if err != nil {
		log.Fatal("Error initializing file processor: %v", zap.Error(err))
	}
	go fp.Start(ctx)

	log.Info("starting app", zap.String("port", config.AppConfig.HTTPServerPort))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server start error", zap.Error(err))
		panic(err)
	}
}

func (s *Server) stop(ctx context.Context) error {
	log.Info("stopping server")
	return s.httpServer.Shutdown(ctx)
}

func newServer() (*Server, error) {
	var err error
	if err = app.Bootstrap(); err != nil {
		return nil, err
	}

	s := &Server{
		httpServer: &http.Server{},
	}
	return s, nil
}
