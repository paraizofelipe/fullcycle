# Full Cycle — Desafios de Go

Repositório com os desafios práticos do curso de Go da plataforma [Full Cycle](https://fullcycle.com.br/). Cada subdiretório é um módulo Go independente correspondente a um desafio.

## Desafios

### 1. `client-server/` — Client-Server com Contextos e Persistência

**Objetivo:** Implementar um sistema cliente-servidor que consulta a cotação do dólar (USD/BRL), persiste os dados e os expõe via HTTP, aplicando corretamente `context.WithTimeout` em cada camada da aplicação.

**Como funciona:**

- O **servidor** expõe o endpoint `GET /cotacao` na porta `8080`
- A cada requisição, busca a cotação atual na [AwesomeAPI](https://docs.awesomeapi.com.br/api-de-moedas), persiste o valor em um banco SQLite e retorna o bid (valor de compra) em JSON
- O **cliente** chama o servidor, recebe a cotação e grava o valor no arquivo `exchange_rates.txt`

**Timeouts configurados:**
| Operação | Timeout |
|---|---|
| Busca na AwesomeAPI | 200ms |
| Escrita no SQLite | 10ms |
| Requisição do cliente ao servidor | 300ms |

**Como executar:**
```bash
# Terminal 1 — subir o servidor
cd client-server && go run cmd/server/main.go

# Terminal 2 — executar o cliente
cd client-server && go run cmd/client/main.go
```

---

### 2. `multithreading/` — Multithreading e Corrida entre APIs

**Objetivo:** Utilizar goroutines e channels para consultar dois provedores de CEP simultaneamente e retornar o resultado da API que responder primeiro, descartando a mais lenta.

**Como funciona:**

- Dispara duas goroutines em paralelo, uma para cada API:
  - [BrasilAPI](https://brasilapi.com.br/) — `https://brasilapi.com.br/api/cep/v1/{cep}`
  - [ViaCEP](https://viacep.com.br/) — `http://viacep.com.br/ws/{cep}/json/`
- Usa um `select` sobre um channel com buffer para capturar a primeira resposta
- Se nenhuma API responder em **1 segundo**, exibe erro de timeout e encerra

**Como executar:**
```bash
cd multithreading && go run main.go <cep>

# Exemplo:
go run main.go 01310100
```

**Exemplo de saída:**
```
API: BrasilAPI
CEP:    01310-100
Rua:    Avenida Paulista
Bairro: Bela Vista
Cidade: São Paulo
Estado: SP
```

---

## Estrutura do Repositório

```
fullcycle/
├── client-server/          # Desafio 1 — cliente, servidor, SQLite e timeouts
│   ├── cmd/
│   │   ├── server/         # Entrypoint do servidor HTTP
│   │   └── client/         # Entrypoint do cliente
│   └── internal/
│       ├── exchange_rates/ # Integração com a AwesomeAPI
│       ├── storage/        # Persistência no SQLite
│       ├── client/         # HTTP client para consumir o servidor
│       └── ctxlog/         # Helper para log de erros de contexto
└── multithreading/         # Desafio 2 — corrida entre BrasilAPI e ViaCEP
    └── main.go
```

## Requisitos

- Go 1.21+
- Conexão com a internet (as APIs são externas)

## Gerenciamento de Módulos

Cada projeto possui seu próprio `go.mod` e deve ser gerenciado de forma independente:

```bash
cd <projeto> && go mod tidy   # instalar dependências
cd <projeto> && go build ./... # compilar
cd <projeto> && go vet ./...   # verificar o código
```
