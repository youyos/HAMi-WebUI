package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type HttpClient struct {
	client *http.Client
}

// NewHttpClient 创建带超时的 HTTP 客户端
func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get 发送 GET 请求
func (hc *HttpClient) Get(ctx context.Context, rawUrl string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawUrl, nil)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// PostJSON 发送 POST 请求，Body 是 JSON
func (hc *HttpClient) PostJSON(ctx context.Context, rawUrl string, data interface{}, headers map[string]string) ([]byte, int, error) {
	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// PostForm 发送 POST 表单请求
func (hc *HttpClient) PostForm(ctx context.Context, rawUrl string, formData map[string]string, headers map[string]string) ([]byte, int, error) {
	data := url.Values{}
	for k, v := range formData {
		data.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

var (
	defaultClient *HttpClient
	once          sync.Once
)

func GetDefaultClient() *HttpClient {
	once.Do(func() {
		defaultClient = NewHttpClient(10 * time.Second)
	})
	return defaultClient
}
