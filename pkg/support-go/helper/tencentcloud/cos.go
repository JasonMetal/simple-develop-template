package tencentcloud

import (
	"context"
	"github.com/JasonMetal/simple-develop-template/pkg/support-go/bootstrap"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"os"
)

func IsExist(object, bucket, region string) bool {
	rs := false
	u, _ := url.Parse("https://" + bucket + ".cos." + region + ".myqcloud.com")

	url := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(url, &http.Client{
		Transport: setTransport(),
	})
	_, err := c.Object.Head(context.Background(), object, nil)
	if err != nil {
		bootstrap.CheckError(err, "cos")
	}

	return rs
}

func setTransport() *cos.AuthorizationTransport {
	return &cos.AuthorizationTransport{
		// 通过环境变量获取密钥
		// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
		SecretID: os.Getenv("COS_SECRETID"),
		// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
		SecretKey: os.Getenv("COS_SECRETKEY"),
		// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
		Transport: &debug.DebugRequestTransport{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
		},
	}
}
