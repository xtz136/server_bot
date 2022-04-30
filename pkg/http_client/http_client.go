package http_client

import (
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
	return nil, nil
}

func Head(hc HttpClientInterface, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return hc.Fetch(req)
}

func Get(hc HttpClientInterface, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return hc.Fetch(req)
}

func PostJson(hc HttpClientInterface, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return hc.Fetch(req)
}

func Post(hc HttpClientInterface, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return hc.Fetch(req)
}

func Delete(hc HttpClientInterface, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return hc.Fetch(req)
}

// 不需要返回值的http客户端，自动关闭连接、丢弃返回值
func NewDumbHttpClient(timeout time.Duration) *DumbHttpClient {
	client := http.Client{Timeout: timeout * time.Second}
	return &DumbHttpClient{
		Client: &client,
	}
}

// 自动关闭连接的http客户端
func NewHttpClient(timeout time.Duration) *HttpClient {
	client := http.Client{Timeout: timeout * time.Second}
	return &HttpClient{
		Client: &client,
	}
}

func init() {
	// 跳过ssl证书验证
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
