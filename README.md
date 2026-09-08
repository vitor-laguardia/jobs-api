# Jobs API

REST API para gerenciamento de jobs e usuários, escrita em Go usando **apenas a standard library** — sem frameworks HTTP, sem ORM.

Projeto construído com foco em **arquitetura limpa**, **testes automatizados** e boas práticas idiomáticas de Go, como exercício prático de transição para desenvolvimento backend.

## Principais características

- **Arquitetura em camadas**: handler → service → repository → domain, com separação clara de responsabilidades
- **Organização por domínio** (package by feature): cada entidade (`job`, `user`) isolada em seu próprio pacote
- **Roteamento nativo** com `http.ServeMux` (Go 1.22+), sem router externo
- **Tratamento de erros idiomático**: erros sentinela (`errors.Is`) e erros customizados (`errors.As`) para diferenciar falhas de domínio, validação e infraestrutura
- **Validação de entrada** via DTOs desacoplados das entidades de domínio, com respostas de erro estruturadas (422) para múltiplos campos inválidos
- **Máquina de estados** para o ciclo de vida dos jobs (`pending → running → done/failed`), com transições validadas no domínio
- **Testes automatizados**: TDD, testes de integração via `httptest`, testes table-driven e cobertura de casos de erro (payload malformado, limites de tamanho, content-type, etc.)

## Stack

Go · `net/http` · `database/sql` (em progresso) · PostgreSQL

## Roadmap

- [ ] Persistência com PostgreSQL via `database/sql` puro (sem ORM/query builder), substituindo o repositório em memória
- [ ] Migrations de schema
- [ ] Processamento assíncrono de jobs com goroutines (execução de uma chamada HTTP real por job, com transição de status refletindo sucesso/falha)
- [ ] Evolução para worker pool (channel + N workers) conforme necessidade de controle de concorrência

## Rodando localmente

```bash
go run cmd/api/main.go
```

## Testes

```bash
go test ./...
```
