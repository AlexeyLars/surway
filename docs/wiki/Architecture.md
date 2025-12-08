# 🏗️ Архитектура Surway

Детальное описание архитектуры проекта, паттернов проектирования и взаимодействия компонентов.

---

## 📐 Общая архитектура

Surway следует принципам **Clean Architecture** и **Domain-Driven Design**, разделяя систему на независимые слои с четкими границами.

```
┌─────────────────────────────────────────────────────────────────┐
│                          Users / Clients                         │
└────────────────────┬─────────────────┬──────────────────────────┘
                     │                 │
        ┌────────────▼──────────┐     │
        │   Frontend (Next.js)  │     │
        │   - React Components  │     │
        │   - API Client        │     │
        │   - SSR/CSR           │     │
        └────────────┬──────────┘     │
                     │                │
                     │ HTTP/REST      │ HTTP/REST
                     │                │
        ┌────────────▼────────────────▼──────────────────┐
        │         Backend (Go + Gin)                     │
        │  ┌──────────────────────────────────────────┐  │
        │  │  Handler Layer (HTTP)                    │  │
        │  │  - Request validation                    │  │
        │  │  - Response formatting                   │  │
        │  │  - Error mapping                         │  │
        │  └──────────────┬───────────────────────────┘  │
        │                 │                               │
        │  ┌──────────────▼───────────────────────────┐  │
        │  │  Service Layer (Business Logic)         │  │
        │  │  - Poll creation logic                  │  │
        │  │  - ID generation                        │  │
        │  │  - URL building                         │  │
        │  └──────────────┬───────────────────────────┘  │
        │                 │                               │
        │  ┌──────────────▼───────────────────────────┐  │
        │  │  Storage Interface                       │  │
        │  └──────────────┬───────────────────────────┘  │
        │                 │                               │
        │  ┌──────────────▼───────────────────────────┐  │
        │  │  Redis Storage Implementation            │  │
        │  │  - CRUD operations                       │  │
        │  │  - Atomic transactions                   │  │
        │  └──────────────┬───────────────────────────┘  │
        └─────────────────┼───────────────────────────────┘
                          │
        ┌─────────────────▼───────────────────────────────┐
        │              Redis Database                     │
        │  - String: poll metadata (JSON + TTL)           │
        │  - Hash: vote counters                          │
        └─────────────────────────────────────────────────┘
```

---

## 🎯 Backend архитектура (Go)

### Слоистая архитектура

#### 1. Handler Layer (`internal/handler/`)

**Ответственность:**
- Обработка HTTP запросов
- Валидация входных данных (через Gin binding)
- Маппинг бизнес-ошибок в HTTP статусы
- Формирование JSON ответов

**Компоненты:**
- `poll.go` — CRUD handlers для опросов
- `router.go` — настройка маршрутов и middleware

**Middleware:**
- `LoggerMiddleware` — структурированное логирование запросов
- `CORSMiddleware` — CORS headers для frontend
- `gin.Recovery()` — восстановление после panic

**Пример handler:**
```go
func (h *PollHandler) Vote(c *gin.Context) {
    pollID := c.Param("id")

    var req model.VoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Валидация
        c.JSON(400, model.ErrorResponse{...})
        return
    }

    // Делегирование в service layer
    err := h.service.Vote(c.Request.Context(), pollID, &req)

    // Маппинг ошибок
    if errors.Is(err, storage.ErrPollNotFound) {
        c.JSON(404, ...)
        return
    }

    c.JSON(200, model.VoteResponse{Success: true})
}
```

#### 2. Service Layer (`internal/service/`)

**Ответственность:**
- Бизнес-логика приложения
- Генерация уникальных ID для опросов
- Построение URL для голосования/результатов
- Логирование бизнес-событий
- Валидация бизнес-правил

**Изоляция:**
- Не знает о HTTP (использует `context.Context`)
- Работает через Storage интерфейс (DI)
- Возвращает бизнес-ошибки (не HTTP)

**Пример:**
```go
func (s *PollService) CreatePoll(ctx context.Context, req *model.CreatePollRequest) (*model.CreatePollResponse, error) {
    // Генерация короткого ID
    pollID := random.NewRandomString(7)

    // Создание Poll entity
    poll := &model.Poll{
        ID:        pollID,
        Title:     req.Title,
        Options:   req.Options,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(s.config.Poll.DefaultTTL),
    }

    // Сохранение через интерфейс
    if err := s.storage.CreatePoll(ctx, poll, s.config.Poll.DefaultTTL); err != nil {
        s.logger.Error("failed to create poll", ...)
        return nil, err
    }

    // Построение response с URL
    return &model.CreatePollResponse{
        PollID:     pollID,
        VoteURL:    fmt.Sprintf("%s/api/v1/polls/%s/vote", baseURL, pollID),
        ResultsURL: fmt.Sprintf("%s/api/v1/polls/%s/results", baseURL, pollID),
    }, nil
}
```

