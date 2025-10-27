# syntax=docker/dockerfile:1.7

# -------- Base & dependency install (with caching) --------
FROM node:20-alpine AS deps
WORKDIR /app
ENV NODE_ENV=development
RUN corepack enable
# 只拷贝依赖声明以最大化缓存
COPY package.json yarn.lock* ./
RUN yarn install --frozen-lockfile

# -------- Build stage --------
FROM node:20-alpine AS build
WORKDIR /app
RUN corepack enable
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN yarn build

# -------- Production runtime --------
FROM node:20-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
RUN adduser -D app && chown app:app /app
RUN corepack enable
COPY package.json yarn.lock* ./
# 仅安装生产依赖：利用 --production (Yarn 1) / 或设置 NODE_ENV=production
RUN yarn install --frozen-lockfile --production --ignore-scripts && yarn cache clean
COPY --from=build /app/dist ./dist
USER app
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD wget -qO- http://127.0.0.1:3000/ping || exit 1
CMD ["node", "dist/main.js"]
