# Payment Service — Техническое Задание v8.0
**Уровень:** Production-Ready Go Backend
**Автор:** Гульжан Каракул | github.com/GulzhanKarakul/payment-service
**Обновлено:** service тесты готовы, coverage 96.3%, идём к handler тестам

---

## Цель проекта

Production-ready микросервис обработки транзакций с бонусной системой.
Основа для реального финтех-стартапа.

**Стек технологий:**
```
Язык:        Go 1.21
HTTP роутер: chi
БД:          PostgreSQL 15
Кэш:         Redis 7
Брокер:      Kafka
Метрики:     Prometheus + Grafana
Логи:        slog (structured JSON)
Контейнеры:  Docker + docker-compose
Тесты:       testify + testcontainers
Деплой:      VPS + Nginx + SSL
```

---

## Что уже сделано

### Инфраструктура
```
✅ git репозиторий — github.com/GulzhanKarakul/payment-service
✅ .gitignore, .dockerignore
✅ Dockerfile — multi-stage build (golang:1.21-alpine → alpine:3.18)
   — CGO_ENABLED=0, GOOS=linux, GOARCH=amd64
   — ldflags="-w -s" (уменьшение бинарника на 30-40%)
   — непривилегированный пользователь appuser
   — ca-certificates, tzdata
✅ docker-compose.yml — для локальной разработки
   — PostgreSQL 15 с healthcheck
   — Redis 7 с healthcheck
   — restart: unless-stopped
   — volumes с персистентностью данных
✅ .env / .env.example
```

### Конфигурация и подключения
```
✅ pkg/config/config.go
   — Структуры: Config, ServerConfig, DatabaseConfig, RedisConfig, KafkaConfig
   — func Load() *Config — из переменных окружения
   — func (c *Config) DSN() string — строка подключения PostgreSQL
   — godotenv.Load() — не падает если .env нет (для продакшна)

✅ pkg/database/postgres.go
   — func NewPostgres(cfg Config) (*sql.DB, error)
   — Настройки пула: MaxOpenConns(25), MaxIdleConns(5)
   — ConnMaxLifetime(5m), ConnMaxIdleTime(30m)
   — db.Ping() для проверки при старте

✅ cmd/main.go
   — Structured logger (slog JSON формат)
   — Загрузка конфига
   — Подключение к PostgreSQL
   — HTTP сервер с таймаутами (Read/Write/Idle)
   — Graceful Shutdown через SIGTERM/SIGINT
   — defer db.Close()
```

