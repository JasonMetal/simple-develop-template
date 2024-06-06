package curl

import (
	"bytes"
	"crypto/tls"
	strLib "github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/strings"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	type args struct {
		req Request
	}
	tests := []struct {
		name     string
		args     args
		wantResp Response
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := Get(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Get() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestHttpBuildQuery(t *testing.T) {
	tests := []struct {
		params map[string]string
		want   string
	}{
		// TODO: Add test cases.
		{
			map[string]string{
				"name": "张三",
				"age":  "27",
			},
			"age=27&name=%E5%BC%A0%E4%B8%89",
		},
		{
			map[string]string{
				"ie": "UTF-8",
				"wd": "facebook",
			},
			"ie=UTF-8&wd=facebook",
		},
	}
	for _, tt := range tests {
		if got := HttpBuildQuery(tt.params); got != tt.want {
			t.Errorf("HttpBuildQuery() = %v, want %v", got, tt.want)
		}
	}
}

func TestPost(t *testing.T) {
	type args struct {
		req Request
	}
	tests := []struct {
		name     string
		args     args
		wantResp Response
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := Post(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Post() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	MaxIdleConns:        0,
	MaxIdleConnsPerHost: 400,
	IdleConnTimeout:     30 * time.Second,
}

func Test_initClient(t *testing.T) {
	req1 := Request{
		Url:    "http://test-site.cc",
		Method: "GET",
		Headers: map[string]string{
			"ContentType": "text/html",
		},
		BodyData: map[string]string{
			"name": "张三",
			"age":  "22",
		},
		ReqTimeOut:  10,
		RespTimeOut: 12,
	}
	want1 := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   400,
			IdleConnTimeout:       30 * time.Second,
			ResponseHeaderTimeout: time.Second * req1.RespTimeOut,
		},
	}
	req2 := Request{
		Url:    "http://test-site.cc",
		Method: "GET",
		Headers: map[string]string{
			"ContentType": "text/html",
		},
		BodyData: map[string]string{
			"name": "张三",
			"age":  "22",
		},
		ReqTimeOut:  10,
		RespTimeOut: 0,
	}
	want2 := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   400,
			IdleConnTimeout:       30 * time.Second,
			ResponseHeaderTimeout: 0,
		},
	}
	tests := []struct {
		request    Request
		wantClient http.Client
	}{
		// TODO: Add test cases.
		{
			req1,
			want1,
		},
		{
			req2,
			want2,
		},
	}
	for _, tt := range tests {
		if gotClient := initClient(tt.request); !reflect.DeepEqual(gotClient, tt.wantClient) {
			t.Errorf("initClient() = %v, want %v", gotClient, tt.wantClient)
		}
	}
}

func Test_send(t *testing.T) {
	test1 := Request{
		Url:         "http://demo.52zsj.com/index/user/go_test",
		Method:      "GET",
		Headers:     nil,
		BodyData:    nil,
		ReqTimeOut:  0,
		RespTimeOut: 0,
	}
	test1R := Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     nil,
		Body:       strLib.StrToBytes("GET"),
	}
	var test1E bool = false

	test2 := Request{
		Url:         "http://demo.52zsj.com/index/user/go_test",
		Method:      "POST",
		Headers:     nil,
		BodyData:    nil,
		ReqTimeOut:  1,
		RespTimeOut: 1,
	}
	test2R := Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     nil,
		Body:       strLib.StrToBytes("POST"),
	}
	var test2E bool = false
	test3 := Request{
		Url:         "http://127.0.0.1/ok",
		Method:      "POST",
		Headers:     nil,
		BodyData:    nil,
		ReqTimeOut:  1,
		RespTimeOut: 1,
	}
	test3R := Response{
		Status:     "404 Not Found ",
		StatusCode: 404,
		Header:     nil,
		Body:       strLib.StrToBytes("无法预料的BODY"),
	}
	var test3E bool = false
	tests := []struct {
		request  Request
		wantResp Response
		wantErr  interface{}
	}{
		// TODO: Add test cases.
		{
			test1,
			test1R,
			test1E,
		},
		{
			test2,
			test2R,
			test2E,
		},
		{
			test3,
			test3R,
			test3E,
		},
	}
	for _, tt := range tests {
		gotResp, err := send(tt.request)
		if (err != nil) != tt.wantErr {
			t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(gotResp, tt.wantResp) {
			t.Errorf("send() gotResp = %v, want %v", gotResp, tt.wantResp)
		}
	}
}

func Test_setBodyData(t *testing.T) {
	str1 := "name=张三&age=27"
	want1 := strings.NewReader(str1)

	str2 := "{\n    \"code\": 1,\n    \"data\": [],\n    \"message\": \"请求成功\"\n}"
	want2 := strings.NewReader(str2)

	byte1 := []byte{}
	wantb1 := bytes.NewBuffer(byte1)
	maps1 := map[string]string{
		"name": "age",
		"age":  "22",
	}
	wantm1 := strings.NewReader(HttpBuildQuery(maps1))
	tests := []struct {
		bodyData interface{}
		want     io.Reader
	}{
		// TODO: Add test cases.
		{
			str1,
			want1,
		},
		{
			str2,
			want2,
		},
		{
			byte1,
			wantb1,
		},
		{
			maps1,
			wantm1,
		},
		{
			12,
			strings.NewReader("123"),
		},
	}
	for _, tt := range tests {
		if got := setBodyData(tt.bodyData); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("setBodyData() = %v, want %v", got, tt.want)
		}
	}
}