#### 3. Storage Layer (`internal/storage/`)

**Ответственность:**
- CRUD операции с данными
- Атомарные транзакции (Redis Pipeline)
- Управление TTL
- Работа с Redis структурами данных

**Storage Interface:**
```go
type Storage interface {
    CreatePoll(ctx context.Context, poll *model.Poll, ttl time.Duration) error
    GetPoll(ctx context.Context, pollID string) (*model.Poll, error)
    Vote(ctx context.Context, pollID string, optionIndices []int) error
    GetResults(ctx context.Context, pollID string) (*model.PollResults, error)
    Close() error
}
```

**Redis Implementation:**
- **Ключи:**
  - `poll:{id}:info` — String с JSON метаданными
  - `poll:{id}:votes` — Hash с счетчиками голосов
- **Атомарность:** Использование Pipeline для batch операций
- **TTL:** Синхронизирован для обоих ключей

**Пример атомарной операции:**
```go
func (s *RedisStorage) Vote(ctx context.Context, pollID string, optionIndices []int) error {
    // Валидация + проверка дубликатов
    poll, err := s.GetPoll(ctx, pollID)
    // ... validation logic ...

    // Атомарное увеличение счетчиков
    pipe := s.client.Pipeline()
    for _, idx := range optionIndices {
        field := fmt.Sprintf("%d", idx)
        pipe.HIncrBy(ctx, pollVotesKey(pollID), field, 1)
    }
    _, err = pipe.Exec(ctx)

    return err
}
```

#### 4. Model Layer (`internal/model/`)

**Domain entities и DTO:**
- `Poll` — основная entity
- `CreatePollRequest/Response` — DTO для создания
- `VoteRequest/Response` — DTO для голосования
- `PollResults` — результаты с агрегированными данными
- `ErrorResponse` — стандартизированные ошибки

**Validation tags:**
```go
type CreatePollRequest struct {
    Title   string   `json:"title" binding:"required,min=3,max=200"`
    Options []string `json:"options" binding:"required,min=2,max=10,dive,required,min=1,max=100"`
}
```

---

## 🎨 Frontend архитектура (Next.js)

### App Router структура

```
app/
├── layout.tsx                 # Root layout
├── page.tsx                   # Home page (/)
├── create/                    # Создание опроса (/create)
│   └── page.tsx
├── [id]/                      # Динамический роутинг (/[id])
│   ├── page.tsx               # Голосование
│   └── results/               # Результаты (/[id]/results)
│       └── page.tsx
├── config/                    # Конфигурация
│   └── api.ts
└── services/                  # API client
    └── pollService.ts
```

### Архитектурные решения

#### Server Components vs Client Components

**Server Components (по умолчанию):**
- Рендеринг на сервере (SSR)
- SEO-friendly
- Уменьшенный bundle size
- Пример: главная страница, страница результатов

**Client Components (`'use client'`):**
- Интерактивные компоненты
- Работа с state и effects
- Обработка событий
- Пример: формы голосования, графики

#### API Client Service

Централизованный сервис для работы с backend:

```typescript
// app/services/pollService.ts
export const pollService = {
  async createPoll(data: CreatePollRequest): Promise<CreatePollResponse> {
    const response = await fetch(`${API_URL}/polls`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  async vote(pollId: string, optionIndices: number[]): Promise<VoteResponse> {
    // ...
  },

  async getResults(pollId: string): Promise<PollResults> {
    // ...
  }
};
```

#### Конфигурация

Использование environment variables:
- **Client-side:** `NEXT_PUBLIC_*` (доступны в браузере)
- **Server-side:** обычные env vars (только на сервере)

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

---

## 🗄️ Структура данных в Redis

### Poll Metadata (String)

**Ключ:** `poll:{poll_id}:info`
**Тип:** String (JSON)
**TTL:** `POLL_DEFAULT_TTL` (168h)

**Структура:**
```json
{
  "id": "abc123",
  "title": "Какие языки программирования вы используете?",
  "options": ["Go", "Python", "JavaScript", "Rust", "TypeScript"],
  "created_at": "2025-12-08T10:00:00Z",
  "expires_at": "2025-12-15T10:00:00Z"
}
```

### Vote Counters (Hash)

**Ключ:** `poll:{poll_id}:votes`
**Тип:** Hash
**TTL:** `POLL_DEFAULT_TTL` (168h, синхронизирован с info)

**Структура:**
```
Field    Value
-----    -----
"0"   -> "42"    # Go: 42 голоса
"1"   -> "28"    # Python: 28 голосов
"2"   -> "35"    # JavaScript: 35 голосов
"3"   -> "15"    # Rust: 15 голосов
"4"   -> "30"    # TypeScript: 30 голосов
```

### Преимущества такой структуры

