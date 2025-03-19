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

1. Request latency for `shorten` handler:

   ```promql
   (
       sum(rate(http_request_duration_seconds_bucket{
            handler="shorten",
            otel_scope_name="cshort",
            le="0.05"
        }[1h]))
        /
        sum(rate(http_request_duration_seconds_count{
            handler="shorten",
            otel_scope_name="cshort"
        }[1h]))
   ) * 100
   ```

2. Request latency for `redirect` handler:
