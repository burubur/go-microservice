package servd

import (
	"os"
	"time"

	"bitbucket.org/kudoindonesia/frontier_biller_sdk/deduplication"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/env"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
	paymentSDK "bitbucket.org/kudoindonesia/frontier_biller_sdk/payment"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/pkg/database"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/tracer"
	"bitbucket.org/kudoindonesia/koolkit/koollog"
	"github.com/burubur/go-microservice/internal/app/payment"
	"github.com/burubur/go-microservice/internal/platform/configuration"
)

// Payment will serve a payment service
func (s *Service) Payment() {
	config := configuration.Values
	productHandlers := initPaymentProductHandler()
	paymentHandler := payment.New(productHandlers)

	envName, err := env.GetEnv(config.App.Env)
	if err != nil {
		log.Fatal("invalid app.env value, please check the yaml config values!", koollog.Err(err))
	}

	paymentConsumerConfig := paymentSDK.Config{
		Env:                     envName,
		AppName:                 config.App.Name,
		EncryptionKey:           config.Redis.EncryptionKey,
		Source:                  config.Kafka.PaymentTopic,
		SourceServer:            config.Kafka.Broker,
		PendingProcessorSource:  config.Kafka.PendingProcessorTopic,
		TransactionUpdateSource: config.Kafka.UpdateTransactionTopic,
		DeduplicationType:       deduplication.PersistanceDeduplication,
		MessageExpiryTime:       config.Worker.Payment.MessageExpireTime * time.Minute,
		TempStorageKeyPrefix:    config.Redis.CachePrefix,
		TempStorageHostServer:   []string{config.Redis.Host},
		CheckStatusTopic:        config.Kafka.CheckStatusTopic,
		Handler:                 paymentHandler.Handle,
		DBConfig: database.SQLConfig{
			Host:                  config.Database.Host,
			Port:                  config.Database.Port,
			Username:              config.Database.User,
			Password:              config.Database.Pass,
			DbName:                config.Database.Name,
			MaxOpenConnection:     config.Database.MaxOpenConnection,
			MaxIdleConnection:     config.Database.MaxIdleConnection,
			MaxConnectionLifeTime: config.Database.MaxConnectionLifeTime,
		},
		TracerType:              tracer.DataDogTracerType,
		TracerAgentAddress:      config.Datadog.ApmAddress,
		TracerRateSampler:       config.Datadog.ApmRateSampler,
		WithoutInquiry:          false,
		NumberOfConsumer:        config.Worker.Payment.NumberOfConsumer,
		NumberOfWorker:          config.Worker.Payment.NumberOfWorker,
		MaxCommitTryCount:       config.Worker.Payment.CommitBackoff,
		NumberOfMessageToCommit: config.Worker.Payment.CommitEvery,
		CommitMessageDuration:   config.Worker.Payment.CommitWithin,
	}

	paymentConsumer := paymentSDK.NewPayment(paymentConsumerConfig)
	paymentChannel := make(chan bool)
	interuptionSignal := make(chan os.Signal, 1)

	log.Info("starting payment service")
	go shutdown(paymentChannel, interuptionSignal)
	paymentConsumer.Run(paymentChannel)
	log.Info("stopped payment service")
}

func initPaymentProductHandler() map[string]payment.Payer {
	// in case you need to retrieve config values, unremark the following line
	// config := configuration.Values

	productHandlers := map[string]payment.Payer{
		// instantigithub.com/burubur/go-microservice/internalatae all available product handler here!
	}

	return productHandlers
}
