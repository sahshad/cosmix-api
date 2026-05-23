# ---------- Build Stage ----------
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

RUN npm install -g pnpm
RUN pnpm install

COPY tsconfig*.json ./
COPY nest-cli.json ./
COPY src ./src

RUN pnpm build

# ---------- Runtime Stage ----------
FROM node:20-alpine

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

RUN npm install -g pnpm
RUN pnpm install --prod

COPY --from=builder /app/dist ./dist

EXPOSE 3000

CMD ["node", "dist/main.js"]