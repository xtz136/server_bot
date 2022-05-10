package http_client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func makeContext() context.Context {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Millisecond)
	return ctx
}

func TestHead(t *testing.T) {
	type args struct {
		ctx context.Context
		hc  HttpClientInterface
		url string
	}
	tests := []struct {
		name      string
		args      args
		sleepTime time.Duration
		want      []byte
		wantErr   error
	}{
		{"basic httpclient", args{makeContext(), NewHttpClient(1 * time.Second), ""}, 0, []byte(""), nil},
		{"basic dumb httpclient", args{makeContext(), NewDumbHttpClient(1 * time.Second), ""}, 0, []byte(""), nil},
		{"timeout httpclient", args{makeContext(), NewHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"timeout dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel httpclient", args{makeContext(), NewHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.Method != "HEAD" {
					http.Error(rw, "Not Found", http.StatusNotFound)
				}

				time.Sleep(tt.sleepTime)
				fmt.Fprint(rw, "OK")
			}))
			defer server.Close()

			tt.args.url = server.URL

			got, err := Head(tt.args.ctx, tt.args.hc, tt.args.url)
			if (err != nil) && errors.Is(tt.wantErr, err) {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
		hc  HttpClientInterface
		url string
	}
	tests := []struct {
		name      string
		args      args
		sleepTime time.Duration
		want      []byte
		wantErr   error
	}{
		{"basic httpclient", args{makeContext(), NewHttpClient(1 * time.Second), ""}, 0, []byte("OK"), nil},
		{"basic dumb httpclient", args{makeContext(), NewDumbHttpClient(1 * time.Second), ""}, 0, []byte(""), nil},
		{"timeout httpclient", args{makeContext(), NewHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"timeout dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel httpclient", args{makeContext(), NewHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.Method != "GET" {
					http.Error(rw, "Not Found", http.StatusNotFound)
				}
				time.Sleep(tt.sleepTime)
				fmt.Fprint(rw, "OK")
			}))
			defer server.Close()
			tt.args.url = server.URL

			got, err := Get(tt.args.ctx, tt.args.hc, tt.args.url)
			if (err != nil) && errors.Is(tt.wantErr, err) {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		hc  HttpClientInterface
		url string
	}
	tests := []struct {
		name      string
		args      args
		sleepTime time.Duration
		want      []byte
		wantErr   error
	}{
		{"basic httpclient", args{makeContext(), NewHttpClient(1 * time.Second), ""}, 0, []byte("OK"), nil},
		{"basic dumb httpclient", args{makeContext(), NewDumbHttpClient(1 * time.Second), ""}, 0, []byte(""), nil},
		{"timeout httpclient", args{makeContext(), NewHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"timeout dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel httpclient", args{makeContext(), NewHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.Method != "DELETE" {
					http.Error(rw, "Not Found", http.StatusNotFound)
				}
				time.Sleep(tt.sleepTime)
				fmt.Fprint(rw, "OK")
			}))
			defer server.Close()
			tt.args.url = server.URL

			got, err := Delete(tt.args.ctx, tt.args.hc, tt.args.url)
			if (err != nil) && errors.Is(tt.wantErr, err) {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		ctx  context.Context
		hc   HttpClientInterface
		body string
	}
	tests := []struct {
		name      string
		args      args
		sleepTime time.Duration
		want      []byte
		wantErr   error
	}{
		{"basic httpclient", args{makeContext(), NewHttpClient(1 * time.Second), ""}, 0, []byte("OK"), nil},
		{"basic dumb httpclient", args{makeContext(), NewDumbHttpClient(1 * time.Second), ""}, 0, []byte(""), nil},
		{"timeout httpclient", args{makeContext(), NewHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"timeout dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel httpclient", args{makeContext(), NewHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.Method != "POST" {
					http.Error(rw, "Not Found", http.StatusNotFound)
				}
				if req.Header["Content-Type"][0] != "application/x-www-form-urlencoded" {
					http.Error(rw, "Invalid", http.StatusInternalServerError)
				}
				time.Sleep(tt.sleepTime)
				fmt.Fprint(rw, "OK")
			}))
			defer server.Close()

			got, err := Post(tt.args.ctx, tt.args.hc, server.URL, bytes.NewBuffer([]byte(tt.args.body)))
			if (err != nil) && errors.Is(tt.wantErr, err) {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostJson(t *testing.T) {
	type args struct {
		ctx  context.Context
		hc   HttpClientInterface
		body string
	}
	tests := []struct {
		name      string
		args      args
		sleepTime time.Duration
		want      []byte
		wantErr   error
	}{
		{"basic httpclient", args{makeContext(), NewHttpClient(1 * time.Second), ""}, 0, []byte("OK"), nil},
		{"basic dumb httpclient", args{makeContext(), NewDumbHttpClient(1 * time.Second), ""}, 0, []byte(""), nil},
		{"timeout httpclient", args{makeContext(), NewHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"timeout dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Millisecond), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel httpclient", args{makeContext(), NewHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
		{"context cancel dumb httpclient", args{makeContext(), NewDumbHttpClient(10 * time.Second), ""}, 20 * time.Millisecond, nil, context.DeadlineExceeded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.Method != "POST" {
					http.Error(rw, "Not Found", http.StatusNotFound)
				}
				if req.Header["Content-Type"][0] != "application/json" {
					http.Error(rw, "Invalid", http.StatusInternalServerError)
				}
				time.Sleep(tt.sleepTime)
				fmt.Fprint(rw, "OK")
			}))
			defer server.Close()

			got, err := PostJson(tt.args.ctx, tt.args.hc, server.URL, bytes.NewBuffer([]byte(tt.args.body)))
			if (err != nil) && errors.Is(tt.wantErr, err) {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
