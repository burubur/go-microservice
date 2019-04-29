package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	validator "gopkg.in/asaskevich/govalidator.v9"
	yaml "gopkg.in/yaml.v2"
)

const defaultPath = "configurations/app.yaml"

// Values is a public var that will keep all the loaded configuration values
var Values *Structure

// Structure is a configuration structure, this structure have a tag for validation
// the validation method is using govalidator library
type Structure struct {
	App struct {
		Name  string `yaml:"name" valid:"required"`
		Env   string `yaml:"env"` // the value should be, "development", "staging", "production" or "kudobox"
		Debug bool   `yaml:"debug"`
	} `yaml:"app"`
	Database struct {
		Host                  string        `yaml:"host" valid:"required"`
		Port                  int           `yaml:"port" valid:"required"`
		User                  string        `yaml:"user" valid:"required"`
		Pass                  string        `yaml:"pass"`
		Name                  string        `yaml:"name" valid:"required"`
		MaxOpenConnection     int           `yaml:"max_open_connection" valid:"required"`
		MaxIdleConnection     int           `yaml:"max_idle_connection" valid:"required"`
		MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time" valid:"required"`
	} `yaml:"database"`
	Kafka struct {
		Broker                 string `yaml:"broker" valid:"required"`
		InquiryTopic           string `yaml:"inquiry_topic" valid:"required"`
		InquiryResultTopic     string `yaml:"inquiry_result_topic" valid:"required"`
		PaymentTopic           string `yaml:"payment_topic" valid:"required"`
		CheckStatusTopic       string `yaml:"check_status_topic" valid:"required"`
		PendingProcessorTopic  string `yaml:"pending_processor_topic" valid:"required"`
		UpdateTransactionTopic string `yaml:"update_transaction_topic" valid:"required"`
	} `yaml:"kafka"`
	Worker struct {
		Inquiry struct {
			NumberOfConsumer int           `yaml:"number_of_consumer" valid:"required"`
			NumberOfWorker   int           `yaml:"number_of_worker" valid:"required"`
			CommitEvery      int           `yaml:"commit_every" valid:"required"`
			CommitBackoff    int           `yaml:"commit_backoff" valid:"required"`
			CommitWithin     time.Duration `yaml:"commit_within" valid:"required"`
		} `yaml:"inquiry"`
		Payment struct {
			NumberOfConsumer  int           `yaml:"number_of_consumer" valid:"required"`
			NumberOfWorker    int           `yaml:"number_of_worker" valid:"required"`
			CommitEvery       int           `yaml:"commit_every" valid:"required"`
			CommitBackoff     int           `yaml:"commit_backoff" valid:"required"`
			CommitWithin      time.Duration `yaml:"commit_within" valid:"required"`
			MessageExpireTime time.Duration `yaml:"message_expire_time" valid:"required"`
		} `yaml:"payment"`
		CheckStatus struct {
			NumberOfConsumer int           `yaml:"number_of_consumer" valid:"required"`
			NumberOfWorker   int           `yaml:"number_of_worker" valid:"required"`
			CommitEvery      int           `yaml:"commit_every" valid:"required"`
			CommitBackoff    int           `yaml:"commit_backoff" valid:"required"`
			CommitWithin     time.Duration `yaml:"commit_within" valid:"required"`
			RetryCount       int           `yaml:"retry_count" valid:"required"`
			RetryEvery       int           `yaml:"retry_every" valid:"required"`
			RetryInDay       int           `yaml:"retry_in_day" valid:"required"`
		} `yaml:"check_status"`
	} `yaml:"worker"`
	Redis struct {
		Host                    string        `yaml:"host" valid:"required"`
		SuccessInquiryRetention time.Duration `yaml:"success_inquiry_retention" valid:"required"`
		FailedInquiryRetention  time.Duration `yaml:"failed_inquiry_retention" valid:"required"`
		CachePrefix             string        `yaml:"cache_prefix" valid:"required"`
		EncryptionKey           string        `yaml:"encryption_key" valid:"required"`
	} `yaml:"redis"`
	HTTP struct {
		MaxIdleConnection    int           `yaml:"max_idle_connection" valid:"required"`
		IdleConnTimeout      time.Duration `yaml:"idle_conn_timeout" valid:"required"`
		MaxConnectionPerHost int           `yaml:"max_connection_per_host" valid:"required"`
		RequestTimeout       time.Duration `yaml:"request_timeout" valid:"required"`
		InsecureSkipVerify   bool          `yaml:"insecure_skip_verify" valid:"required"`
		// CertificatePath      string        `yaml:"certificate_path" valid:"required"`
	} `yaml:"http"`
	Biller struct {
		PrivateKeyPath string `yaml:"private_key_path" valid:"required"`
		EndPoint       string `yaml:"end_point" valid:"required"`
	} `yaml:"biller"`
	Datadog struct {
		ApmAddress     string  `yaml:"apm_address" valid:"required"`
		ApmRateSampler float64 `yaml:"apm_rate_sampler" valid:"required"`
		MetricAddress  string  `yaml:"metric_address" valid:"required"`
	} `yaml:"datadog"`
	ConfigVersion string `yaml:"_config_version"`
}

// Load will cal a c.load()
func Load(path string) error { return load(path) }

// Validate doc
func Validate() error { return validate() }

// Reset doc
func Reset() { reset() }

func load(path string) error {
	if Values != nil {
		err := errors.New("configuration file already loaded")
		return err
	}
	if path == "" {
		path = defaultPath
		ex, err := os.Executable()
		if err != nil {
			return fmt.Errorf("error when os.Executable(), got: %v", err)
		}
		path = filepath.Join(filepath.Dir(ex), filepath.Clean(path))
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("unable to open config file on this path: %s, got: %v", path, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			err = errors.New("cannot close file")
		}
	}()

	var byteBuffer = bytes.NewBuffer([]byte{})

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	p, err := byteBuffer.Write(data)
	if err != nil {
		return err
	}

	if p == 0 {
		return errors.New("failed to read. 0 data is readed")
	}

	return yaml.UnmarshalStrict(byteBuffer.Bytes(), &Values)
}

func validate() (err error) {
	_, err = validator.ValidateStruct(Values)
	return err
}

func reset() {
	Values = nil
}
