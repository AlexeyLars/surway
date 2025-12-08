# 💻 Development Guide

Полное руководство по разработке и контрибьюции в проект Surway.

---

## 📋 Содержание

1. [Начало работы](#начало-работы)
2. [Структура проекта](#структура-проекта)
3. [Backend разработка](#backend-разработка)
4. [Frontend разработка](#frontend-разработка)
5. [Тестирование](#тестирование)
6. [Code Style](#code-style)
7. [Git Workflow](#git-workflow)
8. [Contributing](#contributing)

---

## 🚀 Начало работы

### Prerequisites

- Git
- Go 1.24.2+
- Node.js 20+
- Redis 7+ (или Docker)
- Make (опционально)

### Клонирование и setup

```bash
# 1. Fork репозитория на GitHub

# 2. Клонирование вашего fork
git clone git@github.com:YOUR_USERNAME/surway.git
cd surway

# 3. Добавление upstream remote
git remote add upstream git@github.com:AlexeyLars/surway.git

# 4. Создание ветки для разработки
git checkout -b feature/my-awesome-feature develop
```

### Локальный запуск

**Вариант 1: Docker Compose (рекомендуется)**
```bash
make docker-up
# или
docker compose up -d
```

**Вариант 2: Локально**

*Backend:*
```bash
# Запустите Redis
docker run -d -p 6379:6379 redis:7-alpine

# Установите зависимости
cd backend
go mod download

# Запустите сервер
go run cmd/api/main.go
```

*Frontend:*
```bash
cd frontend
npm install
cp env.example .env.local
npm run dev
```

### Проверка

- Backend: http://localhost:8080/health
- Frontend: http://localhost:3000
- Swagger: http://localhost:8080/swagger/index.html

---

## 📁 Структура проекта

```
surway/
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── api/
│   │       └── main.go         # Entry point
│   ├── internal/               # Приватный код
│   │   ├── handler/            # HTTP handlers
│   │   ├── service/            # Бизнес-логика
│   │   ├── storage/            # Работа с БД
│   │   ├── model/              # Модели данных
│   │   ├── config/             # Конфигурация
│   │   └── lib/                # Утилиты
│   ├── docs/                   # Swagger docs (auto-generated)
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # Next.js frontend
│   ├── app/                    # App Router
│   │   ├── page.tsx
│   │   ├── layout.tsx
│   │   ├── create/
│   │   ├── [id]/
│   │   ├── config/
│   │   └── services/
│   ├── components/             # React компоненты
│   ├── public/
│   ├── Dockerfile
│   └── package.json
│
├── configs/                    # Конфигурационные файлы
├── docs/                       # Документация
│   └── wiki/                   # Wiki страницы
├── Makefile                    # Удобные команды
├── docker-compose.yml          # Dev environment
├── docker-compose.prod.yml     # Prod environment
└── README.md
```

---

## 🔧 Backend разработка (Go)

### Архитектура слоев

**1. Handler Layer** (`internal/handler/`)
- HTTP обработка
- Валидация запросов
- Маппинг ответов

**2. Service Layer** (`internal/service/`)
- Бизнес-логика
- Не знает о HTTP
- Работает через интерфейсы

**3. Storage Layer** (`internal/storage/`)
- CRUD операции
- Работа с Redis
- Реализует Storage интерфейс

### Добавление нового endpoint

**Пример: добавим `GET /api/v1/polls/{id}` для получения метаданных опроса**

**1. Добавить метод в Storage интерфейс:**

```go
// internal/storage/redis.go
type Storage interface {
    CreatePoll(ctx context.Context, poll *model.Poll, ttl time.Duration) error
    GetPoll(ctx context.Context, pollID string) (*model.Poll, error)  // Уже есть
    Vote(ctx context.Context, pollID string, optionIndices []int) error
    GetResults(ctx context.Context, pollID string) (*model.PollResults, error)
    Close() error
}
```

**2. Добавить handler:**

```go
// internal/handler/poll.go

// GetPoll godoc
// @Summary      Get poll metadata
// @Description  Return poll metadata without vote counts
// @Tags         polls
// @Produce      json
// @Param        id path string true "Poll ID"
// @Success      200 {object} model.Poll
// @Failure      404 {object} model.ErrorResponse
// @Router       /polls/{id} [get]
func (h *PollHandler) GetPoll(c *gin.Context) {
    pollID := c.Param("id")

    poll, err := h.service.GetPoll(c.Request.Context(), pollID)
    if err != nil {
        if errors.Is(err, storage.ErrPollNotFound) {
            c.JSON(http.StatusNotFound, model.ErrorResponse{
                Error:   "poll_not_found",
                Message: "Poll not found or expired",
            })
            return
        }

        c.JSON(http.StatusInternalServerError, model.ErrorResponse{
            Error:   "internal_error",
            Message: "Failed to get poll",
        })
        return
    }

    c.JSON(http.StatusOK, poll)
}
```

**3. Добавить метод в Service:**

```go
// internal/service/poll.go

func (s *PollService) GetPoll(ctx context.Context, pollID string) (*model.Poll, error) {
    poll, err := s.storage.GetPoll(ctx, pollID)
    if err != nil {
        if err == storage.ErrPollNotFound {
            s.logger.WarnContext(ctx, "poll not found", slog.String("poll_id", pollID))
            return nil, err
        }

        s.logger.ErrorContext(ctx, "failed to get poll",
            slog.String("poll_id", pollID),
            slog.String("error", err.Error()),
        )
        return nil, fmt.Errorf("failed to get poll: %w", err)
    }

    return poll, nil
}
```

**4. Зарегистрировать роут:**

```go
// internal/handler/router.go

func SetupRouter(handler *PollHandler, logger *slog.Logger, releaseMode bool) *gin.Engine {
    // ... existing code ...

    v1 := router.Group("/api/v1")
    {
        polls := v1.Group("/polls")
        {
            polls.POST("", handler.CreatePoll)
            polls.GET("/:id", handler.GetPoll)           // НОВЫЙ ENDPOINT
            polls.POST("/:id/vote", handler.Vote)
            polls.GET("/:id/results", handler.GetResults)
        }
    }

    return router
}
```

**5. Регенерировать Swagger:**

```bash
cd backend
swag init -g cmd/api/main.go
```

**6. Тестирование:**

```bash
# Запустить сервер
go run cmd/api/main.go

# Тест
curl http://localhost:8080/api/v1/polls/abc123
```

### Логирование

Используйте структурированное логирование:

```go
s.logger.InfoContext(ctx, "poll created",
    slog.String("poll_id", pollID),
    slog.String("title", req.Title),
    slog.Int("options_count", len(req.Options)),
)

s.logger.ErrorContext(ctx, "failed to create poll",
    slog.String("poll_id", pollID),
    slog.String("error", err.Error()),
)
```

### Тестирование

**Unit тесты:**

```go
// internal/service/poll_test.go

func TestPollService_CreatePoll(t *testing.T) {
    // Arrange
    mockStorage := &MockStorage{}
    mockConfig := &config.Config{
        Poll: config.PollConfig{
            DefaultTTL: 168 * time.Hour,
        },
    }
    logger := slog.Default()
    service := NewPollService(mockStorage, mockConfig, logger)

    req := &model.CreatePollRequest{
        Title:   "Test Poll",
        Options: []string{"A", "B", "C"},
    }

    // Act
    resp, err := service.CreatePoll(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.PollID)
    assert.Contains(t, resp.VoteURL, resp.PollID)
}
```

**Запуск тестов:**

```bash
cd backend
go test ./...                          # Все тесты
go test -v ./internal/service/...      # Конкретный пакет
go test -race ./...                    # С race detector
go test -cover ./...                   # С покрытием
```

### Makefile команды

```bash
make help              # Показать команды
make build             # Собрать бинарник
make run               # Запустить локально
make test              # Запустить тесты
make test-coverage     # Покрытие тестами
make fmt               # Форматировать код
make lint              # Линтинг
make deps              # Обновить зависимости
```

---

## 🎨 Frontend разработка (Next.js)

### Структура App Router

```
app/
├── layout.tsx          # Root layout (обертка для всех страниц)
├── page.tsx            # Home page (/)
├── create/
│   └── page.tsx        # Create poll page (/create)
├── [id]/               # Dynamic routing
│   ├── page.tsx        # Vote page (/[id])
│   └── results/
│       └── page.tsx    # Results page (/[id]/results)
├── config/
│   └── api.ts          # API configuration
└── services/
    └── pollService.ts  # API client
```

### Добавление новой страницы

**Пример: страница "О проекте"**

**1. Создайте страницу:**

```tsx
// app/about/page.tsx

export default function AboutPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-4">О проекте Surway</h1>
      <p className="text-lg">
        Surway — современный сервис для создания и проведения опросов...
      </p>
    </div>
  );
}
```

Страница автоматически доступна по `/about`!

### Server vs Client Components

**Server Component (по умолчанию):**

```tsx
// app/[id]/results/page.tsx

export default async function ResultsPage({ params }: { params: { id: string } }) {
  // Данные загружаются на сервере
  const results = await pollService.getResults(params.id);

  return (
    <div>
      <h1>{results.poll.title}</h1>
      <ResultsChart data={results} />
    </div>
  );
}
```

**Client Component:**

```tsx
// components/VoteForm.tsx
'use client';  // ОБЯЗАТЕЛЬНО для интерактивных компонентов

import { useState } from 'react';

export function VoteForm({ options }: { options: string[] }) {
  const [selected, setSelected] = useState<number[]>([]);

  const handleVote = async () => {
    // ... voting logic
  };

  return (
    <form>
      {options.map((option, index) => (
        <label key={index}>
          <input
            type="checkbox"
            checked={selected.includes(index)}
            onChange={() => {/* toggle */}}
          />
          {option}
        </label>
      ))}
      <button onClick={handleVote}>Голосовать</button>
    </form>
  );
}
```

### API Client

Используйте централизованный API client:

```typescript
// app/services/pollService.ts

const API_URL = `${API_CONFIG.protocol}://${API_CONFIG.host}:${API_CONFIG.port}/api/${API_CONFIG.version}`;

export const pollService = {
  async createPoll(data: CreatePollRequest): Promise<CreatePollResponse> {
    const response = await fetch(`${API_URL}/polls`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
      cache: 'no-store',  // Для Server Components
    });

    if (!response.ok) {
      throw new Error('Failed to create poll');
    }

    return response.json();
  },

  // ... другие методы
};
```

### Стилизация (Tailwind CSS)

```tsx
// Пример компонента с Tailwind

export function Button({ children, onClick }: ButtonProps) {
  return (
    <button
      onClick={onClick}
      className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
    >
      {children}
    </button>
  );
}
```

### Анимации (Framer Motion)

```tsx
'use client';

import { motion } from 'framer-motion';

export function AnimatedCard({ children }: { children: React.ReactNode }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="p-6 bg-white rounded-lg shadow-lg"
    >
      {children}
    </motion.div>
  );
}
```

### Тестирование

```bash
cd frontend
npm run lint           # ESLint
npm run build          # Проверка на ошибки сборки
npm run dev            # Dev сервер
```

---

## 🧪 Тестирование

### Backend

**Unit тесты:**
```bash
cd backend
go test ./internal/service/...
go test ./internal/storage/...
```

**Integration тесты:**
```bash
# Требуется запущенный Redis
docker run -d -p 6379:6379 redis:7-alpine
go test -tags=integration ./...
```

**Coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend

**Lint:**
```bash
cd frontend
npm run lint
```

**Build test:**
```bash
npm run build
```

---

## 📏 Code Style

### Go

Следуйте [Effective Go](https://golang.org/doc/effective_go) и:

- **gofmt** — форматирование (автоматически в большинстве IDE)
- **golangci-lint** — комплексный линтер

```bash
# Установка
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Запуск
cd backend
golangci-lint run
```

**Conventions:**
- Публичные функции начинаются с заглавной буквы
- Приватные — со строчной
- Короткие переменные в циклах: `i`, `err`, `ctx`
- Комментарии к публичным функциям
- Обработка ошибок явно (не игнорируйте `err`)

### TypeScript/React

**ESLint конфиг:**
```bash
cd frontend
npm run lint
```

**Conventions:**
- Компоненты в PascalCase: `VoteForm.tsx`
- Функции/переменные в camelCase: `handleVote`
- Константы в UPPER_CASE: `API_URL`
- Используйте TypeScript типы
- Избегайте `any`

---

## 🔀 Git Workflow

### Branching Strategy

```
main (stable, production-ready)
  └── develop (latest development)
        ├── feature/add-websocket
        ├── feature/auth-system
        └── fix/redis-connection
```

### Создание feature ветки

```bash
# 1. Обновите develop
git checkout develop
git pull upstream develop

# 2. Создайте feature ветку
git checkout -b feature/my-awesome-feature

# 3. Разработка
# ... код ...

# 4. Commit
git add .
git commit -m "feat: add awesome feature"

# 5. Push в ваш fork
git push origin feature/my-awesome-feature

# 6. Создайте Pull Request на GitHub
```

### Commit Messages

Используйте [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add WebSocket support for live updates
fix: resolve Redis connection timeout issue
docs: update API documentation
refactor: simplify poll creation logic
test: add unit tests for PollService
chore: update dependencies
```

**Формат:**
```
<type>: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` — новая функциональность
- `fix` — исправление бага
- `docs` — документация
- `refactor` — рефакторинг (без изменения функциональности)
- `test` — добавление/изменение тестов
- `chore` — рутинные задачи (обновление зависимостей и т.д.)
- `perf` — оптимизация производительности
- `ci` — CI/CD изменения

---

## 🤝 Contributing

### Pull Request Process

1. **Fork** репозиторий
2. **Создайте** feature ветку от `develop`
3. **Сделайте** изменения
4. **Добавьте** тесты (если применимо)
5. **Проверьте** что тесты проходят
6. **Запустите** линтер
7. **Commit** с правильным форматом сообщения
8. **Push** в ваш fork
9. **Откройте** Pull Request в `develop` ветку upstream
10. **Дождитесь** code review

### Code Review Checklist

**Reviewer проверяет:**
- [ ] Код следует стандартам проекта
- [ ] Есть тесты для новой функциональности
- [ ] Документация обновлена
- [ ] Нет breaking changes (или они документированы)
- [ ] Commit messages в правильном формате
- [ ] CI/CD пайплайн зеленый

**Author должен:**
- [ ] Ответить на комментарии reviewer
- [ ] Внести правки если нужно
- [ ] Обновить PR после изменений

---

## 🔍 Полезные ресурсы

### Документация

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Redis Documentation](https://redis.io/documentation)

### Наша документация

- [Architecture](Architecture)
- [API Documentation](API-Documentation)
- [Deployment](Deployment)
- [Configuration](Configuration)

---

## 💬 Вопросы?

- **GitHub Issues:** https://github.com/AlexeyLars/surway/issues
- **GitHub Discussions:** https://github.com/AlexeyLars/surway/discussions
- **Author:** [@AlexeyLars](https://github.com/AlexeyLars)

---

_Последнее обновление: December 2025_
