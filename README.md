# server bot

在服务器上运行的机器人，可以通过文字对话机器人，让机器人做出对应的操作。

## 安装

1. 将本项目拷贝到 $GOPATH/src/ 目录下
```bash
git clone https://github.com/xtz136/server_bot $GOPATH/src/server_bot
```

2. 新建配置文件
```bash
cd $GOPATH/src/server_bot
cp configs/config.yaml.example ./config.yaml
```

3. 安装依赖
```bash
go mod tidy
```

## 启动开发模式

1. 安装 air，或者直接看[air的官方文档](https://github.com/cosmtrek/air)
```bash
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

2. 启动
```bash
$GOPATH/bin/air -c .air.conf
```

## 启动生产环境

1. 打包
```
go build -o server_bot .
```

2. 拷贝到生产环境

把打包好的二进制文件`server_bot`和配置文件`config.yaml`都拷贝到自己的服务器上，放在同一个目录下。

过程省略

3. 启动

```bash
GIN_MODE=relase ./server_bot
```

4. 配置nginx，在nginx增加如下配置，确保 proxy_pass 里面的端口和项目配置文件的端口一样

这里使用了二级目录，可以根据实际需要修改

```bash
location /server_bot {
    proxy_set_header	X-Real-IP	$remote_addr;
    proxy_set_header	Host	$http_host;
    proxy_pass          http://127.0.0.1:8080;
}
```

## 运行测试

go test ./...
