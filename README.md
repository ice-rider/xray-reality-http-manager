# Xray Reality Manager

Сервер управления конфигурацией Xray Core с поддержкой VLESS + Reality транспорта и постквантовой криптографии ML-DSA65.

## Возможности

- ✅ Генерация конфигурации `config.json` для Xray Core
- ✅ REST API для добавления клиентов
- ✅ VLESS + Reality транспорт
- ✅ Поддержка постквантовой криптографии ML-DSA65
- ✅ Автоматическая генерация VLESS ссылок для подключения клиентов
- ✅ **Мониторинг трафика клиентов** (статистика по uplink/downlink)
- ✅ **gRPC API для получения статистики** в реальном времени

## Быстрый старт

### 1. Настройка переменных окружения

```bash
cp .env.example .env
```

### 2. Генерация ключей

#### Локально (требуется установленный Xray Core)

```shell
xray x25519
# PrivateKey: qNtTKU81oh3x1nwMPrgl0l5Y559VJfMOdtfjZ_nZdEs
# Password (PublicKey): 4bBkRexEWLhNOzClX_LEwyAnQlE6wEpW12P89g0kpjU
# Hash32: Uh7bCXelOBJ0PCzJF-utt6XvNPqFymTj9Mz9mQRMS70

xray mldsa65
# Seed: Uyamaka3HguRXvo55r2yJoRfu44sKWP1_uaRHHwwrkI
# Verify: UQus-pTcPLuOuCZRo4HdWAJurS8sIzZH2gL...(very long string)
```

#### Через Docker

```bash
# X25519 ключи
docker run --rm -it ghcr.io/ice-rider/xray-docker-alpine:latest xray x25519

# ML-DSA65 (постквантовая подпись)
docker run --rm -it ghcr.io/ice-rider/xray-docker-alpine:latest xray mldsa65
```

### 3. Заполнение .env

```bash
# X25519 ключи для Reality
private_key=ваш_приватный_ключ
public_key=ваш_публичный_ключ

# ML-DSA65 (постквантовая подпись)
mldsa65_seed=ваша_подпись
mldsa65_public=ваш_публичный_ключ

# Short IDs для идентификации клиентов (опционально)
shorts_id=abc123,def456

# IP сервера для генерации ссылок
server_ip=ваш_server_ip
```

### 4. Запуск

#### Локально

```bash
go run cmd/server/main.go
```

#### Через Docker

```bash
# Сборка образа
docker build -t xray-server .

# Запуск контейнера
docker run -d --name xray-server \
  --env-file .env \
  -p 8080:8080 \
  -v $(pwd)/config.json:/app/config.json \
  xray-server
```

Сервер запустится на порту **8080** и создаст `config.json` в текущей директории.

#### Docker Compose (рекомендуется для совместного использования с Xray)

Пример `docker-compose.yml` для совместной работы xray-server и xray-core:

```yaml
services:
  xray-server:
    image: ghcr.io/ice-rider/xray-reality-http-manager:latest
    container_name: xray-server
    restart: unless-stopped
    environment:
      - CONFIG_PATH=/app/config.json
      - private_key=
      - public_key=
      - mldsa65_seed=
      - mldsa65_public=
      - shorts_id=
      - server_ip=
    networks:
      - app
    ports:
      - "8080:8080"
    env_file:
      - .env
    volumes:
      - ./config.json:/app/config.json

  xray:
    image: ghcr.io/ice-rider/xray-docker-alpine:latest
    container_name: xray
    restart: unless-stopped
    networks:
      - app
    ports:
      - "443:443"
    volumes:
      - ./config.json:/etc/xray/config.json:ro
    depends_on:
      - xray-server
    command: ["-c", "/etc/xray/config.json"]

networks:
  monitoring:
    driver: bridge
```

**Важно:** После добавления клиента через API нужно перезагрузить Xray:

```bash
# Перезагрузка Xray для применения новой конфигурации
docker restart xray
```

## Мониторинг (Prometheus + Grafana + HTTP API)

Есть два способа получения статистики:

### 1. HTTP API (рекомендуется для простоты)

Используйте endpoint `/api/v1/stats` для получения статистики через REST API:

```bash
curl -X GET http://localhost:8080/api/v1/stats \
  -H "Authorization: Bearer <token>"
```

Этот метод не требует дополнительных компонентов и работает через gRPC API Xray.

### 2. Prometheus + Grafana (для продвинутого мониторинга)

Для включения мониторинга добавьте в `docker-compose.yml`:

```yaml
services:
  # ... xray-server и xray ...

  xray-exporter:
    image: anatolykopyl/xray-exporter:latest
    container_name: xray-exporter
    restart: unless-stopped
    networks:
      - monitoring
    command:
      - --xray-endpoint=xray:54321
      - --listen=:9550

  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    restart: unless-stopped
    networks:
      - monitoring

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    networks:
      - monitoring
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    networks:
      - monitoring
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"

volumes:
  prometheus_data:
  grafana_data:
```

Создайте `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'xray'
    static_configs:
      - targets: ['xray-exporter:9550']

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
```

**Важно:** API inbound на порт 54321 добавляется автоматически в `config.json`.

