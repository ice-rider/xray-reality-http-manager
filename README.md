# Xray Reality Manager

HTTP API для управления Xray (VLESS + Reality). Динамическое добавление клиентов с генерацией ссылок и поддержкой ML-DSA65.

## Возможности

- ✅ Генерация конфигурации `config.json` для Xray Core
- ✅ REST API для добавления клиентов
- ✅ VLESS + Reality транспорт
- ✅ Поддержка постквантовой криптографии ML-DSA65
- ✅ Автоматическая генерация VLESS ссылок

## Быстрый старт

### 1. Настройка переменных окружения

```bash
cp .env.example .env
```

Заполните `.env`:

```bash
# ML-DSA65
mldsa65_sign=ваша_подпись
mldsa65_public=ваш_публичный_ключ

# X25519 ключи
private_key=ваш_приватный_ключ
public_key=ваш_публичный_ключ

# Short IDs (опционально)
shorts_id=abc123,def456

# IP сервера
server_ip=ваш_server_ip
```

### 2. Запуск

```bash
go run cmd/server/main.go
```

Сервер запустится на порту **8080**.

## API

### Добавить клиента

```bash
curl -X POST http://localhost:8080/api/v1/client \
  -H "Content-Type: application/json" \
  -d '{
    "flow": "xtls-rprx-vision",
    "client_name": "my-client"
  }'
```

**Ответ:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "flow": "xtls-rprx-vision",
  "link": "vless://..."
}
```

## Структура проекта

```
├── cmd/server/main.go       # Точка входа
├── internal/
│   ├── app/                 # Сервис приложения
│   ├── env/                 # Переменные окружения
│   └── http/                # HTTP слой
└── pkg/
    ├── domain/              # Доменные модели
    └── xray/                # Утилиты Xray
```

## Технологии

- Go 1.22
- Xray Core
- VLESS + Reality
- ML-DSA65

## Лицензия

MIT
