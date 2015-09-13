
## Auxiliary Components

* [confd](https://github.com/kelseyhightower/confd) - sync configuration from various k/v stores, out of bound from service.
  You should run this within the same container that runs your service, mounted on your Armor configuration.
* consul-template - same approach
https://github.com/a11r/grpc/blob/doc2/doc/naming.md

## Opinions

- Don't overwhelm with options. There should be one good way to do things.
- Components
  - Use existing proven components. Integration over implementation.
  - Pick best components: quality and performance first.
  - Simple and focused components over complex and flexible ones

## Roadmap

#### v1.0
- Service
  - [x] Web: gin
  - [x] Middleware infra: gin
  - [x] RPC: gRPC
- Middleware
  - [x] Request ID
  - [x] Request tracing
  - [x] Route metrics
- Operational
  - [x] Log: logrus
  - [x] Config: spf13/viper
  - [x] Metrics: armon/go-metrics

#### v2.0
- Developer experience
  - [x] Generator: Armor
  - [] Testing
- Resilience
  - [] Circuit breaker
- Discovery
  - [] Client side LB
  - [] Baked in HTTP client


