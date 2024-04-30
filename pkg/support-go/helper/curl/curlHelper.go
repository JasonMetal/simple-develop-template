package curl

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Url         string
	Method      string
	Headers     map[string]string
	BodyData    interface{}
	ReqTimeOut  time.Duration
	RespTimeOut time.Duration
}
type Response struct {
	Status     string
	StatusCode int
	Header     map[string]string
	Body       []byte
}

// Get HTTP GET 请求
func Get(req Request) (resp Response, err error) {
	req.Method = "GET"
	return send(req)
}

// POST HTTP POST 请求
func Post(req Request) (resp Response, err error) {
	req.Method = "POST"
	return send(req)
}

// send 发送请求
func send(request Request) (resp Response, err error) {
	var client = initClient(request)

	//请求超时时间设置
	var requestTimeout time.Duration = 5
	if request.ReqTimeOut > 0 {
		requestTimeout = 5
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*requestTimeout)
	defer cancel()
	//设置body
	bodyData := setBodyData(request.BodyData)
	// 设置上下文请求
	req, err := http.NewRequestWithContext(ctx, request.Method, request.Url, bodyData)
	if err != nil {
		return resp, err
	}
	//设置请求头
	for k, v := range request.Headers {
		req.Header.Add(k, v)
	}
	response, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer response.Body.Close()
	resp.Status = response.Status
	resp.StatusCode = response.StatusCode
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return resp, err
	}
	resp.Body = body
	header := make(map[string][]string)
	for k, v := range response.Header {
		header[k] = v
	}
	return resp, nil

}

// HttpBuildQuery 格式化url会按照字母进行排序
func HttpBuildQuery(params map[string]string) string {
	var uri url.URL
	q := uri.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	queryStr := q.Encode()
	return queryStr
}

// setBodyData 设置body
func setBodyData(bodyData interface{}) io.Reader {
	var data io.Reader
	switch v := bodyData.(type) {
	case string:
		data = strings.NewReader(v)
	case []byte:
		data = bytes.NewBuffer(v)
	case map[string]string:
		data = strings.NewReader(HttpBuildQuery(bodyData.(map[string]string)))
	}
	return data
}

// initClient 获取HTTP.Client实例
func initClient(request Request) (client http.Client) {
	// 忽略 https 证书校验
	var responseheadertimeout time.Duration = 0
	if request.RespTimeOut > 0 {
		responseheadertimeout = time.Second * request.RespTimeOut
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   400,
		IdleConnTimeout:       30 * time.Second,
		ResponseHeaderTimeout: responseheadertimeout,
	}

	client = http.Client{Transport: transport}

	return client
}
