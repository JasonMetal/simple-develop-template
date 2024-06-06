package bootstrap

import (
	"fmt"
	"github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/config"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type StsClient struct {
	AppId           string
	Client          *sts.Client
	Region          string
	BucketName      string
	DurationSeconds int
	Bucket          string
	Domain          string
}

var StsClientList = make(map[string]StsClient)

func InitSts() {
	initDefault()
}

func initDefault() {
	path := fmt.Sprintf("%sconfig/%s/tencentyun.yml", ProjectPath(), DevEnv)

	stsConfigs, err := config.GetConfig(path)

	configList, err := stsConfigs.Map("tencentyun")
	if err == nil {
		for project, _ := range configList {
			secretId, _ := stsConfigs.String("tencentyun." + project + ".secretId")
			secretKey, _ := stsConfigs.String("tencentyun." + project + ".secretKey")
			appId, _ := stsConfigs.String("tencentyun." + project + ".appId")
			region, _ := stsConfigs.String("tencentyun." + project + ".Region")
			durationSeconds, _ := stsConfigs.Int("tencentyun." + project + ".durationSeconds")
			bucketName, _ := stsConfigs.String("tencentyun." + project + ".bucket")
			domain, _ := stsConfigs.String("tencentyun." + project + ".domain")
			c := sts.NewClient(secretId, secretKey, nil)
			StsClientList[project] = StsClient{
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
