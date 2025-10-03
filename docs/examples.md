
# 示例：ICPC 南京站入群助手部署指南

## 构建和登录icpc-nanjing-mcp

```bash
cd icpc-nanjing/mcp
## 构建
docker build -t icpc-nanjing-mcp:local .
## 运行测试
docker run -d  --name icpc-nanjing-mcp --network onebot-agent -p 3000:3000  -v ./config:/app/config icpc-nanjing-mcp:local
## 部署
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml
kubectl apply -f k8s-secret.yml


## 更新后重启
kubectl rollout restart deployment icpc-nanjing-mcp
kubectl rollout status deployment icpc-nanjing-mcp  
```

## 构建和登录icpc-nanjing-agent

```bash
cd icpc-nanjing/agent
## 构建
docker build -t icpc-nanjing-agent:local .
## 运行测试
docker run -d  --name icpc-nanjing-agent --network onebot-agent -p 3000:3000  -v ./config:/app/config icpc-nanjing-agent:local 
## 部署
kubectl apply -f k8s-secret.yml
kubectl apply -f k8s-deployment.yml
kubectl apply -f k8s-service.yml

## 更新后重启
kubectl rollout restart deployment icpc-nanjing-agent
kubectl rollout status deployment icpc-nanjing-agent
```

