# stress-test

CLI para testes de carga em serviços web, escrita em Go com [Cobra](https://github.com/spf13/cobra) e [Viper](https://github.com/spf13/viper).

## Parâmetros

| Flag            | Descrição                             | Padrão |
|-----------------|---------------------------------------|--------|
| `--url`         | URL do serviço a ser testado          | —      |
| `--requests`    | Número total de requisições           | 100    |
| `--concurrency` | Número de chamadas simultâneas        | 1      |

## Executando localmente

```bash
cd stress-test
go run . --url=http://google.com --requests=100 --concurrency=10
```

## Build da imagem Docker

```bash
cd stress-test
docker build -t stress-test .
```

## Executando via Docker

```bash
docker run stress-test --url=http://httpbin.org/ --requests=1000 --concurrency=10
```

## Exemplo de saída

```
Iniciando stress test...
  URL:          http://httpbin.org/
  Requisições:  1000
  Concorrência: 10

========== Relatório ==========
Tempo total:            4.321s
Total de requisições:   1000
Respostas HTTP 200:     873

Distribuição de outros status:
  HTTP 301: 127
================================
```
