# OpenTelemetry Collector Contrib Distro

This distribution contains all the components from both the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) repository and the [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) repository. This distribution includes open source and vendor supported components.

## Recommendation

As this distribution contains many components, it is a good starting point to try various configurations. However, when running in production, it is recommended to limit the collector to contain only the components necessary for an environment. Some reasons to do this:

- reduce the size of the collector, reducing deployment times for the collector
- improve the security of the collector by reducing the available attack surface area

Building a [custom collector](https://opentelemetry.io/docs/collector/custom-collector/) can be achieved using the [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder).

## Components

The full list of components is available in the [manifest](manifest.yaml)

### Rules for Component Inclusion

- Include all extensions at [Alpha stability](https://github.com/open-telemetry/opentelemetry-collector#alpha) or higher and pipeline components that have at least 1 signal at [Alpha stability](https://github.com/open-telemetry/opentelemetry-collector#alpha) or higher.

Eu posso baixar o otel na versão desejada se eu quiser por exemplo:

```
wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.120.0/otelcol-contrib_0.120.0_linux_amd64.tar.gz

## Descompactar o arquivo
tar xzf otelcol-contrib_0.120.0_linux_amd64.tar.gz

## Executar o binário
./otelcol-contrib

## Mover o binário do otel-contrib para o diretório /bin
 mv otelcol-contrib /bin/ ou mv otelcol-contrib ~/bin/

## Eu posso testar usando o comando a seguir, passando o arquivo de configuração
otelcol-contrib --config simple.yaml

## Depois você pode executar um comando qualquer em outro cmd pra ver os dados chegando no seu colector
go run ./

ou

go run ./cmd/all-in-one/
```

### Agent vs Gateways vs Coletores

**OpenTelemetry Agent**
O Agent é um componente que roda como um processo separado no mesmo host da aplicação:
Características:

Executa como sidecar ou daemon no nó/container
Coleta telemetria localmente da aplicação
Baixa latência entre app e agent
Processamento local dos dados

- Vantagens:

  - Menor impacto na aplicação (offloading de processamento)
  - Resiliência local (buffer local caso backend falhe)
  - Configuração centralizada por nó

- Desvantagens:

  - Overhead de recursos por nó
  - Gerenciamento de configuração distribuído

**OpenTelemetry Gateway**
O Gateway é um componente centralizado que atua como proxy entre múltiplas fontes e backends:
Características:

Instância centralizada recebendo de múltiplos agents/apps
Agregação e roteamento de telemetria
Ponto único de controle e política
Load balancing para backends

- Vantagens:

  - Redução de conexões diretas aos backends
  - Controle centralizado de políticas
  - Agregação e correlação de dados
  - Economia de recursos em escala

- Desvantagens:

  - Ponto único de falha (precisa ser redundante)
  - Maior latência na pipeline

**OpenTelemetry Collector**
O Collector é o componente core que pode funcionar tanto como Agent quanto Gateway:

Arquitetura:
Receivers → Processors → Exporters
Receivers: Coletam dados (OTLP, Jaeger, Zipkin, Prometheus)
Processors: Transformam/filtram dados (sampling, batching)
Exporters: Enviam para backends (Jaeger, Prometheus, ELK)
Padrões de Deployment

1. Agent Pattern
   App → OTel Collector (Agent) → Backend
2. Gateway Pattern
   App → OTel Collector (Agent) → OTel Collector (Gateway) → Backend
3. Direct Pattern
   App → Backend (usando SDK direto)
   Quando usar cada um?
   Use Agent quando:

Aplicações em containers/K8s
Necessita processamento local
Quer desacoplar app do backend

Use Gateway quando:

Ambiente multi-tenant
Necessita controle centralizado
Muitas aplicações → poucos backends

Use Collector quando:

Necessita flexibilidade (pode ser agent ou gateway)
Quer pipeline configurável
Padrão vendor-neutral

![alt text]({A902C855-7374-43F7-93D2-3D995114A0BE}.png)
