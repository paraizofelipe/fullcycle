# Rate Limiter

Middleware de controle de fluxo de requisições HTTP em Go, com suporte a limitação por IP e por token de acesso. Persiste contadores e bloqueios no Redis.

## Funcionamento

O middleware inspeciona cada requisição recebida e decide se ela deve ser permitida ou bloqueada:

1. **Limitação por Token** (precedência): se o header `API_KEY` estiver presente, o token é usado como chave de controle. Tokens podem ter limites individuais configurados.
2. **Limitação por IP**: quando não há token, o IP do cliente é usado. O IP real é extraído de `X-Forwarded-For`, `X-Real-IP` ou `RemoteAddr`, nessa ordem.
3. **Bloqueio**: ao exceder o limite, a chave (IP ou token) fica bloqueada pelo tempo configurado. Durante o bloqueio, todas as requisições são rejeitadas com `HTTP 429`.

### Resposta de bloqueio

```
HTTP/1.1 429 Too Many Requests
Content-Type: text/plain

you have reached the maximum number of requests or actions allowed within a certain time frame
```

## Configuração

Todas as configurações são feitas via variáveis de ambiente (ou arquivo `.env` na raiz do projeto):

| Variável                | Padrão  | Descrição                                                                 |
|-------------------------|---------|---------------------------------------------------------------------------|
| `REDIS_HOST`            | `localhost` | Host do Redis                                                         |
| `REDIS_PORT`            | `6379`  | Porta do Redis                                                            |
| `REDIS_PASSWORD`        | _(vazio)_ | Senha do Redis (se configurada)                                         |
| `IP_RATE_LIMIT`         | `10`    | Máximo de requisições por segundo por IP                                  |
| `TOKEN_RATE_LIMIT`      | `100`   | Máximo de requisições por segundo por token (limite padrão)               |
| `BLOCK_DURATION_SECONDS`| `300`   | Tempo de bloqueio em segundos após exceder o limite                       |
| `TOKENS`                | _(vazio)_ | Limites individuais por token: `"tokenA:200,tokenB:50"`                |

### Exemplo de `.env`

```env
IP_RATE_LIMIT=10
TOKEN_RATE_LIMIT=100
BLOCK_DURATION_SECONDS=300
TOKENS=vip-token:500,admin-token:1000
```

## Executando com Docker Compose

```bash
# Subir a aplicação e o Redis
docker compose up --build

# Testar
curl http://localhost:8080/

# Testar com token
curl -H "API_KEY: meu-token" http://localhost:8080/
```

## Executando localmente

```bash
cd rate-limiter

# Baixar dependências
go mod tidy

# Subir apenas o Redis (necessário)
docker run -d -p 6379:6379 redis:7-alpine

# Iniciar o servidor
go run ./cmd/server
```

## Rodando os testes

Os testes utilizam o `MemoryStorage` (sem Redis), portanto não precisam de infraestrutura externa:

```bash
cd rate-limiter
go test ./...
```

Para rodar via Docker Compose (recomendado para avaliação):

```bash
docker compose run --rm app go test ./...
```

## Arquitetura

```
rate-limiter/
├── cmd/server/             # Entrypoint HTTP
├── config/                 # Carregamento de variáveis de ambiente
└── internal/
    ├── limiter/            # Lógica de negócio do rate limiter
    ├── middleware/         # Middleware HTTP
    └── storage/            # Camada de persistência (Strategy Pattern)
        ├── storage.go      # Interface Strategy
        ├── redis.go        # Implementação Redis (produção)
        └── memory.go       # Implementação em memória (testes)
```

### Design Pattern: Strategy

A interface `storage.Storage` define o contrato de persistência:

```go
type Storage interface {
    Increment(ctx context.Context, key string, window time.Duration) (int64, error)
    IsBlocked(ctx context.Context, key string) (bool, error)
    Block(ctx context.Context, key string, duration time.Duration) error
}
```

Para substituir o Redis por outro backend (ex: Memcached, PostgreSQL), basta:
1. Criar uma nova struct que implemente `Storage`
2. Passar a nova implementação ao construtor `limiter.New(...)`

Nenhuma outra parte do código precisa ser alterada.
