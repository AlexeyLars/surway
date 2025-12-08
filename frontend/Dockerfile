FROM node:20-alpine AS deps
WORKDIR /app
ENV NODE_ENV=development
COPY package*.json ./
RUN npm ci

FROM node:20-alpine AS builder
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# --- ВАЖНО: Объявляем ARG, чтобы Docker принял их при билде ---
ARG NEXT_PUBLIC_API_PROTOCOL
ARG NEXT_PUBLIC_API_HOST
ARG NEXT_PUBLIC_API_PORT
ARG NEXT_PUBLIC_API_VERSION
# --------------------------------------------------------------

# Эти переменные теперь будут доступны команде npm run build
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/.next/standalone ./
# Обрати внимание: для продакшена лучше использовать standalone режим (настраивается в next.config.js),
# но если у тебя стандартный билд, то копируй как было:
# COPY --from=builder /app/.next ./.next
# COPY package*.json ./

EXPOSE 3000

# Если standalone: CMD ["node", "server.js"]
# Если стандартный:
CMD ["npm", "run", "start"]