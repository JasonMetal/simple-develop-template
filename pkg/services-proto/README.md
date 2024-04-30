# services-proto

#### 介绍
proto文件仓库

# go proto 生成
```
1.go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
2.安装protoc代码生成工具 https://github.com/google/protobuf
```

# proto 命名方式
```
rule:[services.项目名]
package services.user[user项目]
package services.XXX[XXX项目]

```

# 由proto文件生成go代码

```shell
cd services-proto/proto/user;
protoc --proto_path=./  --go_out=plugins=grpc:../../pb-go/user/ user.proto 
```

