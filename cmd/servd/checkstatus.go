package servd

import (
	"os"

	csSDK "bitbucket.org/kudoindonesia/frontier_biller_sdk/checkstatus"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/env"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/tracer"
	"bitbucket.org/kudoindonesia/koolkit/koollog"
	"github.com/burubur/go-microservice/internal/app/checkstatus"
	"github.com/burubur/go-microservice/internal/platform/configuration"
)

// CheckStatus will provide a checkstatus service
func (s *Service) CheckStatus() {
	config := configuration.Values

	checkStatusHandlers := initCheckStatusProductHandler()
	checkStatusHandler := checkstatus.New(checkStatusHandlers)

	envName, err := env.GetEnv(config.App.Env)
	if err != nil {
		log.Fatal("invalid app.env value, please check the yaml config values!", koollog.Err(err))
	}

	csConsumerCfg := csSDK.Config{
		Env:                     envName,
		AppName:                 config.App.Name,
		EncryptionKey:           config.Redis.EncryptionKey,
		Source:                  config.Kafka.CheckStatusTopic,
		SourceServer:            config.Kafka.Broker,
		PendingProcessorSource:  config.Kafka.PendingProcessorTopic,
		TransactionUpdateSource: config.Kafka.UpdateTransactionTopic,
		MessageExpiryTime:       config.Worker.Payment.MessageExpireTime,
		TempStorageKeyPrefix:    config.Redis.CachePrefix,
		TempStorageHostServer:   []string{config.Redis.Host},
		Handler:                 checkStatusHandler.Handle,
		TracerType:              tracer.DataDogTracerType,
		TracerAgentAddress:      config.Datadog.ApmAddress,
		TracerRateSampler:       config.Datadog.ApmRateSampler,
		WithoutInquiry:          false,
		NumberOfConsumer:        config.Worker.CheckStatus.NumberOfConsumer,
		NumberOfWorker:          config.Worker.CheckStatus.NumberOfWorker,
		MaxCommitTryCount:       config.Worker.CheckStatus.CommitBackoff,
		NumberOfMessageToCommit: config.Worker.CheckStatus.CommitEvery,
		CommitMessageDuration:   config.Worker.CheckStatus.CommitWithin,
	}

	csConsumer := csSDK.NewCheckstatus(csConsumerCfg)
	csCH := make(chan bool)
	interuptionSignal := make(chan os.Signal, 1)

	log.Info("starting checkstatus service")
	go shutdown(csCH, interuptionSignal)
	csConsumer.Run(csCH)
	log.Info("stopped checkstatus service")
}

func initCheckStatusProductHandler() map[string]checkstatus.StatusChecker {
	// in case you need to retrieve config values, unremark the following line
	// config := configuration.Values

	productHandlers := map[string]checkstatus.StatusChecker{
		// instantigithub.com/burubur/go-microservice/internalatae all available product handler here!
	}

	return productHandlers
}
