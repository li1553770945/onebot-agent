# syntax=docker/dockerfile:1.7

# ========== Base build image ==========
FROM node:20-alpine AS base
WORKDIR /app
ENV PNPM_HOME=/pnpm
ENV PATH=$PNPM_HOME:$PATH
RUN corepack enable
COPY package.json pnpm-lock.yaml* package-lock.json* yarn.lock* ./

# ========== Dependencies layer ==========
FROM base AS deps
# 使用 pnpm 优先，其次 npm/yarn
RUN if [ -f pnpm-lock.yaml ]; then pnpm install --frozen-lockfile; \
    elif [ -f package-lock.json ]; then npm ci; \
    elif [ -f yarn.lock ]; then yarn install --frozen-lockfile; \
    else npm install; fi

# ========== Build layer ==========
FROM deps AS build
COPY tsconfig.json ./
COPY . .
RUN if [ -f pnpm-lock.yaml ]; then pnpm run build; \
    elif [ -f package-lock.json ]; then npm run build; \
    elif [ -f yarn.lock ]; then yarn build; \
    else npm run build; fi

# ========== Production runtime ==========
FROM node:20-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production
# 仅拷贝生产依赖
COPY package.json pnpm-lock.yaml* package-lock.json* yarn.lock* ./
RUN corepack enable && \
    if [ -f pnpm-lock.yaml ]; then pnpm install --prod --frozen-lockfile; \
    elif [ -f package-lock.json ]; then npm ci --omit=dev; \
    elif [ -f yarn.lock ]; then yarn install --production --frozen-lockfile; \
    else npm install --omit=dev; fi

# 拷贝构建产物
COPY --from=build /app/dist ./dist

EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD wget -qO- http://127.0.0.1:3000/ping || exit 1

CMD ["node", "dist/main.js"]
