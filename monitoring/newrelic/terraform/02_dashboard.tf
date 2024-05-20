#################
### Dashboard ###
#################

# Dashboard
resource "newrelic_one_dashboard" "monitoring" {
  name = "Monitoring applications with OpenTelemetry"

  page {
    name = "Metrics"

    ###############
    ### HEADERS ###
    ###############

    # GRPC client
    widget_markdown {
      title  = ""
      column = 1
      row    = 1
      width  = 3
      height = 1

      text = "## GRPC client"
    }

    # GRPC server
    widget_markdown {
      title  = ""
      column = 4
      row    = 1
      width  = 3
      height = 1

      text = "## GRPC server"
    }

    # HTTP client
    widget_markdown {
      title  = ""
      column = 7
      row    = 1
      width  = 3
      height = 1

      text = "## HTTP client"
    }

    # HTTP server
    widget_markdown {
      title  = ""
      column = 10
      row    = 1
      width  = 3
      height = 1

      text = "## HTTP server"
    }

    ##################
    ### THROUGHPUT ###
    ##################

    # Throughput (rpm)
    widget_line {
      title  = "Throughput (rpm)"
      column = 1
      row    = 3
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' TIMESERIES"
      }
    }

    # Throughput (rpm)
    widget_line {
      title  = "Throughput (rpm)"
      column = 4
      row    = 3
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' TIMESERIES"
      }
    }

    # Throughput (rpm)
    widget_line {
      title  = "Throughput (rpm)"
      column = 7
      row    = 3
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' TIMESERIES"
      }
    }

    # Throughput (rpm)
    widget_line {
      title  = "Throughput (rpm)"
      column = 10
      row    = 3
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' TIMESERIES"
      }
    }

    # Throughput per method (rpm)
    widget_line {
      title  = "Throughput per method (rpm)"
      column = 1
      row    = 6
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' FACET rpc.method TIMESERIES"
      }
    }

    # Throughput per method (rpm)
    widget_line {
      title  = "Throughput per method (rpm)"
      column = 4
      row    = 6
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' FACET rpc.method TIMESERIES"
      }
    }

    # Throughput per method (rpm)
    widget_line {
      title  = "Throughput per method (rpm)"
      column = 7
      row    = 6
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' FACET http.method TIMESERIES"
      }
    }

    # Throughput per method (rpm)
    widget_line {
      title  = "Throughput per method (rpm)"
      column = 10
      row    = 6
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' FACET http.method TIMESERIES"
      }
    }

    # Throughput per instance (rpm)
    widget_line {
      title  = "Throughput per instance (rpm)"
      column = 1
      row    = 9
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Throughput per instance (rpm)
    widget_line {
      title  = "Throughput per instance (rpm)"
      column = 4
      row    = 9
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Throughput per instance (rpm)
    widget_line {
      title  = "Throughput per instance (rpm)"
      column = 7
      row    = 9
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Throughput per instance (rpm)
    widget_line {
      title  = "Throughput per instance (rpm)"
      column = 10
      row    = 9
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' FACET k8s.pod.name TIMESERIES"
      }
    }

    ###############
    ### LATENCY ###
    ###############

    # Latency (ms)
    widget_line {
      title  = "Latency (ms)"
      column = 1
      row    = 12
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(rpc.client.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' TIMESERIES"
      }
    }

    # Latency (ms)
    widget_line {
      title  = "Latency (ms)"
      column = 4
      row    = 12
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(rpc.server.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' TIMESERIES"
      }
    }

    # Latency (ms)
    widget_line {
      title  = "Latency (ms)"
      column = 7
      row    = 12
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(http.client.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' TIMESERIES"
      }
    }

    # Latency (ms)
    widget_line {
      title  = "Latency (ms)"
      column = 10
      row    = 12
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(http.server.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' TIMESERIES"
      }
    }

    # Latency per method (ms)
    widget_line {
      title  = "Latency per method (ms)"
      column = 1
      row    = 15
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(rpc.client.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' FACET rpc.method TIMESERIES"
      }
    }

    # Latency per method (ms)
    widget_line {
      title  = "Latency per method (ms)"
      column = 4
      row    = 15
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(rpc.server.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' FACET rpc.method TIMESERIES"
      }
    }

    # Latency per method (ms)
    widget_line {
      title  = "Latency per method (ms)"
      column = 7
      row    = 15
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(http.client.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' FACET http.method TIMESERIES"
      }
    }

    # Latency per method (ms)
    widget_line {
      title  = "Latency per method (ms)"
      column = 10
      row    = 15
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(http.server.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' FACET http.method TIMESERIES"
      }
    }

    # Latency per instance (ms)
    widget_line {
      title  = "Latency per instance (ms)"
      column = 1
      row    = 18
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Latency per instance (ms)
    widget_line {
      title  = "Latency per instance (ms)"
      column = 4
      row    = 18
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(rpc.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Latency per instance (ms)
    widget_line {
      title  = "Latency per instance (ms)"
      column = 7
      row    = 18
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.client.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Latency per instance (ms)
    widget_line {
      title  = "Latency per instance (ms)"
      column = 10
      row    = 18
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(http.server.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' FACET k8s.pod.name TIMESERIES"
      }
    }

    ##################
    ### ERROR RATE ###
    ##################

    # Error rate (%)
    widget_line {
      title  = "Error rate (%)"
      column = 1
      row    = 21
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(rpc.client.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(rpc.grpc.status_code) > 0)/count(rpc.client.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' TIMESERIES"
      }
    }

    # Error rate (%)
    widget_line {
      title  = "Error rate (%)"
      column = 4
      row    = 21
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(rpc.server.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(rpc.grpc.status_code) > 0)/count(rpc.server.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' TIMESERIES"
      }
    }

    # Error rate (%)
    widget_line {
      title  = "Error rate (%)"
      column = 7
      row    = 21
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(http.client.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(http.status_code) >= 300)/count(http.client.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' TIMESERIES"
      }
    }

    # Error rate (%)
    widget_line {
      title  = "Error rate (%)"
      column = 10
      row    = 21
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(http.server.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(http.status_code) >= 300)/count(http.server.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' TIMESERIES"
      }
    }

    # Error rate per method (%)
    widget_line {
      title  = "Error rate per method (%)"
      column = 1
      row    = 24
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(rpc.client.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(rpc.grpc.status_code) > 0)/count(rpc.client.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' FACET rpc.method TIMESERIES"
      }
    }

    # Error rate per method (%)
    widget_line {
      title  = "Error rate per method (%)"
      column = 4
      row    = 24
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(rpc.server.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(rpc.grpc.status_code) > 0)/count(rpc.server.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' FACET rpc.method TIMESERIES"
      }
    }

    # Error rate per method (%)
    widget_line {
      title  = "Error rate per method (%)"
      column = 7
      row    = 24
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(http.client.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(http.status_code) >= 300)/count(http.client.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' FACET http.method TIMESERIES"
      }
    }

    # Error rate per method (%)
    widget_line {
      title  = "Error rate per method (%)"
      column = 10
      row    = 24
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(http.server.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(http.status_code) >= 300)/count(http.server.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' FACET http.method TIMESERIES"
      }
    }

    # Error rate per instance (%)
    widget_line {
      title  = "Error rate per instance (%)"
      column = 1
      row    = 27
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(rpc.client.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(rpc.grpc.status_code) > 0)/count(rpc.client.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Error rate per instance (%)
    widget_line {
      title  = "Error rate per instance (%)"
      column = 4
      row    = 27
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(rpc.server.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(rpc.grpc.status_code) > 0)/count(rpc.server.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Error rate per instance (%)
    widget_line {
      title  = "Error rate per instance (%)"
      column = 7
      row    = 27
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(http.client.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(http.status_code) >= 300)/count(http.client.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Error rate per instance (%)
    widget_line {
      title  = "Error rate per instance (%)"
      column = 10
      row    = 27
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(http.server.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND numeric(http.status_code) >= 300)/count(http.server.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpserver' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Error rate per span (rpm)
    widget_line {
      title  = "Error rate per span (rpm)"
      column = 1
      row    = 30
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT count(*) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcclient' AND otel.status_code = 'ERROR' FACET name TIMESERIES"
      }
    }

    # Error rate per span (rpm)
    widget_line {
      title  = "Error rate per span (rpm)"
      column = 4
      row    = 30
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT count(*) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' AND otel.status_code = 'ERROR' FACET name TIMESERIES"
      }
    }

    # Error rate per span (rpm)
    widget_line {
      title  = "Error rate per rpm (rpm)"
      column = 7
      row    = 30
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT count(*) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'httpclient' AND otel.status_code = 'ERROR' FACET name TIMESERIES"
      }
    }

    # Error rate per span (rpm)
    widget_line {
      title  = "Error rate per span (rpm)"
      column = 10
      row    = 30
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT count(*) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = 'kind-otel' AND service.name = 'grpcserver' AND otel.status_code = 'ERROR' FACET name TIMESERIES"
      }
    }
  }
}
