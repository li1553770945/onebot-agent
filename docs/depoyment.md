

# 部署命令

在使用部署命令之前，请确保所有构建操作已完成，尤其是onebot的登录操作完成，且测试无问题。

## 部署message-dispatch

```bash

cd ./message-dispatch
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
kubectl create configmap message-dispatch-config --from-file=production.yml=config/production.yml
```


## 部署lagrange-onebot

注意：请确保`lagrange-onebot/app`目录存在，并且在测试流程中已经登录，生成了对应的device.json、keystore.json等文件，否则k8s容器启动后还需要手动登录。

由于k8s限制无法使用相对路径，**复制k8s-pv.template.yml为k8s-pv.yml，并且请修改`lagrange-onebot\k8s-pv.yml`中的hostPath路径为你的绝对路径**。如果你在Windows使用并且使用的是Docker Desktop的k8s集群，**不能使用Windows上的绝对路径！**必须使用特定的linux目录写法，不同Docker Desktop版本可能存在差异，可以使用`kubectl debug node/docker-desktop -it --image=busybox`命令启动busybox，想办法找到对应的路径写法。请务必确认路径正确，例如作者电脑上Windows路径为`E:\node-project\onebot-agent\lagrange-onebot\app`，对应的路径写法为：

```yaml
      volumes:
      - name: data-storage
        hostPath: /mnt/host/e/node-project/onebot-agent/lagrange-onebot/app
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