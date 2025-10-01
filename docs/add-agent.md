# 增加一个服务或agent

假设你想增加一个新的服务或agent，步骤如下：

## 修改message-dispatch配置

1. 新建服务

新建一个Agent服务，你可以将test-agent拷贝一份，改个名字，比如icpc-nanjing-agent，然后修改里面的代码实现你想要的功能，打包成镜像并推送到你的镜像仓库。


2. 配置规则

在 `message-dispatch/config/production.yml` 中增加 / 修改一条规则，指定消息的来源和去向。


具体配置规则请参考[message-dispatch的README](../message-dispatch/README.md#路由配置规则)。

3. 更新 ConfigMap

使用 k8s 命令更新 ConfigMap（声明式覆盖）：

```bash
kubectl create configmap message-dispatch-config --from-file=production.yml=config/production.yml -o yaml --dry-run=client | kubectl apply -f -
```

  说明：如果这是第一次用 `kubectl apply` 管理该 ConfigMap，可能看到一个警告：

  > resource configmaps/message-dispatch-config is missing the kubectl.kubernetes.io/last-applied-configuration annotation ...

  这是正常的，表示之前不是用 `apply` 创建的；本次会自动补上 annotation，之后不会再警告。

4. 验证 ConfigMap 是否已更新

```bash
kubectl get configmap message-dispatch-config -o yaml
```

应该能看到配置文件已经被更新。


5. 让运行中的 Pod 生效

目前的设计更新 ConfigMap 本身不会自动重启 Pod。需要手动重启：

```bash
kubectl rollout restart deployment message-dispatch
kubectl rollout status deployment message-dispatch
```
