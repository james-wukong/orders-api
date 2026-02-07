## online ordering system for restaurants

## project structure

online-orders/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── bootstrap.go
│   │   ├── server.go
│   │   └── shutdown.go
│   │
│   ├── config/
│   │   ├── config.go
│   │   └── env.go
│   │
│   ├── domain/
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── errors.go
│   │   └── order/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       └── service.go
│   │
│   ├── infrastructure/
│   │   ├── postgres/
│   │   │   ├── client.go
│   │   │   ├── migrations/
│   │   │   └── user_repository.go
│   │   │
│   │   ├── redis/
│   │   │   ├── client.go
│   │   │   ├── cache.go
│   │   │   └── rate_limiter.go
│   │   │
│   │   └── logger/
│   │       └── zap.go
│   │
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── user_handler.go
│   │   │   │   └── order_handler.go
│   │   │   │
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   └── logging.go
│   │   │   │
│   │   │   └── router.go
│   │   │
│   │   └── grpc/
│   │       └── user_service.go
│   │
│   ├── usecase/
│   │   ├── user/
│   │   │   ├── create_user.go
│   │   │   ├── login_user.go
│   │   │   └── get_user.go
│   │   │
│   │   └── order/
│   │       └── create_order.go
│   │
│   ├── pkg/
│   │   ├── auth/
│   │   │   └── jwt.go
│   │   ├── crypto/
│   │   │   └── password.go
│   │   ├── validator/
│   │   │   └── validator.go
│   │   └── utils/
│   │       └── uuid.go
│   │
│   └── tests/
│       ├── integration/
│       └── unit/
│
├── migrations/
│   └── 0001_init.sql
│
├── scripts/
│   ├── migrate.sh
│   └── seed.sh
│
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       └── deployment.yaml
│
├── go.mod
├── go.sum
└── README.md

### cmd/ - entry points

```sh
cmd/api/main.go
```

- Parses env vars

- Starts HTTP/GRPC server

- Calls internal/app/bootstrap.go

### internal/app/ - application lifecycle

```sh
bootstrap.go   // wiring dependencies
server.go     // HTTP/GRPC server
shutdown.go   // graceful shutdown
```

### internal/config/ - configuration management

```go
type Config struct {
    Postgres PostgresConfig
    Redis    RedisConfig
}
```

- Loads .env

- Handles defaults

- Centralized config logic

### internal/domain/ - business rules

```go
type User struct {
    ID    uuid.UUID
    Email string
}
```

NO database, NO Redis, NO HTTP here. core business logic

- Entities

- Interfaces (repositories)

- Domain errors

### internal/usecase/ - application actions

```go
CreateUser
LoginUser
GetUser
```

Each file: 

- Executes ONE business action

- Coordinates domain + infrastructure

- Think: verbs, not nouns

### internal/infrastructure/ - External systems

```sh
postgres/client.go
postgres/user_repository.go
```

- SQL, GORM, pgx

- Implements domain repository interfaces

#### Redis

```sh
redis/cache.go
redis/rate_limiter.go
```

- Cache

- Locks

- Rate limits

- Session store

Only infra knows Redis/Postgres exist.

### internal/interfaces/ - Adapters

HTTP

```sh
handlers/
middleware/
router.go
```

- Converts HTTP → usecase

- Converts usecase → HTTP response

gRPC (optional)

- Same business logic, different transport

### internal/pkg/ - shared libraries

- JWT

- Password hashing

- Validators

- Utilities

Reusable across services.

### migrations/ - Database versioning

```sh
0001_init.sql
0002_add_users.sql
```

- golang-migrate

- atlas

- goose

### deployments/ - Infrastructure-as-code

- Docker

- Docker Compose

- Kubernetes

## work flow

HTTP Request
   ↓
Handler
   ↓
Usecase
   ↓
Domain logic
   ↓
Postgres / Redis
   ↓
Usecase
   ↓
Handler
   ↓
HTTP Response

