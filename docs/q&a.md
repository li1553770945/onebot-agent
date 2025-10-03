1. 使用了minikube，为什么docker build出来的镜像在minikube里找不到？

因为minikube有自己的docker环境，和宿主机的docker环境是隔离的。可以通过下面的命令让宿主机的docker命令行指向minikube的docker环境：

```bash
eval $(minikube docker-env)
```