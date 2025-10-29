

# 部署命令

在使用部署命令之前，请确保所有构建操作已完成，尤其是onebot的登录操作完成，且测试无问题。

## 部署message-dispatch

```bash

cd ./message-dispatch
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
kubectl create configmap message-dispatch-config --from-file=production.yml=config/production.yml
```

如果后续更新了配置文件，可以使用以下命令更新configmap：

```bash
kubectl delete configmap message-dispatch-config
kubectl create configmap message-dispatch-config --from-file=production.yml=config/production.yml
kubectl rollout restart deployment message-dispatch
```

转发端口到localhost，因为机器人本身是在本机跑的：

```bash
kubectl port-forward service/message-dispatch 30001:15001
```

可以`curl http://localhost:30001/`确认能访问（404是正常的）。

## 部署napcat

下载napcat的最新release版本，根据[官网说明](https://napneko.github.io/guide/boot/Shell)进行部署。



### 创建napcat网络配置

需要创建两个napcat网络配置，分别用于HTTP客户端和HTTP服务端。

1. 服务端：

主机：0.0.0.0
端口：30000
启用CORS：否
启用WebSocket：否

2. 客户端：

上报URL：http://localhost:



## 修改代码后重新部署

如果修改了某个服务的代码，需要重新构建镜像，并更新k8s的deployment，例如：

```bash
docker build -t lagrange-onebot:local .
kubectl rollout restart deployment/lagrange-onebot
```