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



# 部署命令

在使用部署命令之前，请确保上述所有构建操作已完成，尤其是onebot的登录操作完成，且测试无问题。

## 部署message-dispatch

```bash

cd ./message-dispatch
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
kubectl create configmap message-dispatch-config --from-file=production.yml=config/production.yml
```


## 部署lagrange-onebot

注意：请确保`lagrange-onebot/app`目录存在，并且在测试流程中已经登录，生成了对应的device.json、keystore.json等文件，否则k8s容器启动后还需要手动登录。

由于k8s限制无法使用相对路径，请修改`lagrange-onebot\k8s-pv.yml`中的hostPath路径为你的绝对路径，例如：

```yaml
      volumes:
      - name: data-storage
        hostPath: /absolute/path/to/lagrange-onebot/app
```

注意：如果你使用的是Docker Desktop的k8s集群，且你的项目目录在Windows的盘符下（例如E盘），**不能使用Windows上的绝对路径！**，不同Docker Desktop版本可能存在差异，请务必确认路径正确，例如可能为：

```yaml
      volumes:
      - name: data-storage
        hostPath: /run/desktop/mnt/host/e/node-project/onebot-agent/lagrange-onebot/app
```

最好是通过`kubectl apply -f k8s-pv.yml`命令后，查看pv和pvc状态是否为`Bound`来确认路径是否正确。

```bash

cd ./lagrange-onebot
kubectl apply -f k8s-pv.yml
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
```


# 其他

## 修改代码后重新部署

如果修改了某个服务的代码，需要重新构建镜像，并更新k8s的deployment，例如：

```bash
docker build -t lagrange-onebot:local .
kubectl rollout restart deployment/lagrange-onebot
```