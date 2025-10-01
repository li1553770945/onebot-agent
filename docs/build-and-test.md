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

## 构建和登录lagrange-onebot

```bash
cd lagrange-onebot

## 构建
docker build -t lagrange-onebot:local . 


## 运行测试
docker run -d  --name lagrange-onebot --network onebot-agent -p 15000:15000  -v ./app:/app/ lagrange-onebot:local 
```
运行上述命令后，如果之前未登录过，会在`lagrange-onebot/app`目录生成qr-0.png，需要扫描该图片二维码登录，登录成功后会生成keystore.json，也可查看容器log确保登录成功，此时容器的15000端口将对外开放为lagrange.onebot的正向HTTP服务。

## 构建和登录icpc-nanjing-agent

```bash
cd icpc-nanjing-agent
## 构建
docker build -t icpc-nanjing-agent:local .
## 运行测试
docker run -d  --name icpc-nanjing-agent --network onebot-agent -p 3000:3000  -v ./config:/app/config icpc-nanjing-agent:local 
```

