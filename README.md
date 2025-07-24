# Projeto OTel na Prática

Este é o projeto que utilizamos na Especialização em OpenTelemetry no [Dose de Telemetria](https://dosedetelemetria.com). Aqui temos uma aplicação relativamente simples, mas utilizando diversos aspectos de aplicações normais, como conexões HTTP e gRPC entre si, comunicação com banco de dados, envio e recebimento de mensagens via mensageria (message queue).

A aplicação não possui nenhuma instrumentação. Nada. Durante a especialização, vamos utilizar a aplicação para aprender diversos aspectos de observabilidade, com foco em OTel.

---

## **Sumário**

- [Módulos Disponíveis](#módulos-disponíveis)
- [Configuração](#configuração)
- [Como as coisas funcionam](#como-as-coisas-funcionam)
- [Contribuindo](#contribuindo)
- [Licença](#licença)

---

## Módulos Disponíveis

- **`cmd/users`**:
  - **Descrição**: Este módulo contém a aplicação principal para gerenciar usuários. Ele lida com operações como criação, atualização e exclusão de usuários.

- **`cmd/payments`**:
  - **Descrição**: Este módulo é responsável pelo processamento de pagamentos. Ele gerencia transações financeiras e integrações com gateways de pagamento. Ao receber uma requisição para um novo pagamento, coloca a requisição em uma fila de mensagens. Uma rotina na mesma aplicação recebe a mensagem e processa o pagamento, armazenando em um banco de dados SQLLite.

- **`cmd/all-in-one`**:
  - **Descrição**: Este módulo combina todas as funcionalidades em uma única aplicação. Ele é útil para desenvolvimento e testes locais, permitindo executar todos os serviços em um único processo.

- **`cmd/plans`**:
  - **Descrição**: Este módulo gerencia os planos de assinatura disponíveis. Ele lida com a criação, atualização e exclusão de planos. Aceita requisições tanto em HTTP quanto gRPC.

- **`cmd/subscriptions`**:
  - **Descrição**: Este módulo gerencia as assinaturas dos usuários aos planos. Ele lida com a criação, atualização e cancelamento de assinaturas.

---

## Configuração

Por padrão, um arquivo de configuração não é necessário, especialmente ao rodar o "all-in-one". Ao fazer a aplicação rodar separadamente, a maioria dos serviços vai precisar de um arquivo de configuração específico, que segue o seguinte formato:

```yaml
# yaml-language-server: $schema=./config-schema.yaml
payments:
  subscriptions_endpoint: http://localhost:8080/subscriptions
  sqlite:
    dsn: file::memory:?cache=shared
  nats:
    endpoint: nats://localhost:4222
    subject: payment.process
    stream: payments
    consumer_name: payments

subscriptions:
  users_endpoint: http://localhost:8080/users
  plans_endpoint: http://localhost:8080/plans

plans: {}

users: {}

server:
  endpoint:
    grpc: :8081
    http: :8080
```

---

## Como as coisas funcionam

* Os serviços "plans" e "users" não tem dependências com outros serviços. O serviço "subscriptions" precisa fazer conexões com "plans" e "users", enquanto que "payments" faz uma conexão com "subscriptions".

---

## Instalação das ferramentas

### NATS

Instalação do `nats-server` e `nats`

#### macOS via Homebrew
```
## nats-server
brew install nats-server

## NATS Command Line Interface
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```

## Contribuindo

Quer ajudar a melhorar este projeto? Veja como começar no arquivo [CONTRIBUTING.md](CONTRIBUTING.md). O guia explica como criar Issues, enviar Pull Requests e seguir as melhores práticas para contribuir de forma eficiente.

---

## Licença

Este projeto está licenciado sob a licença Apache v2. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

#  Passo a passo aula de sdk-manual

1. Criar a repositório `/internal/telemetry` e depois criar o arquivo `/internal/telemetry/otel.go`
2. Ir até o `cmd/main.go` e fazer a chamada para o opentelemetry funcionar, acrescentar o script abaixo no inicio do código:
```
func main() {
	otelconfigFlag := flag.String("otel", "", "path to the config file")
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	closer, err := telemetry.Setup(context.Background(), *otelConfigFlag)
	if err != nil {
		panic(err)
	}
	defer closer(context.Background())
	c, _ := config.LoadConfig(*configFlag)
```
3. Para fazer o programa funcionar é necessário executar primeiro `nats-server -D -js` em seguida baixe o nats cli `go install github.com/nats-io/natscli/nats@latest` so entao execute o `go run ./cmd/all-in-one/`
4. Se tiver algum problema com o nats talvez seja necessário ajustar sua configuração do go:
```
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```
5. Executar `nats -s localhost:4222 stream create payments --subjects "payment.process" --storage memory --replicas 1 --retention=limits --discard=old --max-msgs 1_000_000 --max-msgs-per-subject 100_000 --max-bytes 4GiB --max-age 1d --max-msg-size 10MiB --dupe-window 2m --allow-rollup --no-deny-delete --no-deny-purge` e testar fazer o run `go run ./cmd/all-in-one/`
6. Testar `curl localhost:8080/payments`
