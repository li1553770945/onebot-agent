

# 部署命令

在使用部署命令之前，请确保所有构建操作已完成，尤其是onebot的登录操作完成，且测试无问题。

## 部署message-dispatch

```bash

cd ./message-dispatch
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
kubectl create configmap message-dispatch-config --from-file=production.yml=config/production.yml
```


## 部署napcat

注意：请确保在测试流程中已经登录，否则k8s容器启动后还需要手动登录。

如果你在Windows使用并且使用的是Docker Desktop的k8s集群，**不能使用Windows上的绝对路径！**必须使用特定的linux目录写法，不同Docker Desktop版本可能存在差异，可以使用`kubectl debug node/docker-desktop -it --image=busybox`命令启动busybox，想办法找到对应的路径写法。请务必确认路径正确，例如作者电脑上Windows路径为`E:\node-project\onebot-agent\napcat-docker`，对应的路径写法为：

```yaml
      volumes:
      - name: data-storage
        hostPath: /mnt/host/e/node-project/onebot-agent/napcat-docker
```

如果使用minikube，需要先使用mount命令将napcat数据目录挂载到minikube虚拟机中，这个命令需要在后台持续运行，例如：

```bash
cd ./napcat-docker
nohup minikube mount /home/liyaning/programming/onebot-agent/napcat-docker/data:/napcat-docker > nohup.log 2>&1 & echo $! > nohup.pid
```

这样会将进程 ID 保存到 `nohup.pid` 文件，方便后续管理。停止 mount 时可以使用：

```bash
cd ./napcat-docker
kill $(cat nohup.pid)
```

然后使用下面的命令启动：

```bash
cd ./napcat-docker
export NAPCAT_DATA_PATH=/napcat-docker  # 修改为你的napcat数据目录路径。这里使用了minikube，因此使用minikube mount的路径而不是真实宿主机路径
bash deploy.sh
```
如果一直处于pending状态，可能是因为PV未正确绑定PVC，请使用`kubectl get pv pvc`检查pv和pvc是否处于Bond状态而不是Pending，检查PV的hostPath路径是否正确。

完成后，请使用`kubectl get pods`查看napcat pod是否启动成功。如果成功，需要使用`kubectl logs <napcat-pod-name>`查看日志，确认napcat是否已经登录。如果未登录，请根据日志中的二维码信息进行登录。

### 创建napcat网络配置

需要创建两个napcat网络配置，分别用于HTTP客户端和HTTP服务端。

1. 服务端：

主机：0.0.0.0
端口：3000
启用CORS：否
启用WebSocket：否

2. 客户端：

主机：message-dispatch.message-dispatch.svc.cluster.local
端口：15001


## 修改代码后重新部署

如果修改了某个服务的代码，需要重新构建镜像，并更新k8s的deployment，例如：

```bash
docker build -t lagrange-onebot:local .
kubectl rollout restart deployment/lagrange-onebot
```