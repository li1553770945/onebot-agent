# 构建命令

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
docker run -d --name message-dispatch --network onebot-agent -p 8080:8080 -v ./config:/app/config message-dispatch:local
```

## 构建和登录lagrange-onebot

```bash
cd lagrange-onebot

## 构建
docker build -t lagrange-onebot:local . 


## 运行测试
docker run -d  --name lagrange-onebot --network onebot-agent -p 8000:8000  -v ./app:/app/ lagrange-onebot:local 
```
运行上述命令后，如果之前未登录过，会在`lagrange-onebot/app`目录生成qr-0.png，需要扫描该图片二维码登录，登录成功后会生成keystore.json，也可查看容器log确保登录成功，此时容器的8000端口将对外开放为lagrange.onebot的正向HTTP服务。



# 部署命令

在使用部署命令之前，请确保上述所有构建操作已完成，尤其是onebot的登录操作完成，且测试无问题。