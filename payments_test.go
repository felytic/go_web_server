package main

import (
	"./payments"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestIsAllResponsesOk(t *testing.T) {
	tests := []struct {
		name      string
		responses []payments.ProviderResponse
		want      bool
	}{
		{"Nil", nil, true},
		{
			"AllOk",
			[]payments.ProviderResponse{
				payments.ProviderResponse{"go1", "http://go1.com", nil},
				payments.ProviderResponse{"go2", "http://go2.com", nil},
			},
			true,
		},
		{
			"OneOk",
			[]payments.ProviderResponse{
				payments.ProviderResponse{"go1", "", nil},
				payments.ProviderResponse{"go2", "http://go2.com", nil},
			},
			false,
		},
		{
			"NotOk",
			[]payments.ProviderResponse{
				payments.ProviderResponse{"go1", "http://error.com", errors.New("Err")},
				payments.ProviderResponse{"go2", "http://go2.com", nil},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := payments.IsAllResponsesOk(tt.responses); got != tt.want {
				t.Errorf("IsAllResponsesOk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_constructProviderResponse(t *testing.T) {

	dummyReq, _ := http.NewRequest("GET", "http://go.com", strings.NewReader(""))

	okResp := new(http.Response)
	okResp.Request = dummyReq
	okResp.Body = ioutil.NopCloser(bytes.NewReader([]byte("http://go.com")))
	okResp.StatusCode = 200

	badResp := new(http.Response)
	badResp.Request = dummyReq
	badResp.StatusCode = 500
	badResp.Status = "500 Internal Server Error"

	type args struct {
		name     string
		response *http.Response
		err      error
	}
	tests := []struct {
		name string
		args args
		want *payments.ProviderResponse
	}{
		{
			"nil",
			args{
				"nil",
				nil,
				errors.New("Bad request"),
			},
			&payments.ProviderResponse{
				"nil",
				"",
				errors.New("Bad request"),
			},
		},
		{
			"OK",
			args{
				"OK",
				okResp,
				nil,
			},
			&payments.ProviderResponse{
				"OK",
				"http://go.com",
				nil,
			},
		},
		{
			"notOK",
			args{
				"notOK",
				okResp,
				errors.New("err"),
			},
			&payments.ProviderResponse{
				"notOK",
				"",
				errors.New("err"),
			},
		},
		{
			"badRespWithNoErr",
			args{
				"badRespWithNoErr",
				badResp,
				nil,
			},
			&payments.ProviderResponse{
				"badRespWithNoErr",
				"",
				errors.New("500 Internal Server Error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := payments.ConstructProviderResponse(tt.args.name, tt.args.response, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructProviderResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
