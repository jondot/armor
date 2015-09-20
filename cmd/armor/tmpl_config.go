package main

const TEMPL_CONF = `

#====================================================================
# Your own hierarchical config
#
# Examples:
#   arm.Config.GetInt("foo.bar")
#   arm.Config.GetString("foo.baz")
#====================================================================
foo:
  bar: 1
  baz: yep!



#====================================================================
# Service
#
# Below are service level / operational knobs.
#====================================================================
#shafile: REVISION
#port: 6060
#interface: 0.0.0.0



#====================================================================
# RPC - we do RPC via gRPC and protocolbuffers (http://grpc.io)
#
# Use RPC in addition to HTTP in cases where you want service-to-
# service communication (microservices), oplog/auditlog, or need to
# handle specialized performance-oriented use cases.
#====================================================================
#rpc:
#  port: 41555
#  interface: 0.0.0.0



#====================================================================
# Middleware
#
# These are the list of prebaked middleware, you can mix and match
# by listing or not listing them here.
#====================================================================
middleware:
  - request_identification
  - request_tracing
  - route_metrics



#====================================================================
# Metrics
# You can have metrics sent to multiple sinks.
#
# See github.com/armon/go-metrics for types of sinks
#====================================================================
#metrics:
#  prefix: dev
#  inmem:
#    interval_secs: 10
#    retain_mins: 1
#  statsd:
#    server: localhost:514



#====================================================================
# Logging
#
# Along with custom hooks, you can config formatting, levels and
# where to dump more critical errors via hooks.
#====================================================================
log:
  console:
    level: error
#  hooks:
#    syslog:
#      transport: udp
#      server: localhost:514
#      level: debug


`

