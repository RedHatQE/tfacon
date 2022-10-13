package common

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestPrintGreen(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintGreen(tt.args.str)
		})
	}
}

func TestPrintRed(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintRed(tt.args.str)
		})
	}
}

func TestPrintHeader(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintHeader(tt.args.version)
		})
	}
}

func TestSendHTTPRequest(t *testing.T) {
	type args struct {
		ctx        context.Context
		method     string
		url        string
		auth_token string
		body       io.Reader
		client     *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := SendHTTPRequest(tt.args.ctx, tt.args.method, tt.args.url, tt.args.auth_token, tt.args.body, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendHTTPRequest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SendHTTPRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleError(tt.args.err, "nopanic")
		})
	}
}
