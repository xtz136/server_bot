package http_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClientInterface interface {
	Fetch(req *http.Request) ([]byte, error)
}

type HttpClient struct {
	Client *http.Client
}

func (hc *HttpClient) Fetch(req *http.Request) ([]byte, error) {
	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s http status code is %v", req.Method, req.URL, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

type DumbHttpClient struct {
	Client *http.Client
}

func (hc *DumbHttpClient) Fetch(req *http.Request) ([]byte, error) {
	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s http status code is %v", req.Method, req.URL, resp.StatusCode)
	}
	return []byte{}, nil
}

func Head(ctx context.Context, hc HttpClientInterface, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return hc.Fetch(req)
}

func Get(ctx context.Context, hc HttpClientInterface, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return hc.Fetch(req)
}

func PostJson(ctx context.Context, hc HttpClientInterface, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	return hc.Fetch(req)
}

func Post(ctx context.Context, hc HttpClientInterface, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)
	return hc.Fetch(req)
}

func Delete(ctx context.Context, hc HttpClientInterface, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return hc.Fetch(req)
}

// 不需要返回值的http客户端，自动关闭连接、丢弃返回值
// 下面例子展示了用法，并且resp永远都是nil
//  dhc := http_client.NewDumbHttpClient(10)
//  resp, err := http_client.PostJson(dhc, url, body)
func NewDumbHttpClient(timeout time.Duration) *DumbHttpClient {
	client := http.Client{Timeout: timeout}
	return &DumbHttpClient{
		Client: &client,
	}
}

// 自动关闭连接的http客户端
// 下面例子展示了用法
//  dhc := http_client.NewHttpClient(10)
//  resp, err := http_client.PostJson(dhc, url, body)
func NewHttpClient(timeout time.Duration) *HttpClient {
	client := http.Client{Timeout: timeout}
	return &HttpClient{
		Client: &client,
	}
}

func init() {
	// 跳过ssl证书验证
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
