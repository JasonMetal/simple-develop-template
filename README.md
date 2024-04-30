# 介绍
`gin-simple-develop-template`

From the previous project, slightly made some challenges, open the box to use, without submodules trouble
Gin-based web backend api framework for business development

# 1. 初始化项目
替换项目中import的 `simple-develop-template` 为项目名称

# 2. go mod 初始化
go mod init 项目名称

# 3. 设置私有仓库
go env -w GOPRIVATE=*.github.com

# 4. 整理go mod
go mod tidy

# 5. 编译部署相关
#### 部署运行
#### 0. go build -o go-test cli.go
#### 测试服范例
#### 1. ./go-test -e test savePageDataCron
#### linux下用 `supervisord` 进行监控相关任务


```shell
#### 查看状态
supervisorctl status|grep goTestDemo
supervisorctl restart goTestDemo
#### conf文件路径
 /etc/supervisord.d/conf/goTestDemo.conf

#### 某台服务
[xxx@test go-websites]# cat /etc/supervisord.d/conf/goTestDemo.conf
[program:goTestDemo]
directory = /home/www/demo
command = /home/www/demo/go-test -e test savePageDataCron
autostart = true
autorestart = true
loglevel = info
stdout_logfile = /var/log/supervisor/goTestDemo.log
stderr_logfile = /var/log/supervisor/goTestDemo_stderr.log
stdout_logfile_maxbytes = 30MB
stdout_logfile_backups = 3
stdout_events_enabled = false
```
