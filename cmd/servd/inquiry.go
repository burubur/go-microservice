package servd

import (
	"os"
	"time"

	"bitbucket.org/kudoindonesia/frontier_biller_sdk/env"
	inquirySDK "bitbucket.org/kudoindonesia/frontier_biller_sdk/inquiry"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/tracer"
	"bitbucket.org/kudoindonesia/koolkit/koollog"
	"github.com/burubur/go-microservice/internal/app/inquiry"
	"github.com/burubur/go-microservice/internal/platform/configuration"
)

// Inquiry will serve an inquiry service
func (s *Service) Inquiry() {
	config := configuration.Values
	productHandlers := initInquiryProductHandler()
	inquiryHandler := inquiry.New(productHandlers)

	envName, err := env.GetEnv(config.App.Env)
	if err != nil {
		log.Fatal("invalid app.env value, please check the yaml config values!", koollog.Err(err))
	}

	inqConsumerCfg := inquirySDK.Config{
		Handler:                      inquiryHandler.Handle,
		Env:                          envName,
		AppName:                      config.App.Name,
		Source:                       config.Kafka.InquiryTopic,
		SourceServer:                 config.Kafka.Broker,
		EncryptionKey:                config.Redis.EncryptionKey,
		SuccessInquiryCacheRetention: config.Redis.SuccessInquiryRetention * time.Minute,
		FailedInquiryCacheRetention:  config.Redis.FailedInquiryRetention * time.Minute,
		TempStorageKeyPrefix:         config.Redis.CachePrefix,
		TempStorageHostServer:        []string{config.Redis.Host},
		StoreInquiryResultSource:     config.Kafka.InquiryResultTopic,
		TracerType:                   tracer.DataDogTracerType,
		TracerAgentAddress:           config.Datadog.ApmAddress,
		TracerRateSampler:            config.Datadog.ApmRateSampler,
		NumberOfConsumer:             config.Worker.Inquiry.NumberOfConsumer,
		NumberOfWorker:               config.Worker.Inquiry.NumberOfWorker,
		MaxCommitTryCount:            config.Worker.Inquiry.CommitBackoff,
		NumberOfMessageToCommit:      config.Worker.Inquiry.CommitEvery,
		CommitMessageDuration:        config.Worker.Inquiry.CommitWithin,
	}

	inquiryConsumer := inquirySDK.NewInquiry(inqConsumerCfg)
	inquiryChannel := make(chan bool)
	interuptionSignal := make(chan os.Signal, 1)

	log.Info("starting inquiry service")
	go shutdown(inquiryChannel, interuptionSignal)
	inquiryConsumer.Run(inquiryChannel)
	log.Info("stopped inquiry service")
}

func initInquiryProductHandler() map[string]inquiry.Inquirer {
	// in case you need to retrieve config values, unremark the following line
	// config := configuration.Values
	productHandlers := map[string]inquiry.Inquirer{
		// instantiatae all available product handler here!
	}
	return productHandlers
}
