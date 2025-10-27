# 构建和测试命令

## 创建测试网络

由于使用了k8s的DNS解析，测试时需要创建一个docker网络，命令如下：
```bash
docker network create onebot-agent
```

## 构建和登录message-dispatch

```bash
cd ./message-dispatch
## 构建
docker build -t message-dispatch:local .
## 运行测试
docker run -d --name message-dispatch --network onebot-agent -p 15001:15001 -v ./config:/app/config message-dispatch:local
```




