package http_client

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type httpBinAnything struct {
	Args struct {
		Name string `json:"name"`
	} `json:"args"`
	Data string `json:"data"`
}

func TestHttpClientGet(t *testing.T) {
	hc := NewHttpClient(3)
	body, err := Get(hc, "https://httpbin.org/anything?name=test")
	if err != nil {
		t.Error(err)
	}

	var anything httpBinAnything
	json.Unmarshal(body, &anything)

	if anything.Args.Name != "test" {
		t.Error("test name, expect test, got", anything.Args.Name)
	}
}

func TestHttpClientHead(t *testing.T) {
	hc := NewHttpClient(3)
	body, err := Head(hc, "https://httpbin.org/anything?name=test")
	if err != nil {
		t.Error(err)
	}

	if len(body) != 0 {
		t.Error("response body should be empty")
	}
}

func TestHttpClientPost(t *testing.T) {
	hc := NewHttpClient(3)
	postData := "{\"post\":\"new data\"}"
	postDataReader := strings.NewReader(postData)
	body, err := PostJson(hc, "https://httpbin.org/anything?name=test", postDataReader)
	if err != nil {
		t.Error(err)
	}

	var anything httpBinAnything
	json.Unmarshal(body, &anything)

	if anything.Args.Name != "test" {
		t.Error("test name, expect test, got", anything.Args.Name)
	}

	if anything.Data != postData {
		t.Error("post data, expect", postData, "got", anything.Data)
	}
}

func TestHttpClientDelete(t *testing.T) {
	hc := NewHttpClient(3)
	_, err := Delete(hc, "https://httpbin.org/anything?name=test")
	if err != nil {
		t.Error(err)
	}
}

func TestHttpClientWithTimeout(t *testing.T) {
	hc := NewHttpClient(1)
	if _, err := Get(hc, "https://httpbin.org/delay/2"); !os.IsTimeout(err) {
		t.Errorf("http client should not timeout, got %v", err)
	}
}

func TestDumbHttpClientGet(t *testing.T) {
	hc := NewDumbHttpClient(3)
	body, err := Get(hc, "https://httpbin.org/anything?name=test")
	if err != nil {
		t.Error(err)
	}

	if len(body) != 0 {
		t.Errorf("response body should be empty, got %s", string(body))
	}
}

func TestDumbHttpClientHead(t *testing.T) {
	hc := NewDumbHttpClient(3)
	body, err := Head(hc, "https://httpbin.org/anything?name=test")
	if err != nil {
		t.Error(err)
	}

	if len(body) != 0 {
		t.Error("response body should be empty")
	}
}

func TestDumbHttpClientPost(t *testing.T) {
	hc := NewDumbHttpClient(3)
	postData := "{\"post\":\"new data\"}"
	postDataReader := strings.NewReader(postData)
	body, err := PostJson(hc, "https://httpbin.org/anything?name=test", postDataReader)
	if err != nil {
		t.Error(err)
	}

	if len(body) != 0 {
		t.Error("response body should be empty")
	}
}

func TestDumbHttpClientDelete(t *testing.T) {
	hc := NewDumbHttpClient(3)
	body, err := Delete(hc, "https://httpbin.org/anything?name=test")
	if err != nil {
		t.Error(err)
	}

	if len(body) != 0 {
		t.Error("response body should be empty")
	}
}

func TestHttpDumbClientWithTimeout(t *testing.T) {
	hc := NewDumbHttpClient(1)
	if _, err := Get(hc, "https://httpbin.org/delay/2"); !os.IsTimeout(err) {
		t.Errorf("http client should not timeout, got %v", err)
	}
}
