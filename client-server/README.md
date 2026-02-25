# Full Cycle Desafio Client/Server

Projeto em Go que expõe um endpoint HTTP para consultar a cotação do dólar (USD/BRL), persistir no SQLite e permitir que um cliente grave a cotação atual em arquivo.

## Como executar

1) Suba o servidor HTTP:

```bash
go run cmd/server/main.go
```

1) Em outro terminal, execute o cliente:

```bash
go run cmd/client/main.go
```

Saídas geradas:

- `exchange_rates.db`: banco SQLite criado pelo servidor
- `exchange_rates.txt`: arquivo criado pelo cliente com o valor atual

## Estrutura e proposito dos pacotes

- `cmd/server`: inicializa o servidor HTTP e registra a rota `/exchange_rates`.
- `cmd/client`: chama o servidor e grava a cotação em `exchange_rates.txt`.
- `internal/exchange_rates`: integra com a AwesomeAPI para buscar a cotação USD/BRL.
- `internal/storage`: persiste a cotação no SQLite.
- `internal/client`: cliente HTTP para consumir o endpoint do servidor.
- `internal/ctxlog`: utilitário para registrar timeout de context.

## Dependências

Dependência direta:

- `modernc.org/sqlite` (driver SQLite puro Go)

Dependências indiretas:

- São listadas em `go.mod` como `// indirect` (necessarias para o driver SQLite).
