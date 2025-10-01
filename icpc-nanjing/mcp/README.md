# MCP Server 模块

提供一个示例 Model Context Protocol (MCP) 服务端，包含一个 `summarize` 工具（调用 MCP LLM sampling）以及 HTTP `/ping` 健康检查端点。

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
docker build -t mcp-server:local .
# 运行
docker run --rm -p 3000:3000 mcp-server:local
curl http://127.0.0.1:3000/ping  # ok
```

## 推送到镜像仓库（示例）
```bash
# 修改 tag
export IMAGE=my.registry.local/mcp-server:1.0.0
# 或 Windows PowerShell: $env:IMAGE="my.registry.local/mcp-server:1.0.0"
docker tag mcp-server:local %IMAGE%
docker push %IMAGE%
```

## Kubernetes 部署（本地 kind / k3s）
```bash
# 使用本地镜像（kind 需先导入）
# kind load docker-image mcp-server:local
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml

kubectl get pods -l app=mcp-server
kubectl logs -f deploy/mcp-server
kubectl port-forward svc/mcp-server 3000:3000
curl http://127.0.0.1:3000/ping
```

## 健康检查
- Readiness: GET /ping -> 200 ok
- Liveness:  GET /ping -> 200 ok

## 环境变量
| 名称 | 默认 | 说明 |
|------|------|------|
| PORT | 3000 | 健康检查 HTTP 端口 |
| NODE_ENV | development/production | 运行模式 |

## 后续可扩展建议
- 增加日志库 (pino / winston)
- 增加 OpenTelemetry 指标与 trace
- 增加更多 MCP Tools（如搜索、向量检索）
- 集成 CI (lint + typecheck + build + docker build)

---
如需进一步定制，继续提需求即可。
