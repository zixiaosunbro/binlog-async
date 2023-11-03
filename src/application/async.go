package application

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"

	"binlog-async/src"
)

var (
	once sync.Once
	// mux maps type to handler
	mux = asynq.NewServeMux()
)

func InitAsyncSvc() {
	srv := asynq.NewServer(src.RedisCfg{}, asynq.Config{
		Concurrency: 2,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})
	log.Info("start async service")
	if err := srv.Run(mux); err != nil {
		panic(fmt.Sprintf("start asynq fail:%s", err.Error()))
	}
	// wait for SIGINT or SIGTERM signals to stop
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	srv.Shutdown()
	log.Println("asynq service stopped")

}

func init() {
	once.Do(func() {
		// TODO: register task type with handler
	})
}
