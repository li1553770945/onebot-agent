
# 示例：入群审核Agent部署

## 构建和登录group-request-handler-mcp

```bash
cd group_request/mcp
## 构建
docker build -t group-request-handler-mcp:local .
## 运行测试
docker run -d  --name group-request-handler-mcp --network onebot-agent -p 3000:3000  -v ./config:/app/config group-request-handler-mcp:local
## 部署
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml


## 更新后重启
kubectl rollout restart deployment group-request-handler-mcp
kubectl rollout status deployment group-request-handler-mcp  
```

## 构建和登录group-request-handler-agent

```bash
cd group_request/agent
## 构建
docker build -t group-request-handler-agent:local .
## 运行测试
docker run -d  --name group-request-handler-agent --network onebot-agent -p 3000:3000  -v ./config:/app/config group-request-handler-agent:local 
## 部署
kubectl apply -f k8s-secret.yml
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml

## 更新后重启
kubectl rollout restart deployment group-request-handler-agent
kubectl rollout status deployment group-request-handler-agent
```

