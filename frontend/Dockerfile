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

RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY package*.json ./

RUN npm ci --omit=dev

EXPOSE 3000

CMD ["npm", "run", "start"]
