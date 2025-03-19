# chaos-shortener

URL shortener pet app with failure simulation tooling

```mermaid
%%{init: {"flowchart": {"defaultRenderer": "elk"}} }%%
flowchart LR
    CLIENT1@{shape: rounded, label: Client}
    INGRESS@{shape: rounded, label: Ingress}
    SVC@{shape: rounded, label: Service}
    SHORT@{shape: processes, label: URL Shortener}
    PGB@{shape: rounded, label: PgBouncer}
    DB@{shape: cyl, label: PostgreSQL}
    PROM@{shape: rounded, label: Prometheus}
    GRAF@{shape: rounded, label: Grafana}

    CLIENT1 -- GET /r/{id} --> INGRESS
    CLIENT1 -- POST /shorten --> INGRESS

    subgraph K8S[Kubernetes Cluster]
        PGB --> DB
        INGRESS --> SVC
        SVC --> SHORT
        SHORT --> PGB

        PROM -. /metrics .-> SVC
        GRAF -.-> PROM
    end
```

## Metrics

Chaos shortener service exposes the following Prometheus metrics:

| Name                          | Type      | Unit    | Description                            |
| ----------------------------- | --------- | ------- | -------------------------------------- |
| http_request_duration_seconds | Histogram | seconds | Duration of HTTP requests              |
| http_responses_total          | Counter   |         | Total number of HTTP requests          |
| shortener_urls_created_total  | Counter   |         | Total number of shortened URLs created |
| shortener_redirects_total     | Counter   |         | Total number of redirects performed    |

## Service Level Indicators

### Request latency

```promql
(
    sum(rate(http_request_duration_seconds_bucket{
        handler="$handler",
        otel_scope_name="cshort",
        le="0.1"
    }[$interval]))
    /
    sum(rate(http_request_duration_seconds_count{
        handler="$handler",
        otel_scope_name="cshort"
    }[$interval]))
) * 100
```

### Error rate

```promql
(
    sum(rate(http_responses_total{
        handler="$handler",
        http_response_status_code=~"^[123][0-9][0-9]$",
        otel_scope_name="cshort"
    }[$interval]))
    /
    sum(rate(http_responses_total{
        handler="$handler",
        otel_scope_name="cshort"
    }[$interval]))
) * 100
```

## Service Level Objectives

Both request latency and error rate SLIs should have values above 99.5%,
in other words, SLO for both SLIs is equal to 99.5%.

## Error budget

Error budget for latency SLO is defined by the following query:

```promql
clamp_min((
    (
        sum(rate(http_request_duration_seconds_bucket{
            handler="$handler",
            otel_scope_name="cshort",
            le="0.1"
        }[$interval])) * 100
        /
        sum(rate(http_request_duration_seconds_count{
            handler="$handler",
            otel_scope_name="cshort"
        }[$interval])) - $target
    )
    /
    (100 - $target)
) * 100, 0)
```

But the error budget for error rate SLO is defined using this query:

```promql
clamp_min((
    (
        sum(rate(http_responses_total{
            handler="$handler",
            http_response_status_code=~"^[123][0-9][0-9]$",
            otel_scope_name="cshort"
        }[$interval])) * 100
        /
        sum(rate(http_responses_total{
            handler="$handler",
            otel_scope_name="cshort"
        }[$interval])) - $target
    )
    /
    (100 - $target)
) * 100, 0)
```

## Alerts

There are 4 alert rules defined in total, 2 rules for each handler. First triggers when handler request latency SLI
drops less then 99.5%, and the second one triggers when hander error rate SLI drops below 99.5%. All alert rules are
defined [here](./charts/chaos-shortener/templates/prometheusrule.yaml).
