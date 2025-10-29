# MCP Server 模块

提供一个QQ群相关MCP服务端，包含几个工具，以及 HTTP `/ping` 健康检查端点。

## 工具列表

1. 接受入群申请
2. 拒绝入群申请
3. 修改群名片
4. 发送群消息

## 目录结构
```
mcp/
  main.ts            # 入口，注册工具与启动健康检查
  Dockerfile         # 多阶段构建镜像
  k8s-deployment.yml # Kubernetes Deployment
  k8s-service.yml    # Kubernetes Service
  tsconfig.json      # TypeScript 配置
  README.md          # 使用说明
```

## 脚本说明 (package.json)
| 脚本 | 作用 |
|------|------|
| dev | 使用 ts-node/ts-node-dev 启动开发（监听重载可自行替换） |
| build | 编译 TypeScript 到 `dist/` |
| start | 以编译后的 JS 启动生产运行 |
| typecheck | 仅做类型检查，不生成文件 |
| clean | 清理 dist 与缓存 |

## 本地开发
```bash
pnpm install
pnpm dev
# 另开终端健康检查
curl http://127.0.0.1:3000/ping   # => ok
```

## 构建产物
```bash
pnpm build
node dist/main.js
```

## Docker 构建与运行
```bash
# 构建镜像 (本地 tag)
docker build -t group-request-handler-mcp:local .
# 运行
docker run --rm -p 3000:3000 group-request-handler-mcp:local
curl http://127.0.0.1:3000/ping  # ok
```

## Kubernetes 部署（本地 kind / k3s）
```bash
# 使用本地镜像(minikube 需先导入或使用eval $(minikube docker-env) )
eval $(minikube docker-env)
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
```

如果是更改了代码重新运行，需要:

```bash
# 重新构建镜像
docker build -t group-request-handler-mcp:local .  
eval $(minikube docker-env)

kubectl rollout restart deployment/group-request-handler-mcp

# 或者先删除后重建
# 删除旧的 deployment
kubectl delete deployment group-request-handler-mcp
# 重新应用 deployment
kubectl apply -f k8s-deployment.yml