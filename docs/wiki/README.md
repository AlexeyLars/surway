# Surway Wiki Content

Этот каталог содержит markdown файлы для GitHub Wiki проекта Surway.

## 📋 Список страниц

1. **Home.md** — главная страница Wiki
2. **Architecture.md** — детальная архитектура системы
3. **API-Documentation.md** — полная документация REST API
4. **Deployment.md** — руководство по деплою в production
5. **Development-Guide.md** — гайд для разработчиков
6. **Configuration.md** — полное описание конфигурации

---

## 🚀 Как загрузить на GitHub Wiki

GitHub Wiki — это отдельный Git репозиторий. Вот как загрузить эти страницы:

### Метод 1: Через Web интерфейс (простой)

1. Перейдите на https://github.com/AlexeyLars/surway/wiki
2. Нажмите "Create the first page" (если Wiki еще не создана)
3. Для каждого файла:
   - Нажмите "New Page"
   - Скопируйте содержимое из соответствующего `.md` файла
   - Укажите название страницы (без расширения .md)
   - Сохраните

### Метод 2: Через Git (рекомендуется)

```bash
# 1. Клонируйте Wiki репозиторий
git clone https://github.com/AlexeyLars/surway.wiki.git

# 2. Скопируйте все wiki файлы
cp docs/wiki/*.md surway.wiki/

# 3. Удалите README.md (не нужен в Wiki)
cd surway.wiki
rm README.md

# 4. Commit и push
git add .
git commit -m "docs: add comprehensive wiki documentation"
git push origin master
```

### Метод 3: Скрипт автоматизации

Создайте скрипт `scripts/deploy-wiki.sh`:

```bash
#!/bin/bash
set -e

WIKI_REPO="https://github.com/AlexeyLars/surway.wiki.git"
WIKI_DIR="temp-wiki"

echo "📚 Deploying Wiki..."

# Клонирование Wiki репозитория
if [ -d "$WIKI_DIR" ]; then
  rm -rf "$WIKI_DIR"
fi

git clone "$WIKI_REPO" "$WIKI_DIR"

# Копирование файлов
cp docs/wiki/*.md "$WIKI_DIR/"

# Удаление README (не нужен в Wiki)
rm -f "$WIKI_DIR/README.md"

# Commit и push
cd "$WIKI_DIR"
git add .
git commit -m "docs: update wiki documentation" || echo "No changes to commit"
git push origin master

# Cleanup
cd ..
rm -rf "$WIKI_DIR"

echo "✅ Wiki deployed successfully!"
echo "🌐 View at: https://github.com/AlexeyLars/surway/wiki"
```

**Использование:**
```bash
chmod +x scripts/deploy-wiki.sh
./scripts/deploy-wiki.sh
```

---

## 📝 Структура Wiki

После загрузки, Wiki будет иметь следующую структуру:

```
Home (Home.md)
├── Quick Start
├── Architecture
│   ├── Backend Architecture
│   ├── Frontend Architecture
│   └── Data Structure
├── API Documentation
│   ├── Endpoints
│   ├── Models
│   └── Examples
├── Development Guide
│   ├── Setup
│   ├── Backend Development
│   ├── Frontend Development
│   └── Contributing
├── Deployment
│   ├── Docker Compose
│   ├── Caddy Setup
│   └── Production Best Practices
└── Configuration
    ├── Backend Config
    ├── Frontend Config
    └── Redis Config
```

---

## 🔗 Внутренние ссылки

В Wiki используйте относительные ссылки без `.md`:

```markdown
См. [Architecture](Architecture) для деталей.
Прочитайте [API Documentation](API-Documentation).
```

GitHub Wiki автоматически создаст правильные ссылки.

---

## ✏️ Редактирование Wiki

### Через Web

1. Перейдите на страницу Wiki
2. Нажмите "Edit" справа вверху
3. Отредактируйте markdown
4. Сохраните

### Через Git

```bash
# 1. Клонируйте Wiki
git clone https://github.com/AlexeyLars/surway.wiki.git
cd surway.wiki

# 2. Редактируйте файлы
nano Home.md

# 3. Commit и push
git add Home.md
git commit -m "docs: update home page"
git push origin master
```

---

## 🎨 Markdown Features

GitHub Wiki поддерживает:

- **Заголовки:** `# H1`, `## H2`, `### H3`
- **Списки:** `- item` или `1. item`
- **Код:** ` ```language ` для блоков
- **Таблицы:** `| Header | Header |`
- **Ссылки:** `[text](url)` или `[[WikiPage]]`
- **Картинки:** `![alt](url)`
- **Emoji:** `:rocket:` → 🚀
- **Task lists:** `- [ ]` и `- [x]`

---

## 🖼️ Добавление изображений

1. **Загрузите** изображения в репозиторий:
   ```
   docs/images/architecture-diagram.png
   ```

2. **Используйте** в Wiki:
   ```markdown
   ![Architecture Diagram](https://raw.githubusercontent.com/AlexeyLars/surway/main/docs/images/architecture-diagram.png)
   ```

Или используйте GitHub Issues для хостинга картинок:
1. Создайте Issue
2. Перетащите картинку
3. Скопируйте сгенерированный URL
4. Используйте в Wiki

---

## 🔄 Синхронизация с основным репозиторием

Рекомендуется:
1. Хранить источники Wiki в `docs/wiki/` основного репозитория
2. При изменениях обновлять оба места
3. Использовать CI/CD для автоматической синхронизации

**GitHub Actions пример:**

```yaml
# .github/workflows/sync-wiki.yml
name: Sync Wiki

on:
  push:
    branches: [main]
    paths:
      - 'docs/wiki/**'

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Deploy to Wiki
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          git clone "https://x-access-token:${GITHUB_TOKEN}@github.com/AlexeyLars/surway.wiki.git" wiki
          cp docs/wiki/*.md wiki/
          cd wiki
          rm README.md
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add .
          git commit -m "docs: auto-sync from main repo" || exit 0
          git push
```

---

## 📞 Вопросы?

- **GitHub Wiki Guide:** https://docs.github.com/en/communities/documenting-your-project-with-wikis
- **Markdown Guide:** https://www.markdownguide.org/
- **Issues:** https://github.com/AlexeyLars/surway/issues

---

_Последнее обновление: December 2025_