Grafana дашборд: импортируйте [16426](https://grafana.com/grafana/dashboards/16426-xray-core/) или аналогичный.

## Авторизация

Все endpoints API защищены JWT аутентификацией, кроме `/api/v1/auth/login`.

### Переменные окружения

```bash
# JWT секрет (опционально, генерируется автоматически если не указан)
jwt_secret=your-secret-key

# Логин/пароль админа (опционально, по умолчанию admin/admin)
admin_username=admin
admin_password=admin
```

### Получить токен

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Ответ:
# {
#   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "expires_in": 86400
# }
```

### Использовать токен

```bash
curl -X POST http://localhost:8080/api/v1/client \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "flow": "xtls-rprx-vision",
    "client_name": "my-client"
  }'
```

**Время жизни токена:** 24 часа.

## API

### Добавить клиента

Для запроса требуется JWT токен в заголовке `Authorization: Bearer <token>`.

```bash
curl -X POST http://localhost:8080/api/v1/client \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "flow": "xtls-rprx-vision",
    "client_name": "my-client"
  }'

# {
#    "id":"49c64dcd-f83c-486a-a590-2a38cf417829",
#    "flow":"xtls-rprx-vision",
#    "link":"vless://49c64dcd-f83c-486a-a590-2a38cf417829@\u003c0.0.0.0\u003e:443?
#      encryption = none \u0026
#      flow = xtls-rprx-vision \u0026
#      fp = firefox \u0026
#      pbk = %3Cpublic_key%3E \u0026
#      pqv = %3Cpublic_mldsa65_verify%3E \u0026
#      security = reality\u0026
#      sid = %3C123abc%3E\u0026
#      sni = www.apple.com#my-client"}
```

`'%3C'`, `'%3E'` - `'<'` и `'>'`
`\u0026` представляет собой символ & в жсон энкодинге
ссылка приведена с переносами для наглядности

### Получить список клиентов

```bash
curl -X GET http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer <token>"

# Ответ:
# {
#   "clients": [
#     {"id": "uuid-1", "flow": "xtls-rprx-vision", "email": "client1@local"},
#     {"id": "uuid-2", "flow": "", "email": "client2@local"}
#   ]
# }
```

### Получить статистику трафика клиентов

Возвращает статистику по всем клиентам с именами (email) и объёмом трафика:

```bash
curl -X GET http://localhost:8080/api/v1/stats \
  -H "Authorization: Bearer <token>"

# Ответ:
# {
#   "clients": [
#     {
#       "id": "uuid-1",
#       "email": "client1@local",
#       "flow": "xtls-rprx-vision",
#       "uplink": 1024567,      # байт отправлено
#       "downlink": 8765432,    # байт получено
#       "total": 9789999        # всего байт
#     },
#     {
#       "id": "uuid-2",
#       "email": "client2@local",
#       "flow": "",
#       "uplink": 0,
#       "downlink": 0,
#       "total": 0
#     }
#   ]
# }
```

**Важно:** Статистика доступна только если в конфигурации Xray включена статистика (включено по умолчанию через `policy.levels.0.statsUserUplink/Downlink`).

### Получить статистику конкретного клиента

```bash
curl -X GET "http://localhost:8080/api/v1/stats?email=client1@local" \
  -H "Authorization: Bearer <token>"

# Ответ:
# {
#   "client": {
#     "id": "uuid-1",
#     "email": "client1@local",
#     "flow": "xtls-rprx-vision",
#     "uplink": 1024567,
#     "downlink": 8765432,
#     "total": 9789999
#   }
# }
```

## Структура проекта

Проект организован в соответствии с принципами Clean Architecture:

```
xray_server/
├── cmd/server/main.go              # Точка входа: сборка зависимостей (DI)
├── internal/
│   ├── domain/                     # Бизнес-модели и интерфейсы
│   │   ├── user.go                 # User, LoginRequest, TokenResponse + UserRepository, JWTService
│   │   ├── config.go               # Config, Inbound, Client + ConfigRepository
│   │   └── stats.go                # UserTraffic + StatsRepository
│   │
│   ├── usecase/                    # Бизнес-правила приложения
│   │   ├── auth/
│   │   │   └── login.go            # LoginUseCase, Execute()
│   │   ├── config/
│   │   │   └── usecase.go          # ConfigUseCase, AddClient(), GetClients()
│   │   └── stats/
│   │       └── usecase.go          # StatsUseCase, GetAllClientsStats()
│   │
│   ├── repository/                 # Реализации портов (БД, gRPC, JWT)
│   │   ├── user_repository.go      # UserRepositorySQLite
│   │   ├── jwt_service.go          # JWTService
│   │   └── stats_repository.go     # StatsRepositorygRPC
│   │
│   └── delivery/                   # HTTP/gRPC адаптеры
│       └── http/
│           ├── router.go           # Маршрутизация, Server
│           ├── auth_handler.go     # Login handler
│           ├── auth_middleware.go  # JWT middleware
│           ├── client_handler.go   # CreateClient, GetClients
│           └── stats_handler.go    # GetClientsStats
│
├── pkg/                            # Общие утилиты (не зависят от бизнеса)
│   ├── config/
│   │   └── env.go                  # Загрузка переменных окружения
│   └── xrayutil/
│       ├── util.go                 # GenerateUUID, ParseShortIds
│       └── vless.go                # GenerateVlessLink
│
├── go.mod
└── go.sum
```

## Технологии

- Go 1.25+
- Xray Core
- VLESS + Reality
- ML-DSA65 (постквантовая криптография)
- gRPC (StatsService API)
- JWT аутентификация
- SQLite
- godotenv

## Лицензия

[MIT](./LICENSE)
