
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

### Cloud Philosophy

"Cloud" infrastructure and all-encompassing microservices. You might have bumped into Netflix' or Spring Boot Cloud,
which by now have gotten years of mileage (at least with Netflix' case).

You would have seen that a microservice in those cases is a monument of capabilities. 
Concerns such as load balancing and service discovery are the responsibility of the microservice; auto detection
of configuration and reconfiguring itself are expected from every microservice.

We feel that this overhead and critical responsibility (too critical) are symptoms of Netflix providing
a solution ahead of its time. They lived it in a container-less world. In a world where monoliths were common.

Times are changing and we believe Docker and cloud support infrastructure such as Kubernetes should solve 
these concerns and microservices should be agnostic of it.

Microservices should only interface to cloud infrastructure concerns (a great example is 12-factor). Examples:

* Logging: STDOUT only. Forget graylog, syslog, and various others.
* Metrics: expose only. Infrastructure should collect metrics from all services by itself. Don't support a ton
of metrics collectors (statsd, influx, etc.)
* Discovery: keep referring to services by their name (should it be DNS name is another question). 
Infrastructure should take care of changing landscape of microservices.
* Services, middleware: should remain as is. It is time proven.
* RPC: we *can* do better than HTTP 1.x.

With Armor, we have adopted these, although we have bypassed some of them to allow an escape hatch. For example:

* Logging supports syslog as well.
* Metrics supports statsd.

These are escape hatches that will remain as long as we don't have a perfect microservice ecosystem,
that is readily available for everyone.

With Armor we have to always walk a fine line between staying simple and not taking too much responsibility
and providing enough headroom to work in a non-ideal infrastructure.


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
  - [x] Config: olebedev/config
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


