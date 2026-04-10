# Economizze-API

API REST de gestão financeira pessoal construída em Go com clean architecture, DDD e CQRS.

---

## Sumário

- [Visão Geral](#visão-geral)
- [Arquitetura](#arquitetura)
- [Tecnologias](#tecnologias)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Domínio](#domínio)
- [Pré-requisitos](#pré-requisitos)
- [Instalação e Execução](#instalação-e-execução)
- [Migrations](#migrations)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Endpoints da API](#endpoints-da-api)
- [Autenticação](#autenticação)
- [Testes](#testes)
- [Observabilidade](#observabilidade)
- [Decisões Técnicas](#decisões-técnicas)

---

## Visão Geral

Sistema de controle financeiro pessoal que permite ao usuário gerenciar contas bancárias, registrar receitas e despesas, definir orçamentos por categoria e acompanhar relatórios mensais.

**Funcionalidades principais:**

- Múltiplas contas (corrente, poupança, carteira, cartão de crédito, investimentos)
- Lançamentos de receita e despesa com categorização
- Transferências entre contas com rastreabilidade
- Orçamentos mensais por categoria com alertas de extrapolação
- Transações recorrentes (salário, aluguel, assinaturas)
- Relatórios: resumo mensal, fluxo de caixa, progresso de orçamentos
- Autenticação via JWT com refresh tokens e OAuth2 (Google)

---

## Arquitetura

O projeto segue **Clean Architecture** combinada com **DDD (Domain-Driven Design)** e **CQRS** no lado de leitura.

```
┌─────────────────────────────────────────────────────┐
│                    HTTP (Gin)                        │  ← Driving Adapter
├─────────────────────────────────────────────────────┤
│                   Use Cases                          │  ← Application Layer
├──────────────┬──────────────────────────────────────┤
│    Domain    │           Ports (interfaces)          │  ← Domain Layer
├──────────────┴──────────────────────────────────────┤
│         Repository │ Database │ Cache │ Email        │  ← Driven Adapters
└─────────────────────────────────────────────────────┘
```

**Regra de dependência:** as setas apontam sempre para dentro. O domínio não importa nada da infraestrutura.

**CQRS:** comandos (escrita) passam pelo domínio com todas as validações e domain events. Queries (leitura) vão direto ao banco com SQL otimizado, sem reconstruir aggregates.

---

## Tecnologias

| Categoria         | Tecnologia                        |
|-------------------|-----------------------------------|
| Linguagem         | Go 1.23                           |
| HTTP Framework    | Gin                               |
| ORM               | GORM                              |
| Banco de Dados    | PostgreSQL 16                     |
| Cache             | Redis 7                           |
| Migrations        | golang-migrate                    |
| Autenticação      | JWT (golang-jwt) + OAuth2 Google  |
| Logging           | Zap (uber-go/zap)                 |
| Métricas          | Prometheus                        |
| Tracing           | OpenTelemetry + Jaeger            |
| Testes            | testify + testcontainers + mockery |
| Containerização   | Docker + Docker Compose           |

---

## Estrutura do Projeto

```
finance-api/
├── cmd/
│   └── api/
│       └── main.go               # bootstrap: DI manual, server, graceful shutdown
│
├── internal/
│   ├── domain/                   # DDD: núcleo do negócio (zero dependências externas)
│   │   ├── account.go            # Aggregate Root: Account
│   │   ├── transaction.go        # Entity: Transaction
│   │   ├── budget.go             # Aggregate Root: Budget
│   │   ├── category.go           # Entity: Category
│   │   ├── recurring.go          # Aggregate Root: RecurringTransaction
│   │   ├── money.go              # Value Object: Money (centavos, imutável)
│   │   ├── period.go             # Value Object: Period (intervalo de datas)
│   │   ├── events.go             # Domain Events: TransactionCreated, BudgetExceeded...
│   │   └── errors.go             # Sentinel errors: ErrNotFound, ErrConflict...
│   │
│   ├── ports/                    # Interfaces (Ports hexagonais)
│   │   ├── account_repository.go
│   │   ├── transaction_repository.go
│   │   ├── budget_repository.go
│   │   ├── category_repository.go
│   │   ├── event_publisher.go
│   │   └── unit_of_work.go       # Coordena transações atômicas entre repositórios
│   │
│   ├── usecase/                  # Application Layer: orquestra domínio + ports
│   │   ├── account/              # CreateAccount, UpdateAccount, DeleteAccount, ListAccounts
│   │   ├── transaction/          # RecordExpense, RecordIncome, Transfer, DeleteTransaction
│   │   ├── budget/               # CreateBudget, UpdateBudget, DeleteBudget
│   │   └── category/             # CreateCategory, ListCategories
│   │
│   ├── query/                    # CQRS: read side — queries sem passar pelo domínio
│   │   ├── summary_query.go      # Resumo mensal (receitas, despesas, saldo)
│   │   ├── cashflow_query.go     # Fluxo de caixa diário
│   │   └── budget_progress_query.go  # Progresso dos orçamentos do mês
│   │
│   └── infra/                    # Adapters: implementações concretas
│       ├── http/
│       │   ├── handler/          # Handlers Gin por recurso
│       │   ├── middleware/       # Auth JWT, Logger, Recovery, CORS, RateLimit
│       │   ├── dto/              # Request/Response types (separados do domínio)
│       │   └── router.go         # Definição de rotas e grupos
│       ├── repository/           # Implementações GORM + modelos de persistência
│       │   ├── models.go         # Structs com tags GORM (separadas das entidades)
│       │   ├── converters.go     # Converte model ↔ entidade de domínio
│       │   └── unit_of_work.go   # Implementação do UoW com transação GORM
│       ├── database/
│       │   ├── postgres.go       # Conexão, pool de conexões, runner de migrations
│       │   └── seed.go           # Categorias padrão do sistema
│       ├── auth/
│       │   ├── jwt_service.go    # Geração e validação de tokens
│       │   ├── bcrypt.go         # Hash de senhas
│       │   └── oauth2.go         # Integração Google OAuth2
│       ├── cache/
│       │   └── redis.go          # Token store para refresh tokens
│       ├── events/
│       │   └── publisher.go      # Publicação de domain events (in-memory)
│       └── logger/
│           └── logger.go         # Configuração do zap + context propagation
│
├── migrations/                   # Arquivos SQL versionados
│   ├── 000001_create_extensions.up.sql
│   ├── 000002_create_users.up.sql
│   ├── 000003_create_categories.up.sql
│   ├── 000004_create_accounts.up.sql
│   ├── 000005_create_transactions.up.sql
│   ├── 000006_create_budgets.up.sql
│   ├── 000007_create_recurring_transactions.up.sql
│   ├── 000008_create_indexes.up.sql
│   └── ...*.down.sql
│
├── config/
│   ├── prometheus.yml
│   └── otel-collector.yaml
│
├── docker-compose.yml
├── docker-compose.observability.yml
├── Makefile
├── .env.example
└── go.mod
```

---

## Domínio

### Aggregates e Value Objects

```
Account (Aggregate Root)
  ├── balance: Money          ← Value Object (centavos + moeda, imutável)
  └── transactions: []Transaction  ← Entities

Budget (Aggregate Root)
  ├── limit: Money
  ├── spent: Money
  └── period: Period          ← Value Object (start + end, imutável)

RecurringTransaction (Aggregate Root)
  └── frequency: Enum (daily/weekly/monthly/yearly)
```

### Domain Events

| Evento                  | Disparado quando                              |
|-------------------------|-----------------------------------------------|
| `TransactionCreated`    | Receita ou despesa registrada                 |
| `AccountBalanceUpdated` | Saldo da conta muda após transação            |
| `BudgetExceeded`        | Gasto do mês ultrapassa o limite do orçamento |

### Regras de negócio no domínio

- Conta corrente e poupança não permitem saldo negativo
- Orçamento excedido dispara evento apenas na primeira extrapolação
- Transferência cria duas transações vinculadas por `transfer_peer_id`
- Conta com saldo positivo não pode ser deletada
- Limite de 10 contas por usuário

---

## Pré-requisitos

- Go 1.23+
- Docker e Docker Compose
- `golang-migrate` CLI (para migrations manuais)
- `mockery` (para gerar mocks em desenvolvimento)

```bash
# Instala as ferramentas necessárias
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/vektra/mockery/v2@latest
```

---

## Instalação e Execução

```bash
# 1. Clone o repositório
git clone https://github.com/seuuser/finance-api.git
cd finance-api

# 2. Copie e configure o arquivo de ambiente
cp .env.example .env
# Edite .env com suas configurações

# 3. Suba os serviços de infraestrutura
make db-up

# 4. Execute as migrations e o seed de categorias
make migrate-up

# 5. Inicie a API
make run

# A API estará disponível em http://localhost:8080
# Documentação Swagger: http://localhost:8080/swagger/index.html
```

### Com Docker Compose completo

```bash
# Sobe tudo: API + banco + redis + observabilidade
docker compose up -d

# Logs da API
docker compose logs -f api
```

---

## Migrations

```bash
# Cria uma nova migration
make migrate-new name=add_column_to_accounts

# Aplica todas as migrations pendentes
make migrate-up

# Reverte a última migration
make migrate-down

# Reverte N migrations
make migrate-down-n n=3

# Versão atual do schema
make migrate-version

# Força versão (após estado dirty)
make migrate-force v=5
```

As migrations ficam em `migrations/` na raiz do projeto. O binário as embute via `//go:embed` no `cmd/api/main.go`, portanto não é necessário montar volumes em produção.

---

## Variáveis de Ambiente

```env
# Servidor
APP_ENV=development          # development | production
APP_PORT=8080
APP_NAME=finance-api
APP_VERSION=1.0.0

# Banco de dados
DATABASE_URL=postgres://finance:finance@localhost:5432/finance_dev?sslmode=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_ACCESS_SECRET=sua-chave-secreta-aqui
JWT_REFRESH_SECRET=outra-chave-secreta-aqui
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h     # 7 dias

# OAuth2 Google
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# Logging
LOG_LEVEL=info               # debug | info | warn | error
LOG_FORMAT=json              # json | console

# Observabilidade
OTLP_ENDPOINT=http://localhost:4318
TRACE_SAMPLE_RATE=0.1        # 10% das requisições em produção
```

---

## Endpoints da API

### Autenticação

| Método | Rota                       | Descrição                          |
|--------|----------------------------|------------------------------------|
| POST   | `/auth/register`           | Cadastro de novo usuário           |
| POST   | `/auth/login`              | Login com email e senha            |
| POST   | `/auth/refresh`            | Renova o access token              |
| POST   | `/auth/logout`             | Logout do dispositivo atual        |
| POST   | `/auth/logout-all`         | Logout de todos os dispositivos    |
| GET    | `/auth/google`             | Inicia fluxo OAuth2 Google         |
| GET    | `/auth/google/callback`    | Callback OAuth2 Google             |

### Contas

| Método | Rota                       | Descrição                          |
|--------|----------------------------|------------------------------------|
| GET    | `/api/v1/accounts`         | Lista todas as contas do usuário   |
| POST   | `/api/v1/accounts`         | Cria nova conta                    |
| GET    | `/api/v1/accounts/:id`     | Detalhes de uma conta              |
| PUT    | `/api/v1/accounts/:id`     | Atualiza nome, cor ou padrão       |
| DELETE | `/api/v1/accounts/:id`     | Desativa a conta                   |

### Transações

| Método | Rota                            | Descrição                        |
|--------|---------------------------------|----------------------------------|
| GET    | `/api/v1/transactions`          | Lista transações com filtros     |
| POST   | `/api/v1/transactions/income`   | Registra receita                 |
| POST   | `/api/v1/transactions/expense`  | Registra despesa                 |
| POST   | `/api/v1/transactions/transfer` | Transferência entre contas       |
| GET    | `/api/v1/transactions/:id`      | Detalhes de uma transação        |
| DELETE | `/api/v1/transactions/:id`      | Remove e reverte o lançamento    |

### Orçamentos

| Método | Rota                       | Descrição                          |
|--------|----------------------------|------------------------------------|
| GET    | `/api/v1/budgets`          | Lista orçamentos do mês            |
| POST   | `/api/v1/budgets`          | Cria orçamento para categoria/mês  |
| PUT    | `/api/v1/budgets/:id`      | Atualiza limite do orçamento       |
| DELETE | `/api/v1/budgets/:id`      | Remove o orçamento                 |

### Categorias

| Método | Rota                       | Descrição                          |
|--------|----------------------------|------------------------------------|
| GET    | `/api/v1/categories`       | Lista categorias (sistema + custom)|
| POST   | `/api/v1/categories`       | Cria categoria personalizada       |
| DELETE | `/api/v1/categories/:id`   | Desativa categoria personalizada   |

### Relatórios (CQRS — read side)

| Método | Rota                            | Descrição                        |
|--------|---------------------------------|----------------------------------|
| GET    | `/api/v1/reports/summary`       | Resumo do mês (receitas/despesas)|
| GET    | `/api/v1/reports/cashflow`      | Fluxo de caixa diário            |
| GET    | `/api/v1/reports/budget-progress` | Progresso dos orçamentos       |

### Sistema

| Método | Rota        | Descrição                              |
|--------|-------------|----------------------------------------|
| GET    | `/health`   | Liveness probe                         |
| GET    | `/ready`    | Readiness probe (verifica DB e Redis)  |
| GET    | `/metrics`  | Métricas Prometheus                    |

---

## Autenticação

Todas as rotas `/api/v1/*` exigem o header:

```
Authorization: Bearer <access_token>
```

O `access_token` expira em 15 minutos. Use `/auth/refresh` com o cookie `refresh_token` para obter um novo par de tokens. O refresh token é armazenado em cookie `HttpOnly` para proteger contra XSS.

**Fluxo completo:**

```
POST /auth/login
  → access_token (15min) + cookie refresh_token (7 dias)

GET /api/v1/accounts
  → Authorization: Bearer <access_token>

POST /auth/refresh          ← quando access_token expirar
  → novo access_token + novo refresh_token (rotação)
```

---

## Testes

```bash
# Todos os testes
make test

# Só unit tests (sem I/O, rápidos)
make test-unit

# Só integration tests (sobe PostgreSQL via testcontainers)
make test-integration

# Com relatório de coverage
make test-coverage
# Abre coverage.html no browser

# Race detector
make test-race

# Benchmarks
make bench
```

**Cobertura mínima:** 80% (verificada no CI).

A camada de domínio é testada com **table-driven tests** sem nenhuma dependência externa. Os repositórios são testados com **testcontainers** (banco PostgreSQL real, efêmero). Os use cases são testados com **mocks gerados pelo mockery**.

---

## Observabilidade

### Logs

Logs estruturados em JSON via `zap`. Cada request recebe um `request_id` único. O `trace_id` do OpenTelemetry é injetado automaticamente em todos os logs do mesmo request.

```bash
# Logs formatados no terminal (desenvolvimento)
make logs-pretty

# Filtrar só erros
make logs-errors

# Seguir um request específico
make logs-request id=<request_id>
```

### Métricas (Prometheus)

Disponíveis em `GET /metrics`. Grafana pré-configurado em `http://localhost:3000`.

Métricas principais:
- `http_requests_total` — total por método, rota e status
- `http_request_duration_seconds` — histograma de latência
- `orders_placed_total` — total de lançamentos por tipo
- `db_query_duration_seconds` — latência das queries

### Tracing (OpenTelemetry + Jaeger)

Interface do Jaeger em `http://localhost:16686`.

Spans automáticos para todas as requests HTTP e queries SQL. Spans manuais nos use cases para operações críticas.

### Stack local completa

```bash
# Sobe Prometheus + Grafana + Jaeger + Loki
docker compose -f docker-compose.observability.yml up -d

# Interfaces:
#   Grafana:    http://localhost:3000
#   Prometheus: http://localhost:9090
#   Jaeger:     http://localhost:16686
```

---

## Decisões Técnicas

### Por que DDD neste projeto?

O domínio financeiro tem regras não-triviais: `Account` precisa impedir saldo negativo, `Budget` precisa detectar a primeira extrapolação e disparar evento, transferências criam dois registros vinculados atomicamente. Essas regras pertencem ao domínio, não ao banco de dados nem à camada HTTP. DDD força essa separação.

### Por que CQRS parcial?

O modelo ideal para criar uma transação (aggregate com validações, events, UoW) é ruim para ler relatórios (JOINs, GROUP BY, window functions). CQRS resolve isso sem duplicar banco: comandos passam pelo aggregate, queries vão direto ao banco com SQL otimizado.

### Por que centavos em vez de decimal?

`float64` tem erros de arredondamento (`0.1 + 0.2 = 0.30000000000001`). `DECIMAL` no banco é preciso mas mais lento. `int64` em centavos é rápido, exato, trivial de serializar e suporta até ~92 trilhões de centavos.

### Por que Unit of Work?

`RecordExpense` precisa atualizar saldo da conta, salvar a transação e atualizar o orçamento atomicamente. Sem UoW, cada repositório tem sua própria conexão e uma falha no meio deixa o banco inconsistente. O UoW injeta a mesma `*gorm.DB` de transação em todos os repositórios dentro do closure.

### Por que embed.FS para migrations?

O binário Go carrega as migrations consigo. O container Docker é autocontido — sem volume montado, sem dependência de filesystem externo. Em testes, passamos `os.DirFS("../../../migrations")` com o path relativo correto.

---

## Licença

MIT
