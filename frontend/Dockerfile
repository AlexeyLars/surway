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


ARG NEXT_PUBLIC_API_PROTOCOL
ARG NEXT_PUBLIC_API_HOST
ARG NEXT_PUBLIC_API_PORT
ARG NEXT_PUBLIC_API_VERSION

ARG NEXT_PUBLIC_FRONTEND_PROTOCOL
ARG NEXT_PUBLIC_FRONTEND_HOST
ARG NEXT_PUBLIC_FRONTEND_PORT

RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/.next/standalone ./

EXPOSE 3000

CMD ["node", "server.js"]