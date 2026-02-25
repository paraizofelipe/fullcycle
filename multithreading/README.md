# CEP Race

Desafio Full Cycle de multithreading em Go.

Consulta simultaneamente duas APIs de CEP e exibe o resultado da que responder primeiro. Se nenhuma responder em 1 segundo, encerra com erro.

## APIs utilizadas

| API | Endpoint |
|-----|----------|
| BrasilAPI | `https://brasilapi.com.br/api/cep/v1/{cep}` |
| ViaCEP | `http://viacep.com.br/ws/{cep}/json/` |

## Como executar

```bash
go run main.go <cep>
```

**Exemplo:**

```bash
go run main.go 01310100
```

**Saída esperada:**

```
API: BrasilAPI
CEP:    01310-100
Rua:    Avenida Paulista
Bairro: Bela Vista
Cidade: São Paulo
Estado: SP
```

## Funcionamento

1. Duas goroutines são disparadas em paralelo, uma para cada API
2. Um `select` aguarda o primeiro resultado chegar no canal
3. Se o contexto expirar antes de qualquer resposta (timeout de 1s), o programa encerra com erro