### База данных
```
✅ 001 — CREATE TABLE clients         (bonus_balance BIGINT)
✅ 002 — CREATE TABLE businesses      (bonus_balance BIGINT)
✅ 003 — CREATE TABLE transactions    (amount BIGINT, bonus_accrued BIGINT)
✅ 004 — CREATE INDEX
✅ 005 — ALTER TABLE transactions ADD COLUMN deleted_at
✅ 006 — CREATE TABLE business_bonus_settings
✅ 007 — ALTER TABLE payments RENAME TO transactions
✅ 008 — DROP UNIQUE на businesses.owner_phone

## Финальная схема БД

```sql
clients (
    id UUID PK DEFAULT gen_random_uuid(),
    phone VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    bonus_balance BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

businesses (
    id UUID PK DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    owner_phone VARCHAR(20) NOT NULL,     -- не UNIQUE, один owner → N бизнесов
    bonus_balance BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

business_bonus_settings (
    id UUID PK DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL UNIQUE REFERENCES businesses(id),
    bonus_percent DECIMAL(5,2) NOT NULL CHECK (bonus_percent > 0),
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

transactions (
    id UUID PK DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id),
    business_id UUID NOT NULL REFERENCES businesses(id),
    amount BIGINT NOT NULL,
    bonus_accrued BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'KZT',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
)
-- Статусы: pending → completed / failed / cancelled
```

Индексы:
— idx_payments_client_id
— idx_payments_business_id
— idx_payments_status
— idx_payments_created_at DESC

COMMENT ON для документации схемы.
```

**Ключевые решения:**
- Все денежные поля в БД — `BIGINT` (тиыны), не DECIMAL
- В Go — `int64`. Никаких ошибок округления
- `owner_phone` не UNIQUE — один owner может иметь несколько бизнесов
- Телефон хранится в E.164 формате: `+77771156580`


### internal/domain/
```
✅ Client, Business, BonusSettings, Transaction — чистые модели
✅ TransactionStatus + константы:
   StatusPending, StatusCompleted, StatusFailed, StatusCancelled
✅ errors.go:
   ErrClientNotFound, ErrClientIsNotActive, ErrClientAlreadyExist
   ErrBusinessNotFound, ErrBusinessIsNotActive, ErrBusinessAlreadyExist
   ErrTransactionNotFound, ErrTransactionCancelled, ErrTransactionProcessed
   ErrInsufficientBonusBalance
   ErrBonusSettingsNotFound, ErrBonusSettingsIsNotActive
   ErrNotFound, ErrInvalidInput, ErrInternalError
```

### internal/repository/
```
✅ client.go
   — Create (pq.Error 23505 → ErrClientAlreadyExist)
   — GetByID, GetByPhone, UpdateBonusBalance

✅ business.go
   — Create, GetByID
   — UpdateBonusBalance (delta int64, может быть отрицательным)

✅ bonus_settings.go
   — Upsert (ON CONFLICT DO UPDATE), GetByBusinessID

✅ transaction.go
   — CreateWithBonus — 4 операции в одной PG транзакции:
     1. INSERT transactions (status=pending)
     2. UPDATE clients bonus_balance += bonus
     3. UPDATE businesses bonus_balance -= bonus WHERE >= bonus
     4. UPDATE transactions SET status=completed, bonus_accrued=bonus
     n==0 на шаге 3 → ErrInsufficientBonusBalance → rollback
   — GetByID, GetByClientID(limit, offset), Cancel
```


### internal/service/
```
✅ repository.go — интерфейсы репозиториев
✅ service.go    — интерфейсы сервисов
✅ client.go     — проверка дублей по телефону
✅ business.go   — Create, GetByID, UpdateBonusBalance
✅ bonus_settings.go — проверка бизнеса + валидация bonusPercent > 0
✅ transaction.go:
   — проверка client.IsActive и business.IsActive
   — resolveBonusPercent — приватный метод
   — нет настроек или is_active=false → bonusPercent=0
   — пагинация: дефолт limit=20, offset=0
```

### internal/handler/ + middleware/
```
✅ handler/helpers.go
   — writeJSON, writeError
   — handleError — маппинг domain ошибок → HTTP коды:
     404: ErrNotFound, ErrClientNotFound, ErrBusinessNotFound, ErrTransactionNotFound
     409: ErrAlreadyExist, ErrTransactionCancelled
     422: ErrIsNotActive, ErrInsufficientBonusBalance
     400: ErrInvalidInput
     500: всё остальное + log.Error

✅ handler/handler.go
   — chi роутер с вложенными subrouters
   — /{id} как отдельный r.Route для корректной маршрутизации

✅ handler/client.go
   — POST   /api/v1/clients
   — GET    /api/v1/clients/{id}
   — GET    /api/v1/clients/phone/{phone}
   — Validate: E.164 формат, длина 10-16, name обязателен

✅ handler/business.go
   — POST   /api/v1/businesses
   — GET    /api/v1/businesses/{id}
   — POST   /api/v1/businesses/{id}/balance  ← пополнение бонусного баланса
   — Request: { "amount": int64 }
   — Validate: amount > 0

✅ handler/bonus_settings.go
   — POST   /api/v1/businesses/{id}/settings (upsert → 200)
   — GET    /api/v1/businesses/{id}/settings

✅ handler/transaction.go
   — POST   /api/v1/transactions
   — GET    /api/v1/transactions/{id}
   — GET    /api/v1/transactions?client_id=X&limit=20&offset=0
   — PATCH  /api/v1/transactions/{id}/cancel → {"message":"transaction cancelled"}

✅ middleware/logger.go
   — генерирует request_id (UUID v4)
   — кладёт в context, X-Request-ID header в ответе
   — логирует incoming request + request completed с duration_ms

✅ middleware/recovery.go
   — перехватывает panic, логирует stack trace
   — возвращает 500 без краша сервиса
   — пишет ответ напрямую (без импорта handler)

✅ Все endpoints проверены через Postman
✅ Postman Environment с переменными (base_url, client_id, business_id, transaction_id)
```

### Тесты — Repository слой ✅
```
✅ testcontainers PostgreSQL — один контейнер на пакет
✅ client, business, bonus_settings, transaction тесты
✅ Rollback проверен: при InsufficientBalance баланс клиента не изменился
✅ Pagination тесты с beyond_last_page кейсом
✅ go test ./internal/repository/... -race → PASS
```

### Тесты — Service слой ✅
```
✅ mockery генерированные моки (with-expecter: true)
✅ .mockery.yaml настроен для всех Repository интерфейсов
✅ helpers_test.go: testLogger, activeClient/inactiveBusiness/etc builders
✅ strPtr helper для *string литералов
✅ Константы testClientID, testBusinessID, testAmount int64

✅ client_test.go
   — Create: отдельные функции (Success, PhoneTaken, RepoSystemError)
   — GetByID: table-driven (3 кейса)
   — GetByPhone: table-driven с wantIsActive *bool паттерном

✅ business_test.go
   — Create: Success, SameOwnerMultipleBusinesses, RepoSystemError
   — GetByID: table-driven с wantIsActive *bool
   — UpdateBonusBalance: 4 отдельные функции (TopUp, Deduct, NotFound, ZeroDelta)

✅ bonus_settings_test.go
   — Upsert: CreateThenUpdate (.Twice()), BusinessNotFound, BusinessInactive,
             InvalidZeroPercent, InvalidNegativePercent, PropagateSystemError
   — GetByBusinessID: table-driven (3 кейса)

✅ transaction_test.go
   — Create: table-driven 10 кейсов (all branches covered)
   — GetByID: table-driven (3 кейса)
   — GetByClientID_Pagination: table-driven с t.Parallel() + tt:=tt
   — Cancel: table-driven (4 кейса)

✅ go test ./internal/service/... -race → PASS
✅ coverage: 96.3%
```

---

## API Endpoints — все рабочие

```
GET    /health                                      → {"status":"ok"}

POST   /api/v1/clients                             → 201 Client
GET    /api/v1/clients/{id}                        → 200 Client
GET    /api/v1/clients/phone/{phone}               → 200 Client

POST   /api/v1/businesses                          → 201 Business
GET    /api/v1/businesses/{id}                     → 200 Business
POST   /api/v1/businesses/{id}/balance             → 200 Business (пополнение баланса)
POST   /api/v1/businesses/{id}/settings            → 200 BonusSettings (upsert)
GET    /api/v1/businesses/{id}/settings            → 200 BonusSettings

POST   /api/v1/transactions                        → 201 Transaction
GET    /api/v1/transactions/{id}                   → 200 Transaction
GET    /api/v1/transactions?client_id=X&limit=20   → 200 []Transaction
PATCH  /api/v1/transactions/{id}/cancel            → 200 {"message":"transaction cancelled"}

System:
GET    /health                        — healthcheck
GET    /metrics                       — Prometheus метрики
```

---

## Бизнес логика — начисление бонусов

```
1. Проверить client.IsActive и business.IsActive
2. resolveBonusPercent:
   — нет настроек → 0
   — is_active=false → 0
   — есть и активны → bonusPercent
3. bonus = int64(float64(amount) * bonusPercent / 100)
4. PostgreSQL транзакция (атомарно):
   INSERT transactions (pending)
   UPDATE clients  bonus_balance += bonus
   UPDATE businesses bonus_balance -= bonus WHERE >= bonus  ← если 0 строк → rollback
   UPDATE transactions status=completed, bonus_accrued=bonus
5. После COMMIT → (будущее) publish в Kafka
```

---

## План разработки

```
✅ Фаза 1 — Инфраструктура + миграции + domain
✅ Фаза 2 — Repository слой
✅ Фаза 3 — Service слой
✅ Фаза 4 — Handler + chi + middleware + Postman
✅ Фаза 5.1 — Repository тесты (testcontainers)
✅ Фаза 5.2 — Service тесты (mockery, 96.3% coverage)

⬜ Фаза 5.3 — Handler тесты (httptest) ← СЕЙЧАС
   [ ] Добавить Service интерфейсы в .mockery.yaml
   [ ] mockery — сгенерировать моки сервисов
   [ ] helpers_test.go (newHandler, TestHandleError_MapsCorrectHTTPCodes)
   [ ] client_test.go
   [ ] business_test.go
   [ ] bonus_settings_test.go
   [ ] transaction_test.go
   [ ] go test ./internal/handler/... -race -cover → цель 80%+

⬜ Фаза 6 — Auth
   [ ] POST /api/v1/auth/register
   [ ] POST /api/v1/auth/login → JWT
   [ ] middleware/auth.go → Bearer token
   [ ] ErrUnauthorized (401), ErrForbidden (403)

⬜ Фаза 7 — Redis
   [ ] cache/transaction.go (Cache-Aside)
   [ ] middleware/rate_limit.go

⬜ Фаза 8 — Kafka
   [ ] broker/producer.go (publish после COMMIT)
   [ ] broker/consumer.go
   [ ] Kafka в docker-compose.yml

⬜ Фаза 9 — Observability
   [ ] metrics/metrics.go (Prometheus)
   [ ] middleware/metrics.go
   [ ] Prometheus + Grafana в docker-compose

⬜ Фаза 10 — Тесты для новых слоёв
   [ ] Cache тесты (Redis testcontainer)
   [ ] Broker тесты (Kafka testcontainer)
   [ ] Integration тесты (полный флоу)
   [ ] go test -race -cover ./... финальный

⬜ Фаза 11 — Деплой
   [ ] docker-compose.prod.yml
   [ ] entrypoint.sh (миграции перед стартом)
   [ ] VPS + Nginx + SSL
   [ ] GitHub Actions CI/CD
```
---

## 🏗️ Архитектура

### Структура папок
```
payment-service/
├── cmd/
│   └── main.go                    ← точка входа

├── internal/                      ← закрытый код
│   ├── domain/                    ← модели и бизнес ошибки
│   │   ├── payment.go
│   │   ├── client.go
│   │   ├── business.go
│   │   └── errors.go
│   │
│   ├── repository/                ← SQL слой (только БД)
│   │   ├── repository.go          ← интерфейсы
│   │   ├── payment.go
│   │   ├── client.go
│   │   └── business.go
│   │
│   ├── service/                   ← бизнес логика
│   │   ├── service.go             ← интерфейсы
│   │   ├── payment.go
│   │   ├── client.go
│   │   └── business.go
│   │
│   ├── handler/                   ← HTTP слой
│   │   ├── handler.go             ← роутер chi
│   │   ├── payment.go
│   │   ├── client.go
│   │   └── business.go
│   │
│   ├── middleware/
│   │   ├── logger.go              ← логирование запросов
│   │   ├── recovery.go            ← перехват паник
│   │   └── metrics.go             ← сбор метрик
│   │
│   ├── cache/
│   │   └── payment.go             ← Redis кэширование
│   │
│   ├── broker/
│   │   ├── producer.go            ← отправка событий в Kafka
│   │   └── consumer.go            ← получение событий
│   │
│   └── metrics/
│       └── metrics.go             ← Prometheus метрики

├── pkg/                           ← переиспользуемый код
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── logger/                    ← структурированный логгер
│   │   └── logger.go
│   └── apperrors/                 ← HTTP ошибки
│       └── errors.go

├── migrations/
├── tests/                         ← интеграционные тесты
│   └── payment_integration_test.go
├── .env
├── .env.example
├── .gitignore
├── .dockerignore
├── Dockerfile
├── docker-compose.yml
├── docker-compose.prod.yml        ← для продакшна с app
└── go.mod
```
### Принципы архитектуры
```
SOLID:
S — Single Responsibility: каждый слой делает одно
O — Open/Closed: новые фичи без изменения существующего
L — Liskov: mockRepo заменяет postgresRepo в тестах
I — Interface Segregation: маленькие интерфейсы
D — Dependency Inversion: зависимости через интерфейсы

Поток данных:
HTTP Request
    ↓
Middleware (logger → recovery → metrics)
    ↓
Handler (JSON, валидация, HTTP коды)
    ↓
Service (бизнес логика, расчёты)
    ↓
Cache (Redis — проверяем кэш)
    ↓
Repository (PostgreSQL)
    ↓
Broker (Kafka — публикуем событие)
    ↓
HTTP Response
```

## Ключевые архитектурные решения

```
— BIGINT для денег (тиыны), не DECIMAL — никаких ошибок округления
— owner_phone не UNIQUE — один owner → N бизнесов
— Интерфейсы репозиториев в service пакете (Go-идиома)
— Fail fast: валидация до IO запросов
— defer tx.Rollback() сразу после BeginTx
— domain ошибки через errors.Is цепочку
— middleware подключается в main.go, не в Routes()
```

## Правила тестов (выработанные)

```
— Отдельные функции когда кейсы проверяют разные ветки логики
— Table-driven когда структура одинакова, меняются только данные
— newXxxSvc/newXxxHandler helpers убирают дублирование
— Константы для testClientID, testAmount int64 — явный тип
— mock.Anything для *string параметров
— tt := tt обязательно при t.Parallel() в цикле
— makeTransactions(n): make([]T, 0, n) + append — не make([]T, n)
— AssertNotCalled через незарегистрированный EXPECT
— 500 не раскрывает текст внутренней ошибки клиенту
```

---

## 📝 Идеи на будущее

### Телефонные номера
```
Сейчас: E.164 (+77771156580), валидация длины
В будущем для мультистрановости:
  — countries (id, iso_code, name, phone_code, currency_code)
  — clients/businesses получают country_id FK
  — библиотека nyaruka/phonenumbers для реальной валидации
```

### Архитектура
```
⭐ owners таблица
   — owners (id, phone UNIQUE, name, created_at)
   — businesses получает owner_id FK, убрать owner_phone

⭐ bonus_transactions ledger
   — вместо просто bonus_balance
   — (id, client_id, transaction_id, amount, type, status, activated_at, expires_at)
   — "250 бонусов сгорит 15 апреля"

⭐ business_bonus_settings: activation_delay_sec, expires_after_sec
   — добавить когда будет bonus_transactions
```

### Надёжность
```
⭐ Idempotency key (X-Idempotency-Key header)
   — защита от двойной оплаты
   — хранить в Redis с TTL

⭐ Outbox pattern для Kafka
   — outbox таблица внутри транзакции
   — воркер публикует в Kafka
   — гарантированная доставка

⭐ pq.Error 23505 в business.Create
   — сейчас только в client.Create
```

### Продукт
```
⭐ Swagger / OpenAPI документация
⭐ Postman collection в репозитории
⭐ Webhook для бизнесов
⭐ История изменений bonus_settings
```

---

## Правила кода

```
Go:
— context первым параметром
— fmt.Errorf("layer.Method: %w", err)
— Аббревиатуры заглавными: GetByID, clientID, businessID
— Pointer receiver везде
— Интерфейсы там где потребляются (в service пакете)
— Константы вместо magic numbers
— Приватные методы для инкапсуляции логики

SQL:
— BIGINT для денег (не DECIMAL)
— Никогда SELECT *
— $1 $2 параметры, никогда fmt.Sprintf
— RETURNING после INSERT/UPDATE
— updated_at = NOW() в SQL
— WHERE deleted_at IS NULL
— bonus_balance >= $1 (не >)
— defer tx.Rollback() сразу после BeginTx

REST:
— URL во множественном числе (/clients, /businesses, /transactions)
— Существительные не глаголы (/cancel, не /cancelTransaction)
— Query params для фильтрации (?client_id=X)
— 201 только при создании, 200 для upsert и update
— UUID для всех публичных ID

Git:
— Conventional Commits: feat/fix/chore/refactor
— Feature ветки → PR → merge в main
```

краткие записи простого обывателя:
// TODO: слабый сервис^ надо рассширить и добавить авторизацию, бонусы за отзывы и рассылку на почту после добавления всех слоев
// примитивные вещи необходимые к добавлению: страна, город, код номера телефона, карта, вход, авторизацаия по воцапу, смс
// бизнес: филиал, сотрудники, категория и подкатегория, вход как сотрудник и проведение транзакции/ вход для хозяина + статистика, инфо о бизнесе(инста, номер филиала, отзывы, рейтинг и тд)
// настройки: скидка, бонус за отзывы, фиксированные/процент, кэшбэк надо отполировать позже все настройки после завершения настройки всех лэйэров, индивидуальные и условные
// транзакции: расширить: при создании выбирать филиал, сотрудника, скидку, бонус и тд, список транзакций: общий(от последних), по дням(итд), по клиентам, по рейтингу, по филиалам, по сотрудникам(много индексации...), рассчет статистики
// клиент: инфо о клиенте, список транзакций, доступные бонусы, перевод бонусов, мои отзывы/рецензии, возможность видеть на карте заведения с общим рейтингом и моим рейтингом поблизости и тд, набор активных индив скидок/бонусов/промок,
// надо вообще в сам сервис и таблице позже переосмыслить по хорошему более строго, но сначала настроим все слои

// веб или мобилка, шекспирский вопрос...
// не хочется даже трогать js так надоел, может просто попрактиковаться и написать пару игр на флаттере(?)
// посмотрим через сколько времени спустя я вернусь сюда
// пет проект выглядит в моих глазах обычно и нет изюминки, даже если пытаться начинать продавать в кордае есть сомнения, слишком заезженно,
// учитывая что я вернусь еще сюда, предлагаю после окончания написания всех лэйеров и тестого запуска, перед расширением САМОГО проекта до РЕАЛЬНОГО написать 5 игр и посмотреть может с рекламы будет больше денжка, а идею пет проекта еще раз обдумать более глубоко и с учетом незаезженности ниши...