1. **Атомарность** — `HINCRBY` атомарен, не нужны локи
2. **Производительность** — O(1) для чтения/записи
3. **TTL** — автоматическое удаление устаревших опросов
4. **Простота** — минимальное количество операций

### Redis Pipeline для атомарности

```go
// Создание опроса - атомарная операция
pipe := redis.Pipeline()
pipe.Set(ctx, pollInfoKey(id), jsonData, ttl)
pipe.HSet(ctx, pollVotesKey(id), initialVotes)
pipe.Expire(ctx, pollVotesKey(id), ttl)
_, err := pipe.Exec(ctx)

// Голосование - атомарное увеличение нескольких счетчиков
pipe := redis.Pipeline()
for _, idx := range optionIndices {
    pipe.HIncrBy(ctx, pollVotesKey(id), fmt.Sprintf("%d", idx), 1)
}
_, err := pipe.Exec(ctx)
```

---

## 🔄 Поток данных

### Создание опроса

```
User → Frontend Form
  ↓
  POST /api/v1/polls { title, options }
  ↓
Handler.CreatePoll
  ↓ validate request
Service.CreatePoll
  ↓ generate ID, build Poll entity
Storage.CreatePoll
  ↓ Redis Pipeline:
    1. SET poll:{id}:info {json} EX 168h
    2. HSET poll:{id}:votes 0 0 1 0 2 0 ...
    3. EXPIRE poll:{id}:votes 168h
  ↓
Response { poll_id, vote_url, results_url }
  ↓
Frontend → Redirect to /[id]
```

### Голосование

```
User → Frontend Vote Form
  ↓
  POST /api/v1/polls/{id}/vote { option_indices: [0, 2] }
  ↓
Handler.Vote
  ↓ validate request, extract id
Service.Vote
  ↓ business logic
Storage.Vote
  ↓ 1. Check poll exists (EXISTS poll:{id}:info)
  ↓ 2. Get poll for validation
  ↓ 3. Validate indices and check duplicates
  ↓ 4. Redis Pipeline:
        HINCRBY poll:{id}:votes "0" 1
        HINCRBY poll:{id}:votes "2" 1
  ↓
Response { success: true, message: "Votes registered successfully (2 options)" }
  ↓
Frontend → Show success, redirect to results
```

### Получение результатов

```
User → Frontend Results Page
  ↓
  GET /api/v1/polls/{id}/results
  ↓
Handler.GetResults
  ↓ extract id
Service.GetResults
  ↓
Storage.GetResults
  ↓ 1. GET poll:{id}:info → parse Poll
  ↓ 2. HGETALL poll:{id}:votes → get all counters
  ↓ 3. Build PollResults { poll, votes map, total }
  ↓
Response {
  poll: { id, title, options, ... },
  votes: { "Go": 42, "Python": 28, ... },
  total: 150
}
  ↓
Frontend → Render charts with Recharts
```

---

## 🔐 Безопасность и надежность

### Обработка ошибок

**Уровни обработки:**
1. **Storage** — возвращает типизированные ошибки (`ErrPollNotFound`, `ErrInvalidOption`)
2. **Service** — логирует и пробрасывает ошибки выше
3. **Handler** — мапит ошибки в HTTP статусы (404, 400, 500)

### Graceful Shutdown

```go
// Listening for SIGINT/SIGTERM
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// Graceful shutdown с таймаутом
ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
defer cancel()

server.Shutdown(ctx)
storage.Close()
```

### Валидация

**Уровень 1: Gin binding tags**
```go
type CreatePollRequest struct {
    Title   string   `binding:"required,min=3,max=200"`
    Options []string `binding:"required,min=2,max=10,dive,required,min=1,max=100"`
}
```

**Уровень 2: Business validation**
```go
// Проверка дубликатов индексов
seen := make(map[int]bool)
for _, idx := range optionIndices {
    if seen[idx] {
        return ErrDuplicateOption
    }
    seen[idx] = true
}
```

### Middleware

- **Recovery** — отлов panic и возврат 500
- **Logger** — структурированное логирование всех запросов
- **CORS** — правильные headers для frontend

---

## 📊 Масштабирование

### Текущая архитектура

- **Stateless backend** — можно горизонтально масштабировать
- **Redis** — single instance (для начала)
- **Frontend** — SSR/SSG через Next.js

### Рекомендации для масштабирования

1. **Redis Cluster** — для высоконагруженных систем
2. **Load Balancer** — перед backend instances (Caddy, nginx, HAProxy)
3. **Redis Sentinel** — для high availability
4. **Кеширование** — CDN для frontend assets
5. **Метрики** — Prometheus для мониторинга

---

## 🔍 Дальнейшее чтение

- [API Documentation](API-Documentation) — детальное описание endpoints
- [Data Structure](Data-Structure) — подробнее о Redis схеме
- [Development Guide](Development-Guide) — как разрабатывать новые фичи
- [Testing](Testing) — тестирование архитектуры

---

_Последнее обновление: December 2025_
