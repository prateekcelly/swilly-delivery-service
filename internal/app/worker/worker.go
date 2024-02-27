package worker

import (
	"context"
	"os"
	"os/signal"
	"swilly-delivery-service/config"
	"swilly-delivery-service/internal/app"
	"swilly-delivery-service/internal/pkg/log"

	"github.com/gocraft/work"
	"go.uber.org/zap"
)

func StartWorker(ctx context.Context) error {
	if err := app.Bootstrap(); err != nil {
		return err
	}

	pool := work.NewWorkerPool(ctx, 10, "delivery", app.AppDependency.Redis)
	pool.JobWithOptions(config.AppConfig.JobName, work.JobOptions{
		MaxFails: 3,
		SkipDead: false,
	}, triggerAlert)
	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	pool.Stop()

	return nil
}

func triggerAlert(job *work.Job) error {
	// Extract arguments from the job
	userID := job.ArgString("userID")
	message := job.ArgString("message")
	if err := job.ArgError(); err != nil {
		return err
	}

	log.Info("Job Arguments", zap.String("userID", userID), zap.String("message", message))
	// Make HTTP call to webhook api. More info in documentation

	return nil
}
