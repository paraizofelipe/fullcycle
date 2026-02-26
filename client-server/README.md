# Full Cycle Desafio Client/Server

Projeto em Go que expõe um endpoint HTTP para consultar a cotação do dólar (USD/BRL), persistir no SQLite e permitir que um cliente grave a cotação atual em arquivo.

## Pré-requisitos

- [Go 1.21+](https://golang.org/dl/)

## Como executar

**1. Inicie o servidor** (deixe rodando em um terminal):

```bash
cd client-server
go run cmd/server/main.go
```

Saída esperada:

```
2024/01/01 00:00:00 listening on :8080
```

**2. Em outro terminal, execute o cliente:**

```bash
cd client-server
go run cmd/client/main.go
```

Saída esperada:

```
2024/01/01 00:00:00 exchange rate received: bid=5.7423
```

## Resultados gerados

| Arquivo | Gerado por | Conteúdo |
|---------|-----------|----------|
| `cotacao.txt` | cliente | `Dólar: 5.7423` |
| `exchange_rates.db` | servidor | banco SQLite com o histórico de cotações |

## Como testar manualmente

Com o servidor rodando, você pode consultar o endpoint diretamente:

```bash
curl http://localhost:8080/cotacao
```

Resposta:

```json
{"bid":"5.7423"}
```

## Estrutura dos pacotes

- `cmd/server`: servidor HTTP na porta 8080, rota `/cotacao`
- `cmd/client`: chama o servidor e grava a cotação em `cotacao.txt`
- `internal/exchange_rates`: integra com a [AwesomeAPI](https://economia.awesomeapi.com.br/json/last/USD-BRL) (timeout: 200ms)
- `internal/storage`: persiste a cotação no SQLite (timeout: 10ms)
- `internal/client`: cliente HTTP para consumir o servidor (timeout: 300ms)
- `internal/ctxlog`: registra erros de timeout de context nos logs
