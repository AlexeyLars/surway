# Survey App - Frontend

Фронтенд приложения для создания и проведения опросов, построенный на Next.js 15 с React 19.

## Конфигурация

Все технические параметры (URL API, порты, хосты) вынесены в конфигурацию.

### Настройка локальной разработки

1. Скопируйте файл `env.example` в `.env.local`:
   ```bash
   cp env.example .env.local
   ```

2. Отредактируйте `.env.local` при необходимости. Значения по умолчанию соответствуют `../configs/frontend.conf`:
   - API: `http://localhost:8080/api/v1`
   - Frontend: `http://localhost:3000`

3. Запустите приложение:
   ```bash
   npm install
   npm run dev
   ```

### Переменные окружения

Все переменные должны начинаться с `NEXT_PUBLIC_` для доступа в браузере:

- `NEXT_PUBLIC_API_PROTOCOL` - протокол API (http/https)
- `NEXT_PUBLIC_API_HOST` - хост API
- `NEXT_PUBLIC_API_PORT` - порт API
- `NEXT_PUBLIC_API_VERSION` - версия API
- `NEXT_PUBLIC_FRONTEND_PROTOCOL` - протокол фронтенда
- `NEXT_PUBLIC_FRONTEND_HOST` - хост фронтенда
- `NEXT_PUBLIC_FRONTEND_PORT` - порт фронтенда

### Развёртывание на Vercel

При развёртывании на Vercel установите переменные окружения в настройках проекта. Переменная `VERCEL_URL` устанавливается автоматически и используется для генерации ссылок.

## Структура

- `/app` - страницы и маршруты приложения
  - `/config` - загрузчик конфигурации
  - `/services` - API сервисы
- `/components` - переиспользуемые компоненты
  - `/analytics` - компоненты аналитики
  - `/ui` - UI компоненты

## Технологии

- Next.js 15.5.4
- React 19.1.0
- TypeScript 5
- Tailwind CSS 4
- Recharts (визуализация)
- Framer Motion (анимации)
