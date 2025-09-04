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
