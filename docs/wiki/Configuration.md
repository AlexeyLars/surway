# ⚙️ Configuration Guide

Полное руководство по конфигурации Surway проекта.

---

## 📋 Оглавление

1. [Backend конфигурация](#backend-конфигурация)
2. [Frontend конфигурация](#frontend-конфигурация)
3. [Redis конфигурация](#redis-конфигурация)
4. [Docker конфигурация](#docker-конфигурация)
5. [Production настройки](#production-настройки)
6. [Environment-specific](#environment-specific)

---

## 🔧 Backend конфигурация

### Environment Variables

Backend использует переменные окружения для конфигурации. Загрузка через [cleanenv](https://github.com/ilyakaznacheev/cleanenv).

**Файл:** `backend/internal/config/config.go`

### Server Configuration

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `SERVER_HOST` | string | `0.0.0.0` | Хост сервера |
| `SERVER_PORT` | int | `8080` | Порт сервера |
| `SERVER_READ_TIMEOUT` | duration | `10s` | Таймаут чтения запроса |
| `SERVER_WRITE_TIMEOUT` | duration | `10s` | Таймаут записи ответа |
| `SERVER_SHUTDOWN_TIMEOUT` | duration | `5s` | Таймаут graceful shutdown |
| `BASE_URL` | string | `http://localhost:8080` | Базовый URL для генерации ссылок |

**Примеры:**

```env
# Development
SERVER_HOST=127.0.0.1
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
BASE_URL=http://localhost:8080

# Production
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
BASE_URL=https://your-domain.com/api
```

### Redis Configuration

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `REDIS_HOST` | string | `localhost` | Хост Redis |
| `REDIS_PORT` | int | `6379` | Порт Redis |
| `REDIS_PASSWORD` | string | _(пусто)_ | Пароль Redis |
| `REDIS_DB` | int | `0` | Номер базы данных (0-15) |

**Примеры:**

```env
# Local development
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Docker
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_secure_password
REDIS_DB=0

# Remote Redis
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=very_secure_password_here
REDIS_DB=1
```

### Poll Configuration

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `POLL_DEFAULT_TTL` | duration | `168h` (7 дней) | TTL опроса по умолчанию |
| `POLL_MAX_TTL` | duration | `720h` (30 дней) | Максимальный TTL опроса |

**Duration format:**
- `h` — часы (hours)
- `m` — минуты (minutes)
- `s` — секунды (seconds)

**Примеры:**
```env
POLL_DEFAULT_TTL=168h    # 7 дней
POLL_DEFAULT_TTL=24h     # 1 день
POLL_DEFAULT_TTL=30m     # 30 минут
POLL_DEFAULT_TTL=3600s   # 1 час

POLL_MAX_TTL=720h        # 30 дней
```

### Environment Mode

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `ENV` | string | `dev` | Окружение: `dev`, `prod` |

**Эффект:**
- `dev` — Gin в debug mode, подробные логи
- `prod` — Gin в release mode, оптимизация

```env
ENV=dev   # Development
ENV=prod  # Production
```

### Полный пример .env для backend

```env
# Environment
ENV=dev

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_SHUTDOWN_TIMEOUT=5s
BASE_URL=http://localhost:8080

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Polls
POLL_DEFAULT_TTL=168h
POLL_MAX_TTL=720h
```

---

## 🎨 Frontend конфигурация

### Environment Variables

Frontend использует Next.js environment variables.

**Важно:**
- `NEXT_PUBLIC_*` — доступны в браузере (client-side)
- Остальные — только на сервере (server-side)

### Client-side Configuration

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `NEXT_PUBLIC_API_PROTOCOL` | string | `http` | Протокол API (http/https) |
| `NEXT_PUBLIC_API_HOST` | string | `localhost` | Хост API |
| `NEXT_PUBLIC_API_PORT` | string | `8080` | Порт API |
| `NEXT_PUBLIC_API_VERSION` | string | `v1` | Версия API |
| `NEXT_PUBLIC_FRONTEND_PROTOCOL` | string | `http` | Протокол frontend |
| `NEXT_PUBLIC_FRONTEND_HOST` | string | `localhost` | Хост frontend |
| `NEXT_PUBLIC_FRONTEND_PORT` | string | `3000` | Порт frontend |

**Использование:**

```typescript
// app/config/api.ts
export const API_CONFIG = {
  protocol: process.env.NEXT_PUBLIC_API_PROTOCOL || 'http',
  host: process.env.NEXT_PUBLIC_API_HOST || 'localhost',
  port: process.env.NEXT_PUBLIC_API_PORT || '8080',
  version: process.env.NEXT_PUBLIC_API_VERSION || 'v1',
};

export const API_URL = `${API_CONFIG.protocol}://${API_CONFIG.host}:${API_CONFIG.port}/api/${API_CONFIG.version}`;
```

### Server-side Configuration (SSR)

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `API_INTERNAL_PROTOCOL` | string | `http` | Протокол для SSR запросов |
| `API_INTERNAL_HOST` | string | `backend` | Хост backend (внутри Docker) |
| `API_INTERNAL_PORT` | string | `8080` | Порт backend |

**Зачем нужно:**
При SSR (Server-Side Rendering) Next.js запрашивает данные с сервера. Внутри Docker сети нужно использовать `backend:8080`, а не `localhost:8080`.

### Next.js Configuration

| Переменная | Тип | По умолчанию | Описание |
|-----------|-----|--------------|----------|
| `NODE_ENV` | string | `development` | Окружение Node.js |
| `HOSTNAME` | string | `0.0.0.0` | Хост для Next.js сервера |
| `PORT` | string | `3000` | Порт Next.js сервера |
| `NEXT_TELEMETRY_DISABLED` | string | `1` | Отключить телеметрию |

### Полный пример .env.local для frontend

```env
# API Configuration (client-side)
NEXT_PUBLIC_API_PROTOCOL=http
NEXT_PUBLIC_API_HOST=localhost
NEXT_PUBLIC_API_PORT=8080
NEXT_PUBLIC_API_VERSION=v1

# Frontend URLs (client-side)
NEXT_PUBLIC_FRONTEND_PROTOCOL=http
NEXT_PUBLIC_FRONTEND_HOST=localhost
NEXT_PUBLIC_FRONTEND_PORT=3000

# SSR Configuration (server-side)
API_INTERNAL_PROTOCOL=http
API_INTERNAL_HOST=localhost
API_INTERNAL_PORT=8080

# Next.js
NODE_ENV=development
HOSTNAME=0.0.0.0
PORT=3000
NEXT_TELEMETRY_DISABLED=1
```

### Docker Build Args

Для production builds некоторые переменные передаются как build arguments:

```dockerfile
# frontend/Dockerfile
ARG NEXT_PUBLIC_API_PROTOCOL
ARG NEXT_PUBLIC_API_HOST
ARG NEXT_PUBLIC_API_PORT
ARG NEXT_PUBLIC_API_VERSION
```

**В docker-compose.prod.yml:**

```yaml
frontend:
  build:
    context: ./frontend
    args:
      - NEXT_PUBLIC_API_PROTOCOL=https
      - NEXT_PUBLIC_API_HOST=your-domain.com
      - NEXT_PUBLIC_API_PORT=443
      - NEXT_PUBLIC_API_VERSION=v1
```

---

## 🗄️ Redis конфигурация

### Production Redis Config

Файл: `configs/redis.conf`

**Основные настройки:**

```conf
# ====================
# PERSISTENCE
# ====================

# RDB снэпшоты
save 900 1          # Сохранить если хотя бы 1 изменение за 15 минут
save 300 10         # Сохранить если хотя бы 10 изменений за 5 минут
save 60 10000       # Сохранить если хотя бы 10000 изменений за 1 минуту

# AOF (Append Only File)
appendonly yes
appendfsync everysec   # Синхронизация каждую секунду (компромисс)

# ====================
# MEMORY MANAGEMENT
# ====================

maxmemory 256mb
maxmemory-policy allkeys-lru  # Удалять least recently used ключи

# ====================
# SECURITY
# ====================

requirepass your_strong_password_here
protected-mode yes

# ====================
# NETWORKING
# ====================

bind 0.0.0.0  # Слушать на всех интерфейсах (в Docker)
port 6379
tcp-backlog 511
timeout 0
tcp-keepalive 300

# ====================
# LOGGING
# ====================

loglevel notice
logfile ""  # Stdout (для Docker logs)

# ====================
# SNAPSHOTTING
# ====================

stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data
```

### Memory Policies

| Policy | Описание |
|--------|----------|
| `noeviction` | Не удалять ключи, возвращать ошибку |
| `allkeys-lru` | Удалять LRU ключи из всех ключей |
| `volatile-lru` | Удалять LRU ключи только с TTL |
| `allkeys-random` | Удалять случайные ключи |
| `volatile-random` | Удалять случайные ключи с TTL |
| `volatile-ttl` | Удалять ключи с наименьшим TTL |

**Рекомендация для Surway:** `allkeys-lru` (все опросы имеют TTL, но LRU более предсказуем)

### Appendfsync Modes

| Mode | Описание | Производительность | Надежность |
|------|----------|-------------------|-----------|
| `always` | Синхронизация каждую операцию | Низкая | Максимальная |
| `everysec` | Синхронизация каждую секунду | Средняя | Хорошая |
| `no` | ОС решает когда синхронизировать | Высокая | Низкая |

**Рекомендация:** `everysec` (хороший баланс)

---

## 🐳 Docker конфигурация

### Development (docker-compose.yml)

```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"  # Доступен извне для debug
    command: redis-server --appendonly yes --appendfsync everysec

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - REDIS_HOST=redis  # Внутри Docker сети
      - REDIS_PORT=6379
    depends_on:
      redis:
        condition: service_healthy

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_HOST=localhost  # Доступ из браузера
      - API_INTERNAL_HOST=backend       # Доступ с SSR
```

### Production (docker-compose.prod.yml)

```yaml
services:
  redis:
    image: redis:7-alpine
    volumes:
      - redis-data:/data
      - ./configs/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    # НЕТ expose портов наружу!

  backend:
    expose:
      - "8080"  # Только внутри Docker сети
    # НЕТ ports!

  frontend:
    expose:
      - "3000"
    # НЕТ ports!

  caddy:
    ports:
      - "80:80"
      - "443:443"  # Только Caddy доступен извне
```

---

## 🚀 Production настройки

### Backend Production

```env
ENV=prod
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=10s
BASE_URL=https://your-domain.com/api

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=very_strong_password_min_32_chars_recommended
REDIS_DB=0

POLL_DEFAULT_TTL=168h
POLL_MAX_TTL=720h
```

### Frontend Production

```env
# Build args (в docker-compose.prod.yml)
NEXT_PUBLIC_API_PROTOCOL=https
NEXT_PUBLIC_API_HOST=your-domain.com
NEXT_PUBLIC_API_PORT=443
NEXT_PUBLIC_API_VERSION=v1

# Runtime
NODE_ENV=production
API_INTERNAL_PROTOCOL=http
API_INTERNAL_HOST=backend
API_INTERNAL_PORT=8080
NEXT_TELEMETRY_DISABLED=1
```

### Redis Production

```conf
maxmemory 512mb
maxmemory-policy allkeys-lru
requirepass very_strong_password_min_32_chars
appendonly yes
appendfsync everysec
save 900 1
save 300 10
save 60 10000
```

---

## 🌍 Environment-specific

### Development

**Цели:**
- Быстрая разработка
- Подробные логи
- Hot reload
- Доступ ко всем портам для debug

**Backend .env:**
```env
ENV=dev
SERVER_HOST=127.0.0.1
SERVER_PORT=8080
BASE_URL=http://localhost:8080
REDIS_HOST=localhost
REDIS_PASSWORD=
```

**Frontend .env.local:**
```env
NEXT_PUBLIC_API_HOST=localhost
NEXT_PUBLIC_API_PORT=8080
```

### Staging

**Цели:**
- Максимально близко к production
- Тестирование перед релизом
- Доступ для QA команды

**Настройки:** почти как production, но может быть:
- Меньше ресурсов (RAM, CPU)
- Тестовые домены (staging.example.com)
- Отключен rate limiting для удобства тестов

### Production

**Цели:**
- Максимальная производительность
- Безопасность
- Надежность
- Мониторинг

**Обязательно:**
- [ ] Сильные пароли Redis
- [ ] HTTPS only
- [ ] Graceful shutdown
- [ ] Health checks
- [ ] Логирование
- [ ] Backup Redis данных
- [ ] Мониторинг метрик

---

## 🔐 Secrets Management

### Не коммитьте секреты!

**Добавьте в .gitignore:**
```
.env
.env.local
.env.production
*.secret
```

### Для production

Используйте:
- **Docker Secrets** — для Docker Swarm
- **Kubernetes Secrets** — для K8s
- **HashiCorp Vault** — enterprise solution
- **AWS Secrets Manager** — для AWS
- **Environment variables** — через CI/CD

**Пример Docker Secrets:**

```yaml
services:
  backend:
    secrets:
      - redis_password
    environment:
      - REDIS_PASSWORD_FILE=/run/secrets/redis_password

secrets:
  redis_password:
    file: ./secrets/redis_password.txt
```

---

## 🧪 Testing Configurations

### Test Environment

```env
ENV=test
SERVER_PORT=8081
REDIS_HOST=localhost
REDIS_PORT=6380
REDIS_DB=15  # Отдельная DB для тестов
POLL_DEFAULT_TTL=1m  # Короткий TTL для тестов
```

---

## 📞 Troubleshooting

### Backend не подключается к Redis

```bash
# Проверьте переменные
echo $REDIS_HOST
echo $REDIS_PORT

# Проверьте доступность Redis
redis-cli -h $REDIS_HOST -p $REDIS_PORT ping

# С паролем
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD ping
```

### Frontend не может достучаться до API

```bash
# Проверьте URL
echo $NEXT_PUBLIC_API_HOST
echo $NEXT_PUBLIC_API_PORT

# Тест из браузера (DevTools Console)
fetch('http://localhost:8080/health')
  .then(r => r.json())
  .then(console.log)
```

### Docker Compose переменные не работают

```bash
# Проверьте что .env файл в корне проекта
ls -la .env

# Проверьте синтаксис
cat .env

# Используйте явный --env-file
docker compose --env-file .env up
```

---

## 📚 Дополнительное чтение

- [cleanenv Documentation](https://github.com/ilyakaznacheev/cleanenv)
- [Next.js Environment Variables](https://nextjs.org/docs/basic-features/environment-variables)
- [Redis Configuration](https://redis.io/topics/config)
- [Docker Compose Environment](https://docs.docker.com/compose/environment-variables/)

---

_Последнее обновление: December 2025_
