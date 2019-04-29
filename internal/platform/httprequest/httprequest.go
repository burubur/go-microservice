package httprequest

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Config is a configuration that will be used when constructing a new instance of Requestor
type Config struct {
	MaxIdleConnection    int
	IdleConnTimeout      time.Duration
	MaxConnectionPerHost int
	HTTPRequestTimeout   time.Duration

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// This should be used only for testing.
	InsecureSkipVerify bool
	Certificate        *tls.Certificate
}

// Client will handle http request response by extending go http.Client package
type Client struct {
	HTTPClient *http.Client
}

// New will construct a customized http client
func New(config Config) *Client {
	transport := &http.Transport{
		MaxIdleConns:    config.MaxIdleConnection,
		IdleConnTimeout: config.IdleConnTimeout * time.Second,
		MaxConnsPerHost: config.MaxConnectionPerHost,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
			Renegotiation:      tls.RenegotiateFreelyAsClient,
		},
	}
	if config.Certificate != nil {
		transport.TLSClientConfig.Certificates = []tls.Certificate{*config.Certificate}
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.HTTPRequestTimeout * time.Second,
	}

	client := &Client{}
	client.HTTPClient = httpClient
	return client
}

// SendRequest will hit a defined endpoint and return a response body in byte format
func (httprequest *Client) SendRequest(url, action string, payload []byte) (int, []byte, error) {
	var statusCode int
	var responseBody []byte

	body := bytes.NewBuffer(payload)
	httpRequest, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return statusCode, responseBody, err
	}

	httpRequest.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpRequest.Header.Add("Accept", "text/xml")
	httpRequest.Header.Set("SOAPAction", action)

	response, err := httprequest.HTTPClient.Do(httpRequest)
	if err != nil {
		return statusCode, responseBody, err
	}

	defer func() {
		if err = response.Body.Close(); err != nil {
			err = errors.New("cannot close response body")
		}
	}()

	statusCode = response.StatusCode

	if statusCode != http.StatusOK {
		err = fmt.Errorf("retrieved a non 200 httpStatusCode, got: %v", statusCode)
		return statusCode, responseBody, err
	}

	responseBody, err = ioutil.ReadAll(response.Body)

	if len(responseBody) == 0 {
		err = errors.New("response body is empty")
	}

	return statusCode, responseBody, err
}
