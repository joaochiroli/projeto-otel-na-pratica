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

ou outro tipo de arquitetura

![alt text]({59AF3D38-70FE-4A74-9E6F-F84B77D220EA}.png)

### Escalonando Collector

Nessa aula vamos falar sobre estratégias para escalonar coletores de telemetria de forma eficiente. Primeiro, abordamos como diferentes tipos de dados – logs, métricas e traces – exigem abordagens distintas. Logs, por exemplo, requerem escalonamento vertical, pois são coletados localmente em cada máquina, enquanto métricas e traces podem ser distribuídos horizontalmente. A forma de ingestão também impacta essa escolha: se os dados chegam via OTLP, podemos adicionar réplicas do coletor para lidar com a carga crescente.

Exploramos também o papel do Target Allocator no OpenTelemetry, um componente essencial para distribuir a carga entre os coletores. Ele garante que cada instância seja responsável por um subconjunto específico de métricas, evitando sobreposição e garantindo que cada série temporal seja única. O princípio do Single Writer é crucial aqui, pois evita que múltiplos coletores gravem os mesmos dados, prevenindo inconsistências.

Por fim, discutimos métricas de escalonamento automático, como o tamanho da fila de processamento e a saturação dos coletores. Uma boa prática é começar com três coletores e monitorar a carga, ajustando dinamicamente conforme necessário. Se a fila ultrapassa 60%, adicionamos mais coletores; se fica abaixo disso, reduzimos, sempre garantindo um mínimo operacional. Esse controle fino otimiza o consumo de recursos e a eficiência do sistema.

### Monitorando Collector

O conceito de filas é muito importante no collector, alguns spans podem falhar ao enviar dados pro backend, depois eles podem ser enviados novamente, mas depende você também pode alterar esse tipo de configuração. Além disso, caso você tenha uma fila muito grande pode ser que ela falhe também então o ideal é achar um número mágico que suporte as aplicações.

Nesta aula, eu mostrei como podemos monitorar o próprio OpenTelemetry Collector. Começamos configurando o Collector para que ele exporte dados de telemetria para um destino específico. Usei como base o repositório otel-call-cookbook, que traz receitas práticas de configuração. Nele, usamos um receiver que aceita protocolos GRPC e HTTP via OTLP, ouvindo nas portas 4317 e 4318. Esse receiver captura dados da aplicação em execução, mas nosso foco aqui não era a aplicação em si.

A parte realmente importante foi entender como coletar a telemetria gerada pelo próprio Collector. Para isso, configuramos três pipelines — uma para logs, outra para métricas e a terceira para rastreamentos — que usam um receiver OTLP e exportam os dados usando um debug exporter. Embora o exporter de debug não seja o mais usado em produção, ele nos serve bem para testes locais, pois permite inspecionar facilmente a saída diretamente no console.

A ideia principal foi mostrar como o próprio processo de observabilidade também pode — e deve — ser observado. Saber como o Collector se comporta nos dá visibilidade crítica sobre o que pode estar acontecendo com a instrumentação. Isso é essencial para quem trabalha com observabilidade em ambientes distribuídos e precisa garantir que tudo esteja fluindo como esperado.

Para essa parte dos estudos estamos usando este repo (https://github.com/jpkrohling/otelcol-cookbook)
