# 🚀 Deployment Guide

Полное руководство по деплою Surway в production окружение.

---

## 📋 Оглавление

1. [Требования](#требования)
2. [Docker Compose Production](#docker-compose-production)
3. [Caddy Setup (SSL)](#caddy-setup)
4. [Environment Variables](#environment-variables)
5. [Database (Redis)](#redis-configuration)
6. [Мониторинг](#мониторинг)
7. [Troubleshooting](#troubleshooting)

---

## 🎯 Требования

### Системные требования

**Минимальные:**
- CPU: 1 vCore
- RAM: 1 GB
- Disk: 10 GB SSD
- OS: Linux (Ubuntu 22.04 LTS рекомендуется)

**Рекомендуемые:**
- CPU: 2+ vCores
- RAM: 2+ GB
- Disk: 20+ GB SSD
- OS: Ubuntu 22.04 LTS

### Софт

- Docker 20.10+
- Docker Compose v2+
- (Опционально) Make для использования Makefile

### Сеть

- Открытые порты: 80 (HTTP), 443 (HTTPS)
- Домен с A-записью, указывающей на ваш сервер

---

## 🐳 Docker Compose Production

### 1. Клонирование репозитория

```bash
# SSH
git clone git@github.com:AlexeyLars/surway.git

# HTTPS
git clone https://github.com/AlexeyLars/surway.git

cd surway
git checkout main  # или develop для latest
```

### 2. Настройка переменных окружения

Создайте `.env` файл в корне проекта:

```bash
nano .env
```

**Минимальная конфигурация:**
```env
# Domain
DOMAIN=your-domain.com

# Environment
ENV=prod

# Backend
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
BASE_URL=https://your-domain.com/api

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_strong_redis_password_here
REDIS_DB=0

# Polls
POLL_DEFAULT_TTL=168h
POLL_MAX_TTL=720h
```

### 3. Настройка Caddyfile

Отредактируйте `Caddyfile`:

```bash
nano Caddyfile
```

**Пример конфигурации:**
```
your-domain.com {
    # Логирование
    log {
        output file /var/log/caddy/access.log
        format json
    }

    # Frontend (Next.js)
    reverse_proxy frontend:3000

    # Backend API
    handle /api/* {
        reverse_proxy backend:8080
    }

    # Swagger (опционально отключить в production)
    handle /swagger/* {
        reverse_proxy backend:8080
    }

    # Health check
    handle /health {
        reverse_proxy backend:8080
    }

    # Security headers
    header {
        # Enable HSTS
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        # Prevent clickjacking
        X-Frame-Options "SAMEORIGIN"
        # XSS protection
        X-Content-Type-Options "nosniff"
        X-XSS-Protection "1; mode=block"
        # CSP (настройте под свои нужды)
        Content-Security-Policy "default-src 'self'"
    }

    # Gzip compression
    encode gzip

    # SSL автоматически через Let's Encrypt
}

# Redirect www to non-www (опционально)
www.your-domain.com {
    redir https://your-domain.com{uri} permanent
}
```

### 4. Настройка docker-compose.prod.yml

Отредактируйте `docker-compose.prod.yml`:

```yaml
services:
  redis:
    image: redis:7-alpine
    container_name: poll-redis
    volumes:
      - redis-data:/data
      - ./configs/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf --requirepass ${REDIS_PASSWORD}
    healthcheck:
      test: ["CMD", "redis-cli", "--pass", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3
    networks:
      - app-net
    restart: always

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    expose:
      - "8080"
    environment:
      - ENV=${ENV}
      - SERVER_HOST=${SERVER_HOST}
      - SERVER_PORT=${SERVER_PORT}
      - BASE_URL=${BASE_URL}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - REDIS_DB=${REDIS_DB}
      - POLL_DEFAULT_TTL=${POLL_DEFAULT_TTL}
    depends_on:
      redis:
        condition: service_healthy
    networks:
      - app-net
    restart: always

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        - NEXT_PUBLIC_API_PROTOCOL=https
        - NEXT_PUBLIC_API_HOST=${DOMAIN}
        - NEXT_PUBLIC_API_PORT=443
        - NEXT_PUBLIC_API_VERSION=v1
    expose:
      - "3000"
    environment:
      - NODE_ENV=production
      - HOSTNAME=0.0.0.0
      - API_INTERNAL_PROTOCOL=http
      - API_INTERNAL_HOST=backend
      - API_INTERNAL_PORT=8080
    depends_on:
      backend:
        condition: service_started
    networks:
      - app-net
    restart: always

  caddy:
    image: caddy:alpine
    container_name: caddy-ingress
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
      - caddy_logs:/var/log/caddy
    networks:
      - app-net
    depends_on:
      - frontend
      - backend

volumes:
  redis-data:
  caddy_data:
  caddy_config:
  caddy_logs:

networks:
  app-net:
    driver: bridge
```

### 5. Запуск

```bash
# Сборка образов
docker compose -f docker-compose.prod.yml build

# Запуск в фоновом режиме
docker compose -f docker-compose.prod.yml up -d

# Проверка логов
docker compose -f docker-compose.prod.yml logs -f
```

### 6. Проверка

```bash
# Health check
curl https://your-domain.com/health

# Проверка API
curl https://your-domain.com/api/v1/polls \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","options":["A","B"]}'
```

---

## 🔒 Caddy Setup (SSL)

### Автоматический SSL от Let's Encrypt

Caddy автоматически получает и обновляет SSL сертификаты!

**Требования:**
1. Домен должен резолвиться в IP вашего сервера
2. Порты 80 и 443 должны быть открыты
3. Caddy должен иметь права на запись в `/data` volume

### Проверка SSL

```bash
# Проверка сертификата
openssl s_client -connect your-domain.com:443 -servername your-domain.com

# Проверка HTTPS
curl -I https://your-domain.com
```

### Custom SSL сертификаты

Если используете свои сертификаты:

```
your-domain.com {
    tls /path/to/cert.pem /path/to/key.pem
    # ... остальная конфигурация
}
```

---

## ⚙️ Environment Variables

### Backend

| Переменная | Обязательна | Пример | Описание |
|-----------|-------------|--------|----------|
| `ENV` | Да | `prod` | Окружение (dev/prod) |
| `SERVER_HOST` | Да | `0.0.0.0` | Хост сервера |
| `SERVER_PORT` | Да | `8080` | Порт сервера |
| `BASE_URL` | Да | `https://domain.com/api` | Базовый URL для ссылок |
| `REDIS_HOST` | Да | `redis` | Хост Redis |
| `REDIS_PORT` | Да | `6379` | Порт Redis |
| `REDIS_PASSWORD` | Нет | `secret` | Пароль Redis |
| `REDIS_DB` | Нет | `0` | Номер БД Redis |
| `POLL_DEFAULT_TTL` | Нет | `168h` | TTL опроса по умолчанию |
| `POLL_MAX_TTL` | Нет | `720h` | Максимальный TTL |

### Frontend

**Build-time (ARG):**
- `NEXT_PUBLIC_API_PROTOCOL` — https
- `NEXT_PUBLIC_API_HOST` — your-domain.com
- `NEXT_PUBLIC_API_PORT` — 443
- `NEXT_PUBLIC_API_VERSION` — v1

**Runtime:**
- `NODE_ENV` — production
- `API_INTERNAL_PROTOCOL` — http
- `API_INTERNAL_HOST` — backend
- `API_INTERNAL_PORT` — 8080

---

## 🗄️ Redis Configuration

### Production настройки

Создайте `configs/redis.conf`:

```conf
# Persistence
save 900 1
save 300 10
save 60 10000
appendonly yes
appendfsync everysec

# Memory management
maxmemory 256mb
maxmemory-policy allkeys-lru

# Security
requirepass your_strong_redis_password_here
protected-mode yes

# Performance
tcp-backlog 511
timeout 0
tcp-keepalive 300

# Logging
loglevel notice
logfile ""

# Снэпшоты
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data
```

### Backup Redis данных

```bash
# Создание backup
docker exec poll-redis redis-cli --pass your_password BGSAVE

# Копирование dump.rdb
docker cp poll-redis:/data/dump.rdb ./backup-$(date +%Y%m%d).rdb

# Восстановление
docker cp backup-20251208.rdb poll-redis:/data/dump.rdb
docker restart poll-redis
```

---

## 📊 Мониторинг

### Логи

**Просмотр логов всех сервисов:**
```bash
docker compose -f docker-compose.prod.yml logs -f
```

**Логи отдельного сервиса:**
```bash
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f redis
docker compose -f docker-compose.prod.yml logs -f caddy
```

**Caddy access logs:**
```bash
docker exec caddy-ingress tail -f /var/log/caddy/access.log
```

### Health Checks

```bash
# Backend health
curl https://your-domain.com/health

# Redis health
docker exec poll-redis redis-cli --pass your_password ping

# Проверка всех контейнеров
docker ps
```

### Метрики (планируется)

В планах интеграция с:
- **Prometheus** — сбор метрик
- **Grafana** — визуализация
- **AlertManager** — алерты

---

## 🔧 Обновление

### Rolling update

```bash
# 1. Pull latest code
cd /path/to/surway
git pull origin main

# 2. Rebuild images
docker compose -f docker-compose.prod.yml build

# 3. Recreate containers
docker compose -f docker-compose.prod.yml up -d

# 4. Проверка логов
docker compose -f docker-compose.prod.yml logs -f
```

### Zero-downtime deployment

Для zero-downtime нужно:
1. Использовать load balancer
2. Запускать несколько инстансов backend
3. Обновлять поочередно

**Пример с 2 backend инстансами:**
```yaml
backend:
  # ... конфигурация
  deploy:
    replicas: 2
```

---

## 🔥 Troubleshooting

### Backend не запускается

**Проблема:** `failed to connect to redis`

**Решение:**
```bash
# Проверьте Redis
docker compose -f docker-compose.prod.yml logs redis

# Проверьте пароль
docker exec poll-redis redis-cli --pass your_password ping

# Проверьте сеть
docker network inspect surway_app-net
```

### Frontend не загружается

**Проблема:** 502 Bad Gateway

**Решение:**
```bash
# Проверьте логи frontend
docker compose -f docker-compose.prod.yml logs frontend

# Перезапустите frontend
docker compose -f docker-compose.prod.yml restart frontend
```

### SSL не работает

**Проблема:** Caddy не может получить сертификат

**Причины:**
1. Домен не резолвится в IP сервера
2. Порты 80/443 закрыты
3. Firewall блокирует запросы

**Решение:**
```bash
# Проверьте DNS
dig your-domain.com

# Проверьте порты
sudo netstat -tulpn | grep -E ':(80|443)'

# Проверьте логи Caddy
docker compose -f docker-compose.prod.yml logs caddy
```

### Redis переполнен

**Проблема:** `OOM command not allowed`

**Решение:**
```bash
# Увеличьте maxmemory в redis.conf
maxmemory 512mb

# Или используйте eviction policy
maxmemory-policy allkeys-lru

# Перезапустите Redis
docker compose -f docker-compose.prod.yml restart redis
```

### Низкая производительность

**Диагностика:**
```bash
# CPU и память контейнеров
docker stats

# Redis статистика
docker exec poll-redis redis-cli --pass your_password INFO stats

# Количество ключей
docker exec poll-redis redis-cli --pass your_password DBSIZE
```

---

## 🔐 Security Checklist

- [ ] Сильный пароль для Redis
- [ ] Firewall настроен (открыты только 80, 443, 22)
- [ ] SSH доступ только по ключам
- [ ] Регулярные обновления системы
- [ ] Backup Redis данных
- [ ] Мониторинг логов
- [ ] Rate limiting (в планах)
- [ ] HTTPS для всех запросов
- [ ] Security headers в Caddy
- [ ] Отключен Swagger в production (опционально)

---

## 📦 Backup Strategy

### Автоматический backup

Создайте cron job:

```bash
# Откройте crontab
crontab -e

# Добавьте строку (backup каждый день в 3 AM)
0 3 * * * /usr/bin/docker exec poll-redis redis-cli --pass your_password BGSAVE && /usr/bin/docker cp poll-redis:/data/dump.rdb /backups/redis-$(date +\%Y\%m\%d).rdb
```

### Хранение backups

Рекомендуется:
1. Локальное хранение (7 дней)
2. Offsite backup (S3, Google Cloud Storage)
3. Регулярная проверка восстановления

---

## 🌐 CDN и масштабирование

### Использование CDN

Для статических файлов frontend:
- Cloudflare
- AWS CloudFront
- Fastly

### Horizontal scaling

**Backend:**
```yaml
backend:
  deploy:
    replicas: 3
```

**Redis:**
- Redis Sentinel для HA
- Redis Cluster для шардинга

---

## 📞 Support

- **Issues:** https://github.com/AlexeyLars/surway/issues
- **Wiki:** https://github.com/AlexeyLars/surway/wiki
- **Author:** [@AlexeyLars](https://github.com/AlexeyLars)

---

_Последнее обновление: December 2025_
