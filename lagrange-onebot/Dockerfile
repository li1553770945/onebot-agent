# 使用 Ubuntu 作为基础镜像
FROM ubuntu:22.04

# 安装必要的依赖
RUN apt-get update && apt-get install -y \
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 安装 .NET 8 运行时
RUN wget https://packages.microsoft.com/config/ubuntu/22.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb \
    && dpkg -i packages-microsoft-prod.deb \
    && rm packages-microsoft-prod.deb \
    && apt-get update \
    && apt-get install -y aspnetcore-runtime-8.0 \
    && rm -rf /var/lib/apt/lists/*

# 创建模板目录来存储原始文件
WORKDIR /template
COPY Lagrange.OneBot .
COPY appsettings.json ./appsettings.json

# 设置执行权限
RUN chmod +x Lagrange.OneBot

# 声明数据卷（这里将作为程序的运行目录）
VOLUME ["/app"]

EXPOSE 8000

# 启动脚本：将模板文件复制到 /app 并运行程序
CMD ["sh", "-c", "cp -a /template/* /app/ && cd /app && ./Lagrange.OneBot"]
