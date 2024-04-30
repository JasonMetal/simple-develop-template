package bootstrap

import (
	"fmt"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/helper/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
)

type CosClient struct {
	AppId           string
	Client          *http.Client
	Region          string
	BucketName      string
	DurationSeconds int
	Bucket          string
	Domain          string
}

var CosClientList = make(map[string]CosClient)

func InitCos() {
	path := fmt.Sprintf("%sconfig/%s/tencentyun.yml", ProjectPath(), DevEnv)

	cosConfigs, err := config.GetConfig(path)

	configList, err := cosConfigs.Map("tencentyun")
	if err == nil {
		for project, _ := range configList {
			secretId, _ := cosConfigs.String("tencentyun." + project + ".secretId")
			secretKey, _ := cosConfigs.String("tencentyun." + project + ".secretKey")
			appId, _ := cosConfigs.String("tencentyun." + project + ".appId")
			region, _ := cosConfigs.String("tencentyun." + project + ".Region")
			durationSeconds, _ := cosConfigs.Int("tencentyun." + project + ".durationSeconds")
			bucketName, _ := cosConfigs.String("tencentyun." + project + ".bucket")
			domain, _ := cosConfigs.String("tencentyun." + project + ".domain")

			c := &http.Client{
				Transport: setTransport(secretId, secretKey),
			}
			CosClientList[project] = CosClient{
				AppId:           appId,
				Client:          c,
				Region:          region,
				DurationSeconds: durationSeconds,
				BucketName:      fmt.Sprintf("%s-%s", bucketName, appId),
				Bucket:          bucketName,
				Domain:          domain,
			}
		}
	}
}

func setTransport(secretID, secretKey string) *cos.AuthorizationTransport {
	return &cos.AuthorizationTransport{
		SecretID:  secretID,
		SecretKey: secretKey,
		// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
		//Transport: &debug.DebugRequestTransport{
		//	RequestHeader:  true,
		//	RequestBody:    true,
		//	ResponseHeader: true,
		//	ResponseBody:   true,
		//},
	}
}
