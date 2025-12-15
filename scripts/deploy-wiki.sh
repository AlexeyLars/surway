#!/bin/bash
set -e

echo "📚 Deploying Surway Wiki documentation..."

WIKI_REPO="https://github.com/AlexeyLars/surway.wiki.git"
WIKI_DIR="temp-wiki"

# Клонирование Wiki репозитория
echo "📥 Cloning Wiki repository..."
if [ -d "$WIKI_DIR" ]; then
  rm -rf "$WIKI_DIR"
fi

git clone "$WIKI_REPO" "$WIKI_DIR"

# Копирование файлов (кроме Home.md и README.md)
echo "📄 Copying wiki pages..."
cp docs/wiki/Architecture.md "$WIKI_DIR/"
cp docs/wiki/API-Documentation.md "$WIKI_DIR/"
cp docs/wiki/Deployment.md "$WIKI_DIR/"
cp docs/wiki/Development-Guide.md "$WIKI_DIR/"
cp docs/wiki/Configuration.md "$WIKI_DIR/"

# Commit и push
echo "💾 Committing changes..."
cd "$WIKI_DIR"
git add .
git commit -m "docs: add comprehensive wiki documentation

- Architecture.md: detailed system architecture with diagrams
- API-Documentation.md: complete REST API reference
- Deployment.md: production deployment guide
- Development-Guide.md: guide for contributors
- Configuration.md: full configuration reference" || echo "No changes to commit"

echo "🚀 Pushing to GitHub..."
git push origin master

# Cleanup
cd ..
rm -rf "$WIKI_DIR"

echo ""
echo "✅ Wiki deployed successfully!"
echo "🌐 View at: https://github.com/AlexeyLars/surway/wiki"
echo ""
echo "📄 Pages added:"
echo "  - Architecture"
echo "  - API-Documentation"
echo "  - Deployment"
echo "  - Development-Guide"
echo "  - Configuration"
