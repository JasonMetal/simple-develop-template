package tencentcloud

import (
	"fmt"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/bootstrap"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/entity"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// GetToken 获取STS临时授权
func GetToken(project string, objectNameList []string) (entity.StsToken, error) {
	token := entity.StsToken{}
	if c, ok := bootstrap.StsClientList[project]; ok {
		opt := &sts.CredentialOptions{
			DurationSeconds: int64(c.DurationSeconds),
			Region:          c.Region,
			Policy:          setPolicy(c.BucketName, c.AppId, objectNameList),
		}
		res, err := c.Client.GetCredential(opt)
		if err != nil {
			bootstrap.CheckError(err, "sts")

			return token, err
		}

		token = entity.StsToken{
			SecretID:     res.Credentials.TmpSecretID,
			SecretKey:    res.Credentials.TmpSecretKey,
			SessionToken: res.Credentials.SessionToken,
			Bucket:       c.BucketName,
			Region:       c.Region,
			Domain:       c.Domain,
		}
	}

	return token, nil
}

// setPolicy 设置STS权限策略
func setPolicy(bucket string, appId string, objectName []string) *sts.CredentialPolicy {
	return &sts.CredentialPolicy{
		Statement: []sts.CredentialPolicyStatement{
			{
				// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
				Action:   getAllowAction(),
				Effect:   "allow",
				Resource: getResource(bucket, appId, objectName),
			},
		},
	}
}

// getResource 生成授权的resource列表
func getResource(bucket string, appId string, objectName []string) []string {
	resource := make([]string, 0)
	for _, v := range objectName {
		object := fmt.Sprintf("qcs::cos:ap-guangzhou:uid/%s:%s/%s", appId, bucket, v)
		resource = append(resource, object)
	}

	return resource
}

// getAllowAction 设置STS允许操作的行为
func getAllowAction() []string {
	return []string{
		// 简单上传
		"cos:PostObject",
		"cos:PutObject",
		// 分片上传
		"cos:InitiateMultipartUpload",
		"cos:ListMultipartUploads",
		"cos:ListParts",
		"cos:UploadPart",
		"cos:CompleteMultipartUpload",
	}
}

// GetBucketByType 根据类型获取bucket
func GetBucketByType(uploadType int) {

}
