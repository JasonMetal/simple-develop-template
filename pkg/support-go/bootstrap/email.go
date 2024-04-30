package bootstrap

import (
	"fmt"
	configHelper "github.com/BioforestChain/simple-develop-template/pkg/support-go/helper/config"
)

type emailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

var EmailConfig = &emailConfig{}

func InitEmail() {
	initDefaultConfigForEmail()
}

// 初始化默认配置
func initDefaultConfigForEmail() {
	configName := "default"
	path := fmt.Sprintf("%sconfig/%s/email.yml", ProjectPath(), DevEnv)
	emailConfigs, err := configHelper.GetConfig(path)
	if err != nil {
		return
	}
	EmailConfig.Host, _ = emailConfigs.String(configName + ".host")
	EmailConfig.Port, _ = emailConfigs.Int(configName + ".port")
	EmailConfig.Username, _ = emailConfigs.String(configName + ".username")
	EmailConfig.Password, _ = emailConfigs.String(configName + ".password")
}
