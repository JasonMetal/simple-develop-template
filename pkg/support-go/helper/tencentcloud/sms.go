package tencentcloud

import (
	"encoding/json"
	"errors"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/bootstrap"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/helper/slices"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"strings"
)

type SmsParams struct {
	TemplateId, SignName   string
	TemplateParams, Mobile []string
}

// Send 发送短信
func (s *SmsParams) Send() (b []byte, err error) {
	c := bootstrap.TencentSmsClient

	// 使用新签名
	if s.SignName != "" {
		c.Request.SignName = common.StringPtr(s.SignName)
	}
	if s.TemplateId == "" {
		err = errors.New("模板ID不能为空!")
		return nil, err
	}
	if len(s.Mobile) == 0 {
		err = errors.New("手机号不能为空!")
		return nil, err
	}
	//防止因切片而修改原手机号的值
	var mobileList = make([]string, 0)

	// 处理手机号 自动兼容 带+号的手机
	for _, value := range s.Mobile {
		//如果存在+号则直接放进去
		if strings.Index(value, "+") >= 0 {
			mobileList = append(mobileList, value)
		} else {
			mobileList = append(mobileList, c.MobileCode+value)
		}
	}
	//去重手机号
	mobileList = slices.Unique(mobileList)
	// 处理模板参数
	if len(s.TemplateParams) > 0 {
		c.Request.TemplateParamSet = common.StringPtrs(s.TemplateParams)
	}
	// 处理模板ID
	c.Request.TemplateId = common.StringPtr(s.TemplateId)
	c.Request.PhoneNumberSet = common.StringPtrs(mobileList)
	response, err := c.Client.SendSms(c.Request)

	b, err = json.Marshal(response.Response)
	return b, err
}
