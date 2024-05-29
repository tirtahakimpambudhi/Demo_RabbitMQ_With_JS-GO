package main

import (
	"context"
	"fmt"
	"go_test_rabbitmq/config"
	"go_test_rabbitmq/internal/entity"
	"go_test_rabbitmq/internal/logger"
	"go_test_rabbitmq/message_broker"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initial config, log, connection, exit, context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cnf := config.NewAppConfig()
	log := logger.NewLogger()
	exit := make(chan os.Signal, 1)
	errChan := make(chan error, 1) // Buffered to avoid blocking
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	connection, err := message_broker.ConnectionRabbitMQ(cnf.GetString("rabbitmq.cloud"))
	failError(err, "failed connecting to RabbitMQ", log)
	defer connection.Close()

	// Initial Exchange, Queue, Bind
	channel, err := connection.Channel()
	failError(err, "failed to create channel", log)
	defer channel.Close()

	err = channel.ExchangeDeclare(
		cnf.GetString("process_service.exchange"),
		amqp091.ExchangeDirect,
		true,  // durable
		false, // auto delete
		false, // no wait
		false, // internal
		nil,   // args
	)
	failError(err, "failed to create exchange", log)

	queue, err := channel.QueueDeclare(
		cnf.GetString("process_service.queue"),
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // args
	)
	failError(err, "failed to create queue", log)

	err = channel.QueueBind(
		queue.Name,
		cnf.GetString("process_service.routing_key"),
		cnf.GetString("process_service.exchange"),
		false, // no wait
		nil,   // args
	)
	failError(err, "failed to create bind", log)

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-exit:
				cancel()
				done <- true
				return
			case err := <-errChan:
				if err != nil {
					log.Error(err.Error())
				}
			default:
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				message := &entity.Message{
					OS:        runtime.GOOS,
					Architecture: runtime.GOARCH,
					CPUs:      runtime.NumCPU(),
					Threads:   runtime.GOMAXPROCS(0),
					GoRoutine: runtime.NumGoroutine(),
					MemoryAllocated: memStats.Alloc,
					TotalMemoryAllocated: memStats.TotalAlloc,
					MemoryObtainedSystem: memStats.Sys,
					GarbaceCollection: memStats.NumGC,
					GoVersion: runtime.GOROOT(),
				}
				json, err := message.ToJSON()
				if err != nil {
					errChan <- err
					continue
				}
				err = channel.PublishWithContext(ctx, cnf.GetString("process_service.exchange"), cnf.GetString("process_service.routing_key"), false, false, amqp091.Publishing{
					ContentType: "application/json",
					Body:        json,
				})
				if err != nil {
					errChan <- err
				}

				time.Sleep(1 * time.Second)
			}
		}
	}()

	<-done
	fmt.Println("Exiting program")
}

func failError(err error, message string, log *logrus.Logger) {
	if err != nil {
		log.Fatalf("%s: %s", message, err.Error())
	}
}
