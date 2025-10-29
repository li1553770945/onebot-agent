#!/bin/bash
# 获取数据路径，默认为脚本所在目录
# 使用方式：
#   ./deploy.sh                    # 使用默认路径 (当前目录)
#   NAPCAT_DATA_PATH=/data ./deploy.sh  # 指定自定义绝对路径
#   NAPCAT_DATA_PATH=/napcat-docker ./deploy.sh  # Minikube 挂载点

export NAPCAT_DATA_PATH=${NAPCAT_DATA_PATH:-/napcat-docker}

echo "Deploying with:"
echo "  NAPCAT_DATA_PATH=$NAPCAT_DATA_PATH"
echo ""

# 使用 envsubst 替换环境变量并应用
export NAPCAT_DATA_PATH

echo "Applying PV/PVC..."
envsubst < k8s-pv-pvc.yml | kubectl apply -f -

echo "Applying Deployment..."
envsubst < k8s-deployment.yml | kubectl apply -f -

echo "Applying Service..."
kubectl apply -f k8s-service.yml

echo ""
echo "✓ Deployment complete!"
echo ""
echo "Check status with:"
echo "  kubectl get pvc"
echo "  kubectl get pod -w"
echo "  kubectl logs -f deployment/napcat"