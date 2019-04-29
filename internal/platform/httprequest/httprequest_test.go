package httprequest_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/burubur/go-microservice/internal/platform/httprequest"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestNew(t *testing.T) {
	got := httprequest.New(httprequest.Config{})
	assert.NotNil(t, got.HTTPClient, "it should have a default http client instance")
	assert.NotNil(t, got.SendRequest, "it should implement a SendRequest method")
}

func TestClient_SendRequest(t *testing.T) {
	httpClient := &http.Client{}
	httpmock.ActivateNonDefault(httpClient)
	defer httpmock.DeactivateAndReset()

	type fields struct {
		HTTPClient *http.Client
	}
	type args struct {
		url     string
		action  string
		payload []byte
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		httpResponder    httpmock.Responder
		wantStatusCode   int
		wantResponseBody []byte
		wantErr          bool
	}{
		{
			name: "0. it should return an error when given an invalid http protocl schema",
			fields: fields{
				HTTPClient: httpClient,
			},
			args: args{
				url:     "url-test", // this is an unsupported protocol schema
				action:  "action-test",
				payload: []byte("payload-test"),
			},
			wantStatusCode: 0,
			wantErr:        true,
		},
		{
			name: "1. it should return an error when http status code is non 200",
			fields: fields{
				HTTPClient: httpClient,
			},
			args: args{
				url:     "http://server.mock",
				action:  "action-test",
				payload: []byte("payload-test"),
			},
			httpResponder:    httpmock.NewStringResponder(http.StatusBadRequest, ""),
			wantStatusCode:   http.StatusBadRequest,
			wantResponseBody: []byte(""),
			wantErr:          true,
		},
		{
			name: "2. it should return an error when the response body from the partner is empty",
			fields: fields{
				HTTPClient: httpClient,
			},
			args: args{
				url:     "http://server.mock",
				action:  "action-test",
				payload: []byte("payload-test"),
			},
			httpResponder:    httpmock.NewStringResponder(http.StatusOK, ""),
			wantStatusCode:   http.StatusOK,
			wantResponseBody: []byte(""),
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			httprequestor := httprequest.Client{
				HTTPClient: httpClient,
			}
			httpmock.RegisterResponder(http.MethodPost, "http://server.mock", tt.httpResponder)
			gotStatusCode, gotResponseBody, err := httprequestor.SendRequest(tt.args.url, tt.args.action, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatusCode != tt.wantStatusCode {
				t.Errorf("Client.SendRequest() gotStatusCode = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
			if (err == nil) == tt.wantErr {
				if !reflect.DeepEqual(gotResponseBody, tt.wantResponseBody) {
					t.Errorf("Client.SendRequest() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
				}
			}
		})
	}
}
