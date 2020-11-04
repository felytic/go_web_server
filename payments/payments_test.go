package payments

import (
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
		responses []ProviderResponse
		want      bool
	}{
		{"Nil", nil, true},
		{
			"AllOk",
			[]ProviderResponse{
				ProviderResponse{"go1", "http://go1.com", nil},
				ProviderResponse{"go2", "http://go2.com", nil},
			},
			true,
		},
		{
			"OneOk",
			[]ProviderResponse{
				ProviderResponse{"go1", "", nil},
				ProviderResponse{"go2", "http://go2.com", nil},
			},
			false,
		},
		{
			"NotOk",
			[]ProviderResponse{
				ProviderResponse{"go1", "http://error.com", errors.New("Err")},
				ProviderResponse{"go2", "http://go2.com", nil},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllResponsesOk(tt.responses); got != tt.want {
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
		want *ProviderResponse
	}{
		{
			"nil",
			args{
				"nil",
				nil,
				errors.New("Bad request"),
			},
			&ProviderResponse{
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
			&ProviderResponse{
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
			&ProviderResponse{
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
			&ProviderResponse{
				"badRespWithNoErr",
				"",
				errors.New("500 Internal Server Error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constructProviderResponse(tt.args.name, tt.args.response, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("constructProviderResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
