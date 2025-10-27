# Multi-stage build for NestJS service
# 1. Builder stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install deps first (better layer caching)
COPY package.json yarn.lock ./
RUN corepack enable && yarn install --frozen-lockfile

# Copy source
COPY nest-cli.json tsconfig*.json ./
COPY src ./src

# Build
RUN yarn build

# 2. Production runtime stage
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Copy only necessary files
COPY package.json yarn.lock ./
RUN corepack enable && yarn install --frozen-lockfile --production

# Copy built dist from builder
COPY --from=builder /app/dist ./dist

# Non-root user (optional)
# RUN addgroup -S app && adduser -S app -G app
# USER app

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD node -e "fetch('http://localhost:3000/ping').then(r=>{if(!r.ok)process.exit(1)}).catch(()=>process.exit(1))" || exit 1

CMD ["node", "dist/main.js"]